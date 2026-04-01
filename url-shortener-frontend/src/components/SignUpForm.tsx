/* eslint-disable @typescript-eslint/no-explicit-any */
import { signupAndRegister } from "../api/signup";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Eye, EyeOff, AlertCircle, CheckCircle, ArrowRight } from "lucide-react";

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

  const strengthColors = ["", "#ef4444", "#f59e0b", "#10b981", "#059669"];
  const strengthBgColors = ["", "rgba(239,68,68,0.1)", "rgba(245,158,11,0.1)", "rgba(16,185,129,0.1)", "rgba(5,150,105,0.1)"];

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

  const labelStyle: React.CSSProperties = {
    display: "block",
    fontSize: "0.8125rem",
    fontWeight: 600,
    color: "var(--text-secondary)",
    marginBottom: "0.5rem",
    letterSpacing: "0.02em",
  };

  return (
    <form style={{ display: "flex", flexDirection: "column", gap: "1.25rem" }} onSubmit={handleSignup}>
      {/* Email */}
      <div>
        <label htmlFor="signup-email" style={labelStyle}>Email address</label>
        <input
          id="signup-email"
          name="email"
          type="email"
          autoComplete="email"
          required
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="input-dark"
          placeholder="you@example.com"
        />
        {email && !isValidEmail(email) && (
          <p style={{ marginTop: "0.375rem", fontSize: "0.75rem", color: "rgb(248,113,113)" }}>
            Please enter a valid email address
          </p>
        )}
      </div>

      {/* Password */}
      <div>
        <label htmlFor="signup-password" style={labelStyle}>Password</label>
        <div style={{ position: "relative" }}>
          <input
            id="signup-password"
            name="password"
            type={showPassword ? "text" : "password"}
            autoComplete="new-password"
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="input-dark"
            placeholder="Minimum 6 characters"
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

        {password && (
          <div style={{ marginTop: "0.625rem", display: "flex", alignItems: "center", gap: "0.625rem" }}>
            <div style={{ flex: 1, height: 3, background: "var(--bg-elevated)", borderRadius: 999, overflow: "hidden" }}>
              <div
                style={{
                  height: "100%",
                  borderRadius: 999,
                  transition: "width 0.3s, background 0.3s",
                  width: `${(passwordStrength.strength / 4) * 100}%`,
                  background: strengthColors[passwordStrength.strength],
                  boxShadow: `0 0 8px ${strengthBgColors[passwordStrength.strength]}`,
                }}
              />
            </div>
            <span style={{ fontSize: "0.75rem", color: strengthColors[passwordStrength.strength] || "var(--text-muted)", fontWeight: 500 }}>
              {passwordStrength.text}
            </span>
          </div>
        )}
      </div>

      {/* Confirm Password */}
      <div>
        <label htmlFor="confirm-password" style={labelStyle}>Confirm password</label>
        <div style={{ position: "relative" }}>
          <input
            id="confirm-password"
            name="confirmPassword"
            type={showConfirmPassword ? "text" : "password"}
            autoComplete="new-password"
            required
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            className="input-dark"
            placeholder="Confirm your password"
            style={{ paddingRight: "2.75rem" }}
          />
          <button
            type="button"
            onClick={() => setShowConfirmPassword(!showConfirmPassword)}
            style={{ position: "absolute", right: "0.875rem", top: "50%", transform: "translateY(-50%)", background: "none", border: "none", cursor: "pointer", color: "var(--text-muted)", display: "flex", alignItems: "center", padding: 0, transition: "color 0.2s" }}
            onMouseEnter={e => (e.currentTarget.style.color = "var(--text-secondary)")}
            onMouseLeave={e => (e.currentTarget.style.color = "var(--text-muted)")}
          >
            {showConfirmPassword ? <EyeOff size={16} /> : <Eye size={16} />}
          </button>
        </div>

        {confirmPassword && (
          <div style={{ marginTop: "0.5rem", display: "flex", alignItems: "center", gap: "0.375rem" }}>
            {passwordsMatch ? (
              <>
                <CheckCircle size={13} color="#10b981" />
                <span style={{ fontSize: "0.75rem", color: "#10b981" }}>Passwords match</span>
              </>
            ) : (
              <>
                <AlertCircle size={13} color="rgb(248,113,113)" />
                <span style={{ fontSize: "0.75rem", color: "rgb(248,113,113)" }}>Passwords do not match</span>
              </>
            )}
          </div>
        )}
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
        disabled={isLoading || !isValidEmail(email) || !isValidPassword(password) || !passwordsMatch}
        className="btn-primary"
        style={{ width: "100%", display: "flex", alignItems: "center", justifyContent: "center", gap: "0.5rem", marginTop: "0.25rem" }}
      >
        {isLoading ? (
          <>
            <svg style={{ animation: "spin 0.7s linear infinite", width: 15, height: 15 }} viewBox="0 0 24 24" fill="none">
              <circle cx="12" cy="12" r="10" stroke="rgba(255,255,255,0.25)" strokeWidth="3" />
              <path d="M12 2a10 10 0 0 1 10 10" stroke="white" strokeWidth="3" strokeLinecap="round" />
            </svg>
            Creating account...
          </>
        ) : (
          <>
            Create account
            <ArrowRight size={15} />
          </>
        )}
      </button>
    </form>
  );
};
