import type { ReactNode, ButtonHTMLAttributes } from "react";

export default function SocialBtn({
    icon,
    label,
    className = "",
    ...props
}: ButtonHTMLAttributes<HTMLButtonElement> & {
    icon: ReactNode;
    label: string;
}) {
    return (
        <button
            type="button"
            {...props}
            className={"flex h-11 w-full items-center justify-center gap-2 rounded-2xl border border-slate-200 bg-white text-sm font-medium text-slate-700 transition hover:bg-slate-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#E91E63]/20 " + className}
        >
            {icon}
            {label}
        </button>
    );
}
