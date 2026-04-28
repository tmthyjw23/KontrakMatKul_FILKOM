Panduan Testing Manual
Prasyarat
Pastikan dua terminal terbuka, masing-masing menjalankan:

Terminal 1 — Backend:


cd backend
go run main.go
# Harus muncul: ✅ Database connection established
#               🚀 Server listening on http://localhost:8080
Terminal 2 — Frontend:


cd frontend
npm run dev
# Harus muncul: ✓ Ready in ... http://localhost:3000
Test 1: Login sebagai Student
Buka http://localhost:3000/login
Pilih Student
Masukkan NIM: 22010001, Password: password123
Klik Login
Ekspektasi: Toast "Welcome, Budi Santoso!" muncul, redirect ke /student
Test 2: Halaman Dashboard Student — Kontrak Mata Kuliah
Setelah login sebagai student, Anda berada di /student
Ekspektasi: Daftar 124 mata kuliah muncul di panel kiri
Klik beberapa mata kuliah — akan masuk ke planner mingguan kanan
Counter SKS di atas grid ikut bertambah
Test 3: Submit Enrollment (Kontrak)
Pilih beberapa mata kuliah dari panel kiri
Pastikan total SKS tidak melebihi 24
Klik Confirm Enrollment
Ekspektasi: Toast berhasil, data terkirim ke backend
Test 4: Halaman Jadwal Student
Pergi ke http://localhost:3000/student/schedule
Ekspektasi: Kursus yang sudah diregistrasi muncul di section "Enrolled Courses"
Section "Registration History" menampilkan semua registrasi beserta statusnya
Test 5: Login sebagai Admin
Kembali ke http://localhost:3000/login
Pilih Admin
Masukkan username: ADMIN001, Password: admin123
Ekspektasi: Redirect ke /admin — dashboard admin dengan 4 menu
Test 6: Admin — Course Management
Masuk ke /admin/courses
Ekspektasi: Semua 124 mata kuliah tampil dalam tabel
Klik + New Course → isi form → Save → kursus baru muncul
Klik Edit pada kursus manapun → ubah data → Save
Klik Delete → konfirmasi → kursus terhapus dari tabel
Test 7: Admin — Contract Monitoring
Masuk ke /admin/monitoring
Ekspektasi: Tabel registrasi muncul dengan status (registered/cancelled)
Filter ke "Pending" — tombol Approve/Reject muncul
Test 8: Admin — User Management
Masuk ke /admin/users
Ekspektasi: Daftar 4 student tampil
Klik + New Student → isi form (NIM unik) → Create
⚠️ Yang belum bisa ditest (Backend masih stub)
Fitur	Status
Update kursus tersimpan ke DB	Stub — response OK tapi DB tidak berubah
Approve/Reject registration tersimpan	Stub
Create student tersimpan ke DB	Stub
Reset Password	404 — belum ada route di backend
Contract period open/close	Stub — hardcoded buka
Jadwal mingguan terisi	Tidak ada data jadwal di course list
Semua item di atas perlu Antigravity menyelesaikan implementasi backend (ubah stub menjadi query DB nyata + tambah route reset-password di main.go).