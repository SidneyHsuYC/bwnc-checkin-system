-- Rename leader_id to student_id for clarity (leader is also a student)
ALTER TABLE classes RENAME COLUMN leader_id TO student_id;

-- Rename the index
DROP INDEX IF EXISTS idx_classes_leader_id;
CREATE INDEX IF NOT EXISTS idx_classes_student_id ON classes(student_id);
