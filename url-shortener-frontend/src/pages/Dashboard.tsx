import React, { useState } from "react";
import {
  ResponsiveContainer,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  PieChart,
  Pie,
  Cell,
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
  Users,
  Smartphone,
  Monitor,
  Tablet,
  Globe,
  ArrowUp,
  ArrowDown,
  Activity,
  Zap,
} from "lucide-react";
import { useAnalyticsDashboard } from "../api/analytics";
import { StatCard } from "../components/StatsCard";
import { MiniShortenerForm } from "../components/MiniShortenForm";
import SignOutButton from "../components/SignOutButton";
import { analyticsColors } from "../constants/colors";

const pieColors = [
  analyticsColors.primary,
  analyticsColors.accent,
  analyticsColors.neutral,
  analyticsColors.info,
];

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
  const [showShortener, setShowShortener] = useState(false);

  if (isLoading) {
    return (
      <div
        className="min-h-screen relative overflow-hidden"
        style={{ backgroundColor: analyticsColors.background }}
      >
        <div className="">
          <div
            className="absolute top-20 left-10 w-72 h-72 rounded-full mix-blend-multiply filter blur-xl animate-pulse"
            style={{
              background: `linear-gradient(135deg, ${analyticsColors.gradient1}, ${analyticsColors.gradient2})`,
            }}
          ></div>
          <div
            className="absolute top-40 right-10 w-96 h-96 rounded-full mix-blend-multiply filter blur-xl animate-pulse delay-700"
            style={{
              background: `linear-gradient(135deg, ${analyticsColors.gradient3}, ${analyticsColors.secondary})`,
            }}
          ></div>
        </div>

        <div className="relative z-10 max-w-7xl mx-auto px-6 py-12">
          <div className="animate-pulse">
            <div className="text-center mb-12">
              <div className="h-12 rounded-lg w-80 mx-auto mb-4 bg-white/50 backdrop-blur-sm"></div>
              <div className="h-5 rounded w-96 mx-auto bg-white/40 backdrop-blur-sm"></div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
              {[...Array(4)].map((_, i) => (
                <div
                  key={i}
                  className="backdrop-blur-lg rounded-3xl p-6 border border-white/20"
                  style={{ backgroundColor: analyticsColors.cardBg }}
                >
                  <div className="h-20"></div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div
        className="min-h-screen relative overflow-hidden"
        style={{ backgroundColor: analyticsColors.background }}
      >
        <div className="relative z-10 max-w-7xl mx-auto px-6 py-12">
          <div className="text-center">
            <div
              className="backdrop-blur-lg rounded-3xl p-8 max-w-md mx-auto shadow-xl border border-white/20"
              style={{ backgroundColor: analyticsColors.cardBg }}
            >
              <div
                className="w-16 h-16 mx-auto mb-4 rounded-2xl flex items-center justify-center"
                style={{
                  background: `linear-gradient(135deg, ${analyticsColors.danger}, ${analyticsColors.warning})`,
                }}
              >
                <BarChart3 className="w-8 h-8 text-white" />
              </div>
              <h2
                className="font-bold text-xl mb-3"
                style={{ color: analyticsColors.text }}
              >
                Error Loading Dashboard
              </h2>
              <p className="mb-6" style={{ color: analyticsColors.textLight }}>
                {error.message}
              </p>
              <button
                onClick={() => refetch()}
                className="px-6 py-3 rounded-xl font-semibold text-white transition-all duration-300 hover:scale-105 hover:shadow-lg"
                style={{
                  background: `linear-gradient(135deg, ${analyticsColors.gradient1}, ${analyticsColors.gradient2})`,
                }}
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
    if (!deviceType)
      return (
        <Globe className="h-5 w-5" style={{ color: analyticsColors.primary }} />
      );

    switch (deviceType.toLowerCase()) {
      case "mobile":
        return (
          <Smartphone
            className="h-5 w-5"
            style={{ color: analyticsColors.primary }}
          />
        );
      case "desktop":
        return (
          <Monitor
            className="h-5 w-5"
            style={{ color: analyticsColors.primary }}
          />
        );
      case "tablet":
        return (
          <Tablet
            className="h-5 w-5"
            style={{ color: analyticsColors.primary }}
          />
        );
      default:
        return (
          <Globe
            className="h-5 w-5"
            style={{ color: analyticsColors.primary }}
          />
        );
    }
  };

  const getTrendIcon = () => {
    if (overview?.trend_direction === "up")
      return (
        <ArrowUp
          className="h-5 w-5"
          style={{ color: analyticsColors.success }}
        />
      );
    if (overview?.trend_direction === "down")
      return (
        <ArrowDown
          className="h-5 w-5"
          style={{ color: analyticsColors.danger }}
        />
      );
    return (
      <Activity
        className="h-5 w-5"
        style={{ color: analyticsColors.textLight }}
      />
    );
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
      className="min-h-screen relative overflow-hidden "
      style={{ backgroundColor: analyticsColors.background }}
    >
      <div className="relative z-10 max-w-7xl mx-auto px-6 py-12 ">
        <header className="text-center mb-12">
          <div
            className="inline-block backdrop-blur-lg rounded-3xl shadow-2xl border border-white/20 p-8 hover:shadow-3xl transition-all duration-500"
            style={{ backgroundColor: analyticsColors.cardBg }}
          >
            <div className="flex items-center justify-center space-x-4 mb-4">
              <h1
                className="text-4xl font-black bg-clip-text text-transparent"
                style={{
                  backgroundImage: `linear-gradient(135deg, ${analyticsColors.gradient1}, ${analyticsColors.gradient2})`,
                }}
              >
                Analytics Dashboard
              </h1>
            </div>
            <p
              className="text-lg font-medium"
              style={{ color: analyticsColors.textLight }}
            >
              Track your URL performance with insights
            </p>
            <div className="flex space-x-4 justify-center relative">
              <button
                onClick={() => setShowShortener((prev) => !prev)}
                className="mt-6 px-5 py-2 rounded-xl text-sm font-semibold text-white transition-all duration-300 hover:scale-105"
                style={{
                  backgroundImage: `linear-gradient(135deg, ${analyticsColors.gradient1}, ${analyticsColors.gradient2})`,
                }}
              >
                {showShortener ? "Hide Shortener" : "Quick Shorten"}
              </button>
              <SignOutButton />
            </div>
          </div>
        </header>
        <div className="transition-all duration-300 min-h-[80px] mt-4">
          {showShortener && (
            <div className="max-w-xl mx-auto">
              <MiniShortenerForm />
            </div>
          )}
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <StatCard
            title="Total URLs"
            value={overview?.total_urls.toLocaleString() ?? 0}
            leftIcon={<Link className="w-5 h-5 text-[rgb(164,25,61)]" />}
            rightIcon={<Zap className="w-5 h-5 text-[#F4A261]" />}
          />
          <StatCard
            title="Total Clicks"
            value={overview?.total_clicks.toLocaleString() ?? 0}
            leftIcon={<MousePointer className="w-5 h-5 text-[#F4A261]" />}
            rightIcon={<Activity className="w-5 h-5 text-[#A4193D]" />}
          />
          <StatCard
            title="Clicks Today"
            value={overview?.clicks_today.toLocaleString() ?? 0}
            leftIcon={<Calendar className="w-5 h-5 text-[#2A9D8F]" />}
            rightIcon={getTrendIcon()}
          />
          <StatCard
            title="Avg Clicks/URL"
            value={overview?.average_clicks.toFixed(1) ?? 0}
            leftIcon={<BarChart3 className="w-5 h-5 text-[#F59E0B]" />}
            rightIcon={<Zap className="w-5 h-5 text-[#F4A261]" />}
          />
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-8">
          <div
            className="lg:col-span-2 backdrop-blur-lg rounded-3xl shadow-xl border border-white/20 p-8 hover:shadow-2xl transition-all duration-500"
            style={{ backgroundColor: analyticsColors.cardBg }}
          >
            <div className="flex items-center justify-between mb-6">
              <div>
                <h2
                  className="text-2xl font-bold"
                  style={{ color: analyticsColors.text }}
                >
                  Click Trends
                </h2>
                <p
                  className="text-sm"
                  style={{ color: analyticsColors.textLight }}
                >
                  Last 7 days
                </p>
              </div>
              <div
                className="p-3 rounded-2xl"
                style={{
                  background: `linear-gradient(135deg, ${analyticsColors.gradient1}, ${analyticsColors.gradient2})`,
                }}
              >
                <TrendingUp className="h-6 w-6 text-white" />
              </div>
            </div>

            {chartData.length > 0 ? (
              <div className="h-80">
                <ResponsiveContainer width="100%" height="100%">
                  <AreaChart data={chartData}>
                    <defs>
                      <linearGradient
                        id="colorGradient"
                        x1="0"
                        y1="0"
                        x2="0"
                        y2="1"
                      >
                        <stop
                          offset="5%"
                          stopColor={analyticsColors.primary}
                          stopOpacity={0.3}
                        />
                        <stop
                          offset="95%"
                          stopColor={analyticsColors.primary}
                          stopOpacity={0}
                        />
                      </linearGradient>
                    </defs>
                    <CartesianGrid strokeDasharray="3 3" stroke="#f1f5f9" />
                    <XAxis
                      dataKey="date"
                      stroke={analyticsColors.textLight}
                      fontSize={12}
                      fontWeight="500"
                    />
                    <YAxis
                      stroke={analyticsColors.textLight}
                      fontSize={12}
                      fontWeight="500"
                    />
                    <Tooltip
                      formatter={(value: number) => [
                        value.toLocaleString(),
                        "Clicks",
                      ]}
                      contentStyle={{
                        backgroundColor: "rgba(255, 255, 255, 0.95)",
                        border: "2px solid rgba(164, 25, 61, 0.2)",
                        borderRadius: "12px",
                        backdropFilter: "blur(10px)",
                        boxShadow: "0 10px 25px -5px rgba(0, 0, 0, 0.1)",
                      }}
                    />
                    <Area
                      type="monotone"
                      dataKey="clicks"
                      stroke={analyticsColors.primary}
                      strokeWidth={3}
                      fill="url(#colorGradient)"
                    />
                  </AreaChart>
                </ResponsiveContainer>
              </div>
            ) : (
              <div className="text-center py-16">
                <BarChart3
                  className="h-16 w-16 mx-auto mb-4"
                  style={{ color: analyticsColors.textLight }}
                />
                <p
                  className="text-lg"
                  style={{ color: analyticsColors.textLight }}
                >
                  No click data available
                </p>
              </div>
            )}
          </div>

          <div
            className="backdrop-blur-lg rounded-3xl shadow-xl border border-white/20 p-8 hover:shadow-2xl transition-all duration-500"
            style={{ backgroundColor: analyticsColors.cardBg }}
          >
            <div className="flex items-center justify-between mb-6">
              <h2
                className="text-xl font-bold"
                style={{ color: analyticsColors.text }}
              >
                Devices
              </h2>
              <Users
                className="h-6 w-6"
                style={{ color: analyticsColors.primary }}
              />
            </div>

            {deviceChartData.length > 0 ? (
              <>
                <div className="h-48 mb-4">
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie
                        data={deviceChartData}
                        cx="50%"
                        cy="50%"
                        innerRadius={40}
                        outerRadius={70}
                        dataKey="value"
                        startAngle={90}
                        endAngle={-270}
                      >
                        {deviceChartData.map((entry, index) => (
                          <Cell
                            key={`${index}-${entry}`}
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
                          backgroundColor: "rgba(255, 255, 255, 0.95)",
                          border: "none",
                          borderRadius: "12px",
                          backdropFilter: "blur(10px)",
                          boxShadow: "0 10px 25px -5px rgba(0, 0, 0, 0.1)",
                        }}
                      />
                    </PieChart>
                  </ResponsiveContainer>
                </div>

                <div className="space-y-3">
                  {device_breakdown.map((device, index) => (
                    <div
                      key={device.device_type || index}
                      className="flex items-center justify-between p-3 rounded-xl bg-white/30 hover:bg-white/50 transition-colors"
                    >
                      <div className="flex items-center space-x-3">
                        <div
                          className="w-3 h-3 rounded-full"
                          style={{
                            backgroundColor:
                              pieColors[index % pieColors.length],
                          }}
                        ></div>
                        {getDeviceIcon(device.device_type)}
                        <span
                          className="font-medium text-sm"
                          style={{ color: analyticsColors.text }}
                        >
                          {formatDeviceType(device.device_type)}
                        </span>
                      </div>
                      <div className="text-right">
                        <span
                          className="font-bold text-sm"
                          style={{ color: analyticsColors.text }}
                        >
                          {device.percentage.toFixed(0)}%
                        </span>
                      </div>
                    </div>
                  ))}
                </div>
              </>
            ) : (
              <div className="text-center py-16">
                <Smartphone
                  className="h-16 w-16 mx-auto mb-4"
                  style={{ color: analyticsColors.textLight }}
                />
                <p
                  className="text-lg"
                  style={{ color: analyticsColors.textLight }}
                >
                  No device data
                </p>
              </div>
            )}
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <div
            className="backdrop-blur-lg rounded-3xl shadow-xl border border-white/20 p-8 hover:shadow-2xl transition-all duration-500"
            style={{ backgroundColor: analyticsColors.cardBg }}
          >
            <div className="flex items-center justify-between mb-6">
              <h2
                className="text-2xl font-bold"
                style={{ color: analyticsColors.text }}
              >
                Top Performing URLs
              </h2>
              <Eye
                className="h-6 w-6"
                style={{ color: analyticsColors.primary }}
              />
            </div>

            {top_urls && top_urls.length > 0 ? (
              <div className="space-y-3">
                {top_urls.slice(0, 5).map((url, index) => (
                  <div
                    key={url.url_id}
                    className="group relative p-4 rounded-2xl bg-white/30 hover:bg-white/50 transition-all duration-300 hover:scale-[1.02] cursor-pointer"
                  >
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-3 min-w-0">
                        <div
                          className="flex-shrink-0 w-8 h-8 rounded-xl flex items-center justify-center font-bold text-white text-sm shadow-md"
                          style={{
                            background: `linear-gradient(135deg, ${analyticsColors.gradient1}, ${analyticsColors.gradient2})`,
                          }}
                        >
                          {index + 1}
                        </div>
                        <div className="min-w-0">
                          <p
                            className="font-bold text-sm truncate"
                            style={{ color: analyticsColors.text }}
                          >
                            /{url.short_code}
                          </p>
                          <p
                            className="text-xs truncate"
                            style={{ color: analyticsColors.textLight }}
                            title={url.original_url}
                          >
                            {url.original_url}
                          </p>
                        </div>
                      </div>
                      <div className="text-right flex-shrink-0">
                        <p
                          className="text-lg font-bold"
                          style={{ color: analyticsColors.primary }}
                        >
                          {url.click_count.toLocaleString()}
                        </p>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-12">
                <Link
                  className="h-16 w-16 mx-auto mb-4"
                  style={{ color: analyticsColors.textLight }}
                />
                <p
                  className="text-lg"
                  style={{ color: analyticsColors.textLight }}
                >
                  No URLs found
                </p>
              </div>
            )}
          </div>

          <div
            className="backdrop-blur-lg rounded-3xl shadow-xl border border-white/20 p-8 hover:shadow-2xl transition-all duration-500"
            style={{ backgroundColor: analyticsColors.cardBg }}
          >
            <div className="flex items-center justify-between mb-6">
              <h2
                className="text-2xl font-bold"
                style={{ color: analyticsColors.text }}
              >
                Traffic Sources
              </h2>
              <Globe
                className="h-6 w-6"
                style={{ color: analyticsColors.primary }}
              />
            </div>

            {top_referrers && top_referrers.length > 0 ? (
              <div className="space-y-3">
                {top_referrers.map((referrer, index) => {
                  const maxClicks = Math.max(
                    ...top_referrers.map((r) => r.clicks)
                  );
                  const percentage = (referrer.clicks / maxClicks) * 100;

                  return (
                    <div
                      key={`${index}-${referrer}`}
                      className="relative p-4 rounded-2xl bg-white/30 hover:bg-white/50 transition-all duration-300"
                    >
                      <div className="absolute inset-0 rounded-2xl overflow-hidden">
                        <div
                          className="h-full opacity-20 transition-all duration-500"
                          style={{
                            width: `${percentage}%`,
                            background: `linear-gradient(90deg, ${analyticsColors.gradient1}, ${analyticsColors.gradient2})`,
                          }}
                        ></div>
                      </div>
                      <div className="relative flex items-center justify-between">
                        <div className="flex items-center space-x-3">
                          <div
                            className="w-2 h-2 rounded-full"
                            style={{
                              backgroundColor:
                                pieColors[index % pieColors.length],
                            }}
                          ></div>
                          <span
                            className="font-medium text-sm"
                            style={{ color: analyticsColors.text }}
                          >
                            {referrer.referrer === "direct"
                              ? "Direct Traffic"
                              : referrer.referrer}
                          </span>
                        </div>
                        <span
                          className="font-bold text-sm"
                          style={{ color: analyticsColors.primary }}
                        >
                          {referrer.clicks.toLocaleString()}
                        </span>
                      </div>
                    </div>
                  );
                })}
              </div>
            ) : (
              <div className="text-center py-12">
                <Globe
                  className="h-16 w-16 mx-auto mb-4"
                  style={{ color: analyticsColors.textLight }}
                />
                <p
                  className="text-lg"
                  style={{ color: analyticsColors.textLight }}
                >
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
