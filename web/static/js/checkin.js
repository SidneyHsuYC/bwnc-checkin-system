// Check-in form handling
const API_BASE_URL = 'http://localhost:8090/api';

let selectedStudentId = null;
let searchTimeout = null;
let selectedEventId = null;

// Load upcoming events on page load
document.addEventListener('DOMContentLoaded', () => {
  loadUpcomingEvents();
  setupEventHandlers();
});

// Load upcoming events
async function loadUpcomingEvents() {
  try {
    const response = await fetch(`${API_BASE_URL}/events/upcoming`);
    if (response.ok) {
      const events = await response.json();
      const eventSelect = document.getElementById('event_id');
      
      events.forEach(event => {
        const option = document.createElement('option');
        option.value = event.id;
        const eventDate = new Date(event.event_time);
        option.textContent = `${event.event_name} - ${eventDate.toLocaleString()}`;
        eventSelect.appendChild(option);
      });
    }
  } catch (error) {
    console.error('Error loading events:', error);
    showMessage('Failed to load events', 'error');
  }
}

// Setup event handlers
function setupEventHandlers() {
  const studentNameInput = document.getElementById('student_name');
  const checkinForm = document.getElementById('checkinForm');
  const continueBtn = document.getElementById('continue_checkin');
  const goHomeBtn = document.getElementById('go_home');

  // Student name search with debounce
  studentNameInput.addEventListener('input', (e) => {
    clearTimeout(searchTimeout);
    const query = e.target.value.trim();
    
    if (query.length < 2) {
      hideAutocomplete();
      return;
    }

    searchTimeout = setTimeout(() => {
      searchStudents(query);
    }, 300);
  });

  // Form submission
  checkinForm.addEventListener('submit', handleCheckinSubmit);

  // Continue to check-in button
  continueBtn.addEventListener('click', () => {
    resetFormForNextStudent();
  });

  // Go home button
  goHomeBtn.addEventListener('click', () => {
    window.location.href = '/';
  });
}

// Search students
async function searchStudents(query) {
  try {
    const response = await fetch(`${API_BASE_URL}/students/search?q=${encodeURIComponent(query)}`);
    if (response.ok) {
      const students = await response.json();
      displayAutocompleteResults(students);
    }
  } catch (error) {
    console.error('Error searching students:', error);
  }
}

// Display autocomplete results
function displayAutocompleteResults(students) {
  const resultsDiv = document.getElementById('autocomplete_results');
  resultsDiv.innerHTML = '';

  if (students.length === 0) {
    resultsDiv.classList.add('hidden');
    return;
  }

  students.forEach(student => {
    const item = document.createElement('div');
    item.className = 'autocomplete-item';
    item.textContent = `${student.first_name} ${student.last_name} (${student.class_info})`;
    item.addEventListener('click', () => selectStudent(student));
    resultsDiv.appendChild(item);
  });

  resultsDiv.classList.remove('hidden');
}

// Select student from autocomplete
function selectStudent(student) {
  selectedStudentId = student.id;
  document.getElementById('student_name').value = `${student.first_name} ${student.last_name}`;
  document.getElementById('student_id').value = student.id;
  document.getElementById('class_info').value = student.class_info;
  document.getElementById('email').value = student.email;
  hideAutocomplete();
}

// Hide autocomplete
function hideAutocomplete() {
  document.getElementById('autocomplete_results').classList.add('hidden');
}

// Handle check-in submission
async function handleCheckinSubmit(e) {
  e.preventDefault();

  const eventId = document.getElementById('event_id').value;
  const studentId = document.getElementById('student_id').value;

  if (!eventId) {
    showMessage('Please select an event', 'error');
    return;
  }

  if (!studentId) {
    showMessage('Please select a student', 'error');
    return;
  }

  const data = {
    student_id: parseInt(studentId),
    event_id: parseInt(eventId)
  };

  try {
    const response = await fetch(`${API_BASE_URL}/checkins`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    });

    if (response.ok) {
      selectedEventId = eventId;
      showSuccessState();
    } else {
      const error = await response.json();
      if (response.status === 409) {
        showMessage('Student already checked in to this event', 'error');
      } else {
        showMessage(error.error || 'Failed to check in', 'error');
      }
    }
  } catch (error) {
    console.error('Error:', error);
    showMessage('Network error. Please try again.', 'error');
  }
}

// Show success state
function showSuccessState() {
  const form = document.getElementById('checkinForm');
  const successActions = document.getElementById('success_actions');
  
  form.classList.add('hidden');
  successActions.classList.remove('hidden');
  showMessage('Check-in successful!', 'success');
}

// Reset form for next student
function resetFormForNextStudent() {
  const form = document.getElementById('checkinForm');
  const successActions = document.getElementById('success_actions');
  const messageDiv = document.getElementById('message');
  
  // Clear student fields
  document.getElementById('student_name').value = '';
  document.getElementById('student_id').value = '';
  document.getElementById('class_info').value = '';
  document.getElementById('email').value = '';
  selectedStudentId = null;
  
  // Keep event selected
  if (selectedEventId) {
    document.getElementById('event_id').value = selectedEventId;
  }
  
  // Show form, hide success actions
  form.classList.remove('hidden');
  successActions.classList.add('hidden');
  messageDiv.classList.add('hidden');
  
  // Focus on student name input
  document.getElementById('student_name').focus();
}

// Show message
function showMessage(text, type) {
  const messageDiv = document.getElementById('message');
  messageDiv.textContent = text;
  messageDiv.className = `message message-${type}`;
  messageDiv.classList.remove('hidden');
}

// Hide autocomplete when clicking outside
document.addEventListener('click', (e) => {
  if (!e.target.closest('.autocomplete-container')) {
    hideAutocomplete();
  }
});
