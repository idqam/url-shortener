import { AnonResult } from "../components/AnonResult";
import { ShortenerForm } from "../components/ShortenerForm";

export const Home = () => {
  return (
    <div className="shortener-container">
      <ShortenerForm />
      <AnonResult />
    </div>
  );
};
