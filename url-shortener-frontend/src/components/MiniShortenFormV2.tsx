//AI generated to test ai capabilities 
/* eslint-disable @typescript-eslint/no-explicit-any */
import { useState } from "react";
import { Link, Copy, Check, ExternalLink, Zap } from "lucide-react";
import { useAuthStore } from "../store/AuthStore";
import { useShortenURL } from "../api/urls";
import { isValidUrl } from "../utils/isValidUrl";

// Using the same color scheme as the dashboard
const analyticsColors = {
  primary: "#A4193D",
  secondary: "#FFDFB9",
  accent: "#F4A261",
  neutral: "#2A9D8F",
  background: "#F2EFDE",
  cardBg: "rgba(255, 255, 255, 0.8)",
  cardBgHover: "rgba(255, 255, 255, 0.95)",
  text: "#111827",
  textLight: "#6B7280",
  textXLight: "#9CA3AF",
  success: "#10B981",
  danger: "#EF4444",
  warning: "#F59E0B",
  gradient1: "#A4193D",
  gradient2: "#F4A261",
  border: "rgba(164, 25, 61, 0.1)",
  borderHover: "rgba(164, 25, 61, 0.2)",
  surfaceLight: "rgba(255, 255, 255, 0.6)",
};

export function MiniShortenerForm() {
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
      console.error("Error:", err);
      setError(err?.message || "An error occurred.");
      setResult(null);
    }
    setUrl("");
  };

  const handleCopy = async () => {
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
    <div
      className="backdrop-blur-lg rounded-3xl shadow-xl border-2 p-8 hover:shadow-2xl transition-all duration-500 max-w-4xl mx-auto"
      style={{
        backgroundColor: analyticsColors.cardBg,
        borderColor: analyticsColors.border,
      }}
    >
      {/* Header */}
      <div className="flex items-center justify-center space-x-4 mb-6">
        <div
          className="p-3 rounded-2xl"
          style={{
            background: `linear-gradient(135deg, ${analyticsColors.gradient1}, ${analyticsColors.gradient2})`,
          }}
        >
          <Zap className="h-6 w-6 text-white" />
        </div>
        <div className="text-center">
          <h3
            className="text-2xl font-black"
            style={{ color: analyticsColors.text }}
          >
            Quick URL Shortener
          </h3>
          <p
            className="text-sm font-medium tracking-wide uppercase"
            style={{ color: analyticsColors.textLight }}
          >
            Create short links instantly
          </p>
        </div>
      </div>

      {/* Form */}
      <div className="space-y-6">
        <div className="flex flex-col lg:flex-row gap-4">
          <div className="flex-1 relative">
            <div className="absolute left-4 top-1/2 transform -translate-y-1/2">
              <Link
                className="h-5 w-5"
                style={{ color: analyticsColors.textLight }}
              />
            </div>
            <input
              type="url"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter") {
                  e.preventDefault();
                  handleSubmit(e);
                }
              }}
              placeholder="Paste your long URL here..."
              required
              className="w-full pl-12 pr-4 py-4 rounded-2xl border-2 text-sm font-medium placeholder-gray-400 focus:outline-none focus:ring-0 transition-all duration-300"
              style={{
                backgroundColor: analyticsColors.surfaceLight,
                borderColor: analyticsColors.border,
                color: analyticsColors.text,
              }}
              onFocus={(e) => {
                e.target.style.borderColor = analyticsColors.borderHover;
                e.target.style.backgroundColor = analyticsColors.cardBgHover;
              }}
              onBlur={(e) => {
                e.target.style.borderColor = analyticsColors.border;
                e.target.style.backgroundColor = analyticsColors.surfaceLight;
              }}
            />
          </div>
          <button
            type="button"
            onClick={handleSubmit}
            disabled={isPending}
            className="px-8 py-4 text-sm font-bold text-white rounded-2xl transition-all duration-300 hover:scale-105 hover:shadow-lg hover:-translate-y-0.5 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100 disabled:hover:translate-y-0"
            style={{
              background: `linear-gradient(135deg, ${analyticsColors.gradient1}, ${analyticsColors.gradient2})`,
            }}
          >
            {isPending ? (
              <span className="flex items-center space-x-2">
                <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                <span>Shortening...</span>
              </span>
            ) : (
              "Shorten URL"
            )}
          </button>
        </div>

        {/* Success Result */}
        {result && (
          <div
            className="p-6 rounded-2xl border-2 transition-all duration-300"
            style={{
              backgroundColor: analyticsColors.surfaceLight,
              borderColor: analyticsColors.success,
            }}
          >
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3 min-w-0 flex-1">
                <div
                  className="p-2 rounded-xl"
                  style={{ backgroundColor: analyticsColors.success }}
                >
                  <Check className="h-4 w-4 text-white" />
                </div>
                <div className="min-w-0 flex-1">
                  <p
                    className="text-xs font-medium uppercase tracking-wide mb-1"
                    style={{ color: analyticsColors.textLight }}
                  >
                    Your short URL
                  </p>
                  <a
                    href={result}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-lg font-bold hover:underline flex items-center space-x-2 group"
                    style={{ color: analyticsColors.primary }}
                  >
                    <span className="truncate">{result}</span>
                    <ExternalLink className="h-4 w-4 group-hover:scale-110 transition-transform" />
                  </a>
                </div>
              </div>
              <button
                onClick={handleCopy}
                className="ml-4 p-3 rounded-xl transition-all duration-300 hover:scale-110 border-2"
                style={{
                  backgroundColor: copied
                    ? analyticsColors.success
                    : analyticsColors.cardBgHover,
                  borderColor: copied
                    ? analyticsColors.success
                    : analyticsColors.border,
                }}
              >
                {copied ? (
                  <Check className="h-5 w-5 text-white" />
                ) : (
                  <Copy
                    className="h-5 w-5"
                    style={{ color: analyticsColors.primary }}
                  />
                )}
              </button>
            </div>
          </div>
        )}

        {/* Error State */}
        {error && (
          <div
            className="p-4 rounded-2xl border-2 text-center"
            style={{
              backgroundColor: "rgba(239, 68, 68, 0.1)",
              borderColor: analyticsColors.danger,
            }}
          >
            <p
              className="text-sm font-semibold"
              style={{ color: analyticsColors.danger }}
            >
              {error}
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
