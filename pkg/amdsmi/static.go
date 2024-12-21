package amdsmi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

// ValueUnit represents a value with its unit
type ValueUnit struct {
	Value any    `json:"value"` // Can be int, float64, or string
	Unit  string `json:"unit"`
}

// CacheInfo represents the details of a cache.
type CacheInfo struct {
	Cache            int       `json:"cache"`
	CacheProperties  []string  `json:"cache_properties"`
	CacheSize        ValueUnit `json:"cache_size"`
	CacheLevel       int       `json:"cache_level"`
	MaxNumCUShared   int       `json:"max_num_cu_shared"`
	NumCacheInstance int       `json:"num_cache_instance"`
}

// VRAM represents the details of video memory.
type VRAM struct {
	Type     string    `json:"type"`
	Vendor   string    `json:"vendor"`
	Size     ValueUnit `json:"size"`
	BitWidth int       `json:"bit_width"`
}

// RAS represents the Reliability, Availability, and Serviceability details.
type RAS struct {
	EEPROMVersion   string `json:"eeprom_version"`
	ParitySchema    string `json:"parity_schema"`
	SingleBitSchema string `json:"single_bit_schema"`
	DoubleBitSchema string `json:"double_bit_schema"`
	PoisonSchema    string `json:"poison_schema"`
	ECCBlockState   string `json:"ecc_block_state"`
}

// Board represents the details of the board.
type Board struct {
	ModelNumber      string `json:"model_number"`
	ProductSerial    string `json:"product_serial"`
	FruID            string `json:"fru_id"`
	ProductName      string `json:"product_name"`
	ManufacturerName string `json:"manufacturer_name"`
}

// Driver represents the driver details.
type Driver struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// VBIOS represents the VBIOS details.
type VBIOS struct {
	Name       string `json:"name"`
	BuildDate  string `json:"build_date"`
	PartNumber string `json:"part_number"`
	Version    string `json:"version"`
}

// Bus represents the details of the bus.
type Bus struct {
	BDF                  string          `json:"bdf"`
	MaxPCIEWidth         any             `json:"max_pcie_width"` // int or "N/A"
	MaxPCIESpeed         json.RawMessage `json:"max_pcie_speed"`
	PCIEInterfaceVersion string          `json:"pcie_interface_version"`
	SlotType             string          `json:"slot_type"`
}

// ASIC represents the details of the ASIC.
type ASIC struct {
	MarketName            string `json:"market_name"`
	VendorID              string `json:"vendor_id"`
	VendorName            string `json:"vendor_name"`
	SubvendorID           string `json:"subvendor_id"`
	DeviceID              string `json:"device_id"`
	SubsystemID           string `json:"subsystem_id"`
	RevID                 string `json:"rev_id"`
	ASICSerial            string `json:"asic_serial"`
	OAMID                 int    `json:"oam_id"`
	NumComputeUnits       int    `json:"num_compute_units,omitempty"`
	TargetGraphicsVersion string `json:"target_graphics_version,omitempty"`
}

// GPU represents the details of a GPU.
type GPU struct {
	GPU              int         `json:"gpu"`
	ASIC             ASIC        `json:"asic"`
	Bus              Bus         `json:"bus"`
	VBIOS            VBIOS       `json:"vbios"`
	Driver           Driver      `json:"driver"`
	Board            Board       `json:"board"`
	RAS              RAS         `json:"ras"`
	ProcessIsolation string      `json:"process_isolation"`
	VRAM             VRAM        `json:"vram"`
	CacheInfo        []CacheInfo `json:"cache_info"`
}

func HandleGetGPUs(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("amd-smi", "static", "--json")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing amd-smi command: %v\n", err)
		return
	}

	// Parse the JSON output into our GPU struct
	var gpus []GPU
	// var gpus json.RawMessage
	err = json.Unmarshal(output, &gpus)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	// Set content type header
	w.Header().Set("Content-Type", "application/json")

	// Encode and send response
	if err := json.NewEncoder(w).Encode(gpus); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}
