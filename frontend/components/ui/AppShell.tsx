import Link from "next/link";
import type { ReactNode } from "react";
import NotificationBell from "@/components/ui/NotificationBell";

const navItems = [
    { href: "/", label: "Home" },
    { href: "/feed", label: "Feed" },
    { href: "/search", label: "Search" },
    { href: "/livestream/status", label: "Live" },
    { href: "/auth/login", label: "Account" },
] as const;

export default function AppShell({ children }: { children: ReactNode }) {
    return (
        <div className="min-h-screen bg-transparent text-zinc-900">
            <header className="sticky top-0 z-30 border-b border-[#f3c0d4] bg-white/90 backdrop-blur">
                <div className="mx-auto flex max-w-7xl items-center justify-between px-4 py-3 sm:px-6 lg:px-8">
                    <Link href="/" className="flex items-center gap-3">
                        <span className="grid h-10 w-10 place-items-center rounded-[20px] bg-[#E91E63] text-sm font-bold text-white">
                            T
                        </span>
                        <div>
                            <p className="text-[11px] uppercase tracking-[0.24em] text-zinc-500">Teamart</p>
                            <p className="text-sm font-semibold text-zinc-900">Social Commerce</p>
                        </div>
                    </Link>

                    <div className="hidden items-center gap-2 md:flex">
                        {navItems.map((item) => (
                            <Link
                                key={item.href}
                                href={item.href}
                                className="rounded-full px-3 py-2 text-sm font-medium text-zinc-600 transition hover:bg-[#FCE4EC] hover:text-[#E91E63]"
                            >
                                {item.label}
                            </Link>
                        ))}
                    </div>

                    <div className="hidden md:block">
                        <NotificationBell />
                    </div>
                </div>
            </header>

            <main className="mx-auto w-full max-w-7xl px-4 pb-24 pt-4 sm:px-6 lg:px-8">{children}</main>

            <nav className="fixed inset-x-0 bottom-0 z-20 border-t border-zinc-200 bg-white/96 backdrop-blur sm:hidden">
                <div className="mx-auto flex max-w-md items-center justify-between px-3 py-2">
                    {navItems.map((item) => (
                        <Link
                            key={item.href}
                            href={item.href}
                            className="flex min-w-0 flex-1 flex-col items-center gap-1 rounded-[18px] px-2 py-2 text-[11px] font-semibold text-zinc-600"
                        >
                            <span>{item.label}</span>
                        </Link>
                    ))}
                </div>
            </nav>

            <footer className="border-t border-[#f3c0d4] bg-white/90">
                <div className="mx-auto max-w-7xl px-4 py-6 text-sm text-zinc-500 sm:px-6 lg:px-8">
                    Built for Teamart’s creator-led commerce experience.
                </div>
            </footer>
        </div>
    );
}
