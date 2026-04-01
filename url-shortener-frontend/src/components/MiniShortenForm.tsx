/* eslint-disable @typescript-eslint/no-explicit-any */
import { useState } from "react";
import { useAuthStore } from "../store/AuthStore";
import { useShortenURL } from "../api/urls";
import { isValidUrl } from "../utils/isValidUrl";
import { ArrowRight, ExternalLink } from "lucide-react";

export function MiniShortenerForm() {
  const [url, setUrl] = useState("");
  const [result, setResult] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
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

  return (
    <div style={{ display: "flex", flexDirection: "column", gap: "0.5rem" }}>
      <form onSubmit={handleSubmit} style={{ display: "flex", gap: "0.5rem", background: "var(--bg-card)", border: "1px solid var(--border)", borderRadius: "0.625rem", padding: "0.3rem" }}>
        <input
          type="url"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="Paste a URL to shorten"
          required
          style={{
            flex: 1, padding: "0.5rem 0.75rem",
            background: "transparent", border: "none", outline: "none",
            color: "var(--text)", fontSize: "0.875rem", fontFamily: "DM Sans, sans-serif",
          }}
        />
        <button
          type="submit"
          disabled={isPending}
          className="btn-primary"
          style={{ display: "flex", alignItems: "center", gap: "0.25rem", padding: "0.5rem 1rem", borderRadius: "0.4rem", fontSize: "0.8125rem" }}
        >
          {isPending ? "..." : <><span>Shorten</span><ArrowRight size={12} /></>}
        </button>
      </form>

      {result && (
        <div style={{ display: "flex", alignItems: "center", gap: "0.375rem", padding: "0.5rem 0.75rem", background: "rgba(34,211,238,0.06)", border: "1px solid rgba(34,211,238,0.2)", borderRadius: "0.5rem" }}>
          <ExternalLink size={12} color="var(--cyan)" style={{ flexShrink: 0 }} />
          <a href={result} target="_blank" rel="noopener noreferrer" style={{ fontSize: "0.8125rem", color: "var(--cyan)", textDecoration: "none", fontFamily: "DM Mono, monospace", wordBreak: "break-all" }}>
            {result}
          </a>
        </div>
      )}

      {error && (
        <p style={{ fontSize: "0.8125rem", color: "rgb(248,113,113)", margin: 0 }}>{error}</p>
      )}
    </div>
  );
}
