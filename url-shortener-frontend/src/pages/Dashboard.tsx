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

const cardStyle: React.CSSProperties = {
  background: "var(--bg-card)",
  border: "1px solid var(--border)",
  borderRadius: "0.875rem",
  padding: "1.5rem",
};

export const AnalyticsDashboard: React.FC<AnalyticsDashboardProps> = ({ token }) => {
  const { data: dashboard, isLoading, error, refetch } = useAnalyticsDashboard(token);
  const [showShortener, setShowShortener] = useState(false);

  if (isLoading) {
    return (
      <div style={{ minHeight: "100vh", background: "var(--bg)" }}>
        <Navbar onToggleShortener={() => {}} showShortener={false} />
        <div style={{ maxWidth: 1280, margin: "0 auto", padding: "2rem 1.5rem" }}>
          <div style={{ opacity: 0.4 }}>
            <div style={{ height: 24, background: "var(--bg-elevated)", borderRadius: 6, width: 140, marginBottom: 8 }} />
            <div style={{ height: 16, background: "var(--bg-elevated)", borderRadius: 6, width: 220, marginBottom: 32 }} />
            <div style={{ display: "grid", gridTemplateColumns: "repeat(4, 1fr)", gap: 16, marginBottom: 32 }}>
              {[...Array(4)].map((_, i) => (
                <div key={i} style={{ background: "var(--bg-card)", border: "1px solid var(--border)", borderRadius: 14, padding: 20, height: 120 }} />
              ))}
            </div>
            <div style={{ background: "var(--bg-card)", border: "1px solid var(--border)", borderRadius: 14, height: 280 }} />
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div style={{ minHeight: "100vh", background: "var(--bg)" }}>
        <Navbar onToggleShortener={() => {}} showShortener={false} />
        <div style={{ display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center", paddingTop: 120 }}>
          <BarChart3 size={40} color="var(--text-muted)" style={{ marginBottom: 16 }} />
          <h2 className="font-display" style={{ fontSize: "1rem", fontWeight: 700, color: "var(--text)", marginBottom: 6 }}>Failed to load dashboard</h2>
          <p style={{ fontSize: "0.875rem", color: "var(--text-muted)", marginBottom: 24 }}>{error.message}</p>
          <button onClick={() => refetch()} className="btn-primary" style={{ padding: "0.625rem 1.25rem", fontSize: "0.875rem" }}>
            Try again
          </button>
        </div>
      </div>
    );
  }

  if (!dashboard) {
    return (
      <div style={{ minHeight: "100vh", background: "var(--bg)", display: "flex", alignItems: "center", justifyContent: "center" }}>
        <p style={{ fontSize: "0.875rem", color: "var(--text-muted)" }}>No analytics data available.</p>
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
      case "mobile":   return <Smartphone size={16} color="var(--text-muted)" />;
      case "desktop":  return <Monitor size={16} color="var(--text-muted)" />;
      case "tablet":   return <Tablet size={16} color="var(--text-muted)" />;
      default:         return <Globe size={16} color="var(--text-muted)" />;
    }
  };

  const chartData = daily_trend.map((day) => ({
    date: new Date(day.date).toLocaleDateString("en-US", { month: "short", day: "numeric" }),
    clicks: day.clicks,
  }));

  return (
    <div style={{ minHeight: "100vh", background: "var(--bg)" }}>
      <Navbar
        onToggleShortener={() => setShowShortener((prev) => !prev)}
        showShortener={showShortener}
      />

      <main style={{ maxWidth: 1280, margin: "0 auto", padding: "2rem 1.5rem" }}>
        {/* Header */}
        <div style={{ marginBottom: "2rem" }}>
          <h1 className="font-display" style={{ fontSize: "1.375rem", fontWeight: 800, color: "var(--text)", letterSpacing: "-0.02em", marginBottom: "0.25rem" }}>
            Analytics
          </h1>
          <p style={{ fontSize: "0.875rem", color: "var(--text-muted)" }}>Your link performance at a glance</p>
        </div>

        {/* Quick shortener */}
        {showShortener && (
          <div style={{ marginBottom: "2rem", maxWidth: 560 }} className="animate-fade-up">
            <MiniShortenerForm />
          </div>
        )}

        {/* Stat cards */}
        <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(200px, 1fr))", gap: "1rem", marginBottom: "1.5rem" }}>
          <StatCard title="Total URLs" value={overview?.total_urls.toLocaleString() ?? 0} icon={<Link size={16} color="var(--accent-bright)" />} />
          <StatCard title="Total Clicks" value={overview?.total_clicks.toLocaleString() ?? 0} icon={<MousePointer size={16} color="var(--accent-bright)" />} />
          <StatCard title="Clicks Today" value={overview?.clicks_today.toLocaleString() ?? 0} icon={<Calendar size={16} color="var(--accent-bright)" />} />
          <StatCard title="Avg / URL" value={overview?.average_clicks.toFixed(1) ?? 0} icon={<BarChart3 size={16} color="var(--accent-bright)" />} />
        </div>

        {/* Charts row */}
        <div style={{ display: "grid", gridTemplateColumns: "2fr 1fr", gap: "1rem", marginBottom: "1rem" }}>
          {/* Click trend */}
          <div style={cardStyle}>
            <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", marginBottom: "1.5rem" }}>
              <div>
                <h2 className="font-display" style={{ fontSize: "0.9375rem", fontWeight: 700, color: "var(--text)", letterSpacing: "-0.01em" }}>Click Trends</h2>
                <p style={{ fontSize: "0.75rem", color: "var(--text-muted)", marginTop: 2 }}>Last 7 days</p>
              </div>
              <TrendingUp size={16} color="var(--text-muted)" />
            </div>

            {chartData.length > 0 ? (
              <div style={{ height: 240 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <AreaChart data={chartData}>
                    <defs>
                      <linearGradient id="clickGradient" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="5%" stopColor="#7c3aed" stopOpacity={0.3} />
                        <stop offset="95%" stopColor="#7c3aed" stopOpacity={0} />
                      </linearGradient>
                    </defs>
                    <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" vertical={false} />
                    <XAxis dataKey="date" stroke="var(--text-muted)" fontSize={11} tickLine={false} axisLine={false} />
                    <YAxis stroke="var(--text-muted)" fontSize={11} tickLine={false} axisLine={false} />
                    <Tooltip
                      formatter={(value: number) => [value.toLocaleString(), "Clicks"]}
                      contentStyle={{
                        backgroundColor: "var(--bg-elevated)",
                        border: "1px solid var(--border-bright)",
                        borderRadius: "8px",
                        color: "var(--text)",
                        fontSize: "12px",
                        boxShadow: "0 4px 20px rgba(0,0,0,0.4)",
                      }}
                      labelStyle={{ color: "var(--text-secondary)" }}
                    />
                    <Area type="monotone" dataKey="clicks" stroke="#9d6efa" strokeWidth={2} fill="url(#clickGradient)" />
                  </AreaChart>
                </ResponsiveContainer>
              </div>
            ) : (
              <div style={{ display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center", height: 240 }}>
                <BarChart3 size={32} color="var(--text-muted)" style={{ marginBottom: 8, opacity: 0.4 }} />
                <p style={{ fontSize: "0.875rem", color: "var(--text-muted)" }}>No click data yet</p>
              </div>
            )}
          </div>

          {/* Devices */}
          <div style={cardStyle}>
            <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", marginBottom: "1.5rem" }}>
              <h2 className="font-display" style={{ fontSize: "0.9375rem", fontWeight: 700, color: "var(--text)", letterSpacing: "-0.01em" }}>Devices</h2>
              <Monitor size={16} color="var(--text-muted)" />
            </div>

            {device_breakdown.length > 0 ? (
              <div style={{ display: "flex", flexDirection: "column", gap: "1rem" }}>
                {device_breakdown.map((device, index) => (
                  <div key={device.device_type || index} style={{ display: "flex", alignItems: "center", gap: "0.75rem" }}>
                    <div style={{ display: "flex", alignItems: "center", gap: "0.5rem", width: 96, flexShrink: 0 }}>
                      {getDeviceIcon(device.device_type)}
                      <span style={{ fontSize: "0.8125rem", color: "var(--text-secondary)" }}>{formatDeviceType(device.device_type)}</span>
                    </div>
                    <div style={{ flex: 1, height: 4, background: "var(--bg-elevated)", borderRadius: 999, overflow: "hidden" }}>
                      <div style={{ height: "100%", background: "var(--accent-bright)", borderRadius: 999, width: `${device.percentage}%`, boxShadow: "0 0 8px var(--accent-glow)" }} />
                    </div>
                    <span style={{ fontSize: "0.8125rem", fontWeight: 600, color: "var(--text)", width: 36, textAlign: "right", flexShrink: 0 }}>
                      {device.percentage.toFixed(0)}%
                    </span>
                  </div>
                ))}
              </div>
            ) : (
              <div style={{ display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center", height: 180 }}>
                <Smartphone size={28} color="var(--text-muted)" style={{ marginBottom: 8, opacity: 0.4 }} />
                <p style={{ fontSize: "0.875rem", color: "var(--text-muted)" }}>No device data yet</p>
              </div>
            )}
          </div>
        </div>

        {/* Lists row */}
        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "1rem" }}>
          {/* Top URLs */}
          <div style={cardStyle}>
            <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", marginBottom: "1rem" }}>
              <h2 className="font-display" style={{ fontSize: "0.9375rem", fontWeight: 700, color: "var(--text)", letterSpacing: "-0.01em" }}>Top URLs</h2>
              <Eye size={16} color="var(--text-muted)" />
            </div>

            {top_urls && top_urls.length > 0 ? (
              <div>
                {top_urls.slice(0, 5).map((url, index) => (
                  <div
                    key={url.url_id}
                    style={{ display: "flex", alignItems: "center", gap: "0.75rem", padding: "0.75rem 0", borderTop: index > 0 ? "1px solid var(--border)" : "none", transition: "background 0.15s", cursor: "default" }}
                    onMouseEnter={e => (e.currentTarget as HTMLDivElement).style.background = "var(--bg-elevated)"}
                    onMouseLeave={e => (e.currentTarget as HTMLDivElement).style.background = "transparent"}
                  >
                    <span style={{ fontSize: "0.75rem", color: "var(--text-muted)", fontFamily: "DM Mono, monospace", width: 16, flexShrink: 0 }}>{index + 1}</span>
                    <div style={{ minWidth: 0, flex: 1 }}>
                      <p style={{ fontSize: "0.875rem", fontWeight: 600, color: "var(--accent-bright)", fontFamily: "DM Mono, monospace", margin: 0 }}>/{url.short_code}</p>
                      <p style={{ fontSize: "0.75rem", color: "var(--text-muted)", margin: 0, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{url.original_url}</p>
                    </div>
                    <span style={{ fontSize: "0.875rem", fontWeight: 700, color: "var(--text)", flexShrink: 0, fontFamily: "DM Mono, monospace" }}>
                      {url.click_count.toLocaleString()}
                    </span>
                  </div>
                ))}
              </div>
            ) : (
              <div style={{ display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center", paddingTop: 40, paddingBottom: 40 }}>
                <Link size={28} color="var(--text-muted)" style={{ marginBottom: 8, opacity: 0.4 }} />
                <p style={{ fontSize: "0.875rem", color: "var(--text-muted)" }}>No URLs yet</p>
              </div>
            )}
          </div>

          {/* Traffic sources */}
          <div style={cardStyle}>
            <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", marginBottom: "1.5rem" }}>
              <h2 className="font-display" style={{ fontSize: "0.9375rem", fontWeight: 700, color: "var(--text)", letterSpacing: "-0.01em" }}>Traffic Sources</h2>
              <Globe size={16} color="var(--text-muted)" />
            </div>

            {top_referrers && top_referrers.length > 0 ? (
              <div style={{ display: "flex", flexDirection: "column", gap: "1rem" }}>
                {top_referrers.map((referrer, index) => {
                  const maxClicks = Math.max(...top_referrers.map((r) => r.clicks));
                  const pct = maxClicks > 0 ? (referrer.clicks / maxClicks) * 100 : 0;
                  return (
                    <div key={index} style={{ display: "flex", alignItems: "center", gap: "0.75rem" }}>
                      <span style={{ fontSize: "0.8125rem", color: "var(--text-secondary)", width: 112, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap", flexShrink: 0 }}>
                        {referrer.referrer === "direct" ? "Direct" : referrer.referrer}
                      </span>
                      <div style={{ flex: 1, height: 4, background: "var(--bg-elevated)", borderRadius: 999, overflow: "hidden" }}>
                        <div style={{ height: "100%", background: "var(--cyan)", borderRadius: 999, width: `${pct}%`, opacity: 0.8 }} />
                      </div>
                      <span style={{ fontSize: "0.8125rem", fontWeight: 600, color: "var(--text)", width: 40, textAlign: "right", flexShrink: 0, fontFamily: "DM Mono, monospace" }}>
                        {referrer.clicks.toLocaleString()}
                      </span>
                    </div>
                  );
                })}
              </div>
            ) : (
              <div style={{ display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center", paddingTop: 40, paddingBottom: 40 }}>
                <Globe size={28} color="var(--text-muted)" style={{ marginBottom: 8, opacity: 0.4 }} />
                <p style={{ fontSize: "0.875rem", color: "var(--text-muted)" }}>No referrer data yet</p>
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  );
};

export default AnalyticsDashboard;
