"use client";

import Link from "next/link";
import useAppStore from "@/store/useAppStore";

export default function NotificationBell() {
    const count = useAppStore((state) => state.notifications.length);
    const status = useAppStore((state) => state.socketStatus);

    return (
        <Link
            href="/feed"
            className="relative inline-flex items-center rounded-full border border-slate-200 bg-white px-3 py-2 text-sm font-medium text-slate-700 shadow-sm transition hover:bg-slate-50"
        >
            <span>Live</span>
            <span className="ml-2 text-xs text-slate-500">{status === "open" ? "connected" : status}</span>
            {count > 0 ? (
                <span className="absolute -right-2 -top-2 inline-flex h-5 min-w-[1.25rem] items-center justify-center rounded-full bg-rose-500 px-1.5 text-[0.65rem] font-semibold text-white">
                    {count}
                </span>
            ) : null}
        </Link>
    );
}
