import React from "react";

interface StatCardProps {
  title: string;
  value: string | number;
  icon: React.ReactNode;
}

export const StatCard: React.FC<StatCardProps> = ({ title, value, icon }) => {
  return (
    <div
      className="card-dark"
      style={{ padding: "1.25rem", transition: "border-color 0.2s, transform 0.2s" }}
      onMouseEnter={e => {
        (e.currentTarget as HTMLDivElement).style.borderColor = "var(--border-bright)";
        (e.currentTarget as HTMLDivElement).style.transform = "translateY(-2px)";
      }}
      onMouseLeave={e => {
        (e.currentTarget as HTMLDivElement).style.borderColor = "var(--border)";
        (e.currentTarget as HTMLDivElement).style.transform = "translateY(0)";
      }}
    >
      <div style={{ width: 36, height: 36, borderRadius: 8, background: "var(--accent-dim)", border: "1px solid rgba(124,58,237,0.2)", display: "flex", alignItems: "center", justifyContent: "center" }}>
        {icon}
      </div>
      <p style={{ fontSize: "0.75rem", fontWeight: 600, color: "var(--text-muted)", marginTop: "0.875rem", letterSpacing: "0.05em", textTransform: "uppercase" }}>{title}</p>
      <p className="font-display" style={{ fontSize: "1.875rem", fontWeight: 800, color: "var(--text)", marginTop: "0.25rem", letterSpacing: "-0.02em" }}>{value}</p>
    </div>
  );
};
