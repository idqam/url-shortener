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
      <div className="text-center mb-4">
        <p className="text-lg font-medium">
          Url Shortener is a free tool to shorten URLs and generate short links
        </p>
        <p className="text-lg text-gray-600">
          URL shortener creates easy to share short links that are protected and safe from malware and viruses.
        </p>
      </div>

      <input
        type="url"
        value={url}
        onChange={(e) => setUrl(e.target.value)}
        placeholder="Paste your URL"
        required
        className="shortener-input"
      />

      <button className="shortener-box shorten-btn" disabled={loading}>
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
