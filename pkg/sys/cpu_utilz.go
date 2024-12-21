package sys

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// ProcStat holds the CPU statistics read from /proc/stat for a single processor.
type ProcStat struct {
	User, Nice, System, Idle, Iowait, Irq, Softirq, Steal, Guest, GuestNice int64
}

// ReadProcStat reads and parses the /proc/stat file to extract CPU usage data.
func ReadProcStat() ([]ProcStat, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cpuStats []ProcStat
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu") && !strings.HasPrefix(line, "cpu ") {
			parts := strings.Fields(line)
			if len(parts) < 8 {
				continue // Skip lines that don't have enough fields
			}

			var stat ProcStat
			stat.User, _ = strconv.ParseInt(parts[1], 10, 64)
			stat.Nice, _ = strconv.ParseInt(parts[2], 10, 64)
			stat.System, _ = strconv.ParseInt(parts[3], 10, 64)
			stat.Idle, _ = strconv.ParseInt(parts[4], 10, 64)
			stat.Iowait, _ = strconv.ParseInt(parts[5], 10, 64)
			stat.Irq, _ = strconv.ParseInt(parts[6], 10, 64)
			stat.Softirq, _ = strconv.ParseInt(parts[7], 10, 64)
			if len(parts) >= 9 {
				stat.Steal, _ = strconv.ParseInt(parts[8], 10, 64)
			}
			if len(parts) >= 10 {
				stat.Guest, _ = strconv.ParseInt(parts[9], 10, 64)
			}
			if len(parts) >= 11 {
				stat.GuestNice, _ = strconv.ParseInt(parts[10], 10, 64)
			}

			cpuStats = append(cpuStats, stat)
		}
	}

	return cpuStats, scanner.Err()
}

// CalcCPUUtilization calculates the CPU utilization percentage.
func CalcCPUUtilization(prevStat, currentStat ProcStat) float64 {
	prevTotal := prevStat.User + prevStat.Nice + prevStat.System + prevStat.Idle + prevStat.Iowait + prevStat.Irq + prevStat.Softirq + prevStat.Steal + prevStat.Guest + prevStat.GuestNice
	currentTotal := currentStat.User + currentStat.Nice + currentStat.System + currentStat.Idle + currentStat.Iowait + currentStat.Irq + currentStat.Softirq + currentStat.Steal + currentStat.Guest + currentStat.GuestNice

	prevIdle := prevStat.Idle + prevStat.Iowait
	currentIdle := currentStat.Idle + currentStat.Iowait

	idleDelta := float64(currentIdle - prevIdle)
	totalDelta := float64(currentTotal - prevTotal)

	cpuUtilization := (totalDelta - idleDelta) / totalDelta * 100.0

	return cpuUtilization
}
