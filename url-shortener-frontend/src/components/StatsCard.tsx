import React from "react";

interface StatCardProps {
  title: string;
  value: string | number;
  leftIcon: React.ReactNode;
  rightIcon: React.ReactNode;
}

export const StatCard: React.FC<StatCardProps> = ({
  title,
  value,
  leftIcon,
  rightIcon,
}) => {
  return (
    <div className="rounded-3xl p-6 bg-white/50 backdrop-blur-md border border-white/20 shadow-xl flex flex-col items-center justify-center space-y-4 text-center min-h-[160px]">
      <div className="flex justify-between items-center w-full px-4">
        <div className="w-10 h-10 flex items-center justify-center bg-white/30 rounded-xl text-xl">
          {leftIcon}
        </div>
        <div className="text-gray-600 text-lg font-semibold">{title}</div>
        <div className="w-10 h-10 flex items-center justify-center bg-white/30 rounded-xl text-xl">
          {rightIcon}
        </div>
      </div>
      <div className="text-4xl font-extrabold text-gray-900">{value}</div>
    </div>
  );
};
