/* eslint-disable @typescript-eslint/no-explicit-any */
import { useState } from "react";
import { useAuthStore } from "../store/AuthStore";
import { useShortenURL } from "../api/urls";
import { isValidUrl } from "../utils/isValidUrl";

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
    <div>
      <form onSubmit={handleSubmit} className="flex flex-col sm:flex-row gap-2">
        <input
          type="url"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="Paste a URL to shorten"
          required
          className="flex-1 px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 placeholder-gray-400"
        />
        <button
          type="submit"
          disabled={isPending}
          className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 transition-colors disabled:opacity-60"
        >
          {isPending ? "Shortening..." : "Shorten"}
        </button>
      </form>

      {result && (
        <div className="mt-2 text-sm bg-gray-50 border border-gray-200 rounded-md px-3 py-2">
          <a href={result} target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline">
            {result}
          </a>
        </div>
      )}

      {error && (
        <p className="mt-2 text-sm text-red-600">{error}</p>
      )}
    </div>
  );
}
