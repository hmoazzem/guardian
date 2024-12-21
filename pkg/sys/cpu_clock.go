package sys

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// CPUClock retrieves the CPU clock speeds from /proc/cpuinfo as a slice of float64.
func CPUClock() ([]float64, error) {
	// Open the /proc/cpuinfo file
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var frequencies []float64
	scanner := bufio.NewScanner(file)

	// Read through the file line by line
	for scanner.Scan() {
		line := scanner.Text()
		// Check for lines containing "cpu MHz"
		if strings.HasPrefix(line, "cpu MHz") {
			// Split the line into key-value parts
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				continue
			}
			// Parse the frequency value
			freqStr := strings.TrimSpace(parts[1])
			freq, err := strconv.ParseFloat(freqStr, 64)
			if err == nil {
				frequencies = append(frequencies, freq)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return frequencies, nil
}

// using shell utils
//
// func CPUClock() ([]float64, error) {
// 	cmd := exec.Command("bash", "-c", "grep 'cpu MHz' /proc/cpuinfo | awk '{print $4}'")
// 	output, err := cmd.Output()
// 	if err != nil {
// 		return nil, err
// 	}
// 	frequencies := []float64{}
// 	for _, line := range strings.Split(string(output), "\n") {
// 		if line == "" {
// 			continue
// 		}
// 		freq, err := strconv.ParseFloat(line, 64)
// 		if err == nil {
// 			frequencies = append(frequencies, freq)
// 		}
// 	}
// 	return frequencies, nil
// }
