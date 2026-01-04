package server

import (
	"context"
	"encoding/json"
	//"encoding/json"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/yerdembek/Assignment_2_BeknurErdembek/internal/model"
	"github.com/yerdembek/Assignment_2_BeknurErdembek/internal/store"
	"github.com/yerdembek/Assignment_2_BeknurErdembek/internal/worker"
)

type Server struct {
	httpServer *http.Server

	store   *store.Store[string, string]
	started time.Time
	reqs    int64

	worker *worker.Worker
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

	w := worker.New(srv)
	w.Start()
	srv.worker = w

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

func (s *Server) handleData(w http.ResponseWriter, r *http.Request) {
	s.incrementRequests()
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		var body model.KeyValue
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if body.Key == "" || body.Value == "" {
			http.Error(w, "key and value are required", http.StatusBadRequest)
			return
		}

		s.store.Set(body.Key, body.Value)

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(body); err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
		}

	case http.MethodGet:
		data := s.store.Snapshot()
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
		}

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleDataKey(w http.ResponseWriter, r *http.Request) {
	s.incrementRequests()
	w.Header().Set("Content-Type", "application/json")

	key := r.URL.Path[len("/data/"):]
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		value, ok := s.store.Get(key)
		if !ok {
			http.NotFound(w, r)
			return
		}
		resp := model.KeyValue{Key: key, Value: value}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
		}

	case http.MethodDelete:
		ok := s.store.Delete(key)
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	s.incrementRequests()
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := model.Stats{
		Requests:      s.Requests(),
		Keys:          s.KeysCount(),
		UptimeSeconds: s.UptimeSeconds(),
	}

	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.worker != nil {
		s.worker.Stop()
	}
	return s.httpServer.Shutdown(ctx)
}
