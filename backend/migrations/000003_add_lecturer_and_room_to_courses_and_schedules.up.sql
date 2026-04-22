ALTER TABLE courses
ADD COLUMN lecturer VARCHAR(120) NOT NULL DEFAULT 'Dosen belum ditentukan' AFTER quota;

ALTER TABLE schedules
ADD COLUMN room VARCHAR(50) NOT NULL DEFAULT 'TBA' AFTER end_time;
