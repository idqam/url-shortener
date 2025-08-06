import type { ErrorResponse } from "react-router-dom";

import { useQuery } from "@tanstack/react-query";
import type {
  AnalyticsDashboard,
  AnalyticsParams,
  TopURLsResponse,
  TopReferrersResponse,
  DeviceBreakdownResponse,
  DailyTrendResponse,
} from "../dtos/analyticsDto";
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

export async function getAnalyticsDashboard(
  token: string
): Promise<AnalyticsDashboard> {
  return fetcher(`${API_BASE}/api/analytics/dashboard`, {
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  });
}

export async function getTopURLs(
  token: string,
  params?: AnalyticsParams
): Promise<TopURLsResponse> {
  const queryParams = new URLSearchParams();
  if (params?.limit) queryParams.append("limit", params.limit.toString());

  const endpoint = `/api/analytics/urls${
    queryParams.toString() ? `?${queryParams}` : ""
  }`;

  return fetcher(`${API_BASE}${endpoint}`, {
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  });
}

export async function getTopReferrers(
  token: string,
  params?: AnalyticsParams
): Promise<TopReferrersResponse> {
  const queryParams = new URLSearchParams();
  if (params?.limit) queryParams.append("limit", params.limit.toString());

  const endpoint = `/api/analytics/referrers${
    queryParams.toString() ? `?${queryParams}` : ""
  }`;

  return fetcher(`${API_BASE}${endpoint}`, {
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  });
}

export async function getDeviceBreakdown(
  token: string
): Promise<DeviceBreakdownResponse> {
  return fetcher(`${API_BASE}/api/analytics/devices`, {
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  });
}

export async function getDailyTrend(
  token: string,
  params?: AnalyticsParams
): Promise<DailyTrendResponse> {
  const queryParams = new URLSearchParams();
  if (params?.days) queryParams.append("days", params.days.toString());

  const endpoint = `/api/analytics/trend${
    queryParams.toString() ? `?${queryParams}` : ""
  }`;

  return fetcher(`${API_BASE}${endpoint}`, {
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  });
}

export function useAnalyticsDashboard(token: string) {
  return useQuery({
    queryKey: ["analytics-dashboard"],
    queryFn: () => getAnalyticsDashboard(token),
    enabled: !!token,
    staleTime: 5 * 60 * 1000,
    refetchInterval: 10 * 60 * 1000,
  });
}

export function useTopURLs(token: string, params?: AnalyticsParams) {
  return useQuery({
    queryKey: ["analytics-top-urls", params],
    queryFn: () => getTopURLs(token, params),
    enabled: !!token,
    staleTime: 2 * 60 * 1000,
  });
}

export function useTopReferrers(token: string, params?: AnalyticsParams) {
  return useQuery({
    queryKey: ["analytics-top-referrers", params],
    queryFn: () => getTopReferrers(token, params),
    enabled: !!token,
    staleTime: 5 * 60 * 1000,
  });
}

export function useDeviceBreakdown(token: string) {
  return useQuery({
    queryKey: ["analytics-device-breakdown"],
    queryFn: () => getDeviceBreakdown(token),
    enabled: !!token,
    staleTime: 10 * 60 * 1000,
  });
}

export function useDailyTrend(token: string, params?: AnalyticsParams) {
  return useQuery({
    queryKey: ["analytics-daily-trend", params],
    queryFn: () => getDailyTrend(token, params),
    enabled: !!token,
    staleTime: 1 * 60 * 1000,
    refetchInterval: 5 * 60 * 1000,
  });
}

export function useAnalyticsOverview(token: string) {
  const query = useAnalyticsDashboard(token);

  return {
    ...query,
    data: query.data?.overview,
  };
}
