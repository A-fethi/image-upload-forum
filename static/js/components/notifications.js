export function showNotification(message, type = 'error') {
    const container = document.getElementById('notification-container');
    
    container.textContent = message;
    container.classList.remove('hidden');
    container.classList.add('show');

    if (type === 'success') {
        container.style.backgroundColor = 'rgba(0, 128, 0, 0.8)';
    } else {
        container.style.backgroundColor = 'rgba(255, 0, 0, 0.8)';
    }

    setTimeout(() => {
        container.classList.remove('show');
        container.classList.add('hidden');
    }, 5000);
}
