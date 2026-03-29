package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"
)

// Mock repository for testing
type mockStudentRepository struct {
	students map[string]*models.Student
	nextID   int
}

func newMockStudentRepository() *mockStudentRepository {
	return &mockStudentRepository{
		students: make(map[string]*models.Student),
		nextID:   1,
	}
}

func (m *mockStudentRepository) Create(ctx context.Context, student *models.Student) error {
	// Check for duplicate email
	for _, s := range m.students {
		if s.Email == student.Email {
			return &duplicateEmailError{}
		}
	}

	student.ID = m.nextID
	m.nextID++
	m.students[student.Email] = student
	return nil
}

func (m *mockStudentRepository) GetByEmail(ctx context.Context, email string) (*models.Student, error) {
	if student, ok := m.students[email]; ok {
		return student, nil
	}
	return nil, nil
}

func (m *mockStudentRepository) GetByID(ctx context.Context, id int) (*models.Student, error) {
	for _, student := range m.students {
		if student.ID == id {
			return student, nil
		}
	}
	return nil, nil
}

func (m *mockStudentRepository) Search(ctx context.Context, query string) ([]*models.StudentSearchResult, error) {
	var results []*models.StudentSearchResult
	for _, student := range m.students {
		results = append(results, &models.StudentSearchResult{
			ID:        student.ID,
			FirstName: student.FirstName,
			LastName:  student.LastName,
			FullName:  student.FirstName + " " + student.LastName,
			ClassInfo: student.ClassInfo,
			Email:     student.Email,
		})
	}
	return results, nil
}

type duplicateEmailError struct{}

func (e *duplicateEmailError) Error() string {
	return "email already exists"
}

func TestCreateStudent(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "valid student with all fields",
			requestBody: map[string]interface{}{
				"first_name": "John",
				"last_name":  "Doe",
				"class_info": "CS 101",
				"email":      "john.doe@example.com",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var student models.Student
				if err := json.NewDecoder(w.Body).Decode(&student); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if student.ID == 0 {
					t.Error("Expected student ID to be set")
				}
				if student.FirstName != "John" {
					t.Errorf("FirstName = %q, want %q", student.FirstName, "John")
				}
			},
		},
		{
			name: "valid student without class info",
			requestBody: map[string]interface{}{
				"first_name": "Jane",
				"last_name":  "Smith",
				"email":      "jane.smith@example.com",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var student models.Student
				if err := json.NewDecoder(w.Body).Decode(&student); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if student.ClassInfo != "" {
					t.Errorf("ClassInfo should be empty, got %q", student.ClassInfo)
				}
			},
		},
		{
			name: "missing first name",
			requestBody: map[string]interface{}{
				"last_name": "Doe",
				"email":     "john.doe@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				json.NewDecoder(w.Body).Decode(&response)
				if _, ok := response["error"]; !ok {
					t.Error("Expected error field in response")
				}
			},
		},
		{
			name: "invalid email format",
			requestBody: map[string]interface{}{
				"first_name": "John",
				"last_name":  "Doe",
				"email":      "invalid-email",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				json.NewDecoder(w.Body).Decode(&response)
				if _, ok := response["error"]; !ok {
					t.Error("Expected error field in response")
				}
			},
		},
		{
			name:           "invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create handler with mock repository
			handler := &StudentHandler{
				repo: newMockStudentRepository(),
			}

			// Create request
			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/api/students", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Call handler
			handler.CreateStudent(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, want %d. Body: %s", w.Code, tt.expectedStatus, w.Body.String())
			}

			// Run additional checks if provided
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestCreateStudent_DuplicateEmail(t *testing.T) {
	handler := &StudentHandler{
		repo: newMockStudentRepository(),
	}

	// Create first student
	student1 := map[string]interface{}{
		"first_name": "John",
		"last_name":  "Doe",
		"email":      "john.doe@example.com",
	}
	body1, _ := json.Marshal(student1)
	req1 := httptest.NewRequest(http.MethodPost, "/api/students", bytes.NewReader(body1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	handler.CreateStudent(w1, req1)

	if w1.Code != http.StatusCreated {
		t.Fatalf("First student creation failed with status %d", w1.Code)
	}

	// Try to create second student with same email
	student2 := map[string]interface{}{
		"first_name": "Jane",
		"last_name":  "Smith",
		"email":      "john.doe@example.com", // Same email
	}
	body2, _ := json.Marshal(student2)
	req2 := httptest.NewRequest(http.MethodPost, "/api/students", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	handler.CreateStudent(w2, req2)

	// Should return conflict
	if w2.Code != http.StatusConflict {
		t.Errorf("Expected status %d for duplicate email, got %d", http.StatusConflict, w2.Code)
	}

	var response map[string]string
	json.NewDecoder(w2.Body).Decode(&response)
	if response["error"] == "" {
		t.Error("Expected error message for duplicate email")
	}
}

func TestSearchStudents(t *testing.T) {
	mockRepo := newMockStudentRepository()

	// Add some test students
	mockRepo.Create(context.Background(), &models.Student{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	})
	mockRepo.Create(context.Background(), &models.Student{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane@example.com",
	})

	handler := &StudentHandler{repo: mockRepo}

	tests := []struct {
		name           string
		query          string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "search with query",
			query:          "john",
			expectedStatus: http.StatusOK,
			expectedCount:  2, // Mock returns all students
		},
		{
			name:           "empty query",
			query:          "",
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/students/search?q="+tt.query, nil)
			w := httptest.NewRecorder()

			handler.SearchStudents(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status = %d, want %d", w.Code, tt.expectedStatus)
			}

			var results []models.StudentSearchResult
			json.NewDecoder(w.Body).Decode(&results)

			if tt.query == "" && len(results) != 0 {
				t.Errorf("Expected empty results for empty query, got %d results", len(results))
			}
		})
	}
}
