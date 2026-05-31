-- seed.sql — example data for local development and demos.
--
-- Safe to run more than once: every statement is idempotent (ON CONFLICT /
-- WHERE NOT EXISTS), so re-running adds nothing and changes nothing. It does
-- NOT delete or reset anything. Run it with ./seed.sh (see that script and the
-- README "Load example data" section).
--
-- What you get: 5 classes, 10 students (some linked to a class, some not,
-- three classes have a student leader), 5 events (dated relative to today so
-- they always show up in the kiosk's recent-events picker), and check-ins
-- against the two past events so an event already has attendance to look at.

BEGIN;

-- 1. Classes ---------------------------------------------------------------
-- class_name is UNIQUE, so ON CONFLICT makes this re-runnable. Leaders are
-- wired up later (step 3) once the students they point at exist.
INSERT INTO classes (class_name, start_date, day_of_week, start_time, end_time)
VALUES
  ('Beginner Yoga',         CURRENT_DATE - 30, 'Mon', '18:00', '19:30'),
  ('Advanced Piano',        CURRENT_DATE - 60, 'Tue', '19:00', '21:00'),
  ('Conversational Spanish', CURRENT_DATE - 20, 'Wed', '18:30', '20:00'),
  ('Watercolor Painting',   CURRENT_DATE - 45, 'Thu', '17:00', '19:00'),
  ('Community Choir',       CURRENT_DATE - 90, 'Sat', '10:00', '12:00')
ON CONFLICT (class_name) DO NOTHING;

-- 2. Students --------------------------------------------------------------
-- email is UNIQUE. class_info is the legacy free-text label shown in the
-- check-in autocomplete; we always write a string ('' when there's no class)
-- because the search query scans it into a non-nullable column. class_id links
-- to the class created above, looked up by name so it works on a re-run too.
INSERT INTO students (first_name, last_name, email, class_info, class_id)
VALUES
  ('Alice',  'Nguyen',    'alice.nguyen@example.com',    'Beginner Yoga (Mon 18:00-19:30)',          (SELECT id FROM classes WHERE class_name = 'Beginner Yoga')),
  ('Bob',    'Martinez',  'bob.martinez@example.com',    'Beginner Yoga (Mon 18:00-19:30)',          (SELECT id FROM classes WHERE class_name = 'Beginner Yoga')),
  ('Carla',  'Singh',     'carla.singh@example.com',     'Advanced Piano (Tue 19:00-21:00)',         (SELECT id FROM classes WHERE class_name = 'Advanced Piano')),
  ('David',  'Okafor',    'david.okafor@example.com',    'Advanced Piano (Tue 19:00-21:00)',         (SELECT id FROM classes WHERE class_name = 'Advanced Piano')),
  ('Emma',   'Johansson', 'emma.johansson@example.com',  'Conversational Spanish (Wed 18:30-20:00)', (SELECT id FROM classes WHERE class_name = 'Conversational Spanish')),
  ('Farid',  'Hassan',    'farid.hassan@example.com',    'Conversational Spanish (Wed 18:30-20:00)', (SELECT id FROM classes WHERE class_name = 'Conversational Spanish')),
  ('Grace',  'Liu',       'grace.liu@example.com',       'Watercolor Painting (Thu 17:00-19:00)',    (SELECT id FROM classes WHERE class_name = 'Watercolor Painting')),
  ('Hiro',   'Tanaka',    'hiro.tanaka@example.com',     'Community Choir (Sat 10:00-12:00)',        (SELECT id FROM classes WHERE class_name = 'Community Choir')),
  ('Isabel', 'Romero',    'isabel.romero@example.com',   'Community Choir (Sat 10:00-12:00)',        (SELECT id FROM classes WHERE class_name = 'Community Choir')),
  ('Jordan', 'Smith',     'jordan.smith@example.com',    '',                                          NULL)
ON CONFLICT (email) DO NOTHING;

-- 3. Class leaders ---------------------------------------------------------
-- A leader is a student (classes.student_id). Set by email so it's idempotent.
UPDATE classes SET student_id = (SELECT id FROM students WHERE email = 'alice.nguyen@example.com')
  WHERE class_name = 'Beginner Yoga';
UPDATE classes SET student_id = (SELECT id FROM students WHERE email = 'carla.singh@example.com')
  WHERE class_name = 'Advanced Piano';
UPDATE classes SET student_id = (SELECT id FROM students WHERE email = 'hiro.tanaka@example.com')
  WHERE class_name = 'Community Choir';

-- 4. Events ----------------------------------------------------------------
-- Times are relative to now so the two past events seed some attendance and
-- the three upcoming ones are checkin-able from the kiosk. Events have no
-- unique key, so we guard on event_name to avoid duplicates on a re-run.
INSERT INTO events (event_name, event_time, event_type)
SELECT v.event_name, v.event_time, v.event_type
FROM (VALUES
  ('Spring Wellness Workshop', NOW() - INTERVAL '14 days', 'Workshop'),
  ('Open Mic Night',           NOW() - INTERVAL '5 days',  'Social'),
  ('Volunteer Orientation',    NOW() + INTERVAL '3 days',  'Orientation'),
  ('Summer Recital',           NOW() + INTERVAL '12 days', 'Performance'),
  ('Community Potluck',         NOW() + INTERVAL '25 days', 'Social')
) AS v(event_name, event_time, event_type)
WHERE NOT EXISTS (SELECT 1 FROM events e WHERE e.event_name = v.event_name);

-- 5. Check-ins -------------------------------------------------------------
-- Populate attendance for the two past events. (student_id, event_id) is
-- unique, so ON CONFLICT keeps this re-runnable. IDs are resolved by natural
-- keys (email / event_name) so we never hard-code serial values.
INSERT INTO check_ins (student_id, event_id)
SELECT s.id, e.id
FROM students s
JOIN events e ON e.event_name = 'Spring Wellness Workshop'
WHERE s.email IN (
  'alice.nguyen@example.com', 'bob.martinez@example.com', 'carla.singh@example.com',
  'emma.johansson@example.com', 'grace.liu@example.com'
)
ON CONFLICT (student_id, event_id) DO NOTHING;

INSERT INTO check_ins (student_id, event_id)
SELECT s.id, e.id
FROM students s
JOIN events e ON e.event_name = 'Open Mic Night'
WHERE s.email IN (
  'carla.singh@example.com', 'david.okafor@example.com',
  'hiro.tanaka@example.com', 'isabel.romero@example.com'
)
ON CONFLICT (student_id, event_id) DO NOTHING;

COMMIT;

-- Summary of what's now in the database (handy when running by hand).
SELECT 'classes'   AS table, COUNT(*) FROM classes
UNION ALL SELECT 'students',  COUNT(*) FROM students
UNION ALL SELECT 'events',    COUNT(*) FROM events
UNION ALL SELECT 'check_ins', COUNT(*) FROM check_ins;
