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

  return (
    <form className="space-y-6" onSubmit={handleSignup}>
      <div>
        <label
          htmlFor="signup-email"
          className="block text-sm font-semibold mb-2"
          style={{ color: "#111827" }}
        >
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
          className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-0 transition-all duration-200"
          onFocus={(e) => (e.target.style.borderColor = "#A4193D")}
          onBlur={(e) => (e.target.style.borderColor = "#E5E7EB")}
          placeholder="Enter your email"
        />
        {email && !isValidEmail(email) && (
          <p className="mt-1 text-sm text-red-600">
            Please enter a valid email address
          </p>
        )}
      </div>

      <div>
        <label
          htmlFor="signup-password"
          className="block text-sm font-semibold mb-2"
          style={{ color: "#111827" }}
        >
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
            className="w-full px-4 py-3 pr-12 border-2 border-gray-200 rounded-xl text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-0 transition-all duration-200"
            onFocus={(e) => (e.target.style.borderColor = "#A4193D")}
            onBlur={(e) => (e.target.style.borderColor = "#E5E7EB")}
            placeholder="Create a password (min 6 characters)"
          />
          <button
            type="button"
            className="absolute inset-y-0 right-0 pr-4 flex items-center"
            onClick={() => setShowPassword(!showPassword)}
          >
            {showPassword ? (
              <EyeOff className="h-5 w-5" style={{ color: "#6B7280" }} />
            ) : (
              <Eye className="h-5 w-5" style={{ color: "#6B7280" }} />
            )}
          </button>
        </div>

        {password && (
          <div className="mt-2">
            <div className="flex items-center space-x-2">
              <div className="flex-1 bg-gray-200 rounded-full h-2">
                <div
                  className="h-2 rounded-full transition-all duration-300"
                  style={{
                    width: `${(passwordStrength.strength / 4) * 100}%`,
                    backgroundColor:
                      passwordStrength.strength === 1
                        ? "#EF4444"
                        : passwordStrength.strength === 2
                        ? "#F59E0B"
                        : passwordStrength.strength === 3
                        ? "#10B981"
                        : passwordStrength.strength === 4
                        ? "#059669"
                        : "#E5E7EB",
                  }}
                />
              </div>
              <span
                className="text-xs font-medium"
                style={{ color: "#6B7280" }}
              >
                {passwordStrength.text}
              </span>
            </div>
          </div>
        )}
      </div>

      <div>
        <label
          htmlFor="confirm-password"
          className="block text-sm font-semibold mb-2"
          style={{ color: "#111827" }}
        >
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
            className="w-full px-4 py-3 pr-12 border-2 border-gray-200 rounded-xl text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-0 transition-all duration-200"
            onFocus={(e) => (e.target.style.borderColor = "#A4193D")}
            onBlur={(e) => (e.target.style.borderColor = "#E5E7EB")}
            placeholder="Confirm your password"
          />
          <button
            type="button"
            className="absolute inset-y-0 right-0 pr-4 flex items-center"
            onClick={() => setShowConfirmPassword(!showConfirmPassword)}
          >
            {showConfirmPassword ? (
              <EyeOff className="h-5 w-5" style={{ color: "#6B7280" }} />
            ) : (
              <Eye className="h-5 w-5" style={{ color: "#6B7280" }} />
            )}
          </button>
        </div>

        {confirmPassword && (
          <div className="mt-2 flex items-center space-x-2">
            {passwordsMatch ? (
              <>
                <CheckCircle className="h-4 w-4 text-green-500" />
                <span className="text-sm text-green-600 font-medium">
                  Passwords match
                </span>
              </>
            ) : (
              <>
                <AlertCircle className="h-4 w-4 text-red-500" />
                <span className="text-sm text-red-600 font-medium">
                  Passwords do not match
                </span>
              </>
            )}
          </div>
        )}
      </div>

      {errorMsg && (
        <div className="flex items-center space-x-2 p-4 rounded-xl bg-red-50 border border-red-200">
          <AlertCircle className="h-5 w-5 text-red-500 flex-shrink-0" />
          <p className="text-sm text-red-700 font-medium">{errorMsg}</p>
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
        className="w-full flex justify-center items-center py-3 px-4 border border-transparent rounded-xl shadow-lg text-white font-semibold text-lg transition-all duration-200 transform hover:scale-105 hover:shadow-xl disabled:opacity-50 disabled:transform-none disabled:hover:shadow-lg"
        style={{
          background:
            isLoading ||
            !isValidEmail(email) ||
            !isValidPassword(password) ||
            !passwordsMatch
              ? "#6B7280"
              : `linear-gradient(135deg, #A4193D 0%, #F4A261 100%)`,
        }}
      >
        {isLoading ? (
          <>
            <svg
              className="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                className="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                strokeWidth="4"
              ></circle>
              <path
                className="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
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
