/* eslint-disable @typescript-eslint/no-explicit-any */
import { useState } from "react";
import { useAuthStore } from "../store/AuthStore";
import { useShortenURL } from "../api/urls";
import { isValidUrl } from "../utils/isValidUrl";

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
    <div className="space-y-3">
      <div className="flex flex-col sm:flex-row gap-2">
        <input
          type="url"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleSubmit(e)}
          placeholder="Paste your URL here"
          required
          className="flex-1 px-4 py-3 border border-gray-300 rounded-md text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 transition-colors"
        />
        <button
          onClick={handleSubmit}
          disabled={isPending}
          className="px-6 py-3 text-sm font-semibold text-white bg-blue-600 rounded-md hover:bg-blue-700 transition-colors disabled:opacity-60 disabled:cursor-not-allowed"
        >
          {isPending ? "Shortening..." : "Shorten"}
        </button>
      </div>

      {result && (
        <div className="flex items-center justify-between gap-3 px-4 py-3 bg-green-50 border border-green-200 rounded-md">
          <a
            href={result}
            target="_blank"
            rel="noopener noreferrer"
            className="text-sm font-medium text-blue-600 hover:underline break-all"
          >
            {result}
          </a>
          <button
            onClick={copyToClipboard}
            className="text-xs font-medium text-gray-600 hover:text-gray-900 flex-shrink-0 border border-gray-300 rounded px-2 py-1 hover:bg-white transition-colors"
          >
            {copied ? "Copied!" : "Copy"}
          </button>
        </div>
      )}

      {error && (
        <p className="text-sm text-red-600 bg-red-50 border border-red-200 rounded-md px-4 py-3">
          {error}
        </p>
      )}
    </div>
  );
}
