import { openAuthModal } from "./auth.js";
import { logout } from "./auth.js";
import { Home } from "./Home.js";
import { SessionCheck } from "./components/sessionChecker.js";
import { filterCat } from "./filter.js";

export const Session = () => {
  const authbtn = document.getElementById("login-register");
  if (authbtn) {
    authbtn.addEventListener("click", (e) => {
      openAuthModal();
    });
  }
};
Home();
Session();
logout();
SessionCheck();
filterCat();
