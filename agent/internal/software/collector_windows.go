//go:build windows

package software

import (
	"sort"

	"golang.org/x/sys/windows/registry"
)

func Collect() ([]Software, error) {

	all := make([]Software, 0)
	seen := make(map[string]bool)

	collectFromKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
		&all, seen,
	)

	collectFromKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`,
		&all, seen,
	)

	collectFromKey(
		registry.CURRENT_USER,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
		&all, seen,
	)

	sort.Slice(all, func(i, j int) bool {
		return all[i].Name < all[j].Name
	})

	return all, nil
}

func collectFromKey(root registry.Key, path string, all *[]Software, seen map[string]bool) {

	k, err := registry.OpenKey(root, path, registry.ENUMERATE_SUB_KEYS|registry.READ)
	if err != nil {
		return
	}
	defer k.Close()

	subKeys, err := k.ReadSubKeyNames(-1)
	if err != nil {
		return
	}

	for _, subKey := range subKeys {

		sk, err := registry.OpenKey(k, subKey, registry.READ)
		if err != nil {
			continue
		}

		displayName, _, err := sk.GetStringValue("DisplayName")
		if err != nil || displayName == "" {
			sk.Close()
			continue
		}

		if seen[displayName] {
			sk.Close()
			continue
		}
		seen[displayName] = true

		version, _, _ := sk.GetStringValue("DisplayVersion")
		publisher, _, _ := sk.GetStringValue("Publisher")
		installDate, _, _ := sk.GetStringValue("InstallDate")
		installLocation, _, _ := sk.GetStringValue("InstallLocation")

		estimatedSize := int64(0)
		size, _, err := sk.GetIntegerValue("EstimatedSize")
		if err == nil {
			estimatedSize = int64(size) * 1024
		}

		*all = append(*all, Software{
			Name:            displayName,
			Version:         version,
			Publisher:       publisher,
			InstallDate:     installDate,
			InstallLocation: installLocation,
			EstimatedSize:   estimatedSize,
		})

		sk.Close()
	}
}
