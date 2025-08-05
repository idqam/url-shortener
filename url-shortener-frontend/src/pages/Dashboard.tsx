import React from "react";
import {
  ResponsiveContainer,
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  PieChart,
  Pie,
  Cell,
} from "recharts";
import {
  Link,
  TrendingUp,
  BarChart3,
  Calendar,
  MousePointer,
  Eye,
  Users,
  Smartphone,
  Monitor,
  Tablet,
  Globe,
} from "lucide-react";
import { useAnalyticsDashboard } from "../api/analytics";
import { colors, pieColors } from "../constants/colors";

interface AnalyticsDashboardProps {
  token: string;
}

export const AnalyticsDashboard: React.FC<AnalyticsDashboardProps> = ({
  token,
}) => {
  const {
    data: dashboard,
    isLoading,
    error,
    refetch,
  } = useAnalyticsDashboard(token);

  if (isLoading) {
    return (
      <div
        className="min-h-screen"
        style={{ backgroundColor: colors.background }}
      >
        <div className="max-w-7xl mx-auto px-6 py-12">
          <div className="animate-pulse">
            <div className="text-center mb-12">
              <div
                className="h-12 rounded-lg w-80 mx-auto mb-4 animate-pulse"
                style={{
                  background: `linear-gradient(90deg, ${colors.secondary} 25%, #f3f4f6 50%, ${colors.secondary} 75%)`,
                }}
              ></div>
              <div
                className="h-5 rounded w-96 mx-auto animate-pulse"
                style={{
                  background: `linear-gradient(90deg, ${colors.secondary} 25%, #f3f4f6 50%, ${colors.secondary} 75%)`,
                }}
              ></div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
              {[...Array(4)].map((_, i) => (
                <div
                  key={i}
                  className="bg-white p-6 rounded-xl shadow-lg border border-gray-100"
                >
                  <div className="flex items-center">
                    <div
                      className="w-14 h-14 rounded-xl animate-pulse"
                      style={{ backgroundColor: colors.secondary }}
                    ></div>
                    <div className="ml-4 flex-1">
                      <div
                        className="h-4 rounded w-20 mb-3 animate-pulse"
                        style={{ backgroundColor: colors.secondary }}
                      ></div>
                      <div
                        className="h-8 rounded w-16 animate-pulse"
                        style={{ backgroundColor: colors.secondary }}
                      ></div>
                    </div>
                  </div>
                </div>
              ))}
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              <div className="bg-white p-8 rounded-xl shadow-lg border border-gray-100">
                <div
                  className="h-7 rounded w-48 mb-6 animate-pulse"
                  style={{ backgroundColor: colors.secondary }}
                ></div>
                <div
                  className="h-80 rounded-lg animate-pulse"
                  style={{ backgroundColor: colors.secondary }}
                ></div>
              </div>
              <div className="bg-white p-8 rounded-xl shadow-lg border border-gray-100">
                <div
                  className="h-7 rounded w-36 mb-6 animate-pulse"
                  style={{ backgroundColor: colors.secondary }}
                ></div>
                <div
                  className="h-80 rounded-lg animate-pulse"
                  style={{ backgroundColor: colors.secondary }}
                ></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div
        className="min-h-screen"
        style={{ backgroundColor: colors.background }}
      >
        <div className="max-w-7xl mx-auto px-6 py-12">
          <div className="text-center">
            <div className="bg-red-50 border border-red-200 rounded-xl p-8 max-w-md mx-auto shadow-lg">
              <h2 className="text-red-800 font-bold text-xl mb-3">
                Error Loading Dashboard
              </h2>
              <p className="text-red-600 mb-6">{error.message}</p>
              <button
                onClick={() => refetch()}
                className="px-6 py-3 rounded-lg font-medium transition-all duration-200 transform hover:scale-105 text-white"
                style={{ backgroundColor: colors.primary }}
              >
                Try Again
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!dashboard) {
    return <div>No analytics data available</div>;
  }

  const { overview, top_urls, top_referrers, device_breakdown, daily_trend } =
    dashboard;

  const formatDeviceType = (deviceType: string | null | undefined): string => {
    if (!deviceType) return "Unknown";
    return deviceType.charAt(0).toUpperCase() + deviceType.slice(1);
  };

  const getDeviceIcon = (deviceType: string | null | undefined) => {
    if (!deviceType)
      return <Globe className="h-5 w-5" style={{ color: colors.primary }} />;

    switch (deviceType.toLowerCase()) {
      case "mobile":
        return (
          <Smartphone className="h-5 w-5" style={{ color: colors.primary }} />
        );
      case "desktop":
        return (
          <Monitor className="h-5 w-5" style={{ color: colors.primary }} />
        );
      case "tablet":
        return <Tablet className="h-5 w-5" style={{ color: colors.primary }} />;
      default:
        return <Globe className="h-5 w-5" style={{ color: colors.primary }} />;
    }
  };

  const getTrendIcon = () => {
    if (overview.trend_direction === "up")
      return <TrendingUp className="h-5 w-5 text-green-500" />;
    if (overview.trend_direction === "down")
      return <TrendingUp className="h-5 w-5 text-red-500 rotate-180" />;
    return <BarChart3 className="h-5 w-5 text-gray-500" />;
  };

  const chartData = daily_trend.map((day) => ({
    date: new Date(day.date).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
    }),
    clicks: day.clicks,
  }));

  const deviceChartData = device_breakdown.map((device) => ({
    name: formatDeviceType(device.device_type),
    value: device.clicks,
    percentage: device.percentage,
  }));

  return (
    <div
      className="min-h-screen"
      style={{ backgroundColor: colors.background }}
    >
      <div className="max-w-7xl mx-auto px-6 py-12">
        <div className="text-center mb-12">
          <h1
            className="text-5xl font-bold mb-4"
            style={{
              background: `linear-gradient(135deg, ${colors.primary} 0%, ${colors.accent} 100%)`,
              WebkitBackgroundClip: "text",
              WebkitTextFillColor: "transparent",
              backgroundClip: "text",
            }}
          >
            Analytics Dashboard
          </h1>
          <p
            className="text-xl max-w-2xl mx-auto"
            style={{ color: colors.textLight }}
          >
            Get comprehensive insights into your URL performance with 7 day
            analytics and metrics
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
          <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-100 transition-all duration-300 hover:shadow-xl hover:-translate-y-1">
            <div className="flex items-center">
              <div
                className="p-3 rounded-xl"
                style={{ backgroundColor: `${colors.secondary}40` }}
              >
                <Link className="h-7 w-7" style={{ color: colors.primary }} />
              </div>
              <div className="ml-4">
                <p
                  className="text-sm font-medium"
                  style={{ color: colors.textLight }}
                >
                  Total URLs
                </p>
                <p
                  className="text-3xl font-bold"
                  style={{ color: colors.text }}
                >
                  {overview.total_urls.toLocaleString()}
                </p>
              </div>
            </div>
          </div>

          <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-100 transition-all duration-300 hover:shadow-xl hover:-translate-y-1">
            <div className="flex items-center">
              <div
                className="p-3 rounded-xl"
                style={{ backgroundColor: `${colors.accent}40` }}
              >
                <MousePointer
                  className="h-7 w-7"
                  style={{ color: colors.accent }}
                />
              </div>
              <div className="ml-4">
                <p
                  className="text-sm font-medium"
                  style={{ color: colors.textLight }}
                >
                  Total Clicks
                </p>
                <p
                  className="text-3xl font-bold"
                  style={{ color: colors.text }}
                >
                  {overview.total_clicks.toLocaleString()}
                </p>
              </div>
            </div>
          </div>

          <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-100 transition-all duration-300 hover:shadow-xl hover:-translate-y-1">
            <div className="flex items-center">
              <div
                className="p-3 rounded-xl"
                style={{ backgroundColor: `${colors.neutral}40` }}
              >
                <Calendar
                  className="h-7 w-7"
                  style={{ color: colors.neutral }}
                />
              </div>
              <div className="ml-4">
                <div className="flex items-center space-x-2">
                  <p
                    className="text-sm font-medium"
                    style={{ color: colors.textLight }}
                  >
                    Clicks Today
                  </p>
                  {getTrendIcon()}
                </div>
                <div className="flex items-baseline space-x-2">
                  <p
                    className="text-3xl font-bold"
                    style={{ color: colors.text }}
                  >
                    {overview.clicks_today.toLocaleString()}
                  </p>
                  <p className="text-sm" style={{ color: colors.textLight }}>
                    vs {overview.clicks_yesterday} yesterday
                  </p>
                </div>
              </div>
            </div>
          </div>

          <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-100 transition-all duration-300 hover:shadow-xl hover:-translate-y-1">
            <div className="flex items-center">
              <div
                className="p-3 rounded-xl"
                style={{ backgroundColor: `${colors.warning}40` }}
              >
                <BarChart3
                  className="h-7 w-7"
                  style={{ color: colors.warning }}
                />
              </div>
              <div className="ml-4">
                <p
                  className="text-sm font-medium"
                  style={{ color: colors.textLight }}
                >
                  Avg Clicks/URL
                </p>
                <p
                  className="text-3xl font-bold"
                  style={{ color: colors.text }}
                >
                  {overview.average_clicks.toFixed(1)}
                </p>
              </div>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-8">
          <div className="lg:col-span-2 bg-white p-8 rounded-xl shadow-lg border border-gray-100">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-2xl font-bold" style={{ color: colors.text }}>
                Top URLs by Clicks
              </h2>
              <Eye className="h-6 w-6" style={{ color: colors.primary }} />
            </div>

            {top_urls && top_urls.length > 0 ? (
              <div className="space-y-4">
                {top_urls.slice(0, 5).map((url, index) => (
                  <div
                    key={url.url_id}
                    className="flex items-center justify-between p-4 rounded-xl border-2 border-gray-100 shadow-md transition-all duration-300 hover:shadow-lg hover:-translate-y-1 hover:border-gray-200 cursor-pointer"
                    style={{ backgroundColor: colors.background }}
                  >
                    <div className="flex items-center space-x-4">
                      <div
                        className="flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center font-bold text-white shadow-md"
                        style={{
                          background: `linear-gradient(135deg, ${colors.primary} 0%, ${colors.accent} 100%)`,
                        }}
                      >
                        <span className="text-lg">{index + 1}</span>
                      </div>
                      <div className="min-w-0">
                        <p
                          className="text-lg font-bold truncate"
                          style={{ color: colors.text }}
                        >
                          /{url.short_code}
                        </p>
                        <p
                          className="text-sm truncate max-w-md"
                          style={{ color: colors.textLight }}
                          title={url.original_url}
                        >
                          {url.original_url}
                        </p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p
                        className="text-2xl font-bold"
                        style={{ color: colors.primary }}
                      >
                        {url.click_count.toLocaleString()}
                      </p>
                      <p
                        className="text-sm font-medium"
                        style={{ color: colors.textLight }}
                      >
                        clicks
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-12">
                <Link
                  className="h-16 w-16 mx-auto mb-4"
                  style={{ color: colors.textLight }}
                />
                <p className="text-lg" style={{ color: colors.textLight }}>
                  No URLs found
                </p>
              </div>
            )}
          </div>

          <div className="bg-white p-8 rounded-xl shadow-lg border border-gray-100">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-2xl font-bold" style={{ color: colors.text }}>
                Device Types
              </h2>
              <Users className="h-6 w-6" style={{ color: colors.primary }} />
            </div>

            {deviceChartData.length > 0 ? (
              <div className="h-80">
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={deviceChartData}
                      cx="50%"
                      cy="50%"
                      outerRadius={80}
                      dataKey="value"
                      label={({ name, percentage }) =>
                        `${name}: ${percentage.toFixed(1)}%`
                      }
                    >
                      {deviceChartData.map((_, index) => (
                        <Cell
                          key={`cell-${index}`}
                          fill={pieColors[index % pieColors.length]}
                        />
                      ))}
                    </Pie>
                    <Tooltip
                      formatter={(value: number) => [
                        value.toLocaleString(),
                        "Clicks",
                      ]}
                      contentStyle={{
                        backgroundColor: "#ffffff",
                        border: "2px solid #e5e7eb",
                        borderRadius: "12px",
                        boxShadow: "0 10px 25px -5px rgba(0, 0, 0, 0.1)",
                      }}
                    />
                  </PieChart>
                </ResponsiveContainer>

                <div className="mt-4 space-y-2">
                  {device_breakdown.map((device, index) => (
                    <div
                      key={device.device_type || index}
                      className="flex items-center justify-between"
                    >
                      <div className="flex items-center space-x-2">
                        {getDeviceIcon(device.device_type)}
                        <span
                          className="font-medium"
                          style={{ color: colors.text }}
                        >
                          {formatDeviceType(device.device_type)}
                        </span>
                      </div>
                      <span
                        className="font-bold"
                        style={{ color: colors.text }}
                      >
                        {device.clicks}
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            ) : (
              <div className="text-center py-16">
                <Smartphone
                  className="h-16 w-16 mx-auto mb-4"
                  style={{ color: colors.textLight }}
                />
                <p className="text-lg" style={{ color: colors.textLight }}>
                  No device data available
                </p>
              </div>
            )}
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <div className="bg-white p-8 rounded-xl shadow-lg border border-gray-100">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-2xl font-bold" style={{ color: colors.text }}>
                7-Day Click Trend
              </h2>
              <TrendingUp
                className="h-6 w-6"
                style={{ color: colors.primary }}
              />
            </div>

            {chartData.length > 0 ? (
              <div className="h-80">
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart data={chartData}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#f1f5f9" />
                    <XAxis
                      dataKey="date"
                      stroke={colors.textLight}
                      fontSize={12}
                      fontWeight="500"
                    />
                    <YAxis
                      stroke={colors.textLight}
                      fontSize={12}
                      fontWeight="500"
                    />
                    <Tooltip
                      formatter={(value: number) => [
                        value.toLocaleString(),
                        "Clicks",
                      ]}
                      contentStyle={{
                        backgroundColor: "#ffffff",
                        border: "2px solid #e5e7eb",
                        borderRadius: "12px",
                        boxShadow: "0 10px 25px -5px rgba(0, 0, 0, 0.1)",
                      }}
                    />
                    <Line
                      type="monotone"
                      dataKey="clicks"
                      stroke={colors.primary}
                      strokeWidth={3}
                      dot={{ fill: colors.primary, strokeWidth: 3, r: 5 }}
                      activeDot={{
                        r: 8,
                        stroke: colors.primary,
                        strokeWidth: 3,
                        fill: "#ffffff",
                      }}
                    />
                  </LineChart>
                </ResponsiveContainer>
              </div>
            ) : (
              <div className="text-center py-16">
                <BarChart3
                  className="h-16 w-16 mx-auto mb-4"
                  style={{ color: colors.textLight }}
                />
                <p className="text-lg" style={{ color: colors.textLight }}>
                  No click data available
                </p>
              </div>
            )}
          </div>

          <div className="bg-white p-8 rounded-xl shadow-lg border border-gray-100">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-2xl font-bold" style={{ color: colors.text }}>
                Top Referrers
              </h2>
              <Globe className="h-6 w-6" style={{ color: colors.primary }} />
            </div>

            {top_referrers && top_referrers.length > 0 ? (
              <div className="space-y-4">
                {top_referrers.map((referrer, index) => (
                  <div
                    key={index}
                    className="flex items-center justify-between p-4 rounded-lg border border-gray-100 hover:bg-gray-50 transition-colors"
                  >
                    <div className="flex items-center space-x-3">
                      <div
                        className="w-8 h-8 rounded-full flex items-center justify-center text-white font-bold text-sm"
                        style={{ backgroundColor: colors.accent }}
                      >
                        {index + 1}
                      </div>
                      <span
                        className="font-medium truncate max-w-xs"
                        style={{ color: colors.text }}
                        title={referrer.referrer}
                      >
                        {referrer.referrer === "direct"
                          ? "Direct Traffic"
                          : referrer.referrer}
                      </span>
                    </div>
                    <span
                      className="font-bold"
                      style={{ color: colors.primary }}
                    >
                      {referrer.clicks.toLocaleString()}
                    </span>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-16">
                <Globe
                  className="h-16 w-16 mx-auto mb-4"
                  style={{ color: colors.textLight }}
                />
                <p className="text-lg" style={{ color: colors.textLight }}>
                  No referrer data available
                </p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default AnalyticsDashboard;
