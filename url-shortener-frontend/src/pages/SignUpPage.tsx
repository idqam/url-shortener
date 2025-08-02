import { SignupForm } from "../components/SignUpForm";
import "./signup-page.css";
export const SignUpPage = () => {
  return (
    <div className="signup-page-container">
      <div className="signup-page">
        <h1>Create Your Account</h1>
        <SignupForm />
      </div>
    </div>
  );
};
