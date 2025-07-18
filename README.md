# Hospital Management System API

Ini adalah **backend service** untuk **Sistem Manajemen Rumah Sakit**. Dibuat dengan bahasa Go dan framework Fiber, layanan ini menyediakan berbagai **endpoint** API untuk mengatur data pasien, dokter, janji temu, dan rekam medis.

## Fitur Utama
  - **Manajemen Pengguna & Authentication**:
      - JWT Authentication
      - Manajemen pengguna (CRUD) khusus untuk Admin.
      - *Role-Based Access Control* untuk membatasi akses berdasarkan role (Admin, Doctor, Nurse, dll.).
  - **CRUD Domains**:
      - Kelola data **Pasien**.
      - Kelola data **Dokter**.
      - Kelola **Janji Temu**.
      - Kelola **Rekam Medis**.
  - **Dashboard & Report**:
      - Endpoint khusus untuk menyajikan data statistik dan ringkasan aktivitas.
  - **Keamanan & Audit**:
      - *Soft Delete* untuk data sensitif (pengguna dinonaktifkan, bukan dihapus).
      - *Audit Trail* untuk melacak siapa yang membuat atau mengubah data.

## Teknologi yang Digunakan
  - **Language**: [Go](https://golang.org/)
  - **Web Framework**: [Fiber v2](https://gofiber.io/)
  - **Database**: [MongoDB](https://www.mongodb.com/) (dengan [Official Go Driver](https://github.com/mongodb/mongo-go-driver))
  - **Authentication**: [JWT (golang-jwt)](https://github.com/golang-jwt/jwt)
  - **Validation**: [go-playground/validator](https://github.com/go-playground/validator)
  - **API Documentation**: [Swaggo](https://github.com/swaggo/swag)
  - **Live Reload**: [Air](https://github.com/cosmtrek/air)

## Get Started
### **Prerequisites**
  - [Go](https://golang.org/dl/) versi 1.21 atau lebih baru.
  - [MongoDB](https://www.mongodb.com/try/download/community) atau akun MongoDB Atlas.
  - [Air](https://github.com/cosmtrek/air) untuk *live-reloading* saat development (opsional, tapi direkomendasikan).
    ```bash
    go install github.com/cosmtrek/air@latest
    ```
  - [Swag CLI](https://github.com/swaggo/swag) untuk generate dokumentasi.
    ```bash
    go install github.com/swaggo/swag/v2/cmd/swag@latest
    ```

### **Installation & Running**

1.  **Clone repository ini:**

    ```bash
    git clone https://github.com/yourusername/hms-api.git
    cd hms-api
    ```

2.  **Buat file `.env`:**
    Salin dari contoh yang sudah ada.

    ```bash
    cp .env.example .env
    ```

3.  **Install dependencies:**

    ```bash
    go mod tidy
    ```

4.  **Generate Swagger documentation:**

    ```bash
    swag init -g cmd/api/main.go
    ```

5.  **Running the Server:**

      - **Dengan Live Reload (Direkomendasikan untuk Development):**
        ```bash
        air
        ```
      - **Secara Manual:**
        ```bash
        go run ./cmd/api/main.go
        ```

API akan berjalan di alamat yang kamu tentukan di `.env` (default: `http://localhost:5021`). Dokumentasi Swagger dapat diakses di `http://localhost:5021/api/docs/index.html`.

### **Running the Database Seeder (Opsional)**

Jika kamu butuh data *dummy* untuk development, jalankan script seeder. Script ini akan mengisi database dengan data pasien, dokter, dll.

```bash
go run ./cmd/seed/main.go
```

## Environment Configuration (`.env`)

| Variabel                 | Deskripsi                                                                 | Contoh Nilai                                          |
| ------------------------ | ------------------------------------------------------------------------- | ----------------------------------------------------- |
| `APP_ADDR`               | Alamat dan port tempat server akan berjalan.                              | `:5021`                                               |
| `CORS_ALLOWED_ORIGINS`   | Daftar origin yang diizinkan untuk mengakses API (dipisahkan koma).       | `http://localhost:5173,http://localhost:5021`         |
| `MONGO_ADDR`             | URI koneksi ke MongoDB.                                                   | `mongodb://user:pass@localhost:27017`                 |
| `MONGO_DB`               | Nama database yang akan digunakan di MongoDB.                             | `hospital-management-system`                          |
| `JWT_SECRET`             | Untuk menandatangani JWT.                                                 | `your-very-strong-and-secret-key`                     |
| `INITIAL_ADMIN_EMAIL`    | Email untuk akun admin pertama yang akan dibuat otomatis.                 | `admin@hospital.com`                                  |
| `INITIAL_ADMIN_PASSWORD` | Password untuk akun admin pertama.                                        | `SuperSecurePassword123!`                             |

## Project Structure

```
.
├── cmd/                # Application entrypoints (main.go)
│   ├── api/
│   └── seed/
├── docs/               # Generated files by Swagger
├── internal/           # Main application source code
│   ├── app/            # Server config, router, middleware
│   ├── domain/         # Structs, entities, and DTOs
│   ├── handlers/       # HTTP handlers for processing requests
│   ├── repository/     # Database interaction logic
│   ├── service/        # Core business logic
│   └── utils/          # Helper functions (validator, password, etc.)
├── .air.toml           # Configuration for Air (live reload)
├── .env.example        # Example environment file
├── go.mod              # Project dependencies
└── README.md           # This file
```
