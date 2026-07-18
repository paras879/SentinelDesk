package deviceid

import (
	"crypto/sha1"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var deviceID string

func Init() error {
	// Derive a stable UUID from this machine's primary MAC address.
	// This is done first so we know what the "correct" ID for this machine is.
	machineID := macDerivedUUID()

	// Check if a persisted .device_id file exists.
	data, err := os.ReadFile(".device_id")
	if err == nil {
		id := strings.TrimSpace(string(data))
		if len(id) == 36 {
			// If the persisted ID matches this machine's MAC-derived ID, use it.
			// If it doesn't match (file was copied from another machine), regenerate.
			if id == machineID {
				deviceID = id
				log.Println("Device ID loaded from .device_id file")
				return nil
			}
			log.Println("Warning: .device_id does not match this machine's MAC — regenerating")
		}
	}

	// Persist the MAC-derived ID so subsequent startups are instant.
	if err := os.WriteFile(".device_id", []byte(machineID), 0644); err != nil {
		log.Println("Warning: could not save .device_id:", err)
	}

	deviceID = machineID
	log.Println("Device ID generated:", machineID)
	return nil
}

func Get() string {
	return deviceID
}

// macDerivedUUID produces a deterministic UUID v5 from the machine's primary
// MAC address. The same MAC always yields the same UUID on any run.
func macDerivedUUID() string {
	mac := primaryMAC()

	// UUID v5 namespace for MAC addresses (custom, fixed namespace).
	namespace := []byte{0x6b, 0xa7, 0xb8, 0x11, 0x9d, 0xad, 0x11, 0xd1,
		0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	h := sha1.New()
	h.Write(namespace)
	h.Write([]byte(mac))
	sum := h.Sum(nil)

	// Set version (5) and variant bits per RFC 4122.
	sum[6] = (sum[6] & 0x0f) | 0x50
	sum[8] = (sum[8] & 0x3f) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		sum[0:4], sum[4:6], sum[6:8], sum[8:10], sum[10:16])
}

// primaryMAC returns the hardware address of the first non-loopback,
// non-virtual network interface that has a MAC address.
func primaryMAC() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "fallback-no-interfaces"
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if len(iface.HardwareAddr) == 0 {
			continue
		}
		return iface.HardwareAddr.String()
	}
	return "fallback-no-mac"
}
