package server

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/Assignment_2_BeknurErdembek/internal/model"
	"github.com/Assignment_2_BeknurErdembek/internal/store"
)

type Server struct {
	httpServer *http.Server

	store   *store.Store[string, string]
	started time.Time
	reqs    int64
}

func New(addr string, s *store.Store[string, string]) *Server {
	srv := &Server{
		store:   s,
		started: time.Now(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/data", srv.handleData)
	mux.HandleFunc("/data/", srv.handleDataKey)
	mux.HandleFunc("/stats", srv.handleStats)

	srv.httpServer = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return srv
}

func (s *Server) incrementRequests() {
	atomic.AddInt64(&s.reqs, 1)
}

func (s *Server) Requests() int64 {
	return atomic.LoadInt64(&s.reqs)
}

func (s *Server) KeysCount() int {
	snap := s.store.Snapshot()
	return len(snap)
}

func (s *Server) UptimeSeconds() int64 {
	return int64(time.Since(s.started).Seconds())
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

// Заглушки для handler'ов — заполним позже
func (s *Server) handleData(w http.ResponseWriter, r *http.Request) {
	s.incrementRequests()
	// TODO: реализовать POST /data и GET /data
}

func (s *Server) handleDataKey(w http.ResponseWriter, r *http.Request) {
	s.incrementRequests()
	// TODO: реализовать GET /data/{key} и DELETE /data/{key}
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	s.incrementRequests()
	// TODO: вернуть JSON model.Stats
}
