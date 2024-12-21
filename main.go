package main

import (
	"embed"
	"flag"
)

var (
	grpcPort    = flag.Int("grcp-port", 30080, "port to listen on")
	httpPort    = flag.Int("http-port", 8081, "port to serve http API and Anuglar web UI")
	directory   = flag.String("dir", ".", "Directory to serve files from")
	spaFallback = flag.Bool("spa", false, "fallback to index.html for not-found files")
	useEmbedded = flag.Bool("embed", false, "use embedded static files")
	envoyConf   = flag.String("envoy-conf", "pkg/config/envoy.yaml", "envoy --config-path")
)

// embeddedFS holds the embedded Angular application files
//
//go:embed gui/dist/gui/browser
var embeddedFS embed.FS

func main() {
	flag.Parse()

	// Built Angular UI
	var fs *embed.FS
	if *useEmbedded {
		fs = &embeddedFS
	}
	go serveHTTP(httpPort, directory, fs)

	// Envoy gRPC gateway
	go runEnvoy(*envoyConf)

	serveGRPC(grpcPort)
}
