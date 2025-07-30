import { AnonResult } from "../components/AnonResult";
import { ShortenerForm } from "../components/ShortenerForm";

export const Home = () => {
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
          <button>Create Account TBA</button>
        </div>
      </div>
    </section>
  );
};
