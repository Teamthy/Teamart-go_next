import type { ReactNode } from "react";

interface TabsProps {
    tabs: { label: string; value: string }[];
    active: string;
    onChange: (value: string) => void;
}

export default function Tabs({ tabs, active, onChange }: TabsProps) {
    return (
        <div className="flex flex-wrap gap-2">
            {tabs.map((tab) => (
                <button
                    key={tab.value}
                    type="button"
                    onClick={() => onChange(tab.value)}
                    className={`rounded-full px-4 py-2 text-sm font-semibold transition ${active === tab.value ? "bg-[#E91E63] text-white" : "bg-white text-zinc-700 border border-zinc-200"
                        }`}
                >
                    {tab.label}
                </button>
            ))}
        </div>
    );
}
