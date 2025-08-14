// TESTING WITH AI
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
  ExternalLink,
} from "lucide-react";
import { useAnalyticsDashboard } from "../api/analytics";
import { StatCard } from "../components/StatsCard";
import SignOutButton from "../components/SignOutButton";
import { MiniShortenerFormV2 } from "../components/MiniShortenFormV2";

// Enhanced colors with additional variations for better visual hierarchy
const analyticsColors = {
  primary: "#A4193D",
  secondary: "#FFDFB9",
  accent: "#F4A261",
  neutral: "#2A9D8F",
  background: "#F2EFDE",
  cardBg: "rgba(255, 255, 255, 0.8)",
  cardBgHover: "rgba(255, 255, 255, 0.95)",
  text: "#111827",
  textLight: "#6B7280",
  textXLight: "#9CA3AF",
  success: "#10B981",
  danger: "#EF4444",
  warning: "#F59E0B",
  info: "#3B82F6",
  gradient1: "#A4193D",
  gradient2: "#F4A261",
  gradient3: "#2A9D8F",
  border: "rgba(164, 25, 61, 0.1)",
  borderHover: "rgba(164, 25, 61, 0.2)",
  surfaceLight: "rgba(255, 255, 255, 0.6)",
  surfaceDark: "rgba(164, 25, 61, 0.05)",
};

const pieColors = [
  analyticsColors.primary,
  analyticsColors.accent,
  analyticsColors.neutral,
  analyticsColors.warning,
];

interface AnalyticsDashboardProps {
  token: string;
}

