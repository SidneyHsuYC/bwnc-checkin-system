// Student form handling
const API_BASE_URL = 'http://localhost:8090/api';

const EMAIL_SUFFIXES = [
  '@gmail.com',
  '@yahoo.com',
  '@outlook.com',
  '@hotmail.com',
  '@icloud.com',
  '@aol.com'
];

let classSearchTimeout = null;

document.addEventListener('DOMContentLoaded', () => {
  setupEmailAutocomplete();
  setupClassAutocomplete();
  setupClassClearOnEdit();
});

// Setup email autocomplete
function setupEmailAutocomplete() {
  const emailInput = document.getElementById('email');
  const suggestionsDiv = document.getElementById('email_suggestions');

  emailInput.addEventListener('input', (e) => {
    const value = e.target.value;
    
    // Show suggestions if there's text before @ or just @
    const atIndex = value.indexOf('@');
    if (atIndex === -1) {
      // No @ yet, show all suggestions with the prefix
      if (value) {
        showEmailSuggestions(value);
      } else {
        suggestionsDiv.classList.add('hidden');
      }
    } else {
      // Has @, show filtered suggestions
      const prefix = value.substring(0, atIndex);
      const suffix = value.substring(atIndex);
      if (prefix) {
        showEmailSuggestionsWithSuffix(prefix, suffix);
      } else {
        suggestionsDiv.classList.add('hidden');
      }
    }
  });

  emailInput.addEventListener('blur', () => {
    // Delay hiding to allow click on suggestion
    setTimeout(() => {
      suggestionsDiv.classList.add('hidden');
    }, 200);
  });
  
  emailInput.addEventListener('focus', (e) => {
    // Show suggestions on focus if there's a value
    const value = e.target.value;
    if (value && !value.includes('@')) {
      showEmailSuggestions(value);
    }
  });
}

// Show email suggestions with suffix filtering
function showEmailSuggestionsWithSuffix(prefix, currentSuffix) {
  const suggestionsDiv = document.getElementById('email_suggestions');
  suggestionsDiv.innerHTML = '';

  const filteredSuffixes = EMAIL_SUFFIXES.filter(suffix => 
    suffix.toLowerCase().startsWith(currentSuffix.toLowerCase())
  );

  if (filteredSuffixes.length === 0) {
    suggestionsDiv.classList.add('hidden');
    return;
  }

  filteredSuffixes.forEach(suffix => {
    const suggestion = document.createElement('div');
    suggestion.className = 'email-suggestion-item';
    suggestion.textContent = prefix + suffix;
    suggestion.addEventListener('click', () => {
      document.getElementById('email').value = prefix + suffix;
      suggestionsDiv.classList.add('hidden');
    });
    suggestionsDiv.appendChild(suggestion);
  });

  suggestionsDiv.classList.remove('hidden');
}

// Show email suggestions (all suffixes)
function showEmailSuggestions(prefix) {
  const suggestionsDiv = document.getElementById('email_suggestions');
  suggestionsDiv.innerHTML = '';

  EMAIL_SUFFIXES.forEach(suffix => {
    const suggestion = document.createElement('div');
    suggestion.className = 'email-suggestion-item';
    suggestion.textContent = prefix + suffix;
    suggestion.addEventListener('click', () => {
      document.getElementById('email').value = prefix + suffix;
      suggestionsDiv.classList.add('hidden');
    });
    suggestionsDiv.appendChild(suggestion);
  });

  suggestionsDiv.classList.remove('hidden');
}

// Setup class autocomplete
function setupClassAutocomplete() {
  const classInput = document.getElementById('class_info');
  const resultsDiv = document.getElementById('class_autocomplete_results');

  classInput.addEventListener('input', (e) => {
    clearTimeout(classSearchTimeout);
    const query = e.target.value.trim();
    
    if (query.length < 2) {
      resultsDiv.classList.add('hidden');
      return;
    }

    classSearchTimeout = setTimeout(() => {
      searchClasses(query);
    }, 300);
  });

  classInput.addEventListener('blur', () => {
    setTimeout(() => {
      resultsDiv.classList.add('hidden');
    }, 200);
  });
}

