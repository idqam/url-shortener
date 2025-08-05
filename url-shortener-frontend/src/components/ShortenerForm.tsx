/* eslint-disable @typescript-eslint/no-explicit-any */
import { useState } from "react";
import { useAuthStore } from "../store/AuthStore";
import { useShortenURL } from "../api/urls";

export function ShortenerForm() {
  const [url, setUrl] = useState("");
  const [result, setResult] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const userId = useAuthStore((state) => state.userId);
  const { mutateAsync: shorten, isPending } = useShortenURL();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!url) return;

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
      console.error("Error:", err);
      setError(err?.message || "An error occurred.");
      setResult(null);
    }
    setUrl("");
  };

  const copyToClipboard = async () => {
    if (result) {
      try {
        await navigator.clipboard.writeText(result);
      } catch (err) {
        console.error("Failed to copy:", err);
      }
    }
  };

  return (
    <div className="space-y-6">
      <div className="space-y-4">
        <div className="flex flex-col sm:flex-row gap-4">
          <div className="flex-1 relative">
            <input
              type="url"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              placeholder="Paste your URL"
              required
              className="w-full px-6 py-4 bg-white/80 backdrop-blur-sm border border-gray-200/50 rounded-xl text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-0 transition-all duration-300 hover:bg-white/90 focus:border-transparent"
              style={{
                boxShadow: "focus:0 0 0 2px rgba(164, 25, 61, 0.2)",
              }}
              onFocus={(e) => {
                e.target.style.borderColor = "#A4193D";
                e.target.style.boxShadow = "0 0 0 3px rgba(164, 25, 61, 0.1)";
              }}
              onBlur={(e) => {
                e.target.style.borderColor = "rgba(229, 231, 235, 0.5)";
                e.target.style.boxShadow = "none";
              }}
              onKeyDown={(e) => {
                if (e.key === "Enter") {
                  handleSubmit(e);
                }
              }}
            />
          </div>
          <button
            onClick={handleSubmit}
            disabled={isPending}
            className="px-8 py-4 font-semibold text-white rounded-xl transition-all duration-300 hover:scale-105 hover:shadow-lg disabled:opacity-70 disabled:cursor-not-allowed disabled:hover:scale-100 flex items-center justify-center min-w-[120px]"
            style={{
              background: isPending ? "#6B7280" : "#6366F1",
              backgroundColor: isPending ? "#6B7280" : "#6366F1",
            }}
            onMouseEnter={(e) => {
              if (!isPending) {
                const target = e.target as HTMLButtonElement;
                target.style.backgroundColor = "A4193D";
              }
            }}
            onMouseLeave={(e) => {
              if (!isPending) {
                const target = e.target as HTMLButtonElement;
                target.style.backgroundColor = "#6366F1";
              }
            }}
          >
            {isPending ? (
              <>
                <svg
                  className="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    className="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    strokeWidth="4"
                  ></circle>
                  <path
                    className="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  ></path>
                </svg>
                Shortening...
              </>
            ) : (
              "Shorten"
            )}
          </button>
        </div>
      </div>

      {result && (
        <div className="bg-white/70 backdrop-blur-sm border border-white/30 rounded-xl p-6 text-center space-y-4">
          <div className="flex items-center justify-center mb-3">
            <div
              className="w-8 h-8 rounded-full flex items-center justify-center mr-3"
              style={{ backgroundColor: "#10B981" }}
            >
              <svg
                className="w-4 h-4 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M5 13l4 4L19 7"
                />
              </svg>
            </div>
            <span
              className="text-lg font-semibold"
              style={{ color: "#10B981" }}
            >
              Success!
            </span>
          </div>

          <div className="bg-white/90 backdrop-blur-sm rounded-lg p-4 border border-gray-200/30">
            <div className="flex items-center justify-between gap-3">
              <a
                href={result}
                target="_blank"
                rel="noopener noreferrer"
                className="font-medium hover:underline break-all flex-1 text-left text-blue-600"
              >
                {result}
              </a>
              <button
                onClick={copyToClipboard}
                className="p-2 rounded-lg transition-all duration-200 hover:scale-110 flex-shrink-0 hover:bg-gray-100"
                title="Copy to clipboard"
              >
                <svg
                  className="w-4 h-4 text-gray-600"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
                  />
                </svg>
              </button>
            </div>
          </div>
        </div>
      )}

      {error && (
        <div className="bg-red-50/90 backdrop-blur-sm border border-red-200/50 rounded-xl p-4 text-center">
          <div className="flex items-center justify-center mb-2">
            <div className="w-6 h-6 bg-red-500 rounded-full flex items-center justify-center mr-2">
              <svg
                className="w-3 h-3 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </div>
            <span className="font-medium text-red-600">Error</span>
          </div>
          <p className="text-red-600 text-sm">{error}</p>
        </div>
      )}
    </div>
  );
}
