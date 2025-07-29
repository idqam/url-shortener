import { useUrlStore } from "../store/UrlStore";
import type { ShortenResponse } from "../types/urls";

export function AnonResult() {
  const result = useUrlStore((s) => s.anonLastResult);
  const error = useUrlStore((s) => s.error);

  if (error) {
    return (
      <div className="shortener-box shortener-result">
        <div className="text-red-500">{error}</div>
      </div>
    );
  }

  if (result) {
    const res = result as ShortenResponse;
    return (
      <div className="shortener-box shortener-result">
        <p className="">Shortened URL:</p>
        <a
          href={res.short_url}
          target="_blank"
          rel="noopener noreferrer"
          className="text-blue-600 underline break-words"
        >
          {res.short_url}
        </a>
      </div>
    );
  }

  return (
    <div className="shortener-box shortener-result flex">
      <div className="flex-1 flex items-center justify-center text-center w-full">
        <p className="text-gray-400 font-bold">No URL generated yet</p>
      </div>
    </div>
  );
}

// <a
//         href={res.short_url}
//         target="_blank"
//         rel="noopener noreferrer"
//         className="text-blue-600 underline break-words"
//       >
//         {res.short_url}
//       </a>
