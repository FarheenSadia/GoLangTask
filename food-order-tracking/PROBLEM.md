# Project: Online Food Delivery Order Tracking System

## Objective
Design and specify a **Food Delivery Order Tracking Service** that simulates customers placing food orders and tracks the lifecycle of each order from creation until delivery, using:
- **PostgreSQL** (persistence)
- **Kafka** (event streaming)
- **Go + Fiber + GORM** (service/API layer)
- **Prometheus** (metrics)
- **Goroutines & Channels** (concurrency)
- **Unit Testing** (quality)
- **Containerization** with **Docker Compose or Podman Compose** (all components must run in containers)

---

## Functional Requirements

1. **Order Creation**
   - Accept new food orders with fields:
     - `customer_name`, `address`, `item` (e.g., pizza/burger/sandwich), `size` (small/medium/large).
   - Generate a unique `order_id` for every order.
   - Newly created orders start with status **PLACED**.
   - Emit an **order-created event** to Kafka immediately upon creation.

2. **Order Status Progression**
   - Orders must automatically progress through these statuses (with simulated delays between stages):
     ```
     PLACED → PREPARING → COOKING → OUT_FOR_DELIVERY → DELIVERED
     ```
   - Each transition must be recorded as an **event** with a timestamp.

3. **Event Tracking**
   - Persist an immutable timeline of order events (`order_events`) tied to each `order_id`.
   - Every status change generates a corresponding event entry.

4. **Persistence**
   - Store orders and events in **PostgreSQL** using a relational schema that supports:
     - Fetching current status for an `order_id`.
     - Fetching the **full event history** (timeline) for an `order_id`.
     - Listing recent orders (optionally by `status`).

5. **APIs (Fiber)**
   - `POST /orders` → create order; respond with `{ order_id }`.
   - `GET /orders/:order_id` → return current order info + timeline of events.
   - `GET /orders?status=...&limit=...` → list recent orders (filter by status).
   - `GET /healthz` → health check.
   - `GET /metrics` → Prometheus exposition format.

6. **Event Streaming (Kafka)**
   - **Producer**: publish order-created messages to topic `orders.v1` (key = `order_id`, value = JSON).
   - **Consumer**: consume from `orders.v1` and ensure DB is consistent (idempotent writes recommended).
   - The service must remain consistent under duplicate or out-of-order messages.

---

## Non-Functional Requirements

- **Concurrency**: Use goroutines and channels for:
  - Simulated order generation (optional).
  - Kafka producer/consumer loops.
  - Status progression workers and/or pipelines.
  - DB batching or timed flush (optional).
- **Observability (Prometheus)**:
  - Expose metrics for HTTP requests (count, latency), Kafka (produced/consumed), DB latency/errors, and status transitions.
- **Resilience**:
  - Graceful shutdown (drain channels; close Kafka and DB clients).
  - Basic error handling and retries for transient failures.
- **Testing**:
  - Unit tests for repository, handlers, and the status progression logic.
  - Mock/stub external dependencies where appropriate.
- **Configuration**:
  - All runtime settings via **environment variables** (e.g., DB DSN, Kafka brokers, topic, HTTP port).

---

## Data Model (Required)

- **orders**
  - `order_id` (string, unique, business key)
  - `customer_name` (string)
  - `address` (string)
  - `item` (string)
  - `size` (string)
  - `status` (enum/string: PLACED, PREPARING, COOKING, OUT_FOR_DELIVERY, DELIVERED)
  - `created_at` (timestamptz)
  - `updated_at` (timestamptz)

- **order_events**
  - `id` (bigserial PK)
  - `order_id` (FK → orders.order_id)
  - `event` (string; e.g., CREATED, PLACED, PREPARING, COOKING, OUT_FOR_DELIVERY, DELIVERED)
  - `timestamp` (timestamptz)
  - `meta` (JSON; optional)

**Indexes**:
- `orders(order_id)` unique
- `order_events(order_id)` index
- Consider `orders(status, updated_at)` for listing/filtering.

---

## Metrics (Minimum Set)

- **HTTP**
  - `http_requests_total{route,method,status}`
  - `http_request_duration_seconds_bucket{route,method,status}`
- **Kafka**
  - `orders_produced_total`
  - `orders_consumed_total`
- **Pipeline**
  - `status_transitions_total{from,to}`
  - `pipeline_step_duration_seconds_bucket{step}`
- **Database**
  - `db_write_latency_seconds_bucket`
  - `db_errors_total`

---

## Required Files & Structure (Suggestion)

/ (repo root)
README.md                       ← overview & run instructions
PROBLEM.md                      ← this specification (deliver with submission)
docker-compose.yml              ← full multi-service stack
.env.example                    ← example env vars for compose
/deploy/prometheus/prometheus.yml  ← scrape config
/migrations/                    ← SQL or auto-migration notes (optional)
/cmd/server/                    ← service entrypoint (code, not required here)
/internal/…                   ← modules (code, not required here)
/tests/…                      ← unit tests (code, not required here)
Dockerfile                      ← service container build

> This structure is indicative; you may organize differently if all **deliverables** below are met.

---

## **Compose File Task (Mandatory)**

Create a **`docker-compose.yml`** (or `podman-compose.yml`) that runs **all components in containers** and can be brought up with a single command.

