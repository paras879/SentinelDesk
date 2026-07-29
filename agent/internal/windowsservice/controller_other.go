//go:build !windows

package windowsservice

func startService(name string) string {
	return "unsupported on this OS"
}

func stopService(name string) string {
	return "unsupported on this OS"
}

func restartService(name string) string {
	return "unsupported on this OS"
}
