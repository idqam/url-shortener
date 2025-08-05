/* eslint-disable @typescript-eslint/no-explicit-any */
import { useEffect, useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { loginUser } from "../api/login";
import { useAuthStore } from "../store/AuthStore";
import { supabase } from "../lib/supabaseClient";
import { Eye, EyeOff, Link as LinkIcon, AlertCircle } from "lucide-react";

const USE_DUMMY_USER = import.meta.env.VITE_USE_DUMMY_USER === "true";

const LoginPage = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [errorMsg, setErrorMsg] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const navigate = useNavigate();
  const { login } = useAuthStore();

  useEffect(() => {
    const maybeAutoLogin = async () => {
      if (!USE_DUMMY_USER) return;

      try {
        const {
          data: { session },
        } = await supabase.auth.getSession();
        if (session?.user?.id && session?.access_token) {
          login(session.user.id, session.access_token);
          navigate("/dashboard");
        }
      } catch (error) {
        console.error("Auto-login failed:", error);
      }
    };

    maybeAutoLogin();
  }, [navigate, login]);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg("");
    setIsLoading(true);

    try {
      const { userId, accessToken } = await loginUser(email, password);

      if (!userId || !accessToken) {
        throw new Error("Failed to get user credentials");
      }

      login(userId, accessToken);
      navigate("/dashboard");
    } catch (err: any) {
      setErrorMsg(err.message || "Login failed");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div
      className="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8"
      style={{ backgroundColor: "#F2EFDE" }}
    >
      <div className="max-w-md w-full space-y-8">
        <div className="text-center">
          <div className="flex justify-center items-center mb-6">
            <div
              className="p-4 rounded-2xl shadow-lg"
              style={{
                background: `linear-gradient(135deg, #A4193D 0%, #F4A261 100%)`,
              }}
            >
              <LinkIcon className="h-12 w-12 text-white" />
            </div>
          </div>
          <h2
            className="text-4xl font-bold mb-2"
            style={{
              background: `linear-gradient(135deg, #A4193D 0%, #F4A261 100%)`,
              WebkitBackgroundClip: "text",
              WebkitTextFillColor: "transparent",
              backgroundClip: "text",
            }}
          >
            Welcome Back
          </h2>
          <p className="text-lg" style={{ color: "#6B7280" }}>
            Sign in to your account to continue
          </p>
        </div>

        <div className="bg-white rounded-2xl shadow-xl p-8 border border-gray-100">
          <form className="space-y-6" onSubmit={handleLogin}>
            <div>
              <label
                htmlFor="email"
                className="block text-sm font-semibold mb-2"
                style={{ color: "#111827" }}
              >
                Email Address
              </label>
              <input
                id="email"
                name="email"
                type="email"
                autoComplete="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full px-4 py-3 border-2 border-gray-200 rounded-xl text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-0 transition-all duration-200"
                style={{
                  borderBlock: "#A4193D",
                }}
                onFocus={(e) => (e.target.style.borderColor = "#A4193D")}
                onBlur={(e) => (e.target.style.borderColor = "#E5E7EB")}
                placeholder="Enter your email"
              />
            </div>

            <div>
              <label
                htmlFor="password"
                className="block text-sm font-semibold mb-2"
                style={{ color: "#111827" }}
              >
                Password
              </label>
              <div className="relative">
                <input
                  id="password"
                  name="password"
                  type={showPassword ? "text" : "password"}
                  autoComplete="current-password"
                  required
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="w-full px-4 py-3 pr-12 border-2 border-gray-200 rounded-xl text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-0 transition-all duration-200"
                  onFocus={(e) => (e.target.style.borderColor = "#A4193D")}
                  onBlur={(e) => (e.target.style.borderColor = "#E5E7EB")}
                  placeholder="Enter your password"
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
            </div>

            {errorMsg && (
              <div className="flex items-center space-x-2 p-4 rounded-xl bg-red-50 border border-red-200">
                <AlertCircle className="h-5 w-5 text-red-500 flex-shrink-0" />
                <p className="text-sm text-red-700 font-medium">{errorMsg}</p>
              </div>
            )}

            <button
              type="submit"
              disabled={isLoading}
              className="w-full flex justify-center items-center py-3 px-4 border border-transparent rounded-xl shadow-lg text-white font-semibold text-lg transition-all duration-200 transform hover:scale-105 hover:shadow-xl disabled:opacity-50 disabled:transform-none disabled:hover:shadow-lg"
              style={{
                background: isLoading
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
                  Signing in...
                </>
              ) : (
                "Sign In"
              )}
            </button>

            <div className="relative my-6">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-gray-200" />
              </div>
              <div className="relative flex justify-center text-sm">
                <span className="px-4 bg-white text-gray-500 font-medium">
                  Don't have an account?
                </span>
              </div>
            </div>

            <div className="text-center">
              <Link
                to="/signup"
                className="inline-flex items-center px-6 py-2 border-2 rounded-xl font-semibold text-sm transition-all duration-200 hover:shadow-md"
                style={{
                  borderColor: "#A4193D",
                  color: "#A4193D",
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.backgroundColor = "#A4193D";
                  e.currentTarget.style.color = "white";
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.backgroundColor = "transparent";
                  e.currentTarget.style.color = "#A4193D";
                }}
              >
                Create New Account
              </Link>
            </div>

            {USE_DUMMY_USER && (
              <div className="text-center">
                <p className="text-xs text-gray-500 italic">
                  Dev Mode: Auto-login enabled
                </p>
              </div>
            )}
          </form>
        </div>

        <div className="text-center">
          <p className="text-sm" style={{ color: "#6B7280" }}>
            Secure login powered by Supabase
          </p>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
