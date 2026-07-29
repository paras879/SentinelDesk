//go:build !windows

package software

func Collect() ([]Software, error) {
	// Not supported on non-Windows
	return []Software{}, nil
}
