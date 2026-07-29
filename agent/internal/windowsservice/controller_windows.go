//go:build windows

package windowsservice

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func startService(name string) string {

	m, err := mgr.Connect()
	if err != nil {
		return err.Error()
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		return err.Error()
	}
	defer s.Close()

	err = s.Start()
	if err != nil {
		return err.Error()
	}

	return ""
}

func stopService(name string) string {

	m, err := mgr.Connect()
	if err != nil {
		return err.Error()
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		return err.Error()
	}
	defer s.Close()

	_, err = s.Control(svc.Stop)
	if err != nil {
		return err.Error()
	}

	return ""
}

func restartService(name string) string {

	m, err := mgr.Connect()
	if err != nil {
		return err.Error()
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		return err.Error()
	}
	defer s.Close()

	_, err = s.Control(svc.Stop)
	if err != nil {
		return err.Error()
	}

	if err := waitForState(s, svc.Stopped, 30*time.Second); err != nil {
		return err.Error()
	}

	err = s.Start()
	if err != nil {
		return err.Error()
	}

	return ""
}

// waitForState polls the service state until it reaches the desired state or the timeout expires.
func waitForState(s *mgr.Service, want svc.State, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		status, err := s.Query()
		if err != nil {
			return err
		}
		if status.State == want {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("service did not reach state %v within %v", want, timeout)
}
