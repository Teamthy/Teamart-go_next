import type { ReactNode } from "react";

export function Skeleton({ className = "" }: { className?: string }) {
    return <div className={`animate-pulse rounded-3xl bg-slate-200/80 dark:bg-slate-700/60 ${className}`} />;
}

export function ProductGridSkeleton({ count = 6 }: { count?: number }) {
    return (
        <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
            {Array.from({ length: count }, (_, index) => (
                <div key={index} className="space-y-4 rounded-3xl border border-slate-200/60 bg-white/90 p-6 shadow-sm shadow-slate-200/10 dark:border-slate-700/60 dark:bg-slate-900/80">
                    <div className="h-56 rounded-3xl bg-slate-200/80 dark:bg-slate-700/60" />
                    <div className="space-y-3">
                        <div className="h-5 w-3/4 rounded-full bg-slate-200/80 dark:bg-slate-700/60" />
                        <div className="h-4 w-1/2 rounded-full bg-slate-200/80 dark:bg-slate-700/60" />
                        <div className="h-4 w-1/3 rounded-full bg-slate-200/80 dark:bg-slate-700/60" />
                    </div>
                </div>
            ))}
        </div>
    );
}

export function SmallSkeletonCard() {
    return (
        <div className="rounded-3xl border border-slate-200/60 bg-white/90 p-5 shadow-sm shadow-slate-200/10 dark:border-slate-700/60 dark:bg-slate-900/80">
            <div className="h-10 w-10 rounded-2xl bg-slate-200/80 dark:bg-slate-700/60 mb-4" />
            <div className="h-4 w-2/3 rounded-full bg-slate-200/80 dark:bg-slate-700/60 mb-2" />
            <div className="h-4 w-1/2 rounded-full bg-slate-200/80 dark:bg-slate-700/60" />
        </div>
    );
}
