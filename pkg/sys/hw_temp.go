package sys

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Hwmon struct {
	Name      string        `json:"name"`
	Composite float64       `json:"composite"`
	Sensors   []HwmonSensor `json:"sensors"`
}

type HwmonSensor struct {
	Name string  `json:"name"`
	Temp float64 `json:"temp"`
}

func HwmonTemp() ([]Hwmon, error) {
	var hwtemps []Hwmon

	// Find all hwmon directories
	hwmonDirs, err := filepath.Glob("/sys/class/hwmon/hwmon*")
	if err != nil {
		return nil, fmt.Errorf("failed to list hwmon directories: %w", err)
	}

	for _, dir := range hwmonDirs {
		// Read the device name
		nameFile := filepath.Join(dir, "name")
		name, err := os.ReadFile(nameFile)
		if err != nil {
			// Skip hwtemps we can't read
			continue
		}

		hwTemp := Hwmon{Name: strings.TrimSpace(string(name))}

		// Find all temp*_input files
		tempFiles, err := filepath.Glob(filepath.Join(dir, "temp*_input"))
		if err != nil {
			return nil, fmt.Errorf("failed to list temp files for %s: %w", hwTemp.Name, err)
		}

		for _, tempFile := range tempFiles {
			// Read temperature value
			tempValue, err := os.ReadFile(tempFile)
			if err != nil {
				// Skip sensors we can't read
				continue
			}
			tempMilliC, err := strconv.Atoi(strings.TrimSpace(string(tempValue)))
			if err != nil {
				// Skip invalid temperature values
				continue
			}
			tempC := float64(tempMilliC) / 1000.0

			// Read the label if available
			labelFile := strings.TrimSuffix(tempFile, "_input") + "_label"
			label := "Unknown"
			if _, err := os.Stat(labelFile); err == nil {
				labelBytes, err := os.ReadFile(labelFile)
				if err == nil {
					label = strings.TrimSpace(string(labelBytes))
				}
			}

			// Append the sensor data
			hwTemp.Sensors = append(hwTemp.Sensors, HwmonSensor{Name: label, Temp: tempC})
		}

		// Optionally calculate Composite if needed
		// Example: Use the average of all sensor temps
		if len(hwTemp.Sensors) > 0 {
			var totalTemp float64
			for _, sensor := range hwTemp.Sensors {
				totalTemp += sensor.Temp
			}
			hwTemp.Composite = totalTemp / float64(len(hwTemp.Sensors))
		}

		hwtemps = append(hwtemps, hwTemp)
	}

	return hwtemps, nil
}
