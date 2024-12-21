package amdsmi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

// Metrics represents the top-level response from amd-smi
type Metrics []GPUMetric

// GPUMetric represents metrics for a single GPU
type GPUMetric struct {
	GPU int `json:"gpu"`
	// Usage Usage `json:"usage"`
	Usage json.RawMessage `json:"usage"`
	// Power Power `json:"power"`
	Power json.RawMessage `json:"power,omitempty"`
	// Clock Clock `json:"clock"`
	Clock       json.RawMessage `json:"clock"`
	Temperature Temp            `json:"temperature"`
	// PCIe        PCIe `json:"pcie"`
	PCIe      json.RawMessage `json:"pcie"`
	ECC       ECC             `json:"ecc"`
	ECCBlocks string          `json:"ecc_blocks"`
	MemUsage  MemUsage        `json:"mem_usage"`
}

// Usage represents GPU usage metrics
type Usage struct {
	GfxActivity  ValueUnit `json:"gfx_activity"`
	UmcActivity  ValueUnit `json:"umc_activity"`
	MmActivity   ValueUnit `json:"mm_activity"`
	VcnActivity  []any     `json:"vcn_activity"`
	JpegActivity []string  `json:"jpeg_activity"`
	GfxBusyInst  string    `json:"gfx_busy_inst"`
	JpegBusy     string    `json:"jpeg_busy"`
	VcnBusy      string    `json:"vcn_busy"`
	GfxBusyAcc   string    `json:"gfx_busy_acc"`
}

// Power represents power-related metrics
type Power struct {
	SocketPower     ValueUnit `json:"socket_power"`
	GfxVoltage      ValueUnit `json:"gfx_voltage"`
	SocVoltage      ValueUnit `json:"soc_voltage"`
	MemVoltage      ValueUnit `json:"mem_voltage"`
	ThrottleStatus  string    `json:"throttle_status"`
	PowerManagement string    `json:"power_management"`
}

// Clock represents clock-related metrics
type Clock struct {
	Gfx  [8]ClockInfo `json:"gfx_0,gfx_1,gfx_2,gfx_3,gfx_4,gfx_5,gfx_6,gfx_7"`
	Mem  [1]ClockInfo `json:"mem_0"`
	VClk [4]ClockInfo `json:"vclk_0,vclk_1,vclk_2,vclk_3"`
	DClk [4]ClockInfo `json:"dclk_0,dclk_1,dclk_2,dclk_3"`
}

// ClockInfo represents clock information for a specific component
type ClockInfo struct {
	Clk       any    `json:"clk"`     // Can be ValueUnit or string
	MinClk    any    `json:"min_clk"` // Can be ValueUnit or string
	MaxClk    any    `json:"max_clk"` // Can be ValueUnit or string
	ClkLocked string `json:"clk_locked"`
	DeepSleep string `json:"deep_sleep"`
}

// Temp represents temperature metrics
type Temp struct {
	Edge    ValueUnit `json:"edge"`
	Hotspot any       `json:"hotspot"` // ValueUnit
	Mem     any       `json:"mem"`     // ValueUnit
}

// PCIe represents PCIe-related metrics
type PCIe struct {
	Width                    int       `json:"width"`
	Speed                    ValueUnit `json:"speed"`
	Bandwidth                string    `json:"bandwidth"`
	ReplayCount              string    `json:"replay_count"`
	L0ToRecoveryCount        string    `json:"l0_to_recovery_count"`
	ReplayRollOverCount      string    `json:"replay_roll_over_count"`
	NakSentCount             string    `json:"nak_sent_count"`
	NakReceivedCount         string    `json:"nak_received_count"`
	CurrentBandwidthSent     string    `json:"current_bandwidth_sent"`
	CurrentBandwidthReceived string    `json:"current_bandwidth_received"`
	MaxPacketSize            string    `json:"max_packet_size"`
	LcPerfOtherEndRecovery   string    `json:"lc_perf_other_end_recovery"`
}

// ECC represents error correction code metrics
type ECC struct {
	TotalCorrectableCount   int    `json:"total_correctable_count"`
	TotalUncorrectableCount int    `json:"total_uncorrectable_count"`
	TotalDeferredCount      int    `json:"total_deferred_count"`
	CacheCorrectableCount   string `json:"cache_correctable_count"`
	CacheUncorrectableCount string `json:"cache_uncorrectable_count"`
}

// MemUsage represents memory usage metrics
type MemUsage struct {
	TotalVRAM        ValueUnit `json:"total_vram"`
	UsedVRAM         ValueUnit `json:"used_vram"`
	FreeVRAM         ValueUnit `json:"free_vram"`
	TotalVisibleVRAM ValueUnit `json:"total_visible_vram"`
	UsedVisibleVRAM  ValueUnit `json:"used_visible_vram"`
	FreeVisibleVRAM  ValueUnit `json:"free_visible_vram"`
	TotalGTT         ValueUnit `json:"total_gtt"`
	UsedGTT          ValueUnit `json:"used_gtt"`
	FreeGTT          ValueUnit `json:"free_gtt"`
}

func getGPUMetrics() (*Metrics, error) {
	// Execute amd-smi command
	cmd := exec.Command("amd-smi", "metric", "--json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute amd-smi: %v", err)
	}

	// Parse the JSON output
	var metrics Metrics
	if err := json.Unmarshal(output, &metrics); err != nil {
		return nil, fmt.Errorf("failed to parse amd-smi output: %v", err)
	}

	return &metrics, nil
}

func HandleMetrics(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get metrics
	metrics, err := getGPUMetrics()
	if err != nil {
		log.Printf("Error getting GPU metrics: %v", err)
		http.Error(w, "Failed to get GPU metrics", http.StatusInternalServerError)
		return
	}

	// Set content type header
	w.Header().Set("Content-Type", "application/json")

	// Encode and send response
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
