package wg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl"
)

type Device struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	PublicKey    string `json:"public_key"`
	PrivateKey   string `json:"private_key,omitempty"`
	FirewallMark int    `json:"firewall_mark,omitempty"`
	ListenPort   int    `json:"listen_port,omitempty"`
	Peers        []Peer `json:"peers,omitempty"`
}

type Peer struct {
	PublicKey         string    `json:"public_key"`
	AllowedIPs        []string  `json:"allowed_ips"`
	Endpoint          string    `json:"endpoint"`
	LastHandshakeTime time.Time `json:"last_handshake_time"`
	ReceiveBytes      int64     `json:"receive_bytes"`
	TransmitBytes     int64     `json:"transmit_bytes"`
}

func HandleGetDevices(w http.ResponseWriter, r *http.Request) {
	client, err := wgctrl.New()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create wgctrl client: %v", err), http.StatusInternalServerError)
		return
	}
	defer client.Close()

	devices, err := client.Devices()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get devices: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert devices to our response format
	response := make([]Device, 0, len(devices))
	for _, device := range devices {
		deviceResp := Device{
			Name:         device.Name,
			Type:         device.Type.String(),
			PublicKey:    device.PublicKey.String(),
			ListenPort:   device.ListenPort,
			Peers:        make([]Peer, 0, len(device.Peers)),
			FirewallMark: device.FirewallMark,
		}

		for _, peer := range device.Peers {
			// Convert AllowedIPs to string representation
			allowedIPs := make([]string, 0, len(peer.AllowedIPs))
			for _, ip := range peer.AllowedIPs {
				allowedIPs = append(allowedIPs, ip.String())
			}

			peerInfo := Peer{
				PublicKey:         peer.PublicKey.String(),
				AllowedIPs:        allowedIPs,
				Endpoint:          peer.Endpoint.String(),
				LastHandshakeTime: peer.LastHandshakeTime,
				ReceiveBytes:      peer.ReceiveBytes,
				TransmitBytes:     peer.TransmitBytes,
			}
			deviceResp.Peers = append(deviceResp.Peers, peerInfo)
		}

		response = append(response, deviceResp)
	}

	// Set content type header
	w.Header().Set("Content-Type", "application/json")

	// Encode and send response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}
