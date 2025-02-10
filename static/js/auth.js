import { showNotification } from "./components/notifications.js";


export function createAuthModal() {
  const authModalHTML = `
    <div class="auth-modal-overlay" id="authModalOverlay">
      <div class="auth-modal" id="authModal">
        <div class="auth-card" id="authCard">
          <div class="auth-content" id="authContent"></div>
        </div>
      </div>
    </div>
  `;

  document.body.insertAdjacentHTML("beforeend", authModalHTML);
  showLoginForm();

  document
    .querySelector('[id="authModalOverlay"]')
    .addEventListener("click", function (event) {
      if (event.target === this) {
        closeAuthModal();
      }
    });
}

export function createLoginForm() {
  const loginFormHTML = `
    <h2>Welcome Back</h2>
    <p>Please log in to continue.</p>
    <form name="auth-form">
        <div class="input-group">
            <input type="text" name="username-email" placeholder="Email or Username" required>
        </div>
        <div class="input-group">
            <input type="password" name="password" placeholder="Password" required>
        </div>
        <button type="submit" class="btn" name="login-btn">Log In</button>
        <p class="redirect">Don't have an account? <a href="#" name="showSignUpLink">Sign Up</a></p>
    </form>
  `;
  document.querySelector('[id="authContent"]').innerHTML = loginFormHTML;

  document
    .querySelector('[name="showSignUpLink"]')
    .addEventListener("click", (e) => {
      e.preventDefault();
      showSignUpForm();
    });

  auth();
}

export function createSignUpForm() {
  const signUpFormHTML = `
    <h2>Create an Account</h2>
    <p>Join us and enjoy all features.</p>
    <form name="auth-form">
        <div class="input-group">
            <input type="text" name="username" placeholder="Username" required>
        </div>
        <div class="input-group">
            <input type="email" name="email" placeholder="Email" required>
        </div>
        <div class="input-group">
            <input type="password" name="password" placeholder="Password" required>
        </div>
        <button type="submit" class="btn" name="signup-btn">Sign Up</button>
        <p class="redirect">Already have an account? <a href="#" name="showLoginLink">Log In</a></p>
    </form>
  `;
  document.querySelector('[id="authContent"]').innerHTML = signUpFormHTML;

  document
    .querySelector('[name="showLoginLink"]')
    .addEventListener("click", (e) => {
      e.preventDefault();
      showLoginForm();

    });

  auth();
}

export function showLoginForm() {
  createLoginForm();
}

export function showSignUpForm() {
  createSignUpForm();
}

export function openAuthModal() {
  if (!document.querySelector('[id="authModalOverlay"]')) {
    createAuthModal();
  }
  document.querySelector('[id="authModalOverlay"]').style.display = "flex";
}

export function closeAuthModal() {
  document.querySelector('[id="authModalOverlay"]').style.display = "none"; // Hide modal
}

export function auth() {
  document
    .querySelector('[name="auth-form"]')
    .addEventListener("submit", async function (e) {
      e.preventDefault();

      const isLogin = document.querySelector('[name="login-btn"]') !== null;
      const data = {
        password: document.querySelector('[name="password"]').value,
      };

      if (isLogin) {
        const inputValue = document
          .querySelector('[name="username-email"]')
          .value.trim();
        if (/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(inputValue)) {
          data.email = inputValue;
        } else {
          data.username = inputValue;
        }
      } else {
        data.username = document.querySelector('[name="username"]').value;
        data.email = document.querySelector('[name="email"]').value;
      }

      try {
        const response = await fetch(
          isLogin ? "/api/auth/login" : "/api/auth/register",
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(data),
            credentials: "include",
          }
        );

        if (response.ok) {
          const message = isLogin ? "Login successful" : "Registration successful";
          showNotification(message, "success");
          window.location.reload();                                                                                                                                                                         
          logout();
      }  else {
          const error = await response.json();
          showNotification(
            `${isLogin ? "Login" : "Registration"} failed: ${error.message}`,
            "error"
          );
        }
      } catch (error) {
        showNotification("An error occurred", "error");
      }
    });
}


export function logout() {
  const logoutBtn = document.getElementById("logout-btn")
  logoutBtn?.addEventListener("click", async (e)=> {
    try {
      const response = await fetch( "/api/auth/logout",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include",
        }
      );

      if (response.ok) {
        showNotification("Logout successful", "success");
        window.location.reload();
    }    
    } catch (error) {
      showNotification("An error occurred", "error");
    }
  })
}

