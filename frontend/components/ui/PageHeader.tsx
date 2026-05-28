import type { ReactNode } from "react";
import Badge from "@/components/ui/badge";
import Button from "@/components/ui/button";

interface PageHeaderProps {
    title: string;
    description: string;
    eyebrow?: string;
    actions?: ReactNode;
}

export default function PageHeader({ title, description, eyebrow, actions }: PageHeaderProps) {
    return (
        <section className="rounded-[32px] bg-white p-5 sm:p-6">
            <div className="flex flex-col gap-5 lg:flex-row lg:items-end lg:justify-between">
                <div className="space-y-4">
                    {eyebrow ? <Badge tone="default">{eyebrow}</Badge> : null}
                    <div className="space-y-2">
                        <p className="text-[11px] uppercase tracking-[0.24em] text-[#E91E63]">Teamart-go_next</p>
                        <div>
                            <h1 className="text-[24px] font-semibold tracking-tight text-zinc-900 sm:text-[28px]">{title}</h1>
                            <p className="mt-2 max-w-2xl text-sm leading-6 text-zinc-600 sm:text-[15px]">{description}</p>
                        </div>
                    </div>
                </div>
                {actions ? <div className="flex flex-wrap gap-3">{actions}</div> : null}
            </div>
        </section>
    );
}
