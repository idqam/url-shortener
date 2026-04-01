/* eslint-disable @typescript-eslint/no-explicit-any */
import { signupAndRegister } from "../api/signup";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Eye, EyeOff, AlertCircle, CheckCircle } from "lucide-react";

export const SignupForm = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const navigate = useNavigate();

  const isValidEmail = (email: string) =>
    /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);

  const isValidPassword = (password: string) => password.length >= 6;

  const passwordsMatch = password === confirmPassword;

  const getPasswordStrength = (password: string) => {
    if (password.length === 0) return { strength: 0, text: "" };
    if (password.length < 6) return { strength: 1, text: "Too short" };
    if (password.length < 8) return { strength: 2, text: "Weak" };
    if (!/[A-Z]/.test(password) || !/[0-9]/.test(password))
      return { strength: 3, text: "Good" };
    return { strength: 4, text: "Strong" };
  };

  const passwordStrength = getPasswordStrength(password);

  const handleSignup = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg(null);

    if (!isValidEmail(email)) {
      setErrorMsg("Please enter a valid email address.");
      return;
    }
    if (!isValidPassword(password)) {
      setErrorMsg("Password must be at least 6 characters long.");
      return;
    }
    if (!passwordsMatch) {
      setErrorMsg("Passwords do not match.");
      return;
    }

    try {
      setIsLoading(true);
      await signupAndRegister(email, password);
      navigate("/login");
    } catch (err: any) {
      setErrorMsg(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  const inputClass =
    "w-full px-4 py-3 border border-gray-300 rounded-md text-gray-900 placeholder-gray-400 focus:outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 transition-colors text-sm";

  return (
    <form className="space-y-5" onSubmit={handleSignup}>
      <div>
        <label htmlFor="signup-email" className="block text-sm font-medium text-gray-700 mb-1.5">
          Email Address
        </label>
        <input
          id="signup-email"
          name="email"
          type="email"
          autoComplete="email"
          required
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className={inputClass}
          placeholder="Enter your email"
        />
        {email && !isValidEmail(email) && (
          <p className="mt-1 text-xs text-red-600">Please enter a valid email address</p>
        )}
      </div>

      <div>
        <label htmlFor="signup-password" className="block text-sm font-medium text-gray-700 mb-1.5">
          Password
        </label>
        <div className="relative">
          <input
            id="signup-password"
            name="password"
            type={showPassword ? "text" : "password"}
            autoComplete="new-password"
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className={`${inputClass} pr-10`}
            placeholder="Minimum 6 characters"
          />
          <button
            type="button"
            className="absolute inset-y-0 right-0 pr-3 flex items-center"
            onClick={() => setShowPassword(!showPassword)}
          >
            {showPassword ? (
              <EyeOff className="h-4 w-4 text-gray-400" />
            ) : (
              <Eye className="h-4 w-4 text-gray-400" />
            )}
          </button>
        </div>

        {password && (
          <div className="mt-2 flex items-center gap-2">
            <div className="flex-1 bg-gray-200 rounded-full h-1.5">
              <div
                className="h-1.5 rounded-full transition-all"
                style={{
                  width: `${(passwordStrength.strength / 4) * 100}%`,
                  backgroundColor:
                    passwordStrength.strength === 1 ? "#EF4444" :
                    passwordStrength.strength === 2 ? "#F59E0B" :
                    passwordStrength.strength === 3 ? "#10B981" : "#059669",
                }}
              />
            </div>
            <span className="text-xs text-gray-500">{passwordStrength.text}</span>
          </div>
        )}
      </div>

      <div>
        <label htmlFor="confirm-password" className="block text-sm font-medium text-gray-700 mb-1.5">
          Confirm Password
        </label>
        <div className="relative">
          <input
            id="confirm-password"
            name="confirmPassword"
            type={showConfirmPassword ? "text" : "password"}
            autoComplete="new-password"
            required
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            className={`${inputClass} pr-10`}
            placeholder="Confirm your password"
          />
          <button
            type="button"
            className="absolute inset-y-0 right-0 pr-3 flex items-center"
            onClick={() => setShowConfirmPassword(!showConfirmPassword)}
          >
            {showConfirmPassword ? (
              <EyeOff className="h-4 w-4 text-gray-400" />
            ) : (
              <Eye className="h-4 w-4 text-gray-400" />
            )}
          </button>
        </div>

        {confirmPassword && (
          <div className="mt-1.5 flex items-center gap-1.5">
            {passwordsMatch ? (
              <>
                <CheckCircle className="h-3.5 w-3.5 text-green-500" />
                <span className="text-xs text-green-600">Passwords match</span>
              </>
            ) : (
              <>
                <AlertCircle className="h-3.5 w-3.5 text-red-500" />
                <span className="text-xs text-red-600">Passwords do not match</span>
              </>
            )}
          </div>
        )}
      </div>

      {errorMsg && (
        <div className="flex items-center gap-2 p-3 rounded-md bg-red-50 border border-red-200">
          <AlertCircle className="h-4 w-4 text-red-500 flex-shrink-0" />
          <p className="text-sm text-red-700">{errorMsg}</p>
        </div>
      )}

      <button
        type="submit"
        disabled={
          isLoading ||
          !isValidEmail(email) ||
          !isValidPassword(password) ||
          !passwordsMatch
        }
        className="w-full py-3 px-4 rounded-md text-white text-sm font-semibold bg-blue-600 hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
      >
        {isLoading ? (
          <>
            <svg className="animate-spin h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
            </svg>
            Creating Account...
          </>
        ) : (
          "Create Account"
        )}
      </button>
    </form>
  );
};
