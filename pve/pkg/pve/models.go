package pve

import "fmt"

// API response types

type APINode struct {
	Node   string `json:"node"`
	Status string `json:"status"`
}

type APINodesResponse struct {
	Data []APINode `json:"data"`
}

type APIVM struct {
	VMID      int64   `json:"vmid"`
	Name      string  `json:"name"`
	Status    string  `json:"status"`
	CPU       float64 `json:"cpu"`
	Cpus      int     `json:"cpus"`
	Mem       int64   `json:"mem"`
	MaxMem    int64   `json:"maxmem"`
	Uptime    int64   `json:"uptime"`
	NetIn     int64   `json:"netin"`
	NetOut    int64   `json:"netout"`
	DiskRead  int64   `json:"diskread"`
	DiskWrite int64   `json:"diskwrite"`
}

type APIVMListResponse struct {
	Data []APIVM `json:"data"`
}

type APINetworkInterface struct {
	Name        string `json:"name"`
	HardwareAddr string `json:"hardware-address"`
	IPAddresses []struct {
		IPAddress string `json:"ip-address"`
		IPType    string `json:"ip-address-type"`
	} `json:"ip-addresses"`
}

type APIQemuAgentNetworkResponse struct {
	Data struct {
		Result []APINetworkInterface `json:"result"`
	} `json:"data"`
}

type APILXCInterfaceResponse struct {
	Data []struct {
		Name   string `json:"name"`
		Inet   string `json:"inet"`
		Inet6  string `json:"inet6"`
		Hwaddr string `json:"hwaddr"`
	} `json:"data"`
}

type APITaskResponse struct {
	Data string `json:"data"`
}

// Output types

type Guest struct {
	VMID   int64  `yaml:"vmid"`
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	Status string `yaml:"status"`
	CPU    int    `yaml:"cpu"`
	Memory int64  `yaml:"memory"`
	Uptime string `yaml:"uptime"`
	IP     string `yaml:"ip"`
}

type ActionResult struct {
	VMID   int64  `yaml:"vmid"`
	Name   string `yaml:"name"`
	Action string `yaml:"action"`
}

// Helper functions

func (vm *APIVM) StatusLabel() string {
	return vm.Status
}

func FormatUptime(seconds int64) string {
	if seconds == 0 {
		return "0s"
	}

	days := seconds / 86400
	seconds %= 86400
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60

	var result string
	if days > 0 {
		result = fmt.Sprintf("%dd", days)
	}
	if hours > 0 {
		if result != "" {
			result += " "
		}
		result += fmt.Sprintf("%dh", hours)
	}
	if minutes > 0 && (days == 0 || hours > 0) {
		if result != "" {
			result += " "
		}
		result += fmt.Sprintf("%dm", minutes)
	}
	if result == "" {
		result = fmt.Sprintf("%ds", seconds)
	}
	return result
}

func (vm *APIVM) ToGuest(vmType, ip string) Guest {
	return Guest{
		VMID:   vm.VMID,
		Name:   vm.Name,
		Type:   vmType,
		Status: vm.StatusLabel(),
		CPU:    vm.Cpus,
		Memory: vm.MaxMem / 1024 / 1024, // Convert to MB
		Uptime: FormatUptime(vm.Uptime),
		IP:     ip,
	}
}
