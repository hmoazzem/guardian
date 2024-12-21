package sys

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type SystemInfo struct {
	Hostname       string `json:"hostname"`
	OS             string `json:"os"`
	CPU            CPU    `json:"cpu"`
	Memory         Memory `json:"memory"`
	Disk           Disk   `json:"disk"`
	Kernel         string `json:"kernel"`
	Uptime         string `json:"uptime"`
	GPUs           []GPU  `json:"gpus"`
	Motherboard    string `json:"motherboard"`
	Virtualization bool   `json:"virtualization_supported"`
}

type CPU struct {
	Model   string  `json:"model"`
	Threads int     `json:"threads"`
	Clock   float64 `json:"clock_ghz"`
}

type Memory struct {
	TotalGB     float64 `json:"total_gb"`
	AvailableGB float64 `json:"available_gb"`
}

type Disk struct {
	Used  string `json:"used"`
	Total string `json:"total"`
	// UsedRatio string `json:"used_ratio"`
}

type GPU struct {
	Name   string `json:"name"`
	Vendor string `json:"vendor"`
	Device string `json:"device"`
}

func Info() (*SystemInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("error reading hostname: %v", err)
	}

	osName, _ := readOSInfo()
	cpu, virt := readCPUInfo()
	memory := readMemoryInfo()
	disk, _ := readDiskInfo()
	kernel, _ := readKernelInfo()
	uptime, _ := readUptime()
	gpus := readGPUInfo()
	motherboard := readMotherboardInfo()

	return &SystemInfo{
		Hostname:       hostname,
		OS:             osName,
		CPU:            cpu,
		Memory:         memory,
		Disk:           disk,
		Kernel:         kernel,
		Uptime:         uptime,
		GPUs:           gpus,
		Motherboard:    motherboard,
		Virtualization: virt,
	}, nil
}

func readOSInfo() (string, error) {
	content, err := readFile("/etc/os-release")
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(line[12:], "\""), nil
		}
	}
	return "Unknown", nil
}

func readCPUInfo() (CPU, bool) {
	content, err := readFile("/proc/cpuinfo")
	if err != nil {
		return CPU{}, false
	}

	model := extractField(content, "model name")
	threads := strings.Count(content, "processor\t:")
	clockSpeed := parseClockSpeed(extractField(content, "cpu MHz"))
	virtSupported := checkVirtualizationSupport(content)

	return CPU{
		Model:   model,
		Threads: threads,
		Clock:   clockSpeed,
	}, virtSupported
}

func readMemoryInfo() Memory {
	content, err := readFile("/proc/meminfo")
	if err != nil {
		return Memory{}
	}
	totalKB := parseValueKB(content, "MemTotal")
	availableKB := parseValueKB(content, "MemAvailable")

	return Memory{
		TotalGB:     kbToGiB(totalKB),
		AvailableGB: kbToGiB(availableKB),
	}
}

func readDiskInfo() (Disk, error) {
	dfCmd := exec.Command("df", "-h", "/")
	output, err := dfCmd.Output()
	if err != nil {
		return Disk{}, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 1 {
		fields := strings.Fields(lines[1])
		return Disk{
			Used:  fields[2],
			Total: fields[1],
		}, nil
	}
	return Disk{}, nil
}

func readKernelInfo() (string, error) {
	content, err := readFile("/proc/version")
	if err != nil {
		return "", err
	}
	fields := strings.Fields(content)
	if len(fields) > 2 {
		return fields[2], nil
	}
	return "Unknown", nil
}

func readUptime() (string, error) {
	content, err := readFile("/proc/uptime")
	if err != nil {
		return "", err
	}
	uptimeSeconds := strings.Split(content, " ")[0]
	seconds, err := time.ParseDuration(fmt.Sprintf("%ss", uptimeSeconds))
	if err != nil {
		return "", err
	}
	return seconds.String(), nil
}

func readGPUInfo() []GPU {
	drmPath := "/sys/class/drm/"
	entries, err := os.ReadDir(drmPath)
	if err != nil {
		return nil
	}

	var gpus []GPU
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "card") && !strings.Contains(entry.Name(), "-") {
			gpuPath := filepath.Join(drmPath, entry.Name(), "device")
			vendor, _ := readFileTrimmed(filepath.Join(gpuPath, "vendor"))
			device, _ := readFileTrimmed(filepath.Join(gpuPath, "device"))
			if vendor != "" && device != "" {
				gpus = append(gpus, GPU{
					Name:   entry.Name(),
					Vendor: vendor,
					Device: device,
				})
			}
		}
	}
	return gpus
}

func readMotherboardInfo() string {
	manufacturer, err := readFileTrimmed("/sys/devices/virtual/dmi/id/board_vendor")
	if err != nil {
		manufacturer = "Unknown"
	}
	model, err := readFileTrimmed("/sys/devices/virtual/dmi/id/board_name")
	if err != nil {
		model = "Unknown"
	}
	return fmt.Sprintf("%s %s", manufacturer, model)
}

// Helper functions
func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func readFileTrimmed(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func parseValueKB(content, field string) int64 {
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, field) {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				var value int64
				fmt.Sscanf(parts[1], "%d", &value)
				return value
			}
		}
	}
	return 0
}

// Helper to convert KB to GiB
func kbToGiB(kb int64) float64 {
	return float64(kb) / (1024 * 1024)
}

func extractField(content, field string) string {
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, field) {
			return strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}
	return "Unknown"
}

func parseClockSpeed(clockSpeed string) float64 {
	var mhz float64
	fmt.Sscanf(clockSpeed, "%f", &mhz)
	return mhz / 1000 // Convert MHz to GHz
}

func checkVirtualizationSupport(cpuInfo string) bool {
	for _, line := range strings.Split(cpuInfo, "\n") {
		if strings.HasPrefix(line, "flags") {
			if strings.Contains(line, "vmx") || strings.Contains(line, "svm") {
				return true
			}
		}
	}
	return false
}
