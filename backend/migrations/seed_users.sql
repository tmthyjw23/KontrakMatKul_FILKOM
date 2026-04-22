-- Seed Data for Users
-- 1 Admin Account
INSERT INTO users (student_number, full_name, email, password, role, max_sks) 
VALUES ('ADMIN001', 'System Administrator', 'admin@filkom.ac.id', 'admin123', 'ADMIN', 24);

-- 4 Student Accounts
INSERT INTO users (student_number, full_name, email, password, role, max_sks) VALUES 
('S1001', 'Timothy Student 1', 'student1@filkom.ac.id', 'pass123', 'STUDENT', 24),
('S1002', 'Timothy Student 2', 'student2@filkom.ac.id', 'pass123', 'STUDENT', 24),
('S1003', 'Timothy Student 3', 'student3@filkom.ac.id', 'pass123', 'STUDENT', 24),
('S1004', 'Timothy Student 4', 'student4@filkom.ac.id', 'pass123', 'STUDENT', 24);
