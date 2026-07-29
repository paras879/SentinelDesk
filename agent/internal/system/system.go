package system

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v3/host"
)

type NetworkAdapter struct {
	Name        string `json:"name"`
	MACAddress  string `json:"mac_address"`
	IPAddresses string `json:"ip_addresses"`
	SubnetMask  string `json:"subnet_mask"`
	IsPrimary   bool   `json:"is_primary"`
}

type SystemInfo struct {
	DeviceName      string
	Hostname        string
	Username        string
	OS              string
	OSVersion       string
	IPAddress       string
	MACAddress      string
	NetworkAdapters []NetworkAdapter
	ConnectedSubnet string
	DefaultGateway  string
	NetworkGroupID  string
}

func GetSystemInfo() (*SystemInfo, error) {

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}

	ip, subnet, adapters := getNetworkInfo()
	gw := getDefaultGateway()
	groupID := generateNetworkGroupID()

	info := &SystemInfo{
		DeviceName:      hostname,
		Hostname:        hostname,
		Username:        currentUser.Username,
		OS:              runtime.GOOS,
		OSVersion:       getOSVersion(),
		IPAddress:       ip,
		MACAddress:      getPrimaryMAC(adapters),
		NetworkAdapters: adapters,
		ConnectedSubnet: subnet,
		DefaultGateway:  gw,
		NetworkGroupID:  groupID,
	}

	return info, nil
}

// getDefaultGateway returns the default gateway IP of the selected adapter.
func getDefaultGateway() string {
	cmd := newExecCommand("powershell", "-NoProfile", "-NonInteractive", "-Command",
		"Get-NetRoute -DestinationPrefix '0.0.0.0/0' | Select-Object -First 1 -ExpandProperty NextHop")
	out, err := cmd.Output()
	if err != nil {
		log.Printf("[DEBUG] Failed to get default gateway: %v", err)
		return ""
	}
	gw := strings.TrimSpace(string(out))
	if gw == "" {
		return ""
	}
	log.Printf("[DEBUG] Default Gateway: %s", gw)
	return gw
}

// generateNetworkGroupID creates a stable group ID based on the physical network fingerprint.
// It hashes (gatewayMAC + "|" + gatewayIP + "|" + dnsServers) to group devices behind the
// same physical router regardless of IP subnet.
func generateNetworkGroupID() string {
	gw := getDefaultGateway()
	if gw == "" {
		log.Printf("[DEBUG] No gateway found, NetworkGroupID will be empty")
		return ""
	}

	// Ping gateway once to populate ARP cache
	newExecCommand("ping", "-n", "1", "-w", "1000", gw).Run()

	// Get gateway MAC from ARP
	gwMAC := ""
	cmd := newExecCommand("powershell", "-NoProfile", "-NonInteractive", "-Command",
		fmt.Sprintf("(Get-NetNeighbor -IPAddress '%s' -AddressFamily IPv4 | Select-Object -First 1).LinkLayerAddress", gw))
	out, err := cmd.Output()
	if err == nil {
		gwMAC = strings.TrimSpace(string(out))
	}

	// Get DNS servers
	dns := ""
	cmd2 := newExecCommand("powershell", "-NoProfile", "-NonInteractive", "-Command",
		"(Get-DnsClientServerAddress -AddressFamily IPv4 | Where-Object {$_.ServerAddresses}).ServerAddresses -join ','")
	out2, err2 := cmd2.Output()
	if err2 == nil {
		dns = strings.TrimSpace(string(out2))
	}

	// Build fingerprint: prefer MAC+IP+DNS, fallback to IP+DNS, then just IP
	raw := gwMAC + "|" + gw + "|" + dns
	log.Printf("[DEBUG] Network fingerprint raw: %s", raw)

	hash := sha256.Sum256([]byte(raw))
	groupID := fmt.Sprintf("%x", hash[:8]) // first 8 bytes = 16 hex chars

	log.Printf("[DEBUG] NetworkGroupID: %s", groupID)
	return groupID
}

