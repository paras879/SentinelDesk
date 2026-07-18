package main

import (
	"log"
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

func main() {
	log.SetPrefix(logPrefix)

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
		if err := register.RegisterDevice(); err != nil {
			log.Println("[ERROR] Registration failed, retrying in 5 seconds:", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	log.Println("[SUCCESS] Device registered")

	interval := config.GetHeartbeatInterval()
	log.Printf("[INFO] Heartbeat interval: %d seconds", interval)

	log.Println("[INFO] Starting heartbeat...")
	if err := heartbeat.SendHeartbeat(); err != nil {
		log.Println("[ERROR] Initial heartbeat failed:", err)
	} else {
		log.Println("[SUCCESS] Heartbeat sent")
	}

	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if err := heartbeat.SendHeartbeat(); err != nil {
				log.Println("[ERROR] Heartbeat failed:", err)
			}
		}
	}()

	log.Println("[INFO] Uploading system info...")
	sendSystemInfo()

	log.Println("[INFO] Uploading software inventory...")
	sendSoftwareInventory()

	log.Println("[INFO] Uploading process inventory...")
	sendProcesses()

	log.Println("[INFO] Uploading Windows services...")
	sendWindowsServices()

	go func() {
		systemInfoTicker := time.NewTicker(time.Duration(interval) * time.Second)
		defer systemInfoTicker.Stop()
		for range systemInfoTicker.C {
			sendSystemInfo()
		}
	}()

	go func() {
		softwareTicker := time.NewTicker(10 * time.Minute)
		defer softwareTicker.Stop()
		for range softwareTicker.C {
			sendSoftwareInventory()
		}
	}()

	go func() {
		processTicker := time.NewTicker(60 * time.Second)
		defer processTicker.Stop()
		for range processTicker.C {
			sendProcesses()
		}
	}()

	go func() {
		serviceTicker := time.NewTicker(5 * time.Minute)
		defer serviceTicker.Stop()
		for range serviceTicker.C {
			sendWindowsServices()
		}
	}()

	go func() {
		pollTicker := time.NewTicker(15 * time.Second)
		defer pollTicker.Stop()
		for range pollTicker.C {
			windowsservice.PollAndExecuteCommands()
		}
	}()

	go liveview.NewStreamer().Start()

	log.Println("[INFO] All services started. Agent is running.")
	select {}
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
