import type { AnalyticsDashboard } from "../dtos/analyticsDto";

export const mockDashboardData: AnalyticsDashboard = {
  overview: {
    total_urls: 42,
    total_clicks: 3817,
    clicks_today: 214,
    clicks_yesterday: 189,
    average_clicks: 90.9,
    trend_direction: "up",
  },
  top_urls: [
    {
      url_id: "dev-1",
      short_code: "gh-repo",
      original_url: "https://github.com/idqam/url-shortener-go-react",
      click_count: 941,
      created_at: "2026-03-01T10:00:00Z",
    },
    {
      url_id: "dev-2",
      short_code: "lnkdn",
      original_url: "https://linkedin.com/in/devuser",
      click_count: 708,
      created_at: "2026-03-05T14:22:00Z",
    },
    {
      url_id: "dev-3",
      short_code: "blog-1",
      original_url: "https://dev.to/devuser/how-i-built-this",
      click_count: 534,
      created_at: "2026-03-10T09:15:00Z",
    },
    {
      url_id: "dev-4",
      short_code: "docs",
      original_url: "https://docs.example.com/api/v2/reference",
      click_count: 312,
      created_at: "2026-03-14T16:40:00Z",
    },
    {
      url_id: "dev-5",
      short_code: "tweet",
      original_url: "https://twitter.com/devuser/status/123456789",
      click_count: 201,
      created_at: "2026-03-20T11:05:00Z",
    },
  ],
  top_referrers: [
    { referrer: "direct", clicks: 1420 },
    { referrer: "twitter.com", clicks: 874 },
    { referrer: "github.com", clicks: 653 },
    { referrer: "reddit.com", clicks: 412 },
    { referrer: "linkedin.com", clicks: 270 },
  ],
  device_breakdown: [
    { device_type: "desktop", clicks: 2140, percentage: 56.1 },
    { device_type: "mobile", clicks: 1298, percentage: 34.0 },
    { device_type: "tablet", clicks: 379, percentage: 9.9 },
  ],
  daily_trend: [
    { date: "2026-03-26T00:00:00Z", clicks: 148 },
    { date: "2026-03-27T00:00:00Z", clicks: 203 },
    { date: "2026-03-28T00:00:00Z", clicks: 175 },
    { date: "2026-03-29T00:00:00Z", clicks: 310 },
    { date: "2026-03-30T00:00:00Z", clicks: 267 },
    { date: "2026-03-31T00:00:00Z", clicks: 189 },
    { date: "2026-04-01T00:00:00Z", clicks: 214 },
  ],
};
