import { useNavigate } from "react-router-dom";
import { ShortenerForm } from "../components/ShortenerForm";
import { Link as LinkIcon, Zap, Globe, BarChart3, QrCode } from "lucide-react";

export const Home = () => {
  const navigate = useNavigate();

  return (
    <div style={{ minHeight: "100vh", background: "var(--bg)", display: "flex", flexDirection: "column", position: "relative", overflow: "hidden" }}>
      {/* Background orbs */}
      <div className="orb" style={{ width: 500, height: 500, background: "rgba(124,58,237,0.12)", top: -200, left: "50%", transform: "translateX(-60%)" }} />
      <div className="orb" style={{ width: 300, height: 300, background: "rgba(34,211,238,0.06)", top: 200, right: -100 }} />

      {/* Grid */}
      <div className="bg-grid" style={{ position: "absolute", inset: 0, opacity: 0.5, maskImage: "radial-gradient(ellipse 80% 60% at 50% 0%, black 40%, transparent 100%)" }} />

      {/* Nav */}
      <header style={{ borderBottom: "1px solid var(--border)", position: "relative", zIndex: 10 }}>
        <div style={{ maxWidth: 960, margin: "0 auto", padding: "0 1.5rem", height: 56, display: "flex", alignItems: "center", justifyContent: "space-between" }}>
          <div style={{ display: "flex", alignItems: "center", gap: "0.5rem" }}>
            <div style={{ width: 28, height: 28, background: "var(--accent)", borderRadius: 8, display: "flex", alignItems: "center", justifyContent: "center", boxShadow: "0 0 12px var(--accent-glow)" }}>
              <LinkIcon size={14} color="#fff" />
            </div>
            <span className="font-display" style={{ fontWeight: 700, color: "var(--text)", fontSize: "0.9375rem", letterSpacing: "-0.01em" }}>Shortlink</span>
          </div>
          <div style={{ display: "flex", alignItems: "center", gap: "0.5rem" }}>
            <button
              onClick={() => navigate("/login")}
              className="btn-ghost"
            >
              Sign in
            </button>
            <button
              onClick={() => navigate("/signup")}
              className="btn-primary"
              style={{ padding: "0.5rem 1.125rem", fontSize: "0.875rem" }}
            >
              Get started
            </button>
          </div>
        </div>
      </header>

      {/* Hero */}
      <main style={{ flex: 1, display: "flex", flexDirection: "column", alignItems: "center", padding: "5rem 1.5rem 4rem", position: "relative", zIndex: 1 }}>
        {/* Badge */}
        <div className="animate-fade-up" style={{ display: "inline-flex", alignItems: "center", gap: "0.375rem", padding: "0.3rem 0.75rem", background: "var(--accent-dim)", border: "1px solid rgba(124,58,237,0.3)", borderRadius: 999, marginBottom: "1.75rem" }}>
          <Zap size={11} color="var(--accent-bright)" />
          <span style={{ fontSize: "0.75rem", fontWeight: 600, color: "var(--accent-bright)", letterSpacing: "0.05em", textTransform: "uppercase" }}>Instant URL Shortening</span>
        </div>

        {/* Headline */}
        <div className="animate-fade-up-d1" style={{ maxWidth: 620, textAlign: "center", marginBottom: "3.5rem" }}>
          <h1 className="font-display text-glow" style={{ fontSize: "clamp(2.25rem, 5vw, 3.5rem)", fontWeight: 800, lineHeight: 1.1, letterSpacing: "-0.03em", color: "var(--text)", marginBottom: "1rem" }}>
            Shorten any URL,{" "}
            <span style={{ color: "var(--accent-bright)" }}>instantly.</span>
          </h1>
          <p style={{ fontSize: "1.0625rem", color: "var(--text-secondary)", lineHeight: 1.6, fontWeight: 300 }}>
            Paste a long link, get a clean short one. Track every click with built-in analytics.
          </p>
        </div>

        {/* Form */}
        <div className="animate-fade-up-d2" style={{ width: "100%", maxWidth: 620 }}>
          <ShortenerForm />
        </div>

        {/* Divider */}
        <div className="animate-fade-up-d3" style={{ width: "100%", maxWidth: 620, margin: "4rem 0 0", position: "relative" }}>
          <div style={{ height: 1, background: "linear-gradient(90deg, transparent, var(--border-bright), transparent)" }} />
          <span style={{ position: "absolute", top: "50%", left: "50%", transform: "translate(-50%, -50%)", background: "var(--bg)", padding: "0 1rem", fontSize: "0.6875rem", fontWeight: 600, letterSpacing: "0.12em", textTransform: "uppercase", color: "var(--text-muted)" }}>Coming soon</span>
        </div>

        {/* Features */}
        <div className="animate-fade-up-d4" style={{ width: "100%", maxWidth: 620, marginTop: "2.5rem", display: "grid", gridTemplateColumns: "repeat(3, 1fr)", gap: "1rem" }}>
          {[
            { icon: Globe, title: "Custom Domains", desc: "Use your own branded domain" },
            { icon: BarChart3, title: "Advanced Analytics", desc: "Fine-grain click insights" },
            { icon: QrCode, title: "QR Codes", desc: "Instant QR generation" },
          ].map(({ icon: Icon, title, desc }) => (
            <div
              key={title}
              className="card-dark"
              style={{ padding: "1.25rem", textAlign: "center", transition: "border-color 0.2s, transform 0.2s" }}
              onMouseEnter={e => {
                (e.currentTarget as HTMLDivElement).style.borderColor = "var(--border-bright)";
                (e.currentTarget as HTMLDivElement).style.transform = "translateY(-2px)";
              }}
              onMouseLeave={e => {
                (e.currentTarget as HTMLDivElement).style.borderColor = "var(--border)";
                (e.currentTarget as HTMLDivElement).style.transform = "translateY(0)";
              }}
            >
              <div style={{ width: 36, height: 36, background: "var(--accent-dim)", border: "1px solid rgba(124,58,237,0.2)", borderRadius: 8, display: "flex", alignItems: "center", justifyContent: "center", margin: "0 auto 0.75rem" }}>
                <Icon size={16} color="var(--accent-bright)" />
              </div>
              <p className="font-display" style={{ fontSize: "0.875rem", fontWeight: 600, color: "var(--text)", marginBottom: "0.25rem" }}>{title}</p>
              <p style={{ fontSize: "0.75rem", color: "var(--text-muted)", lineHeight: 1.5 }}>{desc}</p>
            </div>
          ))}
        </div>
      </main>

      {/* Footer */}
      <footer style={{ borderTop: "1px solid var(--border)", padding: "1rem", position: "relative", zIndex: 1 }}>
        <p style={{ fontSize: "0.75rem", color: "var(--text-muted)", textAlign: "center" }}>
          Basic analytics collected: click counts, referrer sources, device types. No personal data stored with your links.
        </p>
      </footer>
    </div>
  );
};
