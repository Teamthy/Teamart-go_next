import { ReactNode } from "react";

type BadgeVariant =
  | "success"
  | "warning"
  | "danger"
  | "primary"
  | "secondary"
  | "default";

interface BadgeProps {
  children: ReactNode;
  variant?: BadgeVariant;
  className?: string;
}

const variantStyles: Record<BadgeVariant, string> = {
  success:
    "bg-[var(--success-lighter)] text-[var(--success)] border border-[var(--success-light)]",
  warning:
    "bg-[var(--warning-lighter)] text-[var(--warning)] border border-[var(--warning-light)]",
  danger:
    "bg-[var(--danger-lighter)] text-[var(--danger)] border border-[var(--danger-light)]",
  primary:
    "bg-[var(--primary-lighter)] text-[var(--primary)] border border-[var(--primary-light)]",
  secondary:
    "bg-[var(--bg-lighter)] text-[var(--text-secondary)] border border-[var(--border-light)]",
  default:
    "bg-[var(--bg-lighter)] text-[var(--text-primary)] border border-[var(--border-light)]",
};

export default function Badge({
  children,
  variant = "default",
  className = "",
}: BadgeProps) {
  return (
    <span
      className={`inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-medium ${variantStyles[variant]} ${className}`}
    >
      {children}
    </span>
  );
}
