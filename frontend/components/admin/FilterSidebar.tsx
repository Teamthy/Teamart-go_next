import type { ReactNode } from "react";

interface FilterSidebarProps {
    title: string;
    description: string;
    controls: ReactNode;
}

export default function FilterSidebar({ title, description, controls }: FilterSidebarProps) {
    return (
        <aside className="rounded-[24px] border border-slate-200 bg-white p-5 shadow-[0_18px_50px_-28px_rgba(15,23,42,0.2)]">
            <div className="space-y-2">
                <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Filters</p>
                <h2 className="text-lg font-semibold text-slate-900">{title}</h2>
                <p className="text-sm text-slate-600">{description}</p>
            </div>
            <div className="mt-5 space-y-4">{controls}</div>
        </aside>
    );
}
