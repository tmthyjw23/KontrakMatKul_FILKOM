-- =============================================================
-- Migration: Initial Schema for Sistem Kontrak Mata Kuliah
-- Engine: MySQL 8.0+
-- =============================================================

-- Students table
CREATE TABLE IF NOT EXISTS students (
    nim           VARCHAR(20)  NOT NULL,
    name          VARCHAR(100) NOT NULL,
    faculty       VARCHAR(100) NOT NULL,
    study_program VARCHAR(100) NOT NULL,
    cohort_year   INT          NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(20)  NOT NULL DEFAULT 'Student',
    PRIMARY KEY (nim)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Courses (Mata Kuliah) table
CREATE TABLE IF NOT EXISTS courses (
    code          VARCHAR(20)  NOT NULL,
    name          VARCHAR(150) NOT NULL,
    class         VARCHAR(10)  NOT NULL,
    lecturer_name VARCHAR(100) NOT NULL,
    credits       INT          NOT NULL CHECK (credits > 0),
    cohort_target INT          NOT NULL,
    PRIMARY KEY (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Schedules table
CREATE TABLE IF NOT EXISTS schedules (
    id               INT          NOT NULL AUTO_INCREMENT,
    course_code      VARCHAR(20)  NOT NULL,
    day_of_week      VARCHAR(20)  NOT NULL,
    start_time       VARCHAR(10)  NOT NULL, -- e.g. "08:00"
    end_time         VARCHAR(10)  NOT NULL, -- e.g. "10:00"
    room             VARCHAR(50)  NOT NULL,
    additional_notes TEXT,
    PRIMARY KEY (id),
    CONSTRAINT fk_schedules_course FOREIGN KEY (course_code)
        REFERENCES courses(code) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Registrations table
CREATE TABLE IF NOT EXISTS registrations (
    id          INT          NOT NULL AUTO_INCREMENT,
    student_nim VARCHAR(20)  NOT NULL,
    course_code VARCHAR(20)  NOT NULL,
    status      VARCHAR(20)  NOT NULL DEFAULT 'registered',
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT unique_active_registration UNIQUE (student_nim, course_code),
    CONSTRAINT fk_reg_student FOREIGN KEY (student_nim)
        REFERENCES students(nim) ON DELETE CASCADE,
    CONSTRAINT fk_reg_course FOREIGN KEY (course_code)
        REFERENCES courses(code) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
