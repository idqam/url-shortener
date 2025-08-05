/* eslint-disable @typescript-eslint/no-explicit-any */
import { useState } from "react";
import { useAuthStore } from "../store/AuthStore";
import { useShortenURL } from "../api/urls";

export function MiniShortenerForm() {
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

  return (
    <form onSubmit={handleSubmit} className="w-full mt-6 space-y-4">
      <div className="flex flex-col sm:flex-row gap-3">
        <input
          type="url"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="Paste a URL"
          required
          className="flex-1 px-4 py-3 rounded-xl bg-white/80 backdrop-blur-sm border border-white/30 text-sm placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-primary"
        />
        <button
          type="submit"
          disabled={isPending}
          className="px-5 py-3 text-sm font-semibold text-white rounded-xl bg-[#A4193D] hover:scale-105 hover:shadow-md transition disabled:opacity-60"
        >
          {isPending ? "Shortening..." : "Shorten"}
        </button>
      </div>

      {result && (
        <div className="text-sm bg-white/80 backdrop-blur rounded-xl p-3 text-center">
          <a
            href={result}
            target="_blank"
            className="text-blue-600 font-medium hover:underline"
          >
            {result}
          </a>
        </div>
      )}

      {error && (
        <div className="text-sm text-red-600 bg-red-100 rounded-xl p-2 text-center">
          {error}
        </div>
      )}
    </form>
  );
}
