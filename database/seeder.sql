-- Pastikan kita memasukkan data ke database yang benar
USE db_curriculum;

-- Matikan sementara sabuk pengaman relasi (Foreign Key)
SET FOREIGN_KEY_CHECKS = 0;

-- Gunakan DELETE FROM agar lebih aman terhadap relasi tabel
DELETE FROM kontrak_krs;
DELETE FROM jadwal_kelas;
DELETE FROM kurikulum_rekomendasi;
DELETE FROM mata_kuliah;
DELETE FROM mahasiswa;

-- Nyalakan kembali sabuk pengamannya
SET FOREIGN_KEY_CHECKS = 1;

-- ==========================================
-- 1. DATA MAHASISWA DUMMY
-- ==========================================
INSERT INTO mahasiswa (nim, nama, fakultas, program_studi, angkatan) VALUES
('101010', 'Timothy Frontend', 'Fasilkom', 'Informatika', 2023),
('202020', 'Farlen Database', 'Fasilkom', 'Informatika', 2023);

-- ==========================================
-- 2. DATA MATA KULIAH (HANYA 3 KOLOM: Kode, Nama, SKS)
-- ==========================================
INSERT INTO mata_kuliah (kode_mk, nama_mk, sks) VALUES
-- PRA-SEMESTER
('MK001', 'Matematika', 3),
('MK002', 'Keterampilan Komputer Dasar', 2),
('MK003', 'Bahasa Inggris Dasar', 2),
('MK004', 'Pendidikan Keterampilan', 2),
('MK005', 'Physical and Health Education', 2),

-- SEMESTER 1
('MK101', 'Introduction to Programming',  3),
('MK102', 'Introduction to Computational Thinking', 2),
('MK103', 'Dasar Aljabar Linear', 3),
('MK104', 'Pendidikan Pancasila', 2),
('MK105', 'Teladan Kehidupan I', 2),
('MK106', 'Bahasa Inggris Pra Dasar', 2),

-- SEMESTER 2
('MK201', 'Struktur Data dan Algoritma', 3),
('MK202', 'Matematika Diskrit', 3),
('MK203', 'Logika Informatika', 3),
('MK204', 'Bahasa Indonesia', 2),
('MK205', 'Pendidikan Kewarganegaraan', 2),
('MK206', 'Teladan Kehidupan II', 2),
('MK207', 'Filsafat Pendidikan Kristen', 2),
('MK208', 'Bahasa Inggris Dasar', 2),

-- SEMESTER 3
('MK301', 'Organisasi dan Arsitektur Komputer', 3),
('MK302', 'Perancangan Web', 3),
('MK303', 'Pengantar Basisdata', 3),
('MK304', 'Statistik dan Probabilitas', 3),
('MK305', 'Kalkulus', 3),
('MK306', 'Prinsip-Prinsip Nilai Kristiani', 2),
('MK307', 'Bahasa Inggris Pra Menengah', 2),

-- SEMESTER 4
('MK401', 'Jaringan Komputer', 3),
('MK402', 'Teknik Automata dan Kompilasi', 3),
('MK403', 'Pemrograman Berorientasi Objek', 3),
('MK404', 'Sistem Cerdas', 3),
('MK405', 'Sistem Manajemen Basisdata', 3),
('MK406', 'Orang Muda dan Dunia', 2),
('MK407', 'Bahasa Inggris Pra Menengah II', 2),

-- SEMESTER 5
('MK501', 'Jaringan Komputer II', 3),
('MK502', 'Interaksi Manusia dan Komputer', 3),
('MK503', 'Konsep Sistem Operasi', 3),
('MK504', 'Pengembangan Web Front-End', 3),
('MK505', 'Analisis dan Perancangan Sistem', 3),
('MK506', 'Metodologi Penelitian', 2),
('MK507', 'Kehidupan Keluarga', 2),

-- SEMESTER 6
('MK601', 'Pemrograman Visual', 3),
('MK602', 'Pemrograman Sistem', 3),
('MK603', 'Grafika Komputer', 3),
('MK604', 'Pengantar Pengembangan Game', 3),
('MK605', 'Kecerdasan Buatan', 3),
('MK606', 'Penambangan dan Pergudangan Data', 3),
('MK607', 'Pengembangan Web Back-End', 3),
('MK608', 'Etika Komputer', 2),
('MK609', 'Rekayasa Perangkat Lunak', 3),
('MK610', 'Pengembangan Perangkat Bergerak', 3),
('MK611', 'Skripsi I', 2),
('MK612', 'Kehidupan di Akhir Zaman', 2),

-- SEMESTER 7
('MK701', 'Robotika', 3),
('MK702', 'Kewirausahaan (Project Capstone)', 3),
('MK703', 'Skripsi II', 4),

