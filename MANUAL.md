# Sistem Kontrak Mata Kuliah FILKOM

Manual ini adalah panduan umum untuk keseluruhan project. Dokumentasi teknis detail tersedia juga di:

- [Frontend Manual](./frontend/MANUAL.md)
- [Backend Manual](./backend/MANUAL.md)

## Gambaran Umum

Project ini terdiri dari dua aplikasi terpisah:

- `frontend/`: aplikasi Next.js App Router untuk dashboard kontrak mata kuliah
- `backend/`: REST API berbasis Go + Gin + MySQL

Frontend menampilkan:

- daftar mata kuliah
- planner jadwal mingguan
- counter SKS
- tombol konfirmasi enrollment

Backend menangani:

- login dan JWT
- list curriculum dari database
- proses enrollment dengan validasi
- transaction rollback untuk mencegah partial save

## Struktur Folder

```text
SistemKontrakMatKul_FILKOM/
├─ backend/
├─ frontend/
├─ Curriculum.txt
└─ MANUAL.md
```

## Prasyarat

Pastikan tools ini tersedia:

- Node.js 20+ dan npm
- Go 1.22+
- MySQL 8+

Opsional:

- Docker Desktop, jika ingin menjalankan MySQL/API lewat container

## Menjalankan Project

### 1. Backend

Masuk ke folder backend:

```powershell
cd backend
```

Jalankan API:

```powershell
go run ./cmd/api
```

API default berjalan di:

```text
http://localhost:8080
```

Health check:

```text
GET http://localhost:8080/health
```

### 2. Frontend

Masuk ke folder frontend:

```powershell
cd frontend
npm install
npm run dev
```

Frontend default berjalan di:

```text
http://localhost:3000
```

## Alur Sistem

1. Frontend memuat curriculum dari `GET /api/v1/courses`
2. User memilih mata kuliah di panel kiri
3. Jadwal visual muncul realtime di panel kanan
4. Frontend mengirim enrollment ke backend
5. Backend memvalidasi:
   - batas SKS
   - kuota
   - bentrok jadwal
6. Jika valid, enrollment di-commit
7. Jika gagal, transaction di-rollback

## Data Demo

Akun demo backend:

- `student_number`: `22010001`
- `password`: `password123`

Sample course yang sudah disiapkan:

- `IFAP272` Automata
- `IFMI252` Database

## Endpoint Penting

- `POST /api/v1/login`
- `GET /api/v1/courses`
- `POST /api/v1/enrollments`

## Catatan Penting

- Frontend memakai React Query untuk cache curriculum
- Frontend memakai Zustand untuk state pilihan course
- Backend memakai JWT untuk route enrollment
- Backend memakai DB transaction saat proses enrollment

## Troubleshooting Singkat

Jika frontend menampilkan gagal memuat curriculum:

- pastikan backend hidup di `http://localhost:8080`
- pastikan endpoint `GET /api/v1/courses` mengembalikan `200`
- pastikan database MySQL sudah termigrasi dan berisi sample data

Jika backend gagal compile:

- jalankan `go build ./...` di folder `backend`
- periksa `.env`
- pastikan MySQL aktif dan database `sistem_kontrak` bisa diakses
