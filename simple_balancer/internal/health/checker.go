package health

import (
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// отслеживание состояния бэкенда
type HealthChecker struct {
	backends   []*url.URL
	healthPath string
	timeout    time.Duration // таймаут запроса
	interval   time.Duration // интервал между проверками
	mu         sync.RWMutex
	active     map[*url.URL]bool // здоровые бэкенды
}

func NewHealthChecker(backends []*url.URL, healthPath string, timeoutSec, intervalSec int64) *HealthChecker {
	hc := &HealthChecker{
		backends:   backends,
		timeout:    time.Duration(timeoutSec) * time.Second,
		interval:   time.Duration(intervalSec) * time.Second,
		healthPath: healthPath,
		active:     make(map[*url.URL]bool),
	}
	for _, b := range backends {
		hc.active[b] = true
	}
	return hc
}

// периодическая проверка бэкендов
func (hc *HealthChecker) Start() {
	ticker := time.NewTicker(hc.interval)
	go func() {
		for range ticker.C {
			hc.checkAll()
		}
	}()
}

func (hc *HealthChecker) checkAll() {
	for _, backend := range hc.backends {
		go hc.checkOne(backend)
	}
}

func (hc *HealthChecker) checkOne(backend *url.URL) {
	client := http.Client{Timeout: hc.timeout}
	resp, err := client.Get(backend.String() + hc.healthPath)

	hc.mu.Lock()
	defer hc.mu.Unlock()
	if err != nil {
		log.Printf("Health check error for %s: %v", backend, err)
		hc.active[backend] = false
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Health check failed for %s: status %d", backend, resp.StatusCode)
		hc.active[backend] = false
		return
	}
	hc.active[backend] = true
	log.Printf("Backend %s is healthy", backend)
}

func (hc *HealthChecker) GetActive() []*url.URL {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	var actives []*url.URL
	for b, ok := range hc.active {
		if ok {
			actives = append(actives, b)
		}
	}
	return actives
}
