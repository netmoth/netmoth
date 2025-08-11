package web

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/netmoth/netmoth/internal/config"
	"github.com/netmoth/netmoth/internal/storage/redis"
	"github.com/netmoth/netmoth/internal/version"
	redisclient "github.com/redis/go-redis/v9"
)

// CORS middleware
func corsMiddleware(allowedOrigins []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allow := ""
		if len(allowedOrigins) == 0 {
			allow = "*"
		} else {
			for _, ao := range allowedOrigins {
				if ao == origin {
					allow = origin
					break
				}
			}
		}
		if allow != "" {
			w.Header().Set("Access-Control-Allow-Origin", allow)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// WebSocket handler (gorilla/websocket)
func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()
	conn.SetReadLimit(1 << 20) // 1MB
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			_ = conn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(5*time.Second))
		}
	}()
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if err := conn.WriteMessage(mt, message); err != nil {
			break
		}
	}
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

	// Optional profiling only on localhost if enabled
	if os.Getenv("NETMOTH_PPROF") == "1" {
		go func() {
			log.Println("Starting pprof server on 127.0.0.1:6060")
			log.Fatal(http.ListenAndServe("127.0.0.1:6060", nil))
		}()
	}

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
	mux.HandleFunc("/api/agent/register", makeAgentRegistrationHandler(cfg))
	mux.HandleFunc("/api/agent/data", makeAgentDataHandler(cfg))
	mux.HandleFunc("/api/agent/health", makeAgentHealthHandler(cfg))

	// WebSocket route
	mux.HandleFunc("/ws", websocketHandler)

	// Static files (SPA support)
	mux.HandleFunc("/", staticFileHandler)

	// Apply CORS middleware using config origins
	handler := corsMiddleware(cfg.AllowedOrigins, mux)

	// Configure server
	server := &http.Server{
		Addr:         ":3000",
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-shutdownCh
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Netmoth Web Server v%s starting on :3000", version.Version())
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
