export interface AnalyticsOverview {
  total_urls: number;
  total_clicks: number;
  clicks_today: number;
  clicks_yesterday: number;
  average_clicks: number;
  trend_direction: "up" | "down" | "same";
}

export interface TopURL {
  url_id: string;
  short_code: string;
  original_url: string;
  click_count: number;
  created_at: string;
}

export interface ReferrerStat {
  referrer: string;
  clicks: number;
}

export interface DeviceStat {
  device_type: string;
  clicks: number;
  percentage: number;
}

export interface DailyTrend {
  date: string;
  clicks: number;
}

export interface AnalyticsDashboard {
  overview: AnalyticsOverview;
  top_urls: TopURL[];
  top_referrers: ReferrerStat[];
  device_breakdown: DeviceStat[];
  daily_trend: DailyTrend[];
}

export interface TopURLsResponse {
  urls: TopURL[];
}

export interface TopReferrersResponse {
  referrers: ReferrerStat[];
}

export interface DeviceBreakdownResponse {
  devices: DeviceStat[];
}

export interface DailyTrendResponse {
  trend: DailyTrend[];
  days: number;
}

export interface AnalyticsParams {
  limit?: number;
  days?: number;
}

export interface ApiError {
  error: string;
  code?: string;
  field?: string;
}
