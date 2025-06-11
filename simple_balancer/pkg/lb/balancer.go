package lb

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"simple_balancer/internal/health"
	"sync/atomic"
)

type LoadBalancer struct {
	hc           *health.HealthChecker
	proxies      map[string]*httputil.ReverseProxy
	currentIndex int64
}

func NewLoadBalancer(backends []*url.URL, hc *health.HealthChecker) *LoadBalancer {
	proxies := make(map[string]*httputil.ReverseProxy)
	for _, b := range backends {
		proxies[b.String()] = httputil.NewSingleHostReverseProxy(b)
	}
	return &LoadBalancer{hc: hc, proxies: proxies}
}

// выбор следующего доступного бэкенда по алгоритму Round Robin
func (lb *LoadBalancer) getNextBackend() (string, error) {
	backends := lb.hc.GetActive()
	if len(backends) == 0 {
		return "", errors.New("no healthy backends")
	}
	idx := atomic.AddInt64(&lb.currentIndex, 1)
	return backends[int(idx)%len(backends)].String(), nil
}

// обработка входящих запросов (возвращает статус или проксирует)
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.Write([]byte(`{"status":"ok"}`))
		return
	}
	backendURL, err := lb.getNextBackend()
	if err != nil {
		http.Error(w, "No backend available", http.StatusServiceUnavailable)
		log.Printf("Error: %v", err)
		return
	}
	lb.proxies[backendURL].ServeHTTP(w, r)
}
