package systeminfo

import (
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"

	"sentineldesk/agent/internal/deviceid"
)

type SystemInfo struct {
	DeviceID     string `json:"device_id"`
	Hostname     string `json:"hostname"`
	CPUName      string `json:"cpu_name"`
	CPUUsage     float64 `json:"cpu_usage"`
	TotalRAM     uint64 `json:"total_ram"`
	UsedRAM      uint64 `json:"used_ram"`
	FreeRAM      uint64 `json:"free_ram"`
	TotalDisk    uint64 `json:"total_disk"`
	UsedDisk     uint64 `json:"used_disk"`
	FreeDisk     uint64 `json:"free_disk"`
	OSVersion    string `json:"os_version"`
	LoggedInUser string `json:"logged_in_user"`
	LocalIP      string `json:"local_ip"`
	PublicIP     string `json:"public_ip"`
	MACAddress   string `json:"mac_address"`
	Uptime       int64  `json:"uptime"`
	LastBoot     time.Time `json:"last_boot"`
	AgentVersion string `json:"agent_version"`
}

func Collect() (*SystemInfo, error) {

	hostInfo, _ := host.Info()

	vMem, _ := mem.VirtualMemory()

	hostname := ""
	osVersion := ""
	uptime := int64(0)
	lastBoot := time.Time{}

	if hostInfo != nil {
		hostname = hostInfo.Hostname
		osVersion = hostInfo.Platform + " " + hostInfo.PlatformVersion
		uptime = int64(hostInfo.Uptime)
		lastBoot = time.Unix(int64(hostInfo.BootTime), 0)
	}

	totalRAM, usedRAM, freeRAM := uint64(0), uint64(0), uint64(0)
	if vMem != nil {
		totalRAM = vMem.Total
		usedRAM = vMem.Used
		freeRAM = vMem.Free
	}

	totalDisk, usedDisk, freeDisk := getDiskUsage()

	info := &SystemInfo{
		DeviceID:     deviceid.Get(),
		Hostname:     hostname,
		CPUName:      getCPUName(),
		CPUUsage:     getCPUUsage(),
		TotalRAM:     totalRAM,
		UsedRAM:      usedRAM,
		FreeRAM:      freeRAM,
		TotalDisk:    totalDisk,
		UsedDisk:     usedDisk,
		FreeDisk:     freeDisk,
		OSVersion:    osVersion,
		LoggedInUser: GetLoggedInUser(),
		LocalIP:      GetLocalIP(),
		PublicIP:     GetPublicIP(),
		MACAddress:   GetMACAddress(),
		Uptime:       uptime,
		LastBoot:     lastBoot,
		AgentVersion: "1.0.0",
	}

	return info, nil
}

func getCPUName() string {
	info, err := cpu.Info()
	if err != nil || len(info) == 0 {
		return ""
	}
	return info[0].ModelName
}

func getCPUUsage() float64 {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil || len(percent) == 0 {
		return 0
	}
	return percent[0]
}

func getDiskUsage() (total, used, free uint64) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return 0, 0, 0
	}

	for _, p := range partitions {
		if runtime.GOOS == "windows" && !strings.HasPrefix(p.Mountpoint, "C:") {
			continue
		}

		d, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}
		return d.Total, d.Used, d.Free
	}

	if len(partitions) > 0 {
		d, err := disk.Usage(partitions[0].Mountpoint)
		if err == nil {
			return d.Total, d.Used, d.Free
		}
	}

	return 0, 0, 0
}
