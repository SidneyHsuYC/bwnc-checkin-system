// Class form handling
const API_BASE_URL = 'http://localhost:8090/api';

let leaderSearchTimeout = null;

document.addEventListener('DOMContentLoaded', () => {
  setupLeaderAutocomplete();
});

// Setup leader autocomplete
function setupLeaderAutocomplete() {
  const leaderInput = document.getElementById('leader_name');
  const resultsDiv = document.getElementById('leader_autocomplete_results');

  leaderInput.addEventListener('input', (e) => {
    clearTimeout(leaderSearchTimeout);
    const query = e.target.value.trim();
    
    if (query.length < 2) {
      resultsDiv.classList.add('hidden');
      return;
    }

    leaderSearchTimeout = setTimeout(() => {
      searchStudents(query);
    }, 300);
  });

  leaderInput.addEventListener('blur', () => {
    setTimeout(() => {
      resultsDiv.classList.add('hidden');
    }, 200);
  });
}

// Search students for leader
async function searchStudents(query) {
  try {
    const response = await fetch(`${API_BASE_URL}/students/search?q=${encodeURIComponent(query)}`);
    if (response.ok) {
      const students = await response.json();
      displayLeaderResults(students);
    }
  } catch (error) {
    console.error('Error searching students:', error);
  }
}

// Display leader autocomplete results
function displayLeaderResults(students) {
  const resultsDiv = document.getElementById('leader_autocomplete_results');
  resultsDiv.innerHTML = '';

  if (students.length === 0) {
    const noResults = document.createElement('div');
    noResults.className = 'autocomplete-item';
    noResults.textContent = 'No students found';
    noResults.style.color = '#999';
    resultsDiv.appendChild(noResults);
    resultsDiv.classList.remove('hidden');
    return;
  }

  students.forEach(student => {
    const item = document.createElement('div');
    item.className = 'autocomplete-item';
    item.textContent = `${student.first_name} ${student.last_name} (${student.email})`;
    item.addEventListener('click', () => selectLeader(student));
    resultsDiv.appendChild(item);
  });

  resultsDiv.classList.remove('hidden');
}

// Select leader from autocomplete
function selectLeader(student) {
  document.getElementById('leader_name').value = `${student.first_name} ${student.last_name}`;
  document.getElementById('student_id').value = student.id;
  document.getElementById('leader_name').setAttribute('readonly', 'readonly');
  document.getElementById('leader_autocomplete_results').classList.add('hidden');
}

// Clear leader selection when input changes
function setupLeaderClearOnEdit() {
  const leaderInput = document.getElementById('leader_name');
  leaderInput.addEventListener('input', (e) => {
    // If user modifies the text, clear the hidden student_id
    const studentId = document.getElementById('student_id').value;
    if (studentId && e.target.value !== leaderInput.getAttribute('data-selected-name')) {
      document.getElementById('student_id').value = '';
      leaderInput.removeAttribute('readonly');
    }
  });
  
  leaderInput.addEventListener('focus', () => {
    leaderInput.removeAttribute('readonly');
  });
}

document.addEventListener('DOMContentLoaded', () => {
  setupLeaderAutocomplete();
  setupLeaderClearOnEdit();
});

document.getElementById('classForm').addEventListener('submit', async (e) => {
  e.preventDefault();

  const form = e.target;

  // Get form data and format date properly
  const startDate = form.start_date.value; // YYYY-MM-DD from date input
  const data = {
    class_name: form.class_name.value.trim(),
    start_date: startDate + 'T00:00:00Z', // Convert to RFC3339 format
    day_of_week: form.day_of_week.value,
    start_time: form.start_time.value,
    end_time: form.end_time.value
  };

  // Add optional leader - but only if properly selected from dropdown
  const studentId = form.student_id.value;
  const leaderName = form.leader_name.value.trim();
  
  if (leaderName && !studentId) {
    showMessage('Please select a leader from the dropdown list', 'error');
    return;
  }
  
  if (studentId) {
    data.student_id = parseInt(studentId);
  }

  // Validation
  if (!data.class_name || !startDate || !data.day_of_week || !data.start_time || !data.end_time) {
    showMessage('Please fill in all required fields', 'error');
    return;
  }

  // Validate time range
  if (data.start_time >= data.end_time) {
    showMessage('End time must be after start time', 'error');
    return;
  }

  try {
    const response = await fetch(`${API_BASE_URL}/classes`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    });

    if (response.ok) {
      const result = await response.json();
      showMessage(`Class created successfully! ID: ${result.id}`, 'success');
      form.reset();
      document.getElementById('student_id').value = '';
    } else {
      const error = await response.json();
      if (response.status === 409) {
        showMessage('A class with this name already exists', 'error');
      } else {
        showMessage(error.error || 'Failed to create class', 'error');
      }
    }
  } catch (error) {
    console.error('Error:', error);
    showMessage('Network error. Please try again.', 'error');
  }
});

function showMessage(text, type) {
  const messageDiv = document.getElementById('message');
  messageDiv.textContent = text;
  messageDiv.className = `message message-${type}`;
  messageDiv.classList.remove('hidden');

  // Auto-hide success messages after 5 seconds
  if (type === 'success') {
    setTimeout(() => {
      messageDiv.classList.add('hidden');
    }, 5000);
  }
}