**Requirements:**
1. **Services** (minimum):
   - **app**: the Go service container (built from `Dockerfile`).
   - **postgres**: PostgreSQL with persistent volume and healthcheck.
   - **kafka**: Kafka broker (you may use a single-node image with internal ZK or KRaft).
   - **prometheus**: Prometheus server scraping the app’s `/metrics`.
2. **Networking**:
   - All services on an isolated user-defined bridge network.
3. **Volumes**:
   - Persistent named volumes for PostgreSQL data.
   - Bind-mount a `prometheus.yml` config (e.g., `./deploy/prometheus/prometheus.yml`).
4. **Environment**:
   - Use a `.env` file or in-file `environment` sections for:
     - `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`
     - `KAFKA_BROKERS` (e.g., `kafka:9092`)
     - `PG_DSN` (service DSN pointing at **postgres** container)
     - `APP_PORT` (e.g., `8080`)
   - Ensure the app reads these from environment.
5. **Healthchecks**:
   - `postgres`: use `pg_isready`.
   - `app`: HTTP GET `/healthz`.
   - `kafka`: a reasonable TCP or script check (optional but preferred).
6. **Dependency Ordering**:
   - Use `depends_on` to start `postgres` and `kafka` before `app`.
   - (Compose does not guarantee readiness; combine with healthchecks and app-level retries.)
7. **Prometheus Configuration**:
   - `prometheus.yml` must scrape the `app` container at `/metrics` on the configured port.
8. **(Optional but Recommended)**:
   - A Kafka UI (e.g., `kafka-ui`) service for easy topic inspection.
   - Grafana connected to Prometheus for dashboards.

**Acceptance for Compose Task**:
- `compose up` starts all components.
- `app` can create/read orders against `postgres`.
- `kafka` topic `orders.v1` exists (you may create it at startup).
- `prometheus` can scrape `app:/metrics`.
- Stopping and starting the stack preserves PostgreSQL data (via volume).

---

## **Dockerfile Task (Mandatory)**

Create a **`Dockerfile`** that:
- Builds the Go service as a static binary suitable for container execution.
- Uses a multi-stage build (builder + minimal runtime image).
- Exposes the service port (e.g., `8080`).
- Runs the service with environment variables injected at runtime.
- Does not include dev tools in the final image.
- Is compatible with `docker compose` or `podman compose` build flows.

**Acceptance for Dockerfile Task**:
- `docker compose build` (or `podman compose build`) successfully builds the image.
- The container starts and serves `/healthz` and `/metrics` when run within the compose stack.

---

## API Contract (Must-Have)

- `POST /orders`
  - Request: `{ "customer_name": "...", "address": "...", "item": "...", "size": "small|medium|large" }`
  - Response: `201 { "order_id": "ORD-..." }`
- `GET /orders/:order_id`
  - Response: `200 { "order": {...}, "events": [ ... ] }` (or `404` if not found)
- `GET /orders?status=...&limit=...`
  - Response: `200 [ {...}, ... ]`
- `GET /healthz` → `200 {"status":"ok"}`
- `GET /metrics` → Prometheus exposition format

---

## Event Streaming (Must-Have)

- Topic: `orders.v1`
- Key: `order_id`
- Value: JSON message with order details and creation timestamp.
- **Producer**: publish on order creation.
- **Consumer**: consume and ensure DB is up-to-date (idempotent behavior recommended).

---

## Testing (Must-Have)

- Unit tests for:
  - **Repository/DB** operations (create order, append event, update status, fetch order + events).
  - **Handlers** (validate request/response codes and payloads).
  - **Status progression** logic (simulate transitions over time).
- Tests should run locally without requiring the full stack (use mocks/stubs or in-memory DB for unit tests).

---

## Evaluation & Acceptance Criteria

1. **Functionality**
   - Can create orders and retrieve them by `order_id`.
   - Orders progress through all statuses automatically and record events.
   - List endpoint returns filtered results by status.
2. **Data Integrity**
   - Orders and events are persisted in PostgreSQL.
   - Event history is complete and ordered by time.
3. **Streaming**
   - Order creation events are produced to Kafka.
   - Consumer processes events and keeps DB consistent.
4. **Observability**
   - `/metrics` exposes HTTP, Kafka, DB, and pipeline metrics.
   - Prometheus successfully scrapes the app.
5. **Containerization**
   - `docker compose up` (or `podman compose up`) starts all required services.
   - Data persists across restarts via volumes.
6. **Quality**
   - Unit tests exist and pass.
   - Graceful shutdown and basic error handling are implemented.
7. **Configuration**
   - All runtime parameters are configurable via environment variables.

---

## Submission Checklist

- ✅ `PROBLEM.md` (this spec) included in the repo.
- ✅ `docker-compose.yml` with required services and healthchecks.
- ✅ `Dockerfile` (multi-stage).
- ✅ `.env.example` with all necessary variables.
- ✅ `deploy/prometheus/prometheus.yml` scraping the app.
- ✅ README with clear **run instructions** (build, compose up, endpoints).
- ✅ Source code and unit tests (in appropriate directories).
- ✅ Screenshot or logs demonstrating:
  - Successful service startup in compose.
  - Prometheus scraping `/metrics`.
  - Order lifecycle transitions visible via API or logs.

---