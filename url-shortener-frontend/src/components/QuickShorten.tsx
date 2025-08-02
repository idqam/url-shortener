import React, { useState } from "react";
import { Link, Copy, Check } from "lucide-react";
import { useShortenUrlAuth } from "../hooks/useShortenUrlAuth";
import { useAuthStore } from "../store/AuthStore";
import { useUrlStore } from "../store/UrlStore";

export default function QuickShorten() {
  const [url, setUrl] = useState("");
  const [copied, setCopied] = useState(false);

  const { shorten } = useShortenUrlAuth();
  const { userId } = useAuthStore();
  const { loading: shortenLoading, error, lastAuthedResult } = useUrlStore();

  const handleShorten = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!url || !userId) return;

    await shorten({
      url,
      is_public: true,
      user_id: userId,
      code_length: 8,
    });
    setUrl("");
  };

  const handleCopy = async () => {
    if (lastAuthedResult?.short_url) {
      await navigator.clipboard.writeText(lastAuthedResult.short_url);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  const handleReset = () => {
    setUrl("");
    setCopied(false);
  };

  return (
    <div className="max-w-md mx-auto p-6 bg-white rounded-lg shadow-lg border">
      <div className="flex items-center gap-2 mb-4">
        <Link className="w-5 h-5 text-blue-600" />
        <h2 className="text-lg font-semibold text-gray-800">Quick Shorten</h2>
      </div>

      <div className="space-y-4">
        <div>
          <input
            type="url"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            placeholder="Enter your long URL here..."
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            disabled={shortenLoading}
          />
        </div>

        <button
          type="submit"
          onClick={handleShorten}
          disabled={!url || !userId || shortenLoading}
          className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors duration-200"
        >
          {shortenLoading ? "Shortening..." : "Shorten URL"}
        </button>

        {error && (
          <div className="mt-4 p-3 bg-red-50 border border-red-200 rounded-md">
            <p className="text-sm text-red-600">{error}</p>
          </div>
        )}

        {lastAuthedResult && (
          <div className="mt-4 p-3 bg-green-50 border border-green-200 rounded-md">
            <label className="block text-sm font-medium text-green-700 mb-2">
              Shortened URL:
            </label>
            <div className="flex items-center gap-2">
              <input
                type="text"
                value={lastAuthedResult.short_url}
                readOnly
                className="flex-1 px-2 py-1 bg-white border border-green-300 rounded text-sm"
              />
              <button
                onClick={handleCopy}
                className="p-1 text-green-600 hover:text-green-800 transition-colors"
                title="Copy to clipboard"
              >
                {copied ? (
                  <Check className="w-4 h-4 text-green-600" />
                ) : (
                  <Copy className="w-4 h-4" />
                )}
              </button>
            </div>
            {copied && (
              <p className="text-sm text-green-600 mt-1">
                Copied to clipboard!
              </p>
            )}
          </div>
        )}

        {lastAuthedResult && (
          <button
            onClick={handleReset}
            className="w-full text-gray-600 hover:text-gray-800 text-sm underline mt-2"
          >
            Shorten another URL
          </button>
        )}
      </div>
    </div>
  );
}
