import { AnonResult } from "../components/AnonResult";
import { ShortenerForm } from "../components/ShortenerForm";

export const Home = () => {
  return (
    <section className="main-home-container">
      <div className="pad-header">
        <header className="name-header text-red-500">URL SHORTENER</header>
      </div>
      <div className="shortener-container">
        <ShortenerForm />
        <AnonResult />
      </div>
    </section>
  );
};
