//go:build !windows

package main

import "fmt"

func runService() error {
	return fmt.Errorf("Windows Services are not supported on this OS")
}

func installService() error {
	return fmt.Errorf("Windows Services are not supported on this OS")
}

func removeService() error {
	return fmt.Errorf("Windows Services are not supported on this OS")
}

func startService() error {
	return fmt.Errorf("Windows Services are not supported on this OS")
}

func stopService() error {
	return fmt.Errorf("Windows Services are not supported on this OS")
}
