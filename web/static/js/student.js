// Student form handling
const API_BASE_URL = 'http://localhost:8090/api';

document.getElementById('studentForm').addEventListener('submit', async (e) => {
  e.preventDefault();

  const form = e.target;
  const messageDiv = document.getElementById('message');

  // Get form data
  const data = {
    first_name: form.first_name.value.trim(),
    last_name: form.last_name.value.trim(),
    class_info: form.class_info.value.trim(),
    email: form.email.value.trim()
  };

  // Client-side validation
  if (!data.first_name || !data.last_name || !data.class_info || !data.email) {
    showMessage('All fields are required', 'error');
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
