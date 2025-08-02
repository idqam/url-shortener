/* eslint-disable @typescript-eslint/no-explicit-any */
import { signupAndRegister } from "../api/signup";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

export const SignupForm = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const navigate = useNavigate();

  const isValidEmail = (email: string) =>
    /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);

  const isValidPassword = (password: string) => password.length >= 6;

  const handleSignup = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg(null);

    if (!isValidEmail(email)) {
      setErrorMsg("Please enter a valid email.");
      return;
    }

    if (!isValidPassword(password)) {
      setErrorMsg("Password must be at least 6 characters long.");
      return;
    }

    try {
      setIsLoading(true);
      await signupAndRegister(email, password);

      navigate("/");
    } catch (err: any) {
      setErrorMsg(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form className="signup-form" onSubmit={handleSignup}>
      <input
        type="email"
        placeholder="Email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        required
        autoComplete="email"
      />
      <input
        type="password"
        placeholder="Password (min 6 characters)"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        required
        autoComplete="new-password"
      />
      <button type="submit" disabled={isLoading}>
        {isLoading ? "Signing up..." : "Sign Up"}
      </button>

      {errorMsg && <p className="error">{errorMsg}</p>}
    </form>
  );
};
