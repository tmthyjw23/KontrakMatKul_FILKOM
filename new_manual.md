# Panduan Menjalankan Project KontrakMatKul_FILKOM

Panduan ini ditujukan bagi siapa saja (rekan tim, dosen, atau penguji) yang baru melakukan clone repositori ini dan ingin menjalankan program secara lokal di komputer masing-masing.

## 🛠 Prasyarat Sistem
Sebelum memulai, pastikan komputer Anda sudah terinstal:
1. **Golang** (Versi 1.22 atau terbaru) - untuk menjalankan Backend.
2. **Node.js** (Versi 18 atau terbaru) & **npm** - untuk menjalankan Frontend (Next.js).
3. **MySQL Server** - untuk database (bisa menggunakan XAMPP, Laragon, MySQL Installer, atau Docker).
4. **Git** - untuk melakukan clone repositori.

---

## 🚀 Langkah-langkah Menjalankan Program

### 1. Persiapan Database (MySQL)
1. Nyalakan service MySQL di komputer Anda (misalnya klik "Start" MySQL di XAMPP Control Panel).
2. Buka terminal MySQL, phpMyAdmin, atau DBeaver, lalu buat database baru dengan nama `kontrak_matkul`.
   ```sql
   CREATE DATABASE kontrak_matkul;
   ```

### 2. Setup Backend (Golang)
Buka terminal baru dan masuk ke folder `backend`:
```bash
cd backend
```

**A. Konfigurasi Environment:**
Buat file baru bernama `.env` di dalam folder `backend` (sejajar dengan `main.go`). Isi dengan konfigurasi database Anda. Jika Anda menggunakan XAMPP (user `root` tanpa password), isinya seperti ini:
```env
PORT=8080
DATABASE_URL=root:@tcp(127.0.0.1:3306)/kontrak_matkul?parseTime=true
```
*(Catatan: Sesuaikan `root:@` menjadi `root:password_anda@` jika MySQL Anda memiliki password).*

**B. Install Dependencies:**
```bash
go mod tidy
```

**C. Jalankan Seeder & Migrasi Database:**
Perintah ini akan secara otomatis membuat tabel-tabel yang dibutuhkan dan mengisi data awal (seperti daftar mata kuliah, akun admin, dll).
```bash
go run ./cmd/seeder/main.go
```
*Pastikan muncul pesan "🎉 Seeding complete!" tanpa error.*

**D. Jalankan Server Backend:**
```bash
go run main.go
```
*Backend sekarang berjalan di http://localhost:8080*

---

### 3. Setup Frontend (Next.js)
Buka terminal baru (biarkan terminal backend tetap berjalan), lalu masuk ke folder `frontend`:
```bash
cd frontend
```

**A. Konfigurasi Environment (Opsional):**
Secara default, frontend sudah dikonfigurasi untuk menembak ke `http://localhost:8080/api/v1`. Jika Anda ingin mengubahnya, buat file `.env.local` di folder `frontend`:
```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

**B. Install Dependencies:**
```bash
npm install
```

**C. Jalankan Server Frontend:**
```bash
npm run dev
```

---

## 🌐 Mengakses Aplikasi
Buka browser Anda dan kunjungi:
👉 **http://localhost:3000**

### Akun Uji Coba (Testing Accounts)
Karena Anda sudah menjalankan Seeder di langkah 2C, sistem sudah memiliki beberapa akun default yang bisa langsung digunakan:

**Akun Admin:**
- **NIM/Username:** `admin`
- **Password:** `admin123`

**Akun Student (Mahasiswa):**
- **NIM/Username:** `22515040000001`
- **Password:** `password123`

---

## ❓ Troubleshooting (Masalah Umum)

1. **Error: `Port 3000 is in use` saat menjalankan Frontend**
   Ini berarti ada aplikasi lain (atau sisa proses Next.js sebelumnya) yang masih menggunakan port 3000.
   - **Solusi Windows:** Buka terminal dan ketik `taskkill /PID <angka_pid_yang_muncul> /F`
   - **Solusi Alternatif:** Biarkan saja, Next.js akan otomatis menggunakan port 3001 (`http://localhost:3001`).

2. **Error: `Access denied for user 'root'@'localhost'` saat seeder**
   - **Solusi:** Cek kembali file `.env` di folder `backend`. Pastikan username dan password MySQL sudah sesuai dengan komputer yang sedang digunakan.

3. **Backend error kompilasi `declared and not used` atau `undefined`**
   - **Solusi:** Pastikan Anda menarik (pull) branch `jordy-combine` yang paling terbaru, karena bug kompilasi sudah diperbaiki di commit terakhir.

4. **Error: `The system cannot find the path specified` saat menjalankan seeder**
   - **Solusi:** Pastikan Anda menjalankan perintah `go run ./cmd/seeder/main.go` **SAAT POSISI TERMINAL BERADA DI DALAM FOLDER `backend`**, bukan di luar folder utama proyek.

---

## ✨ Fitur Utama yang Tersedia

1. **Student Dashboard:**
   - Pemilihan mata kuliah interaktif dengan sistem drag-and-drop/klik.
   - Deteksi bentrok jadwal secara real-time.
   - Batasan maksimal 24 SKS (termasuk yang sudah didaftarkan sebelumnya).
   - Notifikasi jika periode pendaftaran (Contract Period) sedang ditutup oleh Admin.

2. **Admin Management:**
   - **Monitoring:** Menyetujui (Approve) atau Menolak (Reject) pendaftaran mahasiswa.
   - **User Management:** Membuat akun mahasiswa baru, reset password, dan menghapus user.
   - **Contract Period:** Kontrol penuh untuk membuka atau menutup akses pendaftaran bagi seluruh mahasiswa.

---

## 📝 Alur Kerja Pendaftaran (Workflow)
1. **Admin** membuka periode pendaftaran (Contract Period: ON).
2. **Mahasiswa** login dan memilih mata kuliah di Dashboard sampai menekan "Confirm Enrollment".
3. Status pendaftaran mahasiswa akan menjadi **Pending**.
4. **Admin** masuk ke menu Monitoring dan melakukan **Approve** pada mata kuliah tersebut.
5. Setelah di-approve, mata kuliah akan muncul di **Weekly Schedule** mahasiswa dan SKS-nya akan terkunci (mahasiswa tidak bisa menambah mata kuliah lagi jika sudah mencapai 24 SKS).

---

*Selamat mencoba! Jika menemukan kendala teknis, silakan hubungi tim pengembang.*
