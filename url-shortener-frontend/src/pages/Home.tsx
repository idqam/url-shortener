import { useNavigate } from "react-router-dom";
import { AnonResult } from "../components/AnonResult";
import { ShortenerForm } from "../components/ShortenerForm";
import { SignOutButton } from "../components/SignOutButton";
import { useAuthStore } from "../store/AuthStore";

export const Home = () => {
  const navigate = useNavigate();
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);

  return (
    <section className="main-home-container">
      <div className="pad-header">
        <header className="name-header">URL SHORTENER</header>
        {isAuthenticated && <SignOutButton />}
      </div>

      <div className="shortener-container">
        <ShortenerForm />
        <AnonResult />
        <div className="premium-box">
          <h2>Want More? Try Premium Features!</h2>
          <p>
            Features coming soon — custom domains, analytics, QR codes, and
            more!
          </p>
          <button onClick={() => navigate("/signup")}>Create Account</button>
          <button onClick={() => navigate("/login")}>Sign in</button>
        </div>
      </div>
    </section>
  );
};
