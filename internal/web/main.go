package web

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/netmoth/netmoth/internal/config"
	"github.com/netmoth/netmoth/internal/storage/redis"
	"github.com/netmoth/netmoth/internal/version"
	redisclient "github.com/redis/go-redis/v9"
	"golang.org/x/net/websocket"
)

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// API handlers
func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"version": version.Version(),
	}
	json.NewEncoder(w).Encode(response)
}

// WebSocket handler
func websocketHandler(w http.ResponseWriter, r *http.Request) {
	websocket.Handler(func(ws *websocket.Conn) {
		log.Println("WebSocket connection established")

		for {
			var message string
			err := websocket.Message.Receive(ws, &message)
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				break
			}

			log.Printf("Received: %s", message)

			// Echo the message back
			err = websocket.Message.Send(ws, message)
			if err != nil {
				log.Printf("WebSocket write error: %v", err)
				break
			}
		}
	}).ServeHTTP(w, r)
}

// Static file server with SPA support
func staticFileHandler(w http.ResponseWriter, r *http.Request) {
	// Remove leading slash
	path := strings.TrimPrefix(r.URL.Path, "/")

	// If path is empty, serve index.html
	if path == "" {
		path = "index.html"
	}

	// Construct full file path
	filePath := filepath.Join("./web/dist", path)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// If file doesn't exist, serve index.html for SPA routing
		filePath = filepath.Join("./web/dist", "index.html")
	}

	// Serve the file
	http.ServeFile(w, r, filePath)
}

// New starts the web server
func New(configPath string) {
	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Start profiling in a separate goroutine
	go func() {
		log.Println("Starting pprof server on :6060")
		log.Fatal(http.ListenAndServe(":6060", nil))
	}()

	// Set number of CPUs for optimal performance
	if cfg.MaxCores > 0 && cfg.MaxCores < runtime.NumCPU() {
		runtime.GOMAXPROCS(cfg.MaxCores)
	}

	// Connect to Redis
	redisOpts := &redisclient.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		DB:       0,
	}
	_ = redis.NewRedisHandler(context.Background(), redisOpts)

	// Create mux
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/version", versionHandler)
	mux.HandleFunc("/api/agent/register", agentRegistrationHandler)
	mux.HandleFunc("/api/agent/data", agentDataHandler)
	mux.HandleFunc("/api/agent/health", agentHealthHandler)

	// WebSocket route
	mux.HandleFunc("/ws", websocketHandler)

	// Static files (SPA support)
	mux.HandleFunc("/", staticFileHandler)

	// Apply CORS middleware
	handler := corsMiddleware(mux)

	// Configure server
	server := &http.Server{
		Addr:         ":3000",
		Handler:      handler,
		ReadTimeout:  30,
		WriteTimeout: 30,
		IdleTimeout:  60,
	}

	log.Printf("Netmoth Web Server v%s starting on :3000", version.Version())
	log.Fatal(server.ListenAndServe())
}
