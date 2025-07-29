import { useState } from "react";
import { useShortenUrl } from "../hooks/useShortenUrl";
import { useUrlStore } from "../store/UrlStore";

export function ShortenerForm() {
  const { shorten } = useShortenUrl();
  const loading = useUrlStore((state) => state.loading);
  const [url, setUrl] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!url) return;

    shorten({
      url,
      is_public: true,
      user_id: null,
      code_length: 8,
    });

    setUrl("");
  };

  return (
    <form onSubmit={handleSubmit} className="shortener-box">
      <input
        type="url"
        value={url}
        onChange={(e) => setUrl(e.target.value)}
        placeholder="Paste your URL"
        required
        className="shortener-input"
      />

      <button className="shortener-box shorten-btn" disabled={loading}>
        {" "}
        {loading ? (
          <span className="loading-text">
            Loading
            <span className="dots">
              <span>.</span>
              <span>.</span>
              <span>.</span>
            </span>
          </span>
        ) : (
          "Shorten"
        )}
      </button>
    </form>
  );
}