// Search classes
async function searchClasses(query) {
  try {
    const response = await fetch(`${API_BASE_URL}/classes/search?q=${encodeURIComponent(query)}`);
    if (response.ok) {
      const classes = await response.json();
      displayClassResults(classes);
    }
  } catch (error) {
    console.error('Error searching classes:', error);
  }
}

// Display class autocomplete results
function displayClassResults(classes) {
  const resultsDiv = document.getElementById('class_autocomplete_results');
  resultsDiv.innerHTML = '';

  if (classes.length === 0) {
    const noResults = document.createElement('div');
    noResults.className = 'autocomplete-item';
    noResults.textContent = 'No classes found';
    noResults.style.color = '#999';
    resultsDiv.appendChild(noResults);
    resultsDiv.classList.remove('hidden');
    return;
  }

  classes.forEach(classItem => {
    const item = document.createElement('div');
    item.className = 'autocomplete-item';
    const startDate = new Date(classItem.start_date).toLocaleDateString();
    const leaderName = classItem.leader ? ` - Leader: ${classItem.leader.first_name} ${classItem.leader.last_name}` : '';
    item.textContent = `${classItem.class_name} - ${classItem.day_of_week} ${classItem.start_time}-${classItem.end_time}${leaderName}`;
    item.addEventListener('click', () => selectClass(classItem));
    resultsDiv.appendChild(item);
  });

  resultsDiv.classList.remove('hidden');
}

// Select class from autocomplete
function selectClass(classItem) {
  document.getElementById('class_info').value = classItem.class_name;
  document.getElementById('class_id').value = classItem.id;
  document.getElementById('class_info').setAttribute('readonly', 'readonly');
  document.getElementById('class_autocomplete_results').classList.add('hidden');
}

// Clear class selection when input changes
function setupClassClearOnEdit() {
  const classInput = document.getElementById('class_info');
  classInput.addEventListener('input', (e) => {
    // If user modifies the text, clear the hidden class_id
    const classId = document.getElementById('class_id').value;
    if (classId && e.target.value !== classInput.getAttribute('data-selected-name')) {
      document.getElementById('class_id').value = '';
      classInput.removeAttribute('readonly');
    }
  });
  
  classInput.addEventListener('focus', () => {
    classInput.removeAttribute('readonly');
  });
}

document.getElementById('studentForm').addEventListener('submit', async (e) => {
  e.preventDefault();

  const form = e.target;
  const messageDiv = document.getElementById('message');

  // Get form data
  const data = {
    first_name: form.first_name.value.trim(),
    last_name: form.last_name.value.trim(),
    email: form.email.value.trim()
  };

  // Add optional fields - but validate class selection
  const classInfo = form.class_info.value.trim();
  const classId = form.class_id.value;
  
  // Validate: if there's text in class field, must have selected from dropdown
  if (classInfo && !classId) {
    showMessage('Please select a class from the dropdown list', 'error');
    return;
  }
  
  if (classInfo) {
    data.class_info = classInfo;
  }
  
  if (classId) {
    data.class_id = parseInt(classId);
  }

  // Client-side validation
  if (!data.first_name || !data.last_name || !data.email) {
    showMessage('First name, last name, and email are required', 'error');
    return;
  }

  // Email format validation
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!emailRegex.test(data.email)) {
    showMessage('Please enter a valid email address', 'error');
    return;
  }

  try {
    const response = await fetch(`${API_BASE_URL}/students`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    });

    if (response.ok) {
      const result = await response.json();
      showMessage(`Student account created successfully! ID: ${result.id}`, 'success');
      form.reset();
      document.getElementById('class_id').value = '';
    } else {
      const error = await response.json();
      if (response.status === 409) {
        showMessage('This email is already registered', 'error');
      } else {
        showMessage(error.error || 'Failed to create student account', 'error');
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
