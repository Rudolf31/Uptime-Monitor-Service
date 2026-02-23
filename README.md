# Uptime Monitor Service

A simple, self-hosted uptime monitoring service written in Go.
It provides a REST API to create, list, retrieve and delete monitors, and periodically checks registered URLs for availability and response time.

Monitors are checked concurrently in a worker pool, with next check scheduling and history tracking.
Designed to be lightweight and easy to extend.

---

## Features

* Create, list, get and delete uptime monitors via HTTP API
* Periodic background checks with concurrent workers
* Stores monitor metadata, status, next check time and history
* Thread-safe in-memory storage for concurrent access
* Designed for simplicity and clarity

---

## Architecture Overview

The service consists of multiple logical components:

* **Handlers** — REST API for monitor management
* **Storage** — In-memory storage with RWMutex protection
* **Scheduler** — Dispatches check jobs when monitors are due
* **Worker Pool** — Concurrent workers that perform HTTP checks
* **Checker** — Determines status (`up`/`down`) and response time

---

## API Endpoints

| Method | Path             | Description            |
| ------ | ---------------- | ---------------------- |
| POST   | `/monitors`      | Create a new monitor   |
| GET    | `/monitors`      | List all monitors      |
| GET    | `/monitors/{id}` | Get a monitor by ID    |
| DELETE | `/monitors/{id}` | Delete a monitor by ID |

---

## Example Monitor

```json
POST /monitors
{
  "url": "https://example.com",
  "interval": 30
}
```

Response includes assigned `id`, initial status and scheduling info.

---

## Requirements

* Go 1.20+
* Environment with network access for external HTTP checks

---

## Getting Started

### Clone the repository

```bash
git clone https://github.com/Rudolf31/Uptime-Monitor-Service.git
cd Uptime-Monitor-Service
```

### Build

```bash
go build ./cmd/server
```

This produces a `server` binary.

---

## Running

### Default mode

```bash
./server
```

By default, the server listens on `:8080` (can be changed via `PORT` env).

---

## Environment

| Variable | Default | Description      |
| -------- | ------- | ---------------- |
| `PORT`   | `8080`  | HTTP server port |

Example:

```bash
PORT=9090 ./server
```

---

## Using the API

### Create a Monitor

```bash
curl -X POST http://localhost:8080/monitors \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com","interval":30}'
```

### List Monitors

```bash
curl http://localhost:8080/monitors
```

### Get a Monitor

```bash
curl http://localhost:8080/monitors/1
```

### Delete a Monitor

```bash
curl -X DELETE http://localhost:8080/monitors/1
```

---

## Testing

The project includes:

* API tests for handlers
* Concurrency tests for storage and handlers
* Unit tests for scheduler logic

Run all tests with the race detector:

```bash
go test ./... -race
```

---

## License

This project is open source and available under the MIT License.
