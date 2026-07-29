//go:build !windows

package windowsservice

func Collect() ([]ServiceInfo, error) {
	// Not supported on non-Windows platforms
	return []ServiceInfo{}, nil
}
