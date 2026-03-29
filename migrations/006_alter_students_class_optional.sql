-- Make class_info optional and add class_id foreign key
ALTER TABLE students 
    ALTER COLUMN class_info DROP NOT NULL;

-- Add class_id column
ALTER TABLE students 
    ADD COLUMN class_id INTEGER REFERENCES classes(id) ON DELETE SET NULL;

-- Create index for class_id
CREATE INDEX IF NOT EXISTS idx_students_class_id ON students(class_id);
