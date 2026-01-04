# Generic Concurrent Web Server (Go)
Simple HTTP server in Go using generics, concurrency, background worker, and graceful shutdown as required by the assignment.

## Project structure
```bash
Assignment_2_BeknurErdembek/
├── go.mod
├── main.go
└── internal/
  ├── server/
  ├── store/
  ├── worker/
  └── model/
```

- store: generic in‑memory store Store[K comparable, V any] with mutex and Set/Get/Delete/Snapshot.
- model: request/response models (KeyValue, Stats).
- server: HTTP handlers, request counter, uptime, integration with store and worker.
- worker: background goroutine printing stats every 5 seconds.

## Run
```bash
go run .
```
Server listens on http://localhost:8080

## HTTP API
All data is stored in memory as map[string]string.

### POST /data
*Create or update a key–value pair.*

Request body:
```bash
{
"key": "someKey",
"value": "someValue"
}
```

### GET /data
*Return all stored pairs as a JSON object:*
```bash
{
"k1": "v1",
"k2": "v2"
}
```

### GET /data/{key}
*Return value by key.*
```bash
{
"key": "k1",
"value": "v1"
}
```

### DELETE /data/{key}
*Delete a key–value pair.*

### GET /stats
Return server statistics:
```bash
{
"requests": 100,
"keys": 5,
"uptime_seconds": 42
}
```

## Concurrency, worker, shutdown
- *Store uses sync.RWMutex to protect the internal map from concurrent access.*
- *Request counter is updated atomically so concurrent increments are safe.*
- *Background worker uses time.Ticker (5 seconds) and select on ticker + stop channel to log stats periodically and stop cleanly.*
- *Graceful shutdown catches OS signals (SIGINT, SIGTERM), stops the worker, and calls http.Server.Shutdown with a context timeout.*