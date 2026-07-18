package process

import (
	"sort"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

type ProcessInfo struct {
	PID             int32   `json:"pid"`
	Name            string  `json:"name"`
	ExecutablePath  string  `json:"executable_path"`
	CPUUsage        float64 `json:"cpu_usage"`
	MemoryBytes     uint64  `json:"memory_bytes"`
	MemoryPercent   float32 `json:"memory_percent"`
	ThreadCount     int32   `json:"thread_count"`
	HandleCount     int32   `json:"handle_count"`
	StartTime       string  `json:"start_time"`
	Username        string  `json:"username"`
	Status          string  `json:"status"`
}

func Collect() ([]ProcessInfo, error) {

	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	all := make([]ProcessInfo, 0, len(procs))

	for _, p := range procs {

		info := ProcessInfo{
			PID:    p.Pid,
			Status: "Running",
		}

		name, err := p.Name()
		if err == nil {
			info.Name = name
		}

		exe, err := p.Exe()
		if err == nil {
			info.ExecutablePath = exe
		}

		cpu, err := p.CPUPercent()
		if err == nil {
			info.CPUUsage = cpu
		}

		memInfo, err := p.MemoryInfo()
		if err == nil && memInfo != nil {
			info.MemoryBytes = memInfo.RSS
		}

		memPercent, err := p.MemoryPercent()
		if err == nil {
			info.MemoryPercent = float32(memPercent)
		}

		threads, err := p.NumThreads()
		if err == nil {
			info.ThreadCount = threads
		}

		createTime, err := p.CreateTime()
		if err == nil {
			info.StartTime = time.Unix(createTime/1000, (createTime%1000)*int64(time.Millisecond)).Format(time.RFC3339)
		}

		username, err := p.Username()
		if err == nil {
			info.Username = username
		}

		all = append(all, info)
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].Name < all[j].Name
	})

	return all, nil
}
