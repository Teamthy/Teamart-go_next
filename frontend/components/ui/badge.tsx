import type { ReactNode } from "react";

interface BadgeProps {
    children: ReactNode;
    tone?: "default" | "success" | "warning" | "error" | "info";
}

const toneStyles = {
    default: "bg-[#FCE4EC] text-[#E91E63]",
    success: "bg-emerald-100 text-emerald-800",
    warning: "bg-amber-100 text-amber-800",
    error: "bg-rose-100 text-rose-800",
    info: "bg-sky-100 text-sky-800",
};

export default function Badge({ children, tone = "default" }: BadgeProps) {
    return (
        <span className={`inline-flex items-center rounded-full px-3 py-1 text-[11px] font-semibold ${toneStyles[tone]}`}>
            {children}
        </span>
    );
}
