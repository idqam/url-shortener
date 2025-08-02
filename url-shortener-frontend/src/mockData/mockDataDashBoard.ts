

import type { DashboardStats } from "../types/urls";

export const mockData: DashboardStats = {
  totalUrls: 12,
  totalClicks: 347,
  topUrls: [
    {
      id: "1",
      original_url:
        "https://example.com/very-long-article-about-technology-trends",
      short_code: "tech123",
      click_count: 89,
    },
    {
      id: "2",
      original_url: "https://github.com/username/awesome-project",
      short_code: "git456",
      click_count: 67,
    },
    {
      id: "3",
      original_url: "https://docs.company.com/api/documentation",
      short_code: "docs789",
      click_count: 45,
    },
    {
      id: "4",
      original_url: "https://blog.example.com/how-to-build-better-apis",
      short_code: "blog321",
      click_count: 32,
    },
    {
      id: "5",
      original_url: "https://newsletter.startup.com/weekly-update",
      short_code: "news654",
      click_count: 28,
    },
  ],
  weeklyData: [
    { week: "2024-12-30", clicks: 45 },
    { week: "2025-01-06", clicks: 62 },
    { week: "2025-01-13", clicks: 38 },
    { week: "2025-01-20", clicks: 71 },
    { week: "2025-01-27", clicks: 56 },
    { week: "2025-02-03", clicks: 43 },
    { week: "2025-02-10", clicks: 89 },
    { week: "2025-02-17", clicks: 67 },
  ],
};
