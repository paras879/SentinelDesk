package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"syscall"
	"time"

	"sentineldesk/agent/internal/config"
	"sentineldesk/agent/internal/deviceid"
	"sentineldesk/agent/internal/heartbeat"
	"sentineldesk/agent/internal/liveview"
	"sentineldesk/agent/internal/process"
	"sentineldesk/agent/internal/register"
	"sentineldesk/agent/internal/software"
	"sentineldesk/agent/internal/system"
	"sentineldesk/agent/internal/systeminfo"
	"sentineldesk/agent/internal/windowsservice"
)
const logPrefix = "[SentinelDesk Agent] "

var logFile *os.File

func main() {
	runtime.GOMAXPROCS(1)
	debug.SetGCPercent(20)

	log.SetFlags(log.Ldate | log.Ltime | log.Lmsgprefix)
	log.SetPrefix(logPrefix)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			if err := installService(); err != nil {
				log.Fatalf("Install failed: %v", err)
			}
			log.Println("Service installed successfully")
			return
		case "uninstall":
			if err := removeService(); err != nil {
				log.Fatalf("Uninstall failed: %v", err)
			}
			log.Println("Service uninstalled successfully")
			return
		case "start":
			if err := startService(); err != nil {
				log.Fatalf("Start failed: %v", err)
			}
			log.Println("Service started successfully")
			return
		case "stop":
			if err := stopService(); err != nil {
				log.Fatalf("Stop failed: %v", err)
			}
			log.Println("Service stopped successfully")
			return
		case "run":
			if err := runService(); err != nil {
				log.Fatalf("Service run failed: %v", err)
			}
			return
		}
	}

	if isSvc, err := isWindowsService(); err == nil && isSvc {
		if err := runService(); err != nil {
			log.Fatalf("Service run failed: %v", err)
		}
		return
	}

	setupFileLogging()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("[INFO] Received shutdown signal")
		cancel()
	}()

	runAgent(ctx)
}

func setupFileLogging() {
	progData := os.Getenv("PROGRAMDATA")
	if progData == "" {
		progData = filepath.Join(os.Getenv("SystemDrive"), "\\ProgramData")
	}
	logDir := filepath.Join(progData, "SentinelDesk", "agent", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("[WARN] Failed to create log directory %s: %v", logDir, err)
		return
	}

	logPath := filepath.Join(logDir, "agent.log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("[WARN] Failed to open log file %s: %v", logPath, err)
		return
	}
	logFile = f
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.Printf("[INFO] Logging to %s", logPath)
}

func runAgent(ctx context.Context) {
	log.Println("[INFO] Loading configuration...")
	config.Load()
	cfg := config.Get()
	log.Println("[INFO] Backend:", cfg.ServerURL)

	log.Println("[INFO] Loading Device UUID...")
	if err := deviceid.Init(); err != nil {
		log.Fatal("[FATAL] Failed to initialize device ID:", err)
	}
	log.Println("[SUCCESS] Device UUID loaded:", deviceid.Get())

	sysInfo, err := system.GetSystemInfo()
	if err != nil {
		log.Fatal("[FATAL] Failed to get system info:", err)
	}

	log.Println("[INFO] Device Name:", sysInfo.DeviceName)
	log.Println("[INFO] Hostname:", sysInfo.Hostname)
	log.Println("[INFO] Username:", sysInfo.Username)
	log.Println("[INFO] OS:", sysInfo.OS)
	log.Println("[INFO] OS Version:", sysInfo.OSVersion)
	log.Println("[INFO] IP Address:", sysInfo.IPAddress)
	log.Println("[INFO] MAC Address:", sysInfo.MACAddress)
	log.Println("[INFO] Subnet:", sysInfo.ConnectedSubnet)

	log.Println("[INFO] Registering device...")
	for {
		select {
		case <-ctx.Done():
			log.Println("[INFO] Registration cancelled during shutdown")
			return
		default:
		}
		if err := register.RegisterDevice(); err != nil {
			log.Println("[ERROR] Registration failed, retrying in 5 seconds:", err)
			select {
			case <-ctx.Done():
				log.Println("[INFO] Registration cancelled during shutdown")
				return
			case <-time.After(5 * time.Second):
			}
			continue
		}
		break
	}
	log.Println("[SUCCESS] Device registered")

	interval := config.GetHeartbeatInterval()
	log.Printf("[INFO] Heartbeat interval: %d seconds", interval)

	log.Println("[INFO] Starting WebSocket heartbeat...")
	go heartbeat.StartHeartbeatLoop(ctx, interval)

	go runSystemInfoLoop(ctx, interval)
	go runSoftwareInventoryLoop(ctx)
	go runProcessInventoryLoop(ctx)
	go runWindowsServicesLoop(ctx)
	go runCommandPollLoop(ctx)

	go liveview.NewStreamer().Start()

	log.Println("[INFO] All services started. Agent is running.")
	<-ctx.Done()

	log.Println("[INFO] Agent shutting down gracefully...")
	if logFile != nil {
		logFile.Sync()
	}
	time.Sleep(500 * time.Millisecond)
}

func runSystemInfoLoop(ctx context.Context, interval int) {
	sendSystemInfo()
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sendSystemInfo()
		}
	}
}

func runSoftwareInventoryLoop(ctx context.Context) {
	sendSoftwareInventory()
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sendSoftwareInventory()
		}
	}
}

func runProcessInventoryLoop(ctx context.Context) {
	sendProcesses()
	ticker := time.NewTicker(120 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sendProcesses()
		}
	}
}

func runWindowsServicesLoop(ctx context.Context) {
	sendWindowsServices()
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sendWindowsServices()
		}
	}
}

func runCommandPollLoop(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			windowsservice.PollAndExecuteCommands()
		}
	}
}

func sendSystemInfo() {
	info, err := systeminfo.Collect()
	if err != nil {
		log.Println("[ERROR] SystemInfo collection failed:", err)
		return
	}
	if err := systeminfo.SendSystemInfo(info); err != nil {
		log.Println("[ERROR] SystemInfo upload failed:", err)
	} else {
		log.Println("[SUCCESS] System info uploaded")
	}
}

func sendSoftwareInventory() {
	apps, err := software.Collect()
	if err != nil {
		log.Println("[ERROR] SoftwareInventory collection failed:", err)
		return
	}
	if err := software.SendSoftware(apps); err != nil {
		log.Println("[ERROR] SoftwareInventory upload failed:", err)
	} else {
		log.Printf("[SUCCESS] Software inventory uploaded: %d applications", len(apps))
	}
}

func sendProcesses() {
	procs, err := process.Collect()
	if err != nil {
		log.Println("[ERROR] ProcessInventory collection failed:", err)
		return
	}
	if err := process.SendProcesses(procs); err != nil {
		log.Println("[ERROR] ProcessInventory upload failed:", err)
	} else {
		log.Printf("[SUCCESS] Process inventory uploaded: %d processes", len(procs))
	}
}

func sendWindowsServices() {
	svcs, err := windowsservice.Collect()
	if err != nil {
		log.Println("[ERROR] WindowsServices collection failed:", err)
		return
	}
	if err := windowsservice.SendServices(svcs); err != nil {
		log.Println("[ERROR] WindowsServices upload failed:", err)
	} else {
		log.Printf("[SUCCESS] Windows services uploaded: %d services", len(svcs))
	}
}
