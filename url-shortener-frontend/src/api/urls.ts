import type { ErrorResponse } from "react-router-dom";
import type {
  ShortenURLRequest,
  ShortenURLResponse,
  GetUserURLsResponse,
  GetURLByShortCodeResponse,
} from "../dtos/requests";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { API_BASE } from "../constants/apiBase";

async function fetcher<T>(url: string, options?: RequestInit): Promise<T> {
  const res = await fetch(url, options);
  const contentType = res.headers.get("Content-Type") || "";

  if (!contentType.includes("application/json")) {
    const text = await res.text();
    throw new Error(`Expected JSON, got: ${text}`);
  }

  const data = await res.json();

  if (!res.ok) {
    throw data as ErrorResponse;
  }

  return data as T;
}

export async function shortenURL(
  payload: ShortenURLRequest
): Promise<ShortenURLResponse> {
  return fetcher(`${API_BASE}/api/urls`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
}

export async function getUserURLs(token: string): Promise<GetUserURLsResponse> {
  return fetcher(`${API_BASE}/api/urls`, {
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  });
}

export async function getURLByShortCode(
  shortcode: string
): Promise<GetURLByShortCodeResponse> {
  const encoded = encodeURIComponent(shortcode);
  return fetcher(`${API_BASE}/api/urls/${encoded}`);
}

export function useShortenURL() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: shortenURL,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user-urls"] });
    },
  });
}

export function useUserURLs(token: string) {
  return useQuery({
    queryKey: ["userUrls"],
    queryFn: () => getUserURLs(token),
    enabled: !!token,
  });
}

export function useURLByShortCode(shortcode: string) {
  return useQuery({
    queryKey: ["shortcode", shortcode],
    queryFn: () => getURLByShortCode(shortcode),
    enabled: !!shortcode,
  });
}
