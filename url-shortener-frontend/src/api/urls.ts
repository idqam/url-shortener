import type {
  UrlRequestAnon,
  UrlRequestAuthed,
  ShortenResponse,
  GetUrlsResponse,
} from "../types/urls";

const BASE_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

type APIError = {
  error: string;
  code?: string;
};

async function parseAPIError(res: Response): Promise<Error> {
  try {
    const data: APIError = await res.json();
    return new Error(data.error || "Unexpected error");
  } catch {
    const fallback = await res.text().catch(() => "Unknown error");
    return new Error(fallback || "Unknown error");
  }
}

export async function shortenAnon(
  req: UrlRequestAnon
): Promise<ShortenResponse> {
  const res = await fetch(`${BASE_URL}/api/urls`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });

  if (!res.ok) throw await parseAPIError(res);
  return res.json();
}

export async function shortenAuthed(
  req: UrlRequestAuthed
): Promise<ShortenResponse> {
  const res = await fetch(`${BASE_URL}/api/urls`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });

  if (!res.ok) throw await parseAPIError(res);
  return res.json();
}

export async function fetchUserUrls(userId: string): Promise<GetUrlsResponse> {
  const res = await fetch(`${BASE_URL}/api/urls?user_id=${userId}`, {
    method: "GET",
    headers: { "Content-Type": "application/json" },
  });

  if (!res.ok) throw await parseAPIError(res);
  return res.json();
}
