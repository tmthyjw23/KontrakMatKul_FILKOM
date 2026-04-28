-- =============================================================
-- Migration 002: Add Lecturers, Prerequisites & Curriculums
-- Engine: MySQL 8.0+
-- Run AFTER: 001_initial_schema.sql
-- =============================================================

-- 1. LECTURERS TABLE
--    Central registry of every faculty member (dosen FILKOM).
CREATE TABLE IF NOT EXISTS lecturers (
    id   INT          NOT NULL AUTO_INCREMENT,
    name VARCHAR(150) NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. CURRICULUMS TABLE
--    Describes which study-program curriculum a course belongs to,
--    which semester it is offered in, and what its academic type is.
CREATE TABLE IF NOT EXISTS curriculums (
    id              INT          NOT NULL AUTO_INCREMENT,
    course_code     VARCHAR(20)  NOT NULL,
    study_program   VARCHAR(100) NOT NULL,  -- e.g. 'Informatika', 'Sistem Informasi', 'Teknologi Informasi'
    semester        INT          NOT NULL,  -- 1-9 (0 = pre-requisite / elective pool)
    course_type     VARCHAR(30)  NOT NULL,  -- 'Pre-requisite', 'General', 'Basic', 'Major', 'Elective'
    PRIMARY KEY (id),
    CONSTRAINT fk_curriculum_course FOREIGN KEY (course_code)
        REFERENCES courses(code) ON DELETE CASCADE,
    CONSTRAINT uq_curriculum UNIQUE (course_code, study_program)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. COURSE PREREQUISITES TABLE
--    Many-to-many: a course can have multiple prerequisites,
--    and a single course can be a prerequisite for many others.
CREATE TABLE IF NOT EXISTS course_prerequisites (
    id                  INT         NOT NULL AUTO_INCREMENT,
    course_code         VARCHAR(20) NOT NULL,
    prerequisite_code   VARCHAR(20) NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_prereq_course  FOREIGN KEY (course_code)
        REFERENCES courses(code) ON DELETE CASCADE,
    CONSTRAINT fk_prereq_prereq  FOREIGN KEY (prerequisite_code)
        REFERENCES courses(code) ON DELETE CASCADE,
    CONSTRAINT uq_prereq UNIQUE (course_code, prerequisite_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 4. ALTER COURSES TABLE
--    Add missing columns (class, lecturer_name, cohort_target) to match seed data
ALTER TABLE courses ADD COLUMN IF NOT EXISTS class VARCHAR(10) NOT NULL DEFAULT '-';
ALTER TABLE courses ADD COLUMN IF NOT EXISTS lecturer_name VARCHAR(100) NOT NULL DEFAULT '-';
ALTER TABLE courses ADD COLUMN IF NOT EXISTS cohort_target INT NOT NULL DEFAULT 0;

--    Add lecturer_id column (nullable FK → lecturers).
ALTER TABLE courses ADD COLUMN IF NOT EXISTS lecturer_id INT NULL;
ALTER TABLE courses ADD CONSTRAINT fk_courses_lecturer FOREIGN KEY (lecturer_id) REFERENCES lecturers(id) ON DELETE SET NULL;

-- NOTE: The old `lecturer_name` column is kept for backward compatibility.
--       The seeder will populate `lecturer_id` going forward.
--       Once fully migrated you may DROP COLUMN lecturer_name.
