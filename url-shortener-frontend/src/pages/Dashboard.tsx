import React, { useState } from "react";
import {
  ResponsiveContainer,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  AreaChart,
  Area,
} from "recharts";
import {
  Link,
  TrendingUp,
  BarChart3,
  Calendar,
  MousePointer,
  Eye,
  Globe,
  Smartphone,
  Monitor,
  Tablet,
} from "lucide-react";
import { useAnalyticsDashboard } from "../api/analytics";
import { StatCard } from "../components/StatsCard";
import { MiniShortenerForm } from "../components/MiniShortenForm";
import { Navbar } from "../components/Navbar";

interface AnalyticsDashboardProps {
  token: string;
}

export const AnalyticsDashboard: React.FC<AnalyticsDashboardProps> = ({ token }) => {
  const { data: dashboard, isLoading, error, refetch } = useAnalyticsDashboard(token);
  const [showShortener, setShowShortener] = useState(false);

  if (isLoading) {
    return (
      <div className="min-h-screen bg-white">
        <Navbar onToggleShortener={() => {}} showShortener={false} />
        <div className="max-w-7xl mx-auto px-6 py-8">
          <div className="animate-pulse">
            <div className="h-6 bg-gray-200 rounded w-32 mb-1.5" />
            <div className="h-4 bg-gray-200 rounded w-52 mb-8" />
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
              {[...Array(4)].map((_, i) => (
                <div key={i} className="bg-white border border-gray-200 rounded-lg p-5">
                  <div className="w-9 h-9 bg-gray-200 rounded-lg mb-3" />
                  <div className="h-3 bg-gray-200 rounded w-20 mb-2" />
                  <div className="h-8 bg-gray-200 rounded w-16" />
                </div>
              ))}
            </div>
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              <div className="lg:col-span-2 bg-gray-100 rounded-lg h-72" />
              <div className="bg-gray-100 rounded-lg h-72" />
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-white">
        <Navbar onToggleShortener={() => {}} showShortener={false} />
        <div className="flex flex-col items-center justify-center py-32">
          <BarChart3 className="w-10 h-10 text-gray-300 mb-4" />
          <h2 className="text-base font-semibold text-gray-900 mb-1">Failed to load dashboard</h2>
          <p className="text-sm text-gray-500 mb-6">{error.message}</p>
          <button
            onClick={() => refetch()}
            className="text-sm font-medium bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition-colors"
          >
            Try again
          </button>
        </div>
      </div>
    );
  }

  if (!dashboard) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <p className="text-sm text-gray-500">No analytics data available.</p>
      </div>
    );
  }

  const {
    overview = null,
    top_urls = [],
    top_referrers = [],
    device_breakdown = [],
    daily_trend = [],
  } = dashboard;

  const formatDeviceType = (deviceType: string | null | undefined): string => {
    if (!deviceType) return "Unknown";
    return deviceType.charAt(0).toUpperCase() + deviceType.slice(1);
  };

  const getDeviceIcon = (deviceType: string | null | undefined) => {
    switch (deviceType?.toLowerCase()) {
      case "mobile":
        return <Smartphone className="w-4 h-4 text-gray-400" />;
      case "desktop":
        return <Monitor className="w-4 h-4 text-gray-400" />;
      case "tablet":
        return <Tablet className="w-4 h-4 text-gray-400" />;
      default:
        return <Globe className="w-4 h-4 text-gray-400" />;
    }
  };

  const chartData = daily_trend.map((day) => ({
    date: new Date(day.date).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
    }),
    clicks: day.clicks,
  }));

  return (
    <div className="min-h-screen bg-white">
      <Navbar
        onToggleShortener={() => setShowShortener((prev) => !prev)}
        showShortener={showShortener}
      />

      <main className="max-w-7xl mx-auto px-6 py-8">
        {/* Page header */}
        <div className="mb-8">
          <h1 className="text-xl font-semibold text-gray-900">Analytics</h1>
          <p className="text-sm text-gray-500 mt-0.5">Your link performance at a glance</p>
        </div>

        {/* Quick shortener */}
        {showShortener && (
          <div className="mb-8 max-w-xl">
            <MiniShortenerForm />
          </div>
        )}

        {/* Stat cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          <StatCard
            title="Total URLs"
            value={overview?.total_urls.toLocaleString() ?? 0}
            icon={<Link className="w-4 h-4 text-blue-600" />}
          />
          <StatCard
            title="Total Clicks"
            value={overview?.total_clicks.toLocaleString() ?? 0}
            icon={<MousePointer className="w-4 h-4 text-blue-600" />}
          />
          <StatCard
            title="Clicks Today"
            value={overview?.clicks_today.toLocaleString() ?? 0}
            icon={<Calendar className="w-4 h-4 text-blue-600" />}
          />
          <StatCard
            title="Avg Clicks / URL"
            value={overview?.average_clicks.toFixed(1) ?? 0}
            icon={<BarChart3 className="w-4 h-4 text-blue-600" />}
          />
        </div>

        {/* Charts row */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-6">
          {/* Click trend */}
          <div className="lg:col-span-2 bg-white border border-gray-200 rounded-lg p-6">
            <div className="flex items-center justify-between mb-6">
              <div>
                <h2 className="text-sm font-semibold text-gray-900">Click Trends</h2>
                <p className="text-xs text-gray-500 mt-0.5">Last 7 days</p>
              </div>
              <TrendingUp className="w-4 h-4 text-gray-400" />
            </div>

            {chartData.length > 0 ? (
              <div className="h-64">
                <ResponsiveContainer width="100%" height="100%">
                  <AreaChart data={chartData}>
                    <defs>
                      <linearGradient id="clickGradient" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="5%" stopColor="#2563EB" stopOpacity={0.15} />
                        <stop offset="95%" stopColor="#2563EB" stopOpacity={0} />
                      </linearGradient>
                    </defs>
                    <CartesianGrid strokeDasharray="3 3" stroke="#F3F4F6" />
                    <XAxis dataKey="date" stroke="#9CA3AF" fontSize={11} tickLine={false} />
                    <YAxis stroke="#9CA3AF" fontSize={11} tickLine={false} axisLine={false} />
                    <Tooltip
                      formatter={(value: number) => [value.toLocaleString(), "Clicks"]}
                      contentStyle={{
                        backgroundColor: "#fff",
                        border: "1px solid #E5E7EB",
                        borderRadius: "6px",
                        boxShadow: "0 1px 3px rgba(0,0,0,0.1)",
                        fontSize: "12px",
                      }}
                    />
                    <Area
                      type="monotone"
                      dataKey="clicks"
                      stroke="#2563EB"
                      strokeWidth={2}
                      fill="url(#clickGradient)"
                    />
                  </AreaChart>
                </ResponsiveContainer>
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center h-64">
                <BarChart3 className="w-8 h-8 text-gray-200 mb-2" />
                <p className="text-sm text-gray-400">No click data yet</p>
              </div>
            )}
          </div>

          {/* Devices */}
          <div className="bg-white border border-gray-200 rounded-lg p-6">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-sm font-semibold text-gray-900">Devices</h2>
              <Monitor className="w-4 h-4 text-gray-400" />
            </div>

            {device_breakdown.length > 0 ? (
              <div className="space-y-4">
                {device_breakdown.map((device, index) => (
                  <div key={device.device_type || index} className="flex items-center gap-3">
                    <div className="flex items-center gap-2 w-24 flex-shrink-0">
                      {getDeviceIcon(device.device_type)}
                      <span className="text-sm text-gray-600">
                        {formatDeviceType(device.device_type)}
                      </span>
                    </div>
                    <div className="flex-1 h-1.5 bg-gray-100 rounded-full overflow-hidden">
                      <div
                        className="h-full bg-blue-500 rounded-full"
                        style={{ width: `${device.percentage}%` }}
                      />
                    </div>
                    <span className="text-sm font-medium text-gray-900 w-9 text-right flex-shrink-0">
                      {device.percentage.toFixed(0)}%
                    </span>
                  </div>
                ))}
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center h-48">
                <Smartphone className="w-8 h-8 text-gray-200 mb-2" />
                <p className="text-sm text-gray-400">No device data yet</p>
              </div>
            )}
          </div>
        </div>

        {/* Lists row */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Top URLs */}
          <div className="bg-white border border-gray-200 rounded-lg p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-sm font-semibold text-gray-900">Top Performing URLs</h2>
              <Eye className="w-4 h-4 text-gray-400" />
            </div>

            {top_urls && top_urls.length > 0 ? (
              <div className="divide-y divide-gray-100">
                {top_urls.slice(0, 5).map((url, index) => (
                  <div
                    key={url.url_id}
                    className="flex items-center gap-3 py-3 -mx-6 px-6 hover:bg-gray-50 transition-colors"
                  >
                    <span className="text-xs text-gray-400 font-mono w-4 flex-shrink-0">
                      {index + 1}
                    </span>
                    <div className="min-w-0 flex-1">
                      <p className="text-sm font-medium text-blue-600">/{url.short_code}</p>
                      <p className="text-xs text-gray-400 truncate">{url.original_url}</p>
                    </div>
                    <span className="text-sm font-semibold text-gray-900 flex-shrink-0">
                      {url.click_count.toLocaleString()}
                    </span>
                  </div>
                ))}
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center py-10">
                <Link className="w-8 h-8 text-gray-200 mb-2" />
                <p className="text-sm text-gray-400">No URLs yet</p>
              </div>
            )}
          </div>

          {/* Traffic sources */}
          <div className="bg-white border border-gray-200 rounded-lg p-6">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-sm font-semibold text-gray-900">Traffic Sources</h2>
              <Globe className="w-4 h-4 text-gray-400" />
            </div>

            {top_referrers && top_referrers.length > 0 ? (
              <div className="space-y-4">
                {top_referrers.map((referrer, index) => {
                  const maxClicks = Math.max(...top_referrers.map((r) => r.clicks));
                  const pct = maxClicks > 0 ? (referrer.clicks / maxClicks) * 100 : 0;
                  return (
                    <div key={index} className="flex items-center gap-3">
                      <span className="text-sm text-gray-700 w-28 truncate flex-shrink-0">
                        {referrer.referrer === "direct" ? "Direct" : referrer.referrer}
                      </span>
                      <div className="flex-1 h-1.5 bg-gray-100 rounded-full overflow-hidden">
                        <div
                          className="h-full bg-blue-500 rounded-full"
                          style={{ width: `${pct}%` }}
                        />
                      </div>
                      <span className="text-sm font-medium text-gray-900 w-10 text-right flex-shrink-0">
                        {referrer.clicks.toLocaleString()}
                      </span>
                    </div>
                  );
                })}
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center py-10">
                <Globe className="w-8 h-8 text-gray-200 mb-2" />
                <p className="text-sm text-gray-400">No referrer data yet</p>
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  );
};

export default AnalyticsDashboard;
