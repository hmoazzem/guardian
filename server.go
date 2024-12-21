package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"time"

	"net"

	"github.com/hmoazzem/guardian/pkg/amdsmi"
	"github.com/hmoazzem/guardian/pkg/net/wg"
	"github.com/hmoazzem/guardian/pkg/sys"
	pb "github.com/hmoazzem/guardian/proto/generated"

	mw "github.com/edgeflare/pgo/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server implements the System service
type Server struct {
	pb.UnimplementedMetricsServer
}

func (s *Server) StreamCPUUtilization(empty *pb.Empty, stream pb.Metrics_StreamCPUUtilizationServer) error {
	var prevStats []sys.ProcStat

	for {
		currentStats, err := sys.ReadProcStat()
		if err != nil {
			return err
		}

		// Only calculate utilization if we have previous stats to compare against
		if len(prevStats) > 0 && len(currentStats) == len(prevStats) {
			utilizations := make([]float32, len(currentStats))
			for i, prevStat := range prevStats {
				utilization := sys.CalcCPUUtilization(prevStat, currentStats[i])
				utilizations[i] = float32(utilization) // Convert to float32 for protobuf
			}

			// Send the utilization data
			if err := stream.Send(&pb.CPUUtilization{
				CpuUtilization: utilizations,
			}); err != nil {
				return err
			}
		}

		prevStats = currentStats
		time.Sleep(1 * time.Second)
	}
}

// StreamCPUClock streams CPUClock in GHz
func (s *Server) StreamCPUClock(req *pb.Empty, stream pb.Metrics_StreamCPUClockServer) error {
	for {
		// Simulate fetching CPU clock speeds
		cpuClock, err := sys.CPUClock()
		if err != nil {
			return err
		}

		// Convert to GHz and round to 2 decimal places
		var float32CpuClock []float32
		for _, clockSpeed := range cpuClock {
			ghz := clockSpeed / 1000.0
			roundedGhz := math.Round(ghz*100) / 100
			float32CpuClock = append(float32CpuClock, float32(roundedGhz))
		}

		// Send CPU clock speeds to the client
		err = stream.Send(&pb.CPUClock{CpuClock: float32CpuClock})
		if err != nil {
			return err
		}

		time.Sleep(1 * time.Second) // Stream data every second
	}
}

// StreamHwmon implements the StreamHwmon RPC method
func (s *Server) StreamHwmon(empty *pb.Empty, stream pb.Metrics_StreamHwmonServer) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case <-ticker.C:
			// Get the hardware monitoring data
			hwmons, err := sys.HwmonTemp()
			if err != nil {
				continue // Skip this reading if there's an error, but keep streaming
			}

			// Convert each Hwmon to protobuf message
			for _, hwmon := range hwmons {
				// Create protobuf sensors slice
				pbSensors := make([]*pb.HwmonSensor, len(hwmon.Sensors))
				for i, sensor := range hwmon.Sensors {
					pbSensors[i] = &pb.HwmonSensor{
						Name: sensor.Name,
						Temp: sensor.Temp,
					}
				}

				// Create the protobuf Hwmon message
				pbHwmon := &pb.Hwmon{
					Name:      hwmon.Name,
					Composite: hwmon.Composite,
					Sensors:   pbSensors,
				}

				// Send the data to the client
				if err := stream.Send(pbHwmon); err != nil {
					return err
				}
			}
		}
	}
}

func serveGRPC(grpcPort *int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMetricsServer(grpcServer, &Server{})
	reflection.Register(grpcServer) // For debugging with gRPC reflection tools

	log.Printf("gRPC server listening on :%d", *grpcPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func serveHTTP(port *int, directory *string, fs *embed.FS) {
	mux := http.NewServeMux()

	mux.Handle("GET /sysinfo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sysInfo, _ := sys.Info()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sysInfo)
	}))

	// AMD GPU Handlers
	if _, err := exec.LookPath("amd-smi"); err == nil { // Check if `amd-smi` binary is available
		mux.HandleFunc("GET /amdgpu", amdsmi.HandleGetGPUs)
		mux.HandleFunc("GET /amdgpu-metrics", amdsmi.HandleMetrics)
	}

	// Enable WireGuard routes if running as root
	if os.Geteuid() == 0 {
		mux.HandleFunc("GET /net/wg-devices", wg.HandleGetDevices)
	}

	mux.Handle("GET /", mw.Static(*directory, *spaFallback, fs))

	// Start the server
	log.Printf("starting server on :%d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mw.CORSWithOptions(nil)(mux))) // enable cors in dev; should be disables
}

func runEnvoy(configPath string) {
	cmd := exec.Command("envoy", "--config-path", configPath)

	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start Envoy: %v", err)
	}

	log.Println("Envoy is running...")

	// Wait for the Envoy process to complete
	if err := cmd.Wait(); err != nil {
		log.Printf("Envoy process exited with error: %v", err)
	}
}

// // mockCPUClock function for testing. useful for deving on Mac/Win
// func mockCPUClock() ([]float64, error) {
// 	// Generate random cpuClock for demonstration
// 	numCores := 32
// 	cpuClock := make([]float64, numCores)
// 	for i := range cpuClock {
// 		cpuClock[i] = rand.Float64() * 6.0 // Simulate 0.0 GHz to 6.0 GHz
// 		// cpuClock[i] = 2.5 + rand.Float64()*1.0 // Simulate 2.5 GHz to 3.5 GHz
// 	}
// 	return cpuClock, nil
// }
