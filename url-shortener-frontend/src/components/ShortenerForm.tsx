/* eslint-disable @typescript-eslint/no-explicit-any */
import { useState } from "react";
import { useAuthStore } from "../store/AuthStore";
import { useShortenURL } from "../api/urls";
import { isValidUrl } from "../utils/isValidUrl";
import { Copy, Check, ExternalLink, AlertCircle, ArrowRight } from "lucide-react";

export function ShortenerForm() {
  const [url, setUrl] = useState("");
  const [result, setResult] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const userId = useAuthStore((state) => state.userId);
  const { mutateAsync: shorten, isPending } = useShortenURL();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!url) {
      setError("Please enter a URL.");
      return;
    }
    if (!isValidUrl(url)) {
      setError("Please enter a valid URL (http/https).");
      return;
    }
    try {
      const data = await shorten({
        url,
        code_length: 8,
        is_public: true,
        user_id: userId ?? null,
      });
      setResult(data.short_url);
      setError(null);
    } catch (err: any) {
      setError(err?.message || "An error occurred.");
      setResult(null);
    }
    setUrl("");
  };

  const copyToClipboard = async () => {
    if (result) {
      try {
        await navigator.clipboard.writeText(result);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
      } catch (err) {
        console.error("Failed to copy:", err);
      }
    }
  };

  return (
    <div style={{ display: "flex", flexDirection: "column", gap: "0.75rem" }}>
      {/* Input row */}
      <div style={{ display: "flex", gap: "0.625rem", background: "var(--bg-card)", border: "1px solid var(--border)", borderRadius: "0.75rem", padding: "0.375rem", transition: "border-color 0.2s, box-shadow 0.2s" }}
        onFocusCapture={e => {
          (e.currentTarget as HTMLDivElement).style.borderColor = "var(--accent-bright)";
          (e.currentTarget as HTMLDivElement).style.boxShadow = "0 0 0 3px var(--accent-dim), 0 0 20px var(--accent-glow)";
        }}
        onBlurCapture={e => {
          (e.currentTarget as HTMLDivElement).style.borderColor = "var(--border)";
          (e.currentTarget as HTMLDivElement).style.boxShadow = "none";
        }}
      >
        <input
          type="url"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleSubmit(e)}
          placeholder="Paste your long URL here..."
          required
          style={{
            flex: 1,
            padding: "0.75rem 1rem",
            background: "transparent",
            border: "none",
            outline: "none",
            color: "var(--text)",
            fontSize: "0.9375rem",
            fontFamily: "DM Sans, sans-serif",
          }}
        />
        <button
          onClick={handleSubmit}
          disabled={isPending}
          className="btn-primary"
          style={{ display: "flex", alignItems: "center", gap: "0.375rem", padding: "0.625rem 1.25rem", borderRadius: "0.5rem", fontSize: "0.875rem" }}
        >
          {isPending ? (
            <>
              <svg style={{ animation: "spin 0.7s linear infinite", width: 15, height: 15 }} viewBox="0 0 24 24" fill="none">
                <style>{`@keyframes spin { to { transform: rotate(360deg); } }`}</style>
                <circle cx="12" cy="12" r="10" stroke="rgba(255,255,255,0.25)" strokeWidth="3" />
                <path d="M12 2a10 10 0 0 1 10 10" stroke="white" strokeWidth="3" strokeLinecap="round" />
              </svg>
              Shortening
            </>
          ) : (
            <>
              Shorten
              <ArrowRight size={14} />
            </>
          )}
        </button>
      </div>

      {/* Result */}
      {result && (
        <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", gap: "0.75rem", padding: "0.875rem 1.125rem", background: "rgba(34,211,238,0.06)", border: "1px solid rgba(34,211,238,0.2)", borderRadius: "0.625rem", animation: "fadeUp 0.3s ease" }}>
          <a
            href={result}
            target="_blank"
            rel="noopener noreferrer"
            style={{ display: "flex", alignItems: "center", gap: "0.375rem", fontSize: "0.875rem", fontWeight: 500, color: "var(--cyan)", textDecoration: "none", wordBreak: "break-all", fontFamily: "DM Mono, monospace" }}
          >
            <ExternalLink size={13} style={{ flexShrink: 0 }} />
            {result}
          </a>
          <button
            onClick={copyToClipboard}
            style={{
              display: "flex", alignItems: "center", gap: "0.25rem", padding: "0.375rem 0.75rem",
              background: copied ? "rgba(34,211,238,0.15)" : "rgba(255,255,255,0.05)",
              border: `1px solid ${copied ? "rgba(34,211,238,0.4)" : "var(--border-bright)"}`,
              borderRadius: "0.375rem", color: copied ? "var(--cyan)" : "var(--text-secondary)",
              fontSize: "0.75rem", fontWeight: 600, cursor: "pointer", flexShrink: 0,
              transition: "all 0.2s", fontFamily: "Syne, sans-serif"
            }}
          >
            {copied ? <Check size={12} /> : <Copy size={12} />}
            {copied ? "Copied!" : "Copy"}
          </button>
        </div>
      )}

      {/* Error */}
      {error && (
        <div style={{ display: "flex", alignItems: "center", gap: "0.5rem", padding: "0.75rem 1rem", background: "rgba(239,68,68,0.07)", border: "1px solid rgba(239,68,68,0.25)", borderRadius: "0.625rem", animation: "fadeUp 0.3s ease" }}>
          <AlertCircle size={14} color="rgb(248,113,113)" style={{ flexShrink: 0 }} />
          <p style={{ fontSize: "0.875rem", color: "rgb(248,113,113)", margin: 0 }}>{error}</p>
        </div>
      )}
    </div>
  );
}
