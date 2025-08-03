/* eslint-disable @typescript-eslint/no-unused-vars */
import { useEffect, useState } from "react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { TrendingUp, Link, BarChart3, Calendar } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { useUrlStore } from "../store/UrlStore";
import { useAuthStore } from "../store/AuthStore";
import { SignOutButton } from "../components/SignOutButton";
import { mockData } from "../mockData/mockDataDashBoard";
import type { DashboardStats } from "../types/urls";
import QuickShorten from "../components/QuickShorten";
import { useShortenUrlAuth } from "../hooks/useShortenUrlAuth";

const Dashboard = () => {
  const isAuthenticated = true;
  const userId = "demo-user-id";
  //const { isAuthenticated, userId, logout } = useAuthStore(); TO USE ONCE AUTH API IMPLEMENTED IN GO BACKEND
  const navigate = useNavigate();

  const { shorten } = useShortenUrlAuth();
  const shortenLoading = useUrlStore((state) => state.loading);
  const [url, setUrl] = useState("");

  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!isAuthenticated) {
      navigate("/login");
      return;
    }

    fetchDashboardData();
  }, [isAuthenticated, navigate, userId]);

  const handleShortenSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!url) return;
    shorten({
      url,
      is_public: true,
      user_id: userId,
      code_length: 8,
    });
    setUrl("");
  };

  const fetchDashboardData = async () => {
    try {
      setLoading(true);
      setError(null);

      // Mock data for demonstration
      
      await new Promise((resolve) => setTimeout(resolve, 1000));
      setStats(mockData);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
      console.error("Dashboard fetch error:", err);
    } finally {
      setLoading(false);
    }
  };

  const formatUrl = (url: string, maxLength: number = 40) => {
    if (url.length <= maxLength) return url;
    return url.substring(0, maxLength) + "...";
  };

  const formatWeekLabel = (week: string) => {
    const date = new Date(week);
    return date.toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
    });
  };

  if (loading) {
    return (
      <div className="min-h-screen" style={{ backgroundColor: "#fefefe" }}>
        <div className="max-w-7xl mx-auto px-6 py-12">
          <div className="animate-pulse">
            <div className="text-center mb-12">
              <div className="h-12 bg-gradient-to-r from-gray-200 via-gray-300 to-gray-200 rounded-lg w-80 mx-auto mb-4 animate-pulse"></div>
              <div className="h-5 bg-gradient-to-r from-gray-200 via-gray-300 to-gray-200 rounded w-96 mx-auto animate-pulse"></div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
              {[...Array(4)].map((_, i) => (
                <div
                  key={i}
                  className="bg-white p-6 rounded-xl shadow-lg border border-gray-100"
                >
                  <div className="flex items-center">
                    <div className="w-14 h-14 bg-gradient-to-r from-gray-200 via-gray-300 to-gray-200 rounded-xl animate-pulse"></div>
                    <div className="ml-4 flex-1">
                      <div className="h-4 bg-gradient-to-r from-gray-200 via-gray-300 to-gray-200 rounded w-20 mb-3 animate-pulse"></div>
                      <div className="h-8 bg-gradient-to-r from-gray-200 via-gray-300 to-gray-200 rounded w-16 animate-pulse"></div>
                    </div>
                  </div>
                </div>
              ))}
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              <div className="bg-white p-8 rounded-xl shadow-lg border border-gray-100">
                <div className="h-7 bg-gradient-to-r from-gray-200 via-gray-300 to-gray-200 rounded w-48 mb-6 animate-pulse"></div>
                <div className="space-y-4">
                  {[...Array(5)].map((_, i) => (
                    <div
                      key={i}
                      className="flex items-center justify-between p-4 rounded-xl border-2 border-gray-100 shadow-md"
                      style={{ backgroundColor: "#fefefe" }}
                    >
                      <div className="flex items-center space-x-4">
                        <div className="w-10 h-10 bg-gradient-to-r from-gray-200 via-gray-300 to-gray-200 rounded-full animate-pulse"></div>
                        <div className="min-w-0">
                          <div className="h-5 bg-gradient-to-r from-gray-200 via-gray-300 to-gray-200 rounded w-24 mb-2 animate-pulse"></div>
                          <div className="h-4 bg-gradient-to-r from-gray-200 via-gray-300 to-gray-200 rounded w-48 animate-pulse"></div>
                        </div>
                      </div>
                      <div className="text-right">
                        <div className="h-6 bg-gradient-to-r from-gray-200 via-gray-300 to-gray-200 rounded w-8 mb-1 animate-pulse"></div>
                        <div className="h-3 bg-gradient-to-r from-gray-200 via-gray-300 to-gray-200 rounded w-12 animate-pulse"></div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              <div className="bg-white p-8 rounded-xl shadow-lg border border-gray-100">
                <div className="h-7 bg-gradient-to-r from-gray-200 via-gray-300 to-gray-200 rounded w-36 mb-6 animate-pulse"></div>
                <div className="h-80 bg-gradient-to-r from-gray-100 via-gray-200 to-gray-100 rounded-lg animate-pulse flex items-center justify-center">
                  <div className="flex flex-col items-center space-y-4">
                    <div className="w-16 h-16 bg-gradient-to-r from-gray-200 via-gray-300 to-gray-200 rounded-full animate-pulse"></div>
                    <div className="h-4 bg-gradient-to-r from-gray-200 via-gray-300 to-gray-200 rounded w-32 animate-pulse"></div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen" style={{ backgroundColor: "#fefefe" }}>
        <div className="max-w-7xl mx-auto px-6 py-12">
          <div className="text-center">
            <div className="bg-red-50 border border-red-200 rounded-xl p-8 max-w-md mx-auto shadow-lg">
              <h2 className="text-red-800 font-bold text-xl mb-3">
                Error Loading Dashboard
              </h2>
              <p className="text-red-600 mb-6">{error}</p>
              <button
                onClick={fetchDashboardData}
                className="px-6 py-3 rounded-lg font-medium transition-all duration-200 transform hover:scale-105"
                style={{
                  backgroundColor: "#6366f1",
                  color: "white",
                }}
                onMouseEnter={(e) =>
                  (e.currentTarget.style.backgroundColor = "#4f46e5")
                }
                onMouseLeave={(e) =>
                  (e.currentTarget.style.backgroundColor = "#6366f1")
                }
              >
                Try Again
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen" style={{ backgroundColor: "#fefefe" }}>
      <div className="max-w-7xl mx-auto px-6 py-12">
        <div className="text-center mb-12">
          <h1 className="text-5xl font-bold mb-4" style={{ color: "#111827" }}>
            Dashboard
          </h1>
          <p className="text-xl max-w-2xl mx-auto" style={{ color: "#6b7280" }}>
            Get comprehensive insights into your URL performance with real-time
            analytics and engagement metrics
          </p>
          <SignOutButton />
        </div>
        <div className="mb-10">
          <QuickShorten />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
          <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-100 transition-all duration-300 hover:shadow-xl hover:-translate-y-1">
            <div className="flex items-center">
              <div
                className="p-3 rounded-xl"
                style={{ backgroundColor: "#fcd5ce" }}
              >
                <Link className="h-7 w-7" style={{ color: "#6366f1" }} />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium" style={{ color: "#6b7280" }}>
                  Total URLs
                </p>
                <p className="text-3xl font-bold" style={{ color: "#111827" }}>
                  {stats?.totalUrls || 0}
                </p>
              </div>
            </div>
          </div>

          <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-100 transition-all duration-300 hover:shadow-xl hover:-translate-y-1">
            <div className="flex items-center">
              <div
                className="p-3 rounded-xl"
                style={{ backgroundColor: "#fcd5ce" }}
              >
                <TrendingUp className="h-7 w-7" style={{ color: "#6366f1" }} />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium" style={{ color: "#6b7280" }}>
                  Total Clicks
                </p>
                <p className="text-3xl font-bold" style={{ color: "#111827" }}>
                  {stats?.totalClicks || 0}
                </p>
              </div>
            </div>
          </div>

          <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-100 transition-all duration-300 hover:shadow-xl hover:-translate-y-1">
            <div className="flex items-center">
              <div
                className="p-3 rounded-xl"
                style={{ backgroundColor: "#fcd5ce" }}
              >
                <BarChart3 className="h-7 w-7" style={{ color: "#6366f1" }} />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium" style={{ color: "#6b7280" }}>
                  Avg Clicks/URL
                </p>
                <p className="text-3xl font-bold" style={{ color: "#111827" }}>
                  {stats?.totalUrls
                    ? Math.round((stats.totalClicks || 0) / stats.totalUrls)
                    : 0}
                </p>
              </div>
            </div>
          </div>

          <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-100 transition-all duration-300 hover:shadow-xl hover:-translate-y-1">
            <div className="flex items-center">
              <div
                className="p-3 rounded-xl"
                style={{ backgroundColor: "#fcd5ce" }}
              >
                <Calendar className="h-7 w-7" style={{ color: "#6366f1" }} />
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium" style={{ color: "#6b7280" }}>
                  This Week
                </p>
                <p className="text-3xl font-bold" style={{ color: "#111827" }}>
                  {stats?.weeklyData?.[stats.weeklyData.length - 1]?.clicks ||
                    0}
                </p>
              </div>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <div className="bg-white p-8 rounded-xl shadow-lg border border-gray-100">
            <h2
              className="text-2xl font-bold mb-6"
              style={{ color: "#111827" }}
            >
              Top 5 URLs by Clicks
            </h2>
            {stats?.topUrls && stats.topUrls.length > 0 ? (
              <div className="space-y-4">
                {stats.topUrls.map((url, index) => (
                  <div
                    key={url.id}
                    className="flex items-center justify-between p-4 rounded-xl border-2 border-gray-100 shadow-md transition-all duration-300 hover:shadow-lg hover:-translate-y-1 hover:border-gray-200 cursor-pointer"
                    style={{ backgroundColor: "#fefefe" }}
                  >
                    <div className="flex items-center space-x-4">
                      <div
                        className="flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center font-bold text-white shadow-md"
                        style={{ backgroundColor: "#6366f1" }}
                      >
                        <span className="text-lg">{index + 1}</span>
                      </div>
                      <div className="min-w-0">
                        <p
                          className="text-lg font-bold truncate"
                          style={{ color: "#111827" }}
                        >
                          /{url.short_code}
                        </p>
                        <p
                          className="text-sm truncate"
                          style={{ color: "#6b7280" }}
                          title={url.original_url}
                        >
                          {formatUrl(url.original_url)}
                        </p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p
                        className="text-2xl font-bold"
                        style={{ color: "#6366f1" }}
                      >
                        {url.click_count}
                      </p>
                      <p
                        className="text-sm font-medium"
                        style={{ color: "#6b7280" }}
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
                  style={{ color: "#6b7280" }}
                />
                <p className="text-lg" style={{ color: "#6b7280" }}>
                  No URLs found
                </p>
              </div>
            )}
          </div>

          <div className="bg-white p-8 rounded-xl shadow-lg border border-gray-100">
            <h2
              className="text-2xl font-bold mb-6"
              style={{ color: "#111827" }}
            >
              Weekly Clicks
            </h2>
            {stats?.weeklyData && stats.weeklyData.length > 0 ? (
              <div className="h-80">
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart data={stats.weeklyData}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#f1f5f9" />
                    <XAxis
                      dataKey="week"
                      tickFormatter={formatWeekLabel}
                      stroke="#6b7280"
                      fontSize={12}
                      fontWeight="500"
                    />
                    <YAxis stroke="#6b7280" fontSize={12} fontWeight="500" />
                    <Tooltip
                      labelFormatter={(value) =>
                        `Week of ${formatWeekLabel(value)}`
                      }
                      formatter={(value: number) => [value, "Clicks"]}
                      contentStyle={{
                        backgroundColor: "#ffffff",
                        border: "2px solid #e5e7eb",
                        borderRadius: "12px",
                        boxShadow: "0 10px 25px -5px rgba(0, 0, 0, 0.1)",
                        fontWeight: "500",
                      }}
                    />
                    <Line
                      type="monotone"
                      dataKey="clicks"
                      stroke="#6366f1"
                      strokeWidth={3}
                      dot={{ fill: "#6366f1", strokeWidth: 3, r: 5 }}
                      activeDot={{
                        r: 8,
                        stroke: "#6366f1",
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
                  style={{ color: "#6b7280" }}
                />
                <p className="text-lg" style={{ color: "#6b7280" }}>
                  No click data available
                </p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
