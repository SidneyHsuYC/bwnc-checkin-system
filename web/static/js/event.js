// Event form handling
const API_BASE_URL = 'http://localhost:8090/api';

document.getElementById('eventForm').addEventListener('submit', async (e) => {
  e.preventDefault();

  const form = e.target;
  const messageDiv = document.getElementById('message');

  // Get form data
  const data = {
    event_name: form.event_name.value.trim(),
    event_time: form.event_time.value,
    event_type: form.event_type.value.trim()
  };

  // Client-side validation
  if (!data.event_name || !data.event_time || !data.event_type) {
    showMessage('All fields are required', 'error');
    return;
  }

  // Convert datetime-local to ISO format
  const eventDate = new Date(data.event_time);
  data.event_time = eventDate.toISOString();

  try {
    const response = await fetch(`${API_BASE_URL}/events`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    });

    if (response.ok) {
      const result = await response.json();
      showMessage(`Event created successfully! ID: ${result.id}`, 'success');
      form.reset();
    } else {
      const error = await response.json();
      showMessage(error.error || 'Failed to create event', 'error');
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
