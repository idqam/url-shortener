import { SignupForm } from "../components/SignUpForm";
import { Link, useNavigate } from "react-router-dom";
import { Link as LinkIcon } from "lucide-react";

export const SignUpPage = () => {
  const navigate = useNavigate();

  return (
    <div style={{ minHeight: "100vh", background: "var(--bg)", display: "flex", flexDirection: "column", position: "relative", overflow: "hidden" }}>
      {/* Background orbs */}
      <div className="orb" style={{ width: 450, height: 450, background: "rgba(124,58,237,0.1)", top: -150, left: "50%", transform: "translateX(-50%)" }} />
      <div className="orb" style={{ width: 200, height: 200, background: "rgba(34,211,238,0.06)", bottom: 100, left: -60 }} />
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
              Create an account
            </h1>
            <p style={{ fontSize: "0.9375rem", color: "var(--text-secondary)", fontWeight: 300 }}>
              Start shortening links and tracking analytics
            </p>
          </div>

          {/* Card */}
          <div className="card-dark" style={{ padding: "2rem" }}>
            <SignupForm />

            <p style={{ marginTop: "1.5rem", textAlign: "center", fontSize: "0.875rem", color: "var(--text-muted)" }}>
              Already have an account?{" "}
              <Link to="/login" style={{ color: "var(--accent-bright)", fontWeight: 500, textDecoration: "none" }}>
                Sign in
              </Link>
            </p>
          </div>
        </div>
      </main>
    </div>
  );
};
