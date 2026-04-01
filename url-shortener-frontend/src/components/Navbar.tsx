import { Link as LinkIcon, Zap, LogOut } from "lucide-react";
import { useAuthStore } from "../store/AuthStore";
import { supabase } from "../lib/supabaseClient";

interface NavbarProps {
  onToggleShortener: () => void;
  showShortener: boolean;
}

export function Navbar({ onToggleShortener, showShortener }: NavbarProps) {
  const logout = useAuthStore((state) => state.logout);

  const handleSignOut = async () => {
    await supabase.auth.signOut();
    logout();
  };

  return (
    <header
      style={{
        position: "sticky",
        top: 0,
        zIndex: 10,
        background: "rgba(7,7,15,0.85)",
        backdropFilter: "blur(12px)",
        borderBottom: "1px solid var(--border)",
      }}
    >
      <div
        style={{
          maxWidth: 1280,
          margin: "0 auto",
          padding: "0 1.5rem",
          height: 56,
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
        }}
      >
        <div style={{ display: "flex", alignItems: "center", gap: "0.5rem" }}>
          <div
            style={{
              width: 28,
              height: 28,
              background: "var(--accent)",
              borderRadius: 8,
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              boxShadow: "0 0 12px var(--accent-glow)",
            }}
          >
            <LinkIcon size={14} color="#fff" />
          </div>
          <span
            className="font-display"
            style={{
              fontWeight: 700,
              color: "var(--text)",
              fontSize: "0.9375rem",
              letterSpacing: "-0.01em",
            }}
          >
            LinkForge
          </span>
        </div>

        <div style={{ display: "flex", alignItems: "center", gap: "0.5rem" }}>
          <button
            onClick={onToggleShortener}
            style={{
              display: "flex",
              alignItems: "center",
              gap: "0.375rem",
              padding: "0.4375rem 0.875rem",
              background: showShortener ? "var(--accent-dim)" : "transparent",
              border: `1px solid ${showShortener ? "rgba(124,58,237,0.35)" : "var(--border)"}`,
              borderRadius: "0.5rem",
              color: showShortener
                ? "var(--accent-bright)"
                : "var(--text-secondary)",
              fontSize: "0.8125rem",
              fontWeight: 500,
              cursor: "pointer",
              transition: "all 0.2s",
              fontFamily: "DM Sans, sans-serif",
            }}
            onMouseEnter={(e) => {
              if (!showShortener) {
                (e.currentTarget as HTMLButtonElement).style.borderColor =
                  "var(--border-bright)";
                (e.currentTarget as HTMLButtonElement).style.color =
                  "var(--text)";
              }
            }}
            onMouseLeave={(e) => {
              if (!showShortener) {
                (e.currentTarget as HTMLButtonElement).style.borderColor =
                  "var(--border)";
                (e.currentTarget as HTMLButtonElement).style.color =
                  "var(--text-secondary)";
              }
            }}
          >
            <Zap size={13} />
            {showShortener ? "Hide" : "Quick Shorten"}
          </button>
          <button
            onClick={handleSignOut}
            style={{
              display: "flex",
              alignItems: "center",
              gap: "0.375rem",
              padding: "0.4375rem 0.875rem",
              background: "transparent",
              border: "1px solid var(--border)",
              borderRadius: "0.5rem",
              color: "var(--text-secondary)",
              fontSize: "0.8125rem",
              fontWeight: 500,
              cursor: "pointer",
              transition: "all 0.2s",
              fontFamily: "DM Sans, sans-serif",
            }}
            onMouseEnter={(e) => {
              (e.currentTarget as HTMLButtonElement).style.borderColor =
                "rgba(239,68,68,0.4)";
              (e.currentTarget as HTMLButtonElement).style.color =
                "rgb(248,113,113)";
              (e.currentTarget as HTMLButtonElement).style.background =
                "rgba(239,68,68,0.07)";
            }}
            onMouseLeave={(e) => {
              (e.currentTarget as HTMLButtonElement).style.borderColor =
                "var(--border)";
              (e.currentTarget as HTMLButtonElement).style.color =
                "var(--text-secondary)";
              (e.currentTarget as HTMLButtonElement).style.background =
                "transparent";
            }}
          >
            <LogOut size={13} />
            Sign out
          </button>
        </div>
      </div>
    </header>
  );
}
