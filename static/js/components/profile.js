export function createProfilePopup() {
    const popup = document.createElement('div');
    popup.classList.add('profile-popup');
    popup.style.display = 'none';
  
    const popupContent = document.createElement('div');
    popupContent.classList.add('popup-content');
  
    const userProfile = document.createElement('div');
    userProfile.classList.add('user-profile');
    userProfile.innerHTML = `
      <img src="https://www.w3schools.com/w3images/avatar2.png" alt="Profile" class="profile-image">
      <h3 id="username">Username</h3>
      <p>Member since 2024</p>
    `;
  
    const dropdownMenu = document.createElement('div');
    dropdownMenu.classList.add('dropdown-menu');
    dropdownMenu.innerHTML = `
      <a href="/profile"><i class="fas fa-user"></i> My Profile</a>
      <a href="/settings"><i class="fas fa-cog"></i> Settings</a>
      <a href="/logout"><i class="fas fa-sign-out-alt"></i> Logout</a>
    `;
  
    popupContent.appendChild(userProfile);
    popupContent.appendChild(dropdownMenu);
  
    popup.appendChild(popupContent);
  
    document.body.appendChild(popup);
  
    return popup;
  }
  
  export default createProfilePopup