import { showNotification } from './notifications.js';
import { logout, openAuthModal } from '../auth.js';
import { Home } from '../Home.js';

export function SessionCheck() {
    setInterval(async () => {
        const response = await CheckAuth('/api/session/');
        if (!response) {
            showNotification("Your session has expired.", "error");
            clearCookies();
            openAuthModal();
        } else {
            return true;
        }
    }, 60 * 60 * 1000);
}

async function CheckAuth(url, options = {}) {
    try {
        const response = await fetch(url, options);

        if (response.status === 401) {
            showNotification("Session expired. Please log in again.");
            logout();
            Home();
            openAuthModal();
            return null;
        }
        
        if (!response.ok) {
            showNotification("Something went wrong.");
            return null;
        }

        return await response.json();
    } catch (error) {
        console.error("Request failed:", error);
        showNotification("Something went wrong. Please try again.");
    }
}

function clearCookies() {
    document.cookie = "";
}