export const AnalyticsDashboardV2: React.FC<AnalyticsDashboardProps> = ({
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
      className="min-h-screen relative overflow-hidden"
      style={{ backgroundColor: analyticsColors.background }}
    >
      <div className="relative z-10 max-w-7xl mx-auto px-6 py-8">
        {/* Enhanced Header with Better Typography Hierarchy */}
        <header className="text-center mb-8">
          <div
            className="inline-block backdrop-blur-lg rounded-3xl shadow-2xl border-2 p-8 hover:shadow-3xl transition-all duration-500 hover:-translate-y-1"
            style={{
              backgroundColor: analyticsColors.cardBg,
              borderColor: analyticsColors.border,
            }}
          >
            <div className="flex items-center justify-center space-x-4 mb-6">
              <div
                className="p-3 rounded-2xl"
                style={{
                  background: `linear-gradient(135deg, ${analyticsColors.gradient1}, ${analyticsColors.gradient2})`,
                }}
              >
                <BarChart3 className="h-8 w-8 text-white" />
              </div>
              <h1
                className="text-5xl font-black tracking-tight bg-clip-text text-transparent"
                style={{
                  backgroundImage: `linear-gradient(135deg, ${analyticsColors.gradient1}, ${analyticsColors.gradient2})`,
                }}
              >
                Analytics
              </h1>
            </div>
            <p
              className="text-xl font-medium mb-6"
              style={{ color: analyticsColors.textLight }}
            >
              Track your URL performance with insights
            </p>
            <div className="flex space-x-4 justify-center">
              <button
                onClick={() => setShowShortener((prev) => !prev)}
                className="px-6 py-3 rounded-2xl text-sm font-bold text-white transition-all duration-300 hover:scale-105 hover:shadow-lg hover:-translate-y-0.5"
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

        {/* Shortener Form with smooth transition */}
        <div
          className="transition-all duration-500 ease-out overflow-hidden mb-8"
          style={{
            maxHeight: showShortener ? "200px" : "0px",
            opacity: showShortener ? 1 : 0,
          }}
        >
          <div className="max-w-xl mx-auto pt-4">
            <MiniShortenerFormV2 />
          </div>
        </div>

        {/* Enhanced Stats Grid with varied card heights */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 mb-12">
          <div className="md:col-span-1 lg:col-span-1">
            <StatCard
              title="Total URLs"
              value={overview?.total_urls.toLocaleString() ?? 0}
              leftIcon={
                <Link
                  className="w-6 h-6"
                  style={{ color: analyticsColors.primary }}
                />
              }
              rightIcon={
                <Zap
                  className="w-5 h-5"
                  style={{ color: analyticsColors.accent }}
                />
              }
            />
          </div>
          <div className="md:col-span-1 lg:col-span-2">
            <StatCard
              title="Total Clicks"
              value={overview?.total_clicks.toLocaleString() ?? 0}
              leftIcon={
                <MousePointer
                  className="w-6 h-6"
                  style={{ color: analyticsColors.accent }}
                />
              }
              rightIcon={
                <Activity
                  className="w-5 h-5"
                  style={{ color: analyticsColors.primary }}
                />
              }
            />
          </div>
          <div className="md:col-span-2 lg:col-span-1">
            <StatCard
              title="Today"
              value={overview?.clicks_today.toLocaleString() ?? 0}
              leftIcon={
                <Calendar
                  className="w-6 h-6"
                  style={{ color: analyticsColors.neutral }}
                />
              }
              rightIcon={getTrendIcon()}
            />
          </div>
        </div>

        {/* Enhanced Bento Grid Layout with asymmetric sizing */}
        <div className="grid grid-cols-1 lg:grid-cols-5 gap-8 mb-12">
          {/* Click Trends - Takes 3 columns */}
          <div
            className="lg:col-span-3 backdrop-blur-lg rounded-3xl shadow-xl border-2 p-8 hover:shadow-2xl transition-all duration-500 hover:-translate-y-1 group"
            style={{
              backgroundColor: analyticsColors.cardBg,
              borderColor: analyticsColors.border,
            }}
          >
            <div className="flex items-start justify-between mb-8">
              <div>
                <h2
                  className="text-3xl font-black mb-2"
                  style={{ color: analyticsColors.text }}
                >
                  Click Trends
                </h2>
                <p
                  className="text-sm font-medium tracking-wide uppercase"
                  style={{ color: analyticsColors.textLight }}
                >
                  Last 7 days
                </p>
              </div>
              <div
                className="p-4 rounded-2xl transition-transform duration-300 group-hover:scale-110"
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
                      fontWeight="600"
                    />
                    <YAxis
                      stroke={analyticsColors.textLight}
                      fontSize={12}
                      fontWeight="600"
                    />
                    <Tooltip
                      formatter={(value: number) => [
                        value.toLocaleString(),
                        "Clicks",
                      ]}
                      contentStyle={{
                        backgroundColor: analyticsColors.cardBgHover,
                        border: `2px solid ${analyticsColors.borderHover}`,
                        borderRadius: "16px",
                        backdropFilter: "blur(10px)",
                        boxShadow: "0 25px 50px -12px rgba(0, 0, 0, 0.15)",
                        fontWeight: "600",
                      }}
                    />
                    <Area
                      type="monotone"
                      dataKey="clicks"
                      stroke={analyticsColors.primary}
                      strokeWidth={4}
                      fill="url(#colorGradient)"
                    />
                  </AreaChart>
                </ResponsiveContainer>
              </div>
            ) : (
              <div className="text-center py-20">
                <BarChart3
                  className="h-20 w-20 mx-auto mb-6"
                  style={{ color: analyticsColors.textXLight }}
                />
                <p
                  className="text-xl font-semibold"
                  style={{ color: analyticsColors.textLight }}
                >
                  No click data available
                </p>
              </div>
            )}
          </div>

          {/* Devices Card - Takes 2 columns */}
          <div
            className="lg:col-span-2 backdrop-blur-lg rounded-3xl shadow-xl border-2 p-8 hover:shadow-2xl transition-all duration-500 hover:-translate-y-1 group"
            style={{
              backgroundColor: analyticsColors.cardBg,
              borderColor: analyticsColors.border,
            }}
          >
            <div className="flex items-center justify-between mb-8">
              <h2
                className="text-2xl font-black"
                style={{ color: analyticsColors.text }}
              >
                Devices
              </h2>
              <div
                className="p-3 rounded-xl transition-transform duration-300 group-hover:scale-110"
                style={{ backgroundColor: analyticsColors.surfaceLight }}
              >
                <Users
                  className="h-6 w-6"
                  style={{ color: analyticsColors.primary }}
                />
              </div>
            </div>

            {deviceChartData.length > 0 ? (
              <>
                <div className="h-48 mb-6">
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie
                        data={deviceChartData}
                        cx="50%"
                        cy="50%"
                        innerRadius={45}
                        outerRadius={75}
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
                          backgroundColor: analyticsColors.cardBgHover,
                          border: "none",
                          borderRadius: "16px",
                          backdropFilter: "blur(10px)",
                          boxShadow: "0 25px 50px -12px rgba(0, 0, 0, 0.15)",
                          fontWeight: "600",
                        }}
                      />
                    </PieChart>
                  </ResponsiveContainer>
                </div>

                <div className="space-y-4">
                  {device_breakdown.map((device, index) => (
                    <div
                      key={device.device_type || index}
                      className="flex items-center justify-between p-4 rounded-2xl transition-all duration-300 hover:scale-[1.02] cursor-pointer border-2"
                      style={{
                        backgroundColor: analyticsColors.surfaceLight,
                        borderColor: "transparent",
                      }}
                      onMouseEnter={(e) => {
                        e.currentTarget.style.borderColor =
                          analyticsColors.borderHover;
                        e.currentTarget.style.backgroundColor =
                          analyticsColors.cardBgHover;
                      }}
                      onMouseLeave={(e) => {
                        e.currentTarget.style.borderColor = "transparent";
                        e.currentTarget.style.backgroundColor =
                          analyticsColors.surfaceLight;
                      }}
                    >
                      <div className="flex items-center space-x-4">
                        <div
                          className="w-4 h-4 rounded-full"
                          style={{
                            backgroundColor:
                              pieColors[index % pieColors.length],
                          }}
                        ></div>
                        {getDeviceIcon(device.device_type)}
                        <span
                          className="font-bold text-sm"
                          style={{ color: analyticsColors.text }}
                        >
                          {formatDeviceType(device.device_type)}
                        </span>
                      </div>
                      <div className="text-right">
                        <span
                          className="font-black text-lg"
                          style={{ color: analyticsColors.primary }}
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

        {/* Bottom Row - Asymmetric Layout */}
        <div className="grid grid-cols-1 lg:grid-cols-7 gap-8">
          {/* Top URLs - Takes 4 columns */}
          <div
            className="lg:col-span-4 backdrop-blur-lg rounded-3xl shadow-xl border-2 p-8 hover:shadow-2xl transition-all duration-500 hover:-translate-y-1 group"
            style={{
              backgroundColor: analyticsColors.cardBg,
              borderColor: analyticsColors.border,
            }}
          >
            <div className="flex items-center justify-between mb-8">
              <div>
                <h2
                  className="text-3xl font-black mb-2"
                  style={{ color: analyticsColors.text }}
                >
                  Top Performing URLs
                </h2>
                <p
                  className="text-sm font-medium tracking-wide uppercase"
                  style={{ color: analyticsColors.textLight }}
                >
                  Most clicked links
                </p>
              </div>
              <div
                className="p-4 rounded-2xl transition-transform duration-300 group-hover:scale-110"
                style={{
                  background: `linear-gradient(135deg, ${analyticsColors.neutral}, ${analyticsColors.accent})`,
                }}
              >
                <Eye className="h-6 w-6 text-white" />
              </div>
            </div>

            {top_urls && top_urls.length > 0 ? (
              <div className="space-y-4">
                {top_urls.slice(0, 5).map((url, index) => (
                  <div
                    key={url.url_id}
                    className="group relative p-6 rounded-3xl transition-all duration-300 hover:scale-[1.02] cursor-pointer border-2"
                    style={{
                      backgroundColor:
                        index === 0
                          ? analyticsColors.surfaceDark
                          : analyticsColors.surfaceLight,
                      borderColor: "transparent",
                    }}
                    onMouseEnter={(e) => {
                      e.currentTarget.style.borderColor =
                        analyticsColors.borderHover;
                      e.currentTarget.style.backgroundColor =
                        analyticsColors.cardBgHover;
                    }}
                    onMouseLeave={(e) => {
                      e.currentTarget.style.borderColor = "transparent";
                      e.currentTarget.style.backgroundColor =
                        index === 0
                          ? analyticsColors.surfaceDark
                          : analyticsColors.surfaceLight;
                    }}
                  >
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-4 min-w-0 flex-1">
                        <div
                          className={`flex-shrink-0 w-10 h-10 rounded-2xl flex items-center justify-center font-black text-white ${
                            index === 0 ? "text-lg" : "text-sm"
                          } shadow-lg`}
                          style={{
                            background:
                              index === 0
                                ? `linear-gradient(135deg, ${analyticsColors.gradient1}, ${analyticsColors.gradient2})`
                                : `linear-gradient(135deg, ${analyticsColors.textLight}, ${analyticsColors.neutral})`,
                          }}
                        >
                          {index + 1}
                        </div>
                        <div className="min-w-0 flex-1">
                          <div className="flex items-center space-x-2 mb-1">
                            <p
                              className={`font-black truncate ${
                                index === 0 ? "text-lg" : "text-sm"
                              }`}
                              style={{ color: analyticsColors.text }}
                            >
                              /{url.short_code}
                            </p>
                            <ExternalLink
                              className="w-3 h-3"
                              style={{ color: analyticsColors.textXLight }}
                            />
                          </div>
                          <p
                            className="text-xs truncate"
                            style={{ color: analyticsColors.textLight }}
                            title={url.original_url}
                          >
                            {url.original_url}
                          </p>
                        </div>
                      </div>
                      <div className="text-right flex-shrink-0 ml-4">
                        <p
                          className={`font-black ${
                            index === 0 ? "text-2xl" : "text-lg"
                          }`}
                          style={{ color: analyticsColors.primary }}
                        >
                          {url.click_count.toLocaleString()}
                        </p>
                        <p
                          className="text-xs font-medium uppercase tracking-wide"
                          style={{ color: analyticsColors.textLight }}
                        >
                          clicks
                        </p>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-16">
                <Link
                  className="h-20 w-20 mx-auto mb-6"
                  style={{ color: analyticsColors.textXLight }}
                />
                <p
                  className="text-xl font-semibold"
                  style={{ color: analyticsColors.textLight }}
                >
                  No URLs found
                </p>
              </div>
            )}
          </div>

          {/* Traffic Sources - Takes 3 columns */}
          <div
            className="lg:col-span-3 backdrop-blur-lg rounded-3xl shadow-xl border-2 p-8 hover:shadow-2xl transition-all duration-500 hover:-translate-y-1 group"
            style={{
              backgroundColor: analyticsColors.cardBg,
              borderColor: analyticsColors.border,
            }}
          >
            <div className="flex items-center justify-between mb-8">
              <div>
                <h2
                  className="text-2xl font-black mb-2"
                  style={{ color: analyticsColors.text }}
                >
                  Traffic Sources
                </h2>
                <p
                  className="text-sm font-medium tracking-wide uppercase"
                  style={{ color: analyticsColors.textLight }}
                >
                  Where clicks come from
                </p>
              </div>
              <div
                className="p-3 rounded-xl transition-transform duration-300 group-hover:scale-110"
                style={{
                  background: `linear-gradient(135deg, ${analyticsColors.info}, ${analyticsColors.neutral})`,
                }}
              >
                <Globe className="h-6 w-6 text-white" />
              </div>
            </div>

            {top_referrers && top_referrers.length > 0 ? (
              <div className="space-y-4">
                {top_referrers.map((referrer, index) => {
                  const maxClicks = Math.max(
                    ...top_referrers.map((r) => r.clicks)
                  );
                  const percentage = (referrer.clicks / maxClicks) * 100;

                  return (
                    <div
                      key={`${index}-${referrer}`}
                      className="relative p-5 rounded-2xl transition-all duration-300 hover:scale-[1.02] cursor-pointer border-2"
                      style={{
                        backgroundColor: analyticsColors.surfaceLight,
                        borderColor: "transparent",
                      }}
                      onMouseEnter={(e) => {
                        e.currentTarget.style.borderColor =
                          analyticsColors.borderHover;
                        e.currentTarget.style.backgroundColor =
                          analyticsColors.cardBgHover;
                      }}
                      onMouseLeave={(e) => {
                        e.currentTarget.style.borderColor = "transparent";
                        e.currentTarget.style.backgroundColor =
                          analyticsColors.surfaceLight;
                      }}
                    >
                      <div className="absolute inset-0 rounded-2xl overflow-hidden">
                        <div
                          className="h-full opacity-10 transition-all duration-700"
                          style={{
                            width: `${percentage}%`,
                            background: `linear-gradient(90deg, ${analyticsColors.gradient1}, ${analyticsColors.gradient2})`,
                          }}
                        ></div>
                      </div>
                      <div className="relative flex items-center justify-between">
                        <div className="flex items-center space-x-4">
                          <div
                            className="w-3 h-3 rounded-full"
                            style={{
                              backgroundColor:
                                pieColors[index % pieColors.length],
                            }}
                          ></div>
                          <span
                            className="font-bold text-sm"
                            style={{ color: analyticsColors.text }}
                          >
                            {referrer.referrer === "direct"
                              ? "Direct Traffic"
                              : referrer.referrer}
                          </span>
                        </div>
                        <div className="text-right">
                          <span
                            className="font-black text-lg"
                            style={{ color: analyticsColors.primary }}
                          >
                            {referrer.clicks.toLocaleString()}
                          </span>
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            ) : (
              <div className="text-center py-16">
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

        {/* Average Clicks Card - Full width accent */}
        <div className="mt-8">
          <div
            className="backdrop-blur-lg rounded-3xl shadow-xl border-2 p-8 hover:shadow-2xl transition-all duration-500 hover:-translate-y-1 group"
            style={{
              backgroundColor: analyticsColors.cardBg,
              borderColor: analyticsColors.border,
            }}
          >
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-6">
                <div
                  className="p-4 rounded-2xl transition-transform duration-300 group-hover:scale-110"
                  style={{
                    background: `linear-gradient(135deg, ${analyticsColors.warning}, ${analyticsColors.accent})`,
                  }}
                >
                  <BarChart3 className="h-8 w-8 text-white" />
                </div>
                <div>
                  <h3
                    className="text-2xl font-black mb-1"
                    style={{ color: analyticsColors.text }}
                  >
                    Average Clicks per URL
                  </h3>
                  <p
                    className="text-sm font-medium tracking-wide uppercase"
                    style={{ color: analyticsColors.textLight }}
                  >
                    Performance metric
                  </p>
                </div>
              </div>
              <div className="text-right">
                <span
                  className="text-5xl font-black"
                  style={{ color: analyticsColors.primary }}
                >
                  {overview?.average_clicks?.toFixed(1) ?? "0.0"}
                </span>
                <p
                  className="text-sm font-medium uppercase tracking-wide mt-1"
                  style={{ color: analyticsColors.textLight }}
                >
                  avg clicks
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AnalyticsDashboardV2;
