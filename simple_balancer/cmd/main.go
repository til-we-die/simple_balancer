package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"

	"simple_balancer/internal/health"
	"simple_balancer/pkg/lb"
)

// структура конфигурации приложения
type Config struct {
	ListenPort string   `json:"listen_port"`
	Backends   []string `json:"backends"`
}

func loadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	configFile := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	var backendURLs []*url.URL
	for _, b := range config.Backends {
		parsed, err := url.Parse(b)
		if err != nil {
			log.Fatalf("Invalid backend URL %s: %v", b, err)
		}
		backendURLs = append(backendURLs, parsed)
		log.Printf("Configured backend: %s", parsed.String())
	}

	checker := health.NewHealthChecker(backendURLs, "/health", 5, 10)
	checker.Start()

	loadBalancer := lb.NewLoadBalancer(backendURLs, checker)

	server := http.Server{
		Addr:    ":" + config.ListenPort,
		Handler: loadBalancer,
	}

	log.Printf("Load balancer started on port %s", config.ListenPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
