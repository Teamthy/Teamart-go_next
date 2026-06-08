import { ReactNode } from "react";
import { TrendingUp, TrendingDown } from "lucide-react";

interface MetricCardProps {
  label: string;
  value: string | number;
  subtext?: string;
  trend?: {
    value: number;
    direction: "up" | "down";
    period: string;
  };
  icon?: ReactNode;
  onClick?: () => void;
  className?: string;
}

export default function MetricCard({
  label,
  value,
  subtext,
  trend,
  icon,
  onClick,
  className = "",
}: MetricCardProps) {
  return (
    <div
      onClick={onClick}
      className={`bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6 hover:shadow-md transition-shadow ${
        onClick ? "cursor-pointer" : ""
      } ${className}`}
    >
      <div className="flex items-start justify-between mb-4">
        <div>
          <p className="text-[var(--text-tertiary)] text-sm font-medium mb-1">
            {label}
          </p>
          <div className="flex items-baseline gap-2">
            <p className="text-3xl font-bold text-[var(--text-primary)]">
              {value}
            </p>
            {subtext && (
              <p className="text-sm text-[var(--text-tertiary)]">{subtext}</p>
            )}
          </div>
        </div>
        {icon && (
          <div className="w-10 h-10 bg-[var(--primary-lighter)] rounded-lg flex items-center justify-center text-[var(--primary)]">
            {icon}
          </div>
        )}
      </div>

      {trend && (
        <div className="flex items-center gap-1 text-sm">
          {trend.direction === "up" ? (
            <>
              <TrendingUp className="w-4 h-4 text-[var(--success)]" />
              <span className="text-[var(--success)] font-medium">
                +{trend.value}%
              </span>
            </>
          ) : (
            <>
              <TrendingDown className="w-4 h-4 text-[var(--danger)]" />
              <span className="text-[var(--danger)] font-medium">
                -{trend.value}%
              </span>
            </>
          )}
          <span className="text-[var(--text-tertiary)]">
            vs {trend.period}
          </span>
        </div>
      )}
    </div>
  );
}
