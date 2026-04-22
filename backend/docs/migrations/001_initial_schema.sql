-- =============================================================
-- Migration: Initial Schema for Sistem Kontrak Mata Kuliah
-- =============================================================

-- Students table
CREATE TABLE IF NOT EXISTS students (
    nim          VARCHAR(20)  PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    faculty      VARCHAR(100) NOT NULL,
    study_program VARCHAR(100) NOT NULL,
    cohort_year  INT          NOT NULL,
    role         VARCHAR(20)  NOT NULL DEFAULT 'Student'
);

-- Courses (Mata Kuliah) table
CREATE TABLE IF NOT EXISTS courses (
    code          VARCHAR(20)  PRIMARY KEY,
    name          VARCHAR(150) NOT NULL,
    class         VARCHAR(10)  NOT NULL,
    lecturer_name VARCHAR(100) NOT NULL,
    credits       INT          NOT NULL CHECK (credits > 0),
    cohort_target INT          NOT NULL
);

-- Schedules table
CREATE TABLE IF NOT EXISTS schedules (
    id               SERIAL       PRIMARY KEY,
    course_code      VARCHAR(20)  NOT NULL REFERENCES courses(code) ON DELETE CASCADE,
    day_of_week      VARCHAR(20)  NOT NULL,
    start_time       VARCHAR(10)  NOT NULL, -- e.g. "08:00"
    end_time         VARCHAR(10)  NOT NULL, -- e.g. "10:00"
    room             VARCHAR(50)  NOT NULL,
    additional_notes TEXT
);

-- Registrations table
CREATE TABLE IF NOT EXISTS registrations (
    id          SERIAL       PRIMARY KEY,
    student_nim VARCHAR(20)  NOT NULL REFERENCES students(nim) ON DELETE CASCADE,
    course_code VARCHAR(20)  NOT NULL REFERENCES courses(code) ON DELETE CASCADE,
    status      VARCHAR(20)  NOT NULL DEFAULT 'registered',
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_active_registration UNIQUE (student_nim, course_code)
);
