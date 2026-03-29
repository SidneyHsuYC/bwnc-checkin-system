-- Add class_name column to classes table (nullable first)
ALTER TABLE classes ADD COLUMN IF NOT EXISTS class_name VARCHAR(255);

-- Update existing rows with unique placeholder names based on their ID
UPDATE classes SET class_name = 'Class-' || id WHERE class_name IS NULL;

-- Now make it NOT NULL
ALTER TABLE classes ALTER COLUMN class_name SET NOT NULL;

-- Add unique constraint on class_name
ALTER TABLE classes ADD CONSTRAINT unique_class_name UNIQUE (class_name);

-- Create index on class_name for faster searches
CREATE INDEX IF NOT EXISTS idx_classes_class_name ON classes(class_name);
