# Backend Manual

Dokumentasi ini khusus untuk aplikasi backend di folder `backend/`.

## Stack

- Go 1.22+
- Gin
- MySQL
- Zap Logger
- JWT

## Tujuan Backend

Backend menyediakan REST API untuk:

- login user
- mengambil daftar course
- menyimpan enrollment
- menjalankan validasi bisnis akademik

## Struktur Utama

```text
backend/
├─ cmd/api/
├─ internal/
│  ├─ config/
│  ├─ delivery/http/
│  ├─ domain/
│  ├─ repository/mysql/
│  └─ usecase/
├─ migrations/
├─ pkg/
└─ MANUAL.md
```

Makna layer:

- `config`: load environment
- `delivery/http`: handler, middleware, routes
- `domain`: model, contract repository, contract service
- `repository/mysql`: query database
- `usecase`: business logic
- `pkg/database`: bootstrap connection pool
- `pkg/logger`: bootstrap Zap

## Cara Menjalankan

```powershell
cd backend
go run ./cmd/api
```

Untuk compile check:

```powershell
go build ./...
```

## Environment

File yang dipakai:

- `.env`
- `.env.example`

Variabel penting:

- `APP_PORT`
- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `JWT_SECRET`

## Database

Database default:

```text
sistem_kontrak
```

Migration yang sudah ada:

- `000001_create_core_tables`
- `000002_add_role_and_password_to_users`
- `000003_add_lecturer_and_room_to_courses_and_schedules`

## Tabel Inti

- `users`
- `courses`
- `schedules`
- `enrollments`

## Endpoint

### 1. Health Check

```http
GET /health
```

### 2. Login

```http
POST /api/v1/login
```

Payload:

```json
{
  "student_number": "22010001",
  "password": "password123"
}
```

### 3. List Courses

```http
GET /api/v1/courses
```

Ini adalah endpoint yang dipakai frontend untuk memuat curriculum.

### 4. Enroll Course

```http
POST /api/v1/enrollments
```

Header:

```text
Authorization: Bearer <token>
```

Payload:

```json
{
  "course_id": 1
}
```

## Business Logic Enrollment

Usecase enrollment melakukan:

- mulai transaction DB
- cek user dan course
- cek duplicate enrollment
- cek total SKS vs `max_sks`
- cek kuota course
- cek bentrok jadwal
- insert enrollment
- commit atau rollback

Jika salah satu validasi gagal, perubahan tidak disimpan.

## Middleware

### CORS

Mengizinkan origin:

```text
http://localhost:3000
```

### JWT

Middleware JWT:

- membaca `Authorization: Bearer <token>`
- memvalidasi signature
- menyimpan `user_id` ke Gin context

## Logging

Backend memakai Zap production logger.

File:

- `pkg/logger/zap.go`

## Sample User

User demo:

- `student_number`: `22010001`
- `password`: `password123`

## Troubleshooting

### `GET /api/v1/courses` error 500

Biasanya salah satu dari:

- tabel database belum ada
- migration belum diterapkan
- database masih kosong

### JWT unauthorized

Periksa:

- token ada
- `JWT_SECRET` backend sesuai saat token dibuat

### Backend tidak bisa konek ke MySQL

Periksa:

- MySQL aktif
- `.env` benar
- user punya akses ke database

## Testing Manual

### Login

```powershell
curl.exe -X POST http://localhost:8080/api/v1/login `
  -H "Content-Type: application/json" `
  -d "{\"student_number\":\"22010001\",\"password\":\"password123\"}"
```

### Courses

```powershell
curl.exe http://localhost:8080/api/v1/courses
```

### Enrollment

```powershell
curl.exe -X POST http://localhost:8080/api/v1/enrollments `
  -H "Content-Type: application/json" `
  -H "Authorization: Bearer <TOKEN>" `
  -d "{\"course_id\":1}"
```
