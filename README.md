# Industrial Fleet Management System

Sistem backend berbasis **Modular Monolith** menggunakan **Go (Golang)** yang dirancang untuk mengelola aset alat berat industri pertambangan dan memproses data telemetri secara *reliable* dan *asynchronous*.

## 📌 Latar Belakang

Perusahaan pertambangan mengoperasikan banyak alat berat di berbagai lokasi kerja. Setiap unit alat berat menghasilkan data operasional penting yang perlu dicatat dan dipantau secara berkala, antara lain:
* **Operating hours** (hour meter)
* **Fuel level**
* **Engine temperature**
* **Equipment status** (Active, Idle, Breakdown, Maintenance)
* **Lokasi equipment** (koordinat GPS / area kerja)
* **Maintenance history**

**Masalah Utama:** Bagaimana membangun infrastruktur *backend* yang dapat mengelola identitas alat berat serta menerima dan memproses ingestion data telemetri berfrekuensi tinggi secara aman, handal (*reliable*), dan terstruktur.

## 🎯 Project Objectives

Membangun backend *modular monolith* yang mampu:
1. **Mengelola Equipment:** Menangani identitas unik, status, dan lokasi setiap unit alat berat.
2. **Telemetry Ingestion:** Menerima data telemetri dari alat berat atau simulator melalui HTTP API.
3. **Asynchronous Processing:** Memvalidasi dan memproses pipeline data telemetri secara *asynchronous* menggunakan Go worker pool.
4. **State & Historical Management:** Menyimpan data histori telemetri sekaligus kondisi terkini (*current state*) dari setiap equipment.
5. **Maintenance Management:** Mengelola jadwal, pesanan kerja (*work order*), dan histori perawatan alat berat.
6. **Auth & Access Control:** Menyediakan mekanisme autentikasi dan otorisasi pengguna (RBAC).
7. **Real-time Data Delivery:** Menyediakan dukungan *real-time updates* menggunakan Redis dan WebSocket pada fase pengembangan berikutnya.

## 🚀 Scope & Tech Stack

Proyek ini berfokus pada **Backend-First Development** dengan pilihan teknologi sebagai berikut:

| Komponen | Keputusan Teknologi / Arsitektur |
| :--- | :--- |
| **Language** | Go (Golang) |
| **Architecture** | Modular Monolith |
| **HTTP Router** | Chi + `net/http` |
| **Database** | PostgreSQL |
| **Database Driver** | `pgx` |
| **Cache & Coordination** | Redis |
| **Authentication** | JWT / Session (ditentukan pada tahap desain auth) |
| **Telemetry Transport** | HTTP Ingestion API |
| **Async Processing** | Go Worker Pool |
| **Testing** | Go standard `testing` package |
| **Containerization** | Docker & Docker Compose |
| **Frontend** | React.js (Fase Berikutnya) |
| **API Documentation** | OpenAPI (Swagger) |

### ⚠️ Batasan Penting (Out of Scope Phase 1)
Untuk tetap fokus pada *backend fundamentals* dan *system reliability*:
* **TIDAK** menggunakan *Microservices architecture*.
* **TIDAK** menggunakan message broker kompleks seperti *Apache Kafka* di tahap awal.
* **TIDAK** menerapkan *Kubernetes (K8s)*.
* **TIDAK** tersambung ke perangkat hardware IoT fisik sungguhan. Data telemetri diproduksi menggunakan **simulator internal**.

## 📂 Project Structure

Struktur direktori proyek mengikuti pola *Modular Monolith* berbasis domain internal:

```text
industrial-fleet-management/
│
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
│
├── internal/                 # Private application and code logic
│   ├── equipment/            # Domain: Equipment Management
│   │   ├── domain/           # Entities & Interface contracts
│   │   ├── service/          # Business logic
│   │   ├── repository/       # Database queries & data access
│   │   └── handler/          # HTTP Handlers / Controllers
│   │
│   ├── telemetry/            # Domain: Telemetry Ingestion & Processing
│   │   ├── domain/
│   │   ├── service/
│   │   ├── repository/
│   │   ├── handler/
│   │   └── worker/           # Async Go Worker Pool for telemetry
│   │
│   ├── maintenance/          # Domain: Maintenance & Work Orders
│   │   ├── domain/
│   │   ├── service/
│   │   ├── repository/
│   │   └── handler/
│   │
│   ├── platform/             # Infrastructure components & shared packages
│   │   ├── database/         # PostgreSQL connection (`pgx`)
│   │   ├── redis/            # Redis client wrapper
│   │   ├── config/           # Environment config loader
│   │   └── logger/           # Structured logging engine
│   │
│   └── app/                  # Application bootstrap & DI container
│       └── app.go
│
├── migrations/               # Database migration files (.sql)
│
├── deployments/              # Deployment configurations
│   └── docker/               # Extra Docker configs
│
├── docs/                     # Comprehensive documentation
│   ├── architecture.md       # Architecture specification
│   ├── domain-model.md       # Domain context & entities diagram
│   └── api-specification.md  # OpenAPI specification
│
├── scripts/                  # Helper scripts (seeding, data simulator)
│
├── .env.example              # Environment variables template
├── .gitignore
├── Dockerfile                # Multi-stage Go Dockerfile
├── docker-compose.yml        # PostgreSQL, Redis & App setup
├── Makefile                  # Build & run automation tasks
├── README.md
└── go.mod
```

## 🛠️ Quick Start

### Prerequisites
* Go 1.27+
* Docker & Docker Compose
* Make (opsional)

### Setup & Running Locally
1. **Clone repository:**
   ```bash
   git clone https://github.com/your-username/industrial-fleet-management.git
   cd industrial-fleet-management
   ```

2. **Salin environment file:**
   ```bash
   cp .env.example .env
   ```

3. **Jalankan infrastruktur (PostgreSQL & Redis):**
   ```bash
   docker compose up -d
   ```

4. **Jalankan aplikasi Go:**
   ```bash
   go run cmd/api/main.go
   ```