// getDefaultRouteIP returns the local IP used for the default route
// by dialing an external UDP address. Returns nil if unsuccessful.
func getDefaultRouteIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP
}

func getNetworkInfo() (string, string, []NetworkAdapter) {
	primaryIP := ""
	subnet := ""
	var adapters []NetworkAdapter

	defaultRouteIP := getDefaultRouteIP()
	if defaultRouteIP != nil {
		log.Printf("[DEBUG] Default route local IP: %s", defaultRouteIP.String())
	} else {
		log.Printf("[DEBUG] Could not determine default route IP, will use first adapter as fallback")
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return "", "", nil
	}

	// First pass: find the interface index that owns the default route IP
	gwIndex := -1
	if defaultRouteIP != nil {
		for _, iface := range interfaces {
			if iface.Flags&net.FlagUp == 0 {
				continue
			}
			if iface.Flags&net.FlagLoopback != 0 {
				continue
			}
			if iface.Flags&net.FlagPointToPoint != 0 {
				continue
			}
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				ipnet, ok := addr.(*net.IPNet)
				if !ok || ipnet.IP.To4() == nil {
					continue
				}
				if ipnet.IP.Equal(defaultRouteIP) {
					gwIndex = iface.Index
					break
				}
			}
			if gwIndex != -1 {
				break
			}
		}
	}

	// Second pass: build adapter list and set primary from gateway interface
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if iface.Flags&net.FlagPointToPoint != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		var ipList []string
		var mask string
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			if ipnet.IP.To4() == nil {
				continue
			}
			ip := ipnet.IP.String()
			ipList = append(ipList, ip)
			if mask == "" {
				mask = net.IP(ipnet.Mask).String()
			}
			// Prefer the default gateway interface
			if iface.Index == gwIndex {
				primaryIP = ip
				networkIP := ipnet.IP.Mask(ipnet.Mask)
				ones, _ := ipnet.Mask.Size()
				subnet = fmt.Sprintf("%s/%d", networkIP.String(), ones)
				log.Printf("[DEBUG] ========== Selected Interface ==========")
				log.Printf("[DEBUG] Interface Name: %s", iface.Name)
				log.Printf("[DEBUG] Hardware Address: %s", iface.HardwareAddr.String())
				log.Printf("[DEBUG] Local IP: %s", primaryIP)
				log.Printf("[DEBUG] Subnet Mask: %s", mask)
				log.Printf("[DEBUG] ConnectedSubnet: %s", subnet)
				log.Printf("[DEBUG] ========================================")
			}
			// Fallback: first interface found
			if primaryIP == "" {
				primaryIP = ip
				networkIP := ipnet.IP.Mask(ipnet.Mask)
				ones, _ := ipnet.Mask.Size()
				subnet = fmt.Sprintf("%s/%d", networkIP.String(), ones)
			}
		}

		if len(ipList) > 0 {
			mac := iface.HardwareAddr.String()
			adapter := NetworkAdapter{
				Name:        iface.Name,
				MACAddress:  mac,
				IPAddresses: fmt.Sprintf("%v", ipList),
				SubnetMask:  mask,
				IsPrimary:   iface.Index == gwIndex,
			}
			adapters = append(adapters, adapter)
		}
	}

	return primaryIP, subnet, adapters
}

func getOSVersion() string {
	info, err := host.Info()
	if err != nil || info == nil {
		return runtime.GOOS
	}
	return info.Platform + " " + info.PlatformVersion
}

func getPrimaryMAC(adapters []NetworkAdapter) string {
	for _, a := range adapters {
		if a.IsPrimary {
			return a.MACAddress
		}
	}
	if len(adapters) > 0 {
		return adapters[0].MACAddress
	}
	return ""
}

func AdaptersToJSON(adapters []NetworkAdapter) string {
	b, _ := json.Marshal(adapters)
	return string(b)
}
