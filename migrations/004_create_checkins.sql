-- Create check_ins table
CREATE TABLE IF NOT EXISTS check_ins (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL REFERENCES students(id) ON DELETE RESTRICT,
    event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE RESTRICT,
    checked_in_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_checkins_student ON check_ins(student_id);
CREATE INDEX IF NOT EXISTS idx_checkins_event ON check_ins(event_id);
CREATE INDEX IF NOT EXISTS idx_checkins_time ON check_ins(checked_in_at);

-- Create unique constraint to prevent duplicate check-ins for the same student and event
CREATE UNIQUE INDEX IF NOT EXISTS idx_checkins_student_event ON check_ins(student_id, event_id);
