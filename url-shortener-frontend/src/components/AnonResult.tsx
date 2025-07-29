import { useUrlStore } from "../store/UrlStore";

export function AnonResult() {
  const result = useUrlStore((s) => s.anonLastResult);
  const error = useUrlStore((s) => s.error);

  if (error) {
    return <div className="text-red-500 shortener-box">{error}</div>;
  }

  if (!result) return null;

  return (
    <div className="shortener-box shortener-result">
      <p className="text-lg font-medium mb-2">Shortened URL:</p>
      <a
        href={result.short_url}
        target="_blank"
        rel="noopener noreferrer"
        className="text-blue-600 underline break-words"
      >
        {result.short_url}
      </a>
    </div>
  );
}
