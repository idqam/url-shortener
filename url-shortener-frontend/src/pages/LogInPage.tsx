/* eslint-disable @typescript-eslint/no-explicit-any */
import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { loginUser } from "../api/login";
import { useAuthStore } from "../store/AuthStore";
import { Eye, EyeOff, AlertCircle, Link as LinkIcon, ArrowRight } from "lucide-react";

const LoginPage = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [errorMsg, setErrorMsg] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const navigate = useNavigate();
  const { login } = useAuthStore();

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
    <div style={{ minHeight: "100vh", background: "var(--bg)", display: "flex", flexDirection: "column", position: "relative", overflow: "hidden" }}>
      {/* Background orbs */}
      <div className="orb" style={{ width: 450, height: 450, background: "rgba(124,58,237,0.1)", top: -150, left: "50%", transform: "translateX(-50%)" }} />
      <div className="orb" style={{ width: 250, height: 250, background: "rgba(34,211,238,0.06)", bottom: 50, right: -80 }} />
      <div className="bg-grid" style={{ position: "absolute", inset: 0, opacity: 0.4, maskImage: "radial-gradient(ellipse 70% 50% at 50% 0%, black 30%, transparent 100%)" }} />

      {/* Nav */}
      <header style={{ borderBottom: "1px solid var(--border)", position: "relative", zIndex: 10 }}>
        <div style={{ maxWidth: 960, margin: "0 auto", padding: "0 1.5rem", height: 56, display: "flex", alignItems: "center" }}>
          <div
            style={{ display: "flex", alignItems: "center", gap: "0.5rem", cursor: "pointer" }}
            onClick={() => navigate("/")}
          >
            <div style={{ width: 28, height: 28, background: "var(--accent)", borderRadius: 8, display: "flex", alignItems: "center", justifyContent: "center", boxShadow: "0 0 12px var(--accent-glow)" }}>
              <LinkIcon size={14} color="#fff" />
            </div>
            <span className="font-display" style={{ fontWeight: 700, color: "var(--text)", fontSize: "0.9375rem", letterSpacing: "-0.01em" }}>Shortlink</span>
          </div>
        </div>
      </header>

      {/* Main */}
      <main style={{ flex: 1, display: "flex", alignItems: "center", justifyContent: "center", padding: "2rem 1.5rem", position: "relative", zIndex: 1 }}>
        <div style={{ width: "100%", maxWidth: 420 }} className="animate-fade-up">
          {/* Heading */}
          <div style={{ marginBottom: "2rem" }}>
            <h1 className="font-display" style={{ fontSize: "2rem", fontWeight: 800, letterSpacing: "-0.03em", color: "var(--text)", marginBottom: "0.375rem" }}>
              Welcome back
            </h1>
            <p style={{ fontSize: "0.9375rem", color: "var(--text-secondary)", fontWeight: 300 }}>
              Sign in to your account to continue
            </p>
          </div>

          {/* Card */}
          <div className="card-dark" style={{ padding: "2rem" }}>
            <form style={{ display: "flex", flexDirection: "column", gap: "1.25rem" }} onSubmit={handleLogin}>
              {/* Email */}
              <div>
                <label htmlFor="email" style={{ display: "block", fontSize: "0.8125rem", fontWeight: 600, color: "var(--text-secondary)", marginBottom: "0.5rem", letterSpacing: "0.02em" }}>
                  Email address
                </label>
                <input
                  id="email"
                  name="email"
                  type="email"
                  autoComplete="email"
                  required
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="input-dark"
                  placeholder="you@example.com"
                />
              </div>

              {/* Password */}
              <div>
                <label htmlFor="password" style={{ display: "block", fontSize: "0.8125rem", fontWeight: 600, color: "var(--text-secondary)", marginBottom: "0.5rem", letterSpacing: "0.02em" }}>
                  Password
                </label>
                <div style={{ position: "relative" }}>
                  <input
                    id="password"
                    name="password"
                    type={showPassword ? "text" : "password"}
                    autoComplete="current-password"
                    required
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="input-dark"
                    placeholder="Enter your password"
                    style={{ paddingRight: "2.75rem" }}
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    style={{ position: "absolute", right: "0.875rem", top: "50%", transform: "translateY(-50%)", background: "none", border: "none", cursor: "pointer", color: "var(--text-muted)", display: "flex", alignItems: "center", padding: 0, transition: "color 0.2s" }}
                    onMouseEnter={e => (e.currentTarget.style.color = "var(--text-secondary)")}
                    onMouseLeave={e => (e.currentTarget.style.color = "var(--text-muted)")}
                  >
                    {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                  </button>
                </div>
              </div>

              {/* Error */}
              {errorMsg && (
                <div style={{ display: "flex", alignItems: "center", gap: "0.5rem", padding: "0.75rem 1rem", background: "rgba(239,68,68,0.07)", border: "1px solid rgba(239,68,68,0.25)", borderRadius: "0.5rem" }}>
                  <AlertCircle size={14} color="rgb(248,113,113)" style={{ flexShrink: 0 }} />
                  <p style={{ fontSize: "0.875rem", color: "rgb(248,113,113)", margin: 0 }}>{errorMsg}</p>
                </div>
              )}

              {/* Submit */}
              <button
                type="submit"
                disabled={isLoading}
                className="btn-primary"
                style={{ width: "100%", display: "flex", alignItems: "center", justifyContent: "center", gap: "0.5rem", marginTop: "0.25rem" }}
              >
                {isLoading ? (
                  <>
                    <svg style={{ animation: "spin 0.7s linear infinite", width: 15, height: 15 }} viewBox="0 0 24 24" fill="none">
                      <circle cx="12" cy="12" r="10" stroke="rgba(255,255,255,0.25)" strokeWidth="3" />
                      <path d="M12 2a10 10 0 0 1 10 10" stroke="white" strokeWidth="3" strokeLinecap="round" />
                    </svg>
                    Signing in...
                  </>
                ) : (
                  <>
                    Sign in
                    <ArrowRight size={15} />
                  </>
                )}
              </button>

              {/* Signup link */}
              <p style={{ textAlign: "center", fontSize: "0.875rem", color: "var(--text-muted)", margin: 0 }}>
                Don't have an account?{" "}
                <Link to="/signup" style={{ color: "var(--accent-bright)", fontWeight: 500, textDecoration: "none" }}>
                  Create one
                </Link>
              </p>
            </form>
          </div>
        </div>
      </main>
    </div>
  );
};

export default LoginPage;
