-- Buat database jika belum ada
CREATE DATABASE IF NOT EXISTS db_curriculum;

-- Gunakan database tersebut
USE db_curriculum;

-- Hapus tabel jika sudah ada agar tidak error saat di-import ulang
DROP TABLE IF EXISTS kontrak_krs;
DROP TABLE IF EXISTS jadwal_kelas;
DROP TABLE IF EXISTS kurikulum_rekomendasi;
DROP TABLE IF EXISTS mata_kuliah;
DROP TABLE IF EXISTS mahasiswa;

-- =========================================
-- TABEL DATA MASTER
-- =========================================

-- Tabel Mahasiswa
CREATE TABLE mahasiswa (
    nim VARCHAR(20) PRIMARY KEY,
    nama VARCHAR(100) NOT NULL,
    fakultas VARCHAR(100) NOT NULL,
    program_studi VARCHAR(100) NOT NULL,
    angkatan INT NOT NULL
);

-- Tabel Mata Kuliah
CREATE TABLE mata_kuliah (
    kode_mk VARCHAR(20) PRIMARY KEY,
    nama_mk VARCHAR(150) NOT NULL,
    sks INT NOT NULL
);

-- Tabel Kurikulum Rekomendasi
CREATE TABLE kurikulum_rekomendasi (
    id INT AUTO_INCREMENT PRIMARY KEY,
    program_studi VARCHAR(100) NOT NULL,
    angkatan INT NOT NULL,
    semester_wajib INT NOT NULL,
    kode_mk VARCHAR(20) NOT NULL,
    FOREIGN KEY (kode_mk) REFERENCES mata_kuliah(kode_mk)
);

-- Tabel Jadwal Kelas
CREATE TABLE jadwal_kelas (
    id_jadwal INT AUTO_INCREMENT PRIMARY KEY,
    kode_mk VARCHAR(20) NOT NULL,
    kelas VARCHAR(10) NOT NULL,
    nama_dosen VARCHAR(100) NOT NULL,
    hari INT NOT NULL, -- 1=Senin, 2=Selasa, dst
    jam_mulai TIME NOT NULL, 
    jam_selesai TIME NOT NULL, 
    tempat VARCHAR(100),
    keterangan VARCHAR(255),
    FOREIGN KEY (kode_mk) REFERENCES mata_kuliah(kode_mk)
);

-- =========================================
-- TABEL DATA TRANSAKSI
-- =========================================

-- Tabel Kontrak KRS
CREATE TABLE kontrak_krs (
    id_kontrak INT AUTO_INCREMENT PRIMARY KEY,
    nim VARCHAR(20) NOT NULL,
    id_jadwal INT NOT NULL,
    semester_aktif VARCHAR(50) NOT NULL,
    FOREIGN KEY (nim) REFERENCES mahasiswa(nim),
    FOREIGN KEY (id_jadwal) REFERENCES jadwal_kelas(id_jadwal)
);