-- SEMESTER 8 (Core Courses & Tracks)
('MK801', 'Prinsip-Prinsip Desain Kreatif', 3),
('MK802', 'Keamanan Komputer', 3),
('MK803', 'Desain untuk Visualisasi dan Komputer', 3),
('MK804', 'Pengantar Animasi', 3),
('MK805', 'Internet of Things', 3),
('MK806', 'Industrial Experience in IT', 4),
('MK807', 'Pemrosesan Bahasa Alami', 3),
('MK808', 'Pembelajaran Mesin', 3),
('MK809', 'Rekayasa DevOps', 3),
('MK810', 'Manajemen Proyek', 3);

-- ==========================================
-- 3. DATA REKOMENDASI KURIKULUM (Mapping ke Semester)
-- ==========================================
INSERT INTO kurikulum_rekomendasi (program_studi, angkatan, semester_wajib, kode_mk) VALUES
('Informatika', 2023, 1, 'MK101'), ('Informatika', 2023, 1, 'MK102'), ('Informatika', 2023, 1, 'MK103'),
('Informatika', 2023, 1, 'MK104'), ('Informatika', 2023, 1, 'MK105'), ('Informatika', 2023, 1, 'MK106'),
('Informatika', 2023, 2, 'MK201'), ('Informatika', 2023, 2, 'MK202'), ('Informatika', 2023, 2, 'MK203'),
('Informatika', 2023, 2, 'MK204'), ('Informatika', 2023, 2, 'MK205'), ('Informatika', 2023, 2, 'MK206'),
('Informatika', 2023, 2, 'MK207'), ('Informatika', 2023, 2, 'MK208'),
('Informatika', 2023, 3, 'MK301'), ('Informatika', 2023, 3, 'MK302'), ('Informatika', 2023, 3, 'MK303'),
('Informatika', 2023, 3, 'MK304'), ('Informatika', 2023, 3, 'MK305'), ('Informatika', 2023, 3, 'MK306'),
('Informatika', 2023, 3, 'MK307');

-- ==========================================
-- 4. DATA JADWAL KELAS (NAMA DOSEN DIMASUKKAN KE SINI)
-- ==========================================
INSERT INTO jadwal_kelas (kode_mk, kelas, nama_dosen, hari, jam_mulai, jam_selesai, tempat) VALUES
-- Dosen Semester 1
('MK103', 'A', 'Sahulata, Reynoldus', 'Senin', '08:00:00', '10:00:00', 'Gedung A'),
('MK104', 'A', 'Silinaung, Max', 'Senin', '10:00:00', '12:00:00', 'Gedung A'),
('MK105', 'A', 'Watopa, James', 'Selasa', '08:00:00', '10:00:00', 'Gedung B'),
('MK106', 'A', 'Lumoindong, Boy', 'Selasa', '10:00:00', '12:00:00', 'Gedung B'),

-- Dosen Semester 2
('MK201', 'A', 'Laoh, Lidya', 'Rabu', '08:00:00', '10:00:00', 'Lab Komputer 1'),
('MK203', 'A', 'Najoan, Regi', 'Rabu', '10:00:00', '12:00:00', 'Lab Komputer 2'),
('MK204', 'A', 'Edson Yahuda Putra', 'Kamis', '08:00:00', '10:00:00', 'Gedung C'),

-- Dosen Semester 3
('MK301', 'A', 'Waworundeng, Jacquline', 'Senin', '13:00:00', '15:00:00', 'Gedung C'),
('MK302', 'A', 'Najoan, Regi', 'Selasa', '13:00:00', '15:00:00', 'Lab Komputer 3'),
('MK303', 'A', 'Mandias, Green', 'Rabu', '13:00:00', '15:00:00', 'Lab Basis Data'),
('MK305', 'A', 'Sahulata, Reynoldus', 'Kamis', '13:00:00', '15:00:00', 'Gedung A'),
('MK306', 'A', 'Woy, Lauda', 'Jumat', '08:00:00', '10:00:00', 'Gedung B'),

-- Dosen Semester 4 (Sesuai Gambar KRS)
('MK401', 'D', 'Mokodaser, Wilsen', 'Senin, Rabu', '14:40:00', '16:00:00', 'GK1 - 405'),
('MK402', 'A', 'Sahulata, Reynoldus', 'Senin, Rabu', '16:10:00', '17:30:00', 'GK1 - 303'),
('MK404', 'B', 'Sandag, Green A.', 'Selasa, Kamis', '14:40:00', '16:00:00', 'GK1 - 503 Software Lab'),
('MK405', 'B', 'Mandias, Green Ferry', 'Selasa, Kamis', '07:10:00', '08:30:00', 'GK1 - 503 Software Lab'),
('MK406', 'B', 'Lumingkewas, Edwin Melky', 'Selasa', '16:10:00', '18:00:00', 'GA - English Lab 1'),

-- Dosen Semester 5
('MK502', 'B', 'Tangka, George', 'Selasa, Kamis', '13:10:00', '14:30:00', 'GK1 - 303'),

-- Dosen Semester 6
('MK607', 'B', 'Djimesha, Enrico', 'Senin, Rabu', '08:40:00', '10:00:00', 'GA - English Lab 3');