package windowsservice

import (
	"sort"

	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

type ServiceInfo struct {
	ServiceName    string `json:"service_name"`
	DisplayName    string `json:"display_name"`
	Status         string `json:"status"`
	StartType      string `json:"start_type"`
	ExecutablePath string `json:"executable_path"`
	PID            int32  `json:"pid"`
	ServiceAccount string `json:"service_account"`
	Description    string `json:"description"`
	CanStop        bool   `json:"can_stop"`
	CanPause       bool   `json:"can_pause"`
	AcceptShutdown bool   `json:"accept_shutdown"`
}

func Collect() ([]ServiceInfo, error) {

	m, err := mgr.Connect()
	if err != nil {
		return nil, err
	}
	defer m.Disconnect()

	names, err := m.ListServices()
	if err != nil {
		return nil, err
	}

	all := make([]ServiceInfo, 0, len(names))

	for _, name := range names {

		s, err := m.OpenService(name)
		if err != nil {
			continue
		}

		config, err := s.Config()
		if err != nil {
			s.Close()
			continue
		}

		status, err := s.Query()
		if err != nil {
			s.Close()
			continue
		}

		info := ServiceInfo{
			ServiceName:    name,
			DisplayName:    config.DisplayName,
			Status:         stateToString(status.State),
			StartType:      startTypeToString(config.StartType, name),
			ExecutablePath: config.BinaryPathName,
			PID:            0,
			ServiceAccount: config.ServiceStartName,
			Description:    config.Description,
			CanStop:        status.Accepts&svc.AcceptStop != 0,
			CanPause:       status.Accepts&svc.AcceptPauseAndContinue != 0,
			AcceptShutdown: status.Accepts&svc.AcceptShutdown != 0,
		}

		all = append(all, info)
		s.Close()
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].DisplayName < all[j].DisplayName
	})

	return all, nil
}

func stateToString(state svc.State) string {
	switch state {
	case svc.Running:
		return "Running"
	case svc.Stopped:
		return "Stopped"
	case svc.Paused:
		return "Paused"
	case svc.StartPending:
		return "Start Pending"
	case svc.StopPending:
		return "Stop Pending"
	case svc.PausePending:
		return "Pause Pending"
	case svc.ContinuePending:
		return "Continue Pending"
	default:
		return "Unknown"
	}
}

func startTypeToString(startType uint32, serviceName string) string {
	switch startType {
	case mgr.StartAutomatic:
		if isDelayedAutoStart(serviceName) {
			return "Automatic (Delayed)"
		}
		return "Automatic"
	case mgr.StartManual:
		return "Manual"
	case mgr.StartDisabled:
		return "Disabled"
	default:
		return "Unknown"
	}
}

func isDelayedAutoStart(serviceName string) bool {

	path := `SYSTEM\CurrentControlSet\Services\` + serviceName

	k, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		path,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return false
	}
	defer k.Close()

	val, _, err := k.GetIntegerValue("DelayedAutostart")
	if err != nil {
		return false
	}

	return val == 1
}
