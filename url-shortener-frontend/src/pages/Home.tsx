import { useNavigate } from "react-router-dom";
import { AnonResult } from "../components/AnonResult";
import { ShortenerForm } from "../components/ShortenerForm";

export const Home = () => {
  const navigate = useNavigate();

  return (
    <section className="main-home-container">
      <div className="pad-header">
        <header className="name-header">URL SHORTENER</header>
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
          <div className="signup-log-container">
            <button onClick={() => navigate("/signup")}>Create Account</button>
            <button onClick={() => navigate("/login")}>Sign in</button>
          </div>
        </div>
      </div>
    </section>
  );
};
