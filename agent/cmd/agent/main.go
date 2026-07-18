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

func main() {

	config.Load()

	log.Println("SentinelDesk Agent Started")

	if err := deviceid.Init(); err != nil {
		log.Fatal("Failed to initialize device ID:", err)
	}

	log.Println("Device ID :", deviceid.Get())

	info, err := system.GetSystemInfo()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Device Name :", info.DeviceName)
	log.Println("Hostname    :", info.Hostname)
	log.Println("Username    :", info.Username)
	log.Println("OS          :", info.OS)
	log.Println("IP Address  :", info.IPAddress)
	log.Println("MAC Address :", info.MACAddress)
	log.Println("Subnet      :", info.ConnectedSubnet)

	// ====================================================
	// Phase 1: Register Device
	// Retry every 5 seconds until the backend confirms
	// the device record exists.
	// ====================================================
	for {
		if err := register.RegisterDevice(); err != nil {
			log.Println("Registration failed, retrying in 5 seconds:", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	log.Println("Device registration completed.")

	interval := config.GetHeartbeatInterval()

	log.Printf("Service Interval : %d seconds\n", interval)
	log.Println("Starting services...")

	// ====================================================
	// Phase 2: Start all background services
	// Only reached after registration succeeds.
	// ====================================================

	// Heartbeat immediately + periodic
	if err := heartbeat.SendHeartbeat(); err != nil {
		log.Println("Heartbeat Error:", err)
	} else {
		log.Println("Heartbeat sent successfully.")
	}

	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if err := heartbeat.SendHeartbeat(); err != nil {
				log.Println("Heartbeat Error:", err)
			} else {
				log.Println("Heartbeat sent successfully.")
			}
		}
	}()

	// System Info immediately + periodic
	sendSystemInfo()
	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			sendSystemInfo()
		}
	}()

	// Software Inventory immediately + periodic
	sendSoftwareInventory()
	go func() {
		softwareTicker := time.NewTicker(10 * time.Minute)
		defer softwareTicker.Stop()

		for range softwareTicker.C {
			sendSoftwareInventory()
		}
	}()

	// Process List immediately + periodic
	sendProcesses()
	go func() {
		processTicker := time.NewTicker(60 * time.Second)
		defer processTicker.Stop()

		for range processTicker.C {
			sendProcesses()
		}
	}()

	// Windows Services immediately + periodic
	sendWindowsServices()
	go func() {
		serviceTicker := time.NewTicker(5 * time.Minute)
		defer serviceTicker.Stop()

		for range serviceTicker.C {
			sendWindowsServices()
		}
	}()

	// Service Command Polling immediately + periodic
	pollServiceCommands()
	go func() {
		pollTicker := time.NewTicker(15 * time.Second)
		defer pollTicker.Stop()

		for range pollTicker.C {
			pollServiceCommands()
		}
	}()

	// Live View Streaming
	go liveview.NewStreamer().Start()

	select {}
}

func sendWindowsServices() {
	svcs, err := windowsservice.Collect()
	if err != nil {
		log.Println("WindowsServices Collection Error:", err)
		return
	}

	if err := windowsservice.SendServices(svcs); err != nil {
		log.Println("WindowsServices Send Error:", err)
	} else {
		log.Printf("Windows services uploaded: %d services.", len(svcs))
	}
}

func pollServiceCommands() {
	windowsservice.PollAndExecuteCommands()
}

func sendProcesses() {
	procs, err := process.Collect()
	if err != nil {
		log.Println("ProcessInventory Collection Error:", err)
		return
	}

	if err := process.SendProcesses(procs); err != nil {
		log.Println("ProcessInventory Send Error:", err)
	} else {
		log.Printf("Process inventory sent: %d processes.", len(procs))
	}
}

func sendSoftwareInventory() {
	apps, err := software.Collect()
	if err != nil {
		log.Println("SoftwareInventory Collection Error:", err)
		return
	}

	if err := software.SendSoftware(apps); err != nil {
		log.Println("SoftwareInventory Send Error:", err)
	} else {
		log.Printf("Software inventory sent: %d applications.", len(apps))
	}
}

func sendSystemInfo() {
	info, err := systeminfo.Collect()
	if err != nil {
		log.Println("SystemInfo Collection Error:", err)
		return
	}

	if err := systeminfo.SendSystemInfo(info); err != nil {
		log.Println("SystemInfo Send Error:", err)
	} else {
		log.Println("System information sent successfully.")
	}
}
