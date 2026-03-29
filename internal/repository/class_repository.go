package repository

import (
	"context"
	"database/sql"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"
)

// ClassRepository defines the interface for class data operations
type ClassRepository interface {
	Create(ctx context.Context, class *models.Class) error
	GetByID(ctx context.Context, id int) (*models.ClassWithLeader, error)
	List(ctx context.Context) ([]models.ClassWithLeader, error)
	Search(ctx context.Context, query string) ([]models.ClassWithLeader, error)
	Update(ctx context.Context, class *models.Class) error
	Delete(ctx context.Context, id int) error
}

// PostgresClassRepository implements ClassRepository for PostgreSQL
type PostgresClassRepository struct {
	db *sql.DB
}

// NewPostgresClassRepository creates a new PostgreSQL class repository
func NewPostgresClassRepository(db *sql.DB) ClassRepository {
	return &PostgresClassRepository{db: db}
}

// Create inserts a new class
func (r *PostgresClassRepository) Create(ctx context.Context, class *models.Class) error {
	query := `
		INSERT INTO classes (class_name, start_date, day_of_week, start_time, end_time, student_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(
		ctx, query,
		class.ClassName, class.StartDate, class.DayOfWeek, class.StartTime, class.EndTime, class.StudentID,
	).Scan(&class.ID, &class.CreatedAt, &class.UpdatedAt)
}

// GetByID retrieves a class by ID with leader details
func (r *PostgresClassRepository) GetByID(ctx context.Context, id int) (*models.ClassWithLeader, error) {
	query := `
		SELECT 
			c.id, c.class_name, c.start_date, c.day_of_week, c.start_time, c.end_time, 
			c.student_id, c.created_at, c.updated_at,
			s.id, s.first_name, s.last_name, s.email
		FROM classes c
		LEFT JOIN students s ON c.student_id = s.id
		WHERE c.id = $1`

	var class models.ClassWithLeader
	var studentID sql.NullInt64
	var leaderStudentID sql.NullInt64
	var firstName, lastName, email sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&class.ID, &class.ClassName, &class.StartDate, &class.DayOfWeek, &class.StartTime, &class.EndTime,
		&studentID, &class.CreatedAt, &class.UpdatedAt,
		&leaderStudentID, &firstName, &lastName, &email,
	)

	if err != nil {
		return nil, err
	}

	if studentID.Valid {
		sid := int(studentID.Int64)
		class.StudentID = &sid

		if leaderStudentID.Valid {
			class.Leader = &models.Student{
				ID:        int(leaderStudentID.Int64),
				FirstName: firstName.String,
				LastName:  lastName.String,
				Email:     email.String,
			}
		}
	}

	return &class, nil
}

// List retrieves all classes with leader details
func (r *PostgresClassRepository) List(ctx context.Context) ([]models.ClassWithLeader, error) {
	query := `
		SELECT 
			c.id, c.class_name, c.start_date, c.day_of_week, c.start_time, c.end_time, 
			c.student_id, c.created_at, c.updated_at,
			s.id, s.first_name, s.last_name, s.email
		FROM classes c
		LEFT JOIN students s ON c.student_id = s.id
		ORDER BY c.class_name, c.day_of_week, c.start_time`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []models.ClassWithLeader
	for rows.Next() {
		var class models.ClassWithLeader
		var studentID sql.NullInt64
		var leaderStudentID sql.NullInt64
		var firstName, lastName, email sql.NullString

		err := rows.Scan(
			&class.ID, &class.ClassName, &class.StartDate, &class.DayOfWeek, &class.StartTime, &class.EndTime,
			&studentID, &class.CreatedAt, &class.UpdatedAt,
			&leaderStudentID, &firstName, &lastName, &email,
		)
		if err != nil {
			return nil, err
		}

		if studentID.Valid {
			sid := int(studentID.Int64)
			class.StudentID = &sid

			if leaderStudentID.Valid {
				class.Leader = &models.Student{
					ID:        int(leaderStudentID.Int64),
					FirstName: firstName.String,
					LastName:  lastName.String,
					Email:     email.String,
				}
			}
		}

		classes = append(classes, class)
	}

	return classes, rows.Err()
}

// Search searches for classes by class name, day of week or time
func (r *PostgresClassRepository) Search(ctx context.Context, query string) ([]models.ClassWithLeader, error) {
	sqlQuery := `
		SELECT 
			c.id, c.class_name, c.start_date, c.day_of_week, c.start_time, c.end_time, 
			c.student_id, c.created_at, c.updated_at,
			s.id, s.first_name, s.last_name, s.email
		FROM classes c
		LEFT JOIN students s ON c.student_id = s.id
		WHERE c.class_name ILIKE $1 OR c.day_of_week ILIKE $1 OR c.start_time::text ILIKE $1
		ORDER BY c.class_name, c.day_of_week, c.start_time`

	rows, err := r.db.QueryContext(ctx, sqlQuery, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []models.ClassWithLeader
	for rows.Next() {
		var class models.ClassWithLeader
		var studentID sql.NullInt64
		var leaderStudentID sql.NullInt64
		var firstName, lastName, email sql.NullString

		err := rows.Scan(
			&class.ID, &class.ClassName, &class.StartDate, &class.DayOfWeek, &class.StartTime, &class.EndTime,
			&studentID, &class.CreatedAt, &class.UpdatedAt,
			&leaderStudentID, &firstName, &lastName, &email,
		)
		if err != nil {
			return nil, err
		}

		if studentID.Valid {
			sid := int(studentID.Int64)
			class.StudentID = &sid

			if leaderStudentID.Valid {
				class.Leader = &models.Student{
					ID:        int(leaderStudentID.Int64),
					FirstName: firstName.String,
					LastName:  lastName.String,
					Email:     email.String,
				}
			}
		}

		classes = append(classes, class)
	}

	return classes, rows.Err()
}

// Update updates a class
func (r *PostgresClassRepository) Update(ctx context.Context, class *models.Class) error {
	query := `
		UPDATE classes 
		SET class_name = $1, start_date = $2, day_of_week = $3, start_time = $4, end_time = $5, student_id = $6
		WHERE id = $7
		RETURNING updated_at`

	return r.db.QueryRowContext(
		ctx, query,
		class.ClassName, class.StartDate, class.DayOfWeek, class.StartTime, class.EndTime, class.StudentID, class.ID,
	).Scan(&class.UpdatedAt)
}

// Delete deletes a class
func (r *PostgresClassRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM classes WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
