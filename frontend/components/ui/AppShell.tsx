import Link from "next/link";
import type { ReactNode } from "react";
import NotificationBell from "@/components/ui/NotificationBell";

const navItems = [
    { href: "/", label: "Home" },
    { href: "/seller", label: "Seller Dashboard" },
    { href: "/auth/login", label: "Login" },
    { href: "/auth/register", label: "Register" },
    { href: "/feed", label: "Feed" },
    { href: "/search", label: "Search" },
    { href: "/cart", label: "Cart" },
    { href: "/checkout", label: "Checkout" },
    { href: "/livestream/status", label: "Livestream" },
] as const;

export default function AppShell({ children }: { children: ReactNode }) {
    return (
        <div className="min-h-screen bg-slate-50 text-slate-900">
            <header className="border-b border-slate-200 bg-white/90 backdrop-blur-sm">
                <div className="mx-auto flex max-w-7xl items-center justify-between px-4 py-4 sm:px-6 lg:px-8">
                    <Link href="/" className="text-lg font-semibold tracking-tight text-slate-900">
                        Teamart Customer Web
                    </Link>
                    <div className="hidden items-center gap-3 md:flex">
                        <nav className="flex items-center gap-3">
                            {navItems.map((item) => (
                                <Link
                                    key={item.href}
                                    href={item.href}
                                    className="rounded-full px-3 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-100 hover:text-slate-900"
                                >
                                    {item.label}
                                </Link>
                            ))}
                        </nav>
                        <NotificationBell />
                    </div>
                </div>
            </header>
            <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">{children}</main>
            <footer className="border-t border-slate-200 bg-white/90">
                <div className="mx-auto max-w-7xl px-4 py-6 text-sm text-slate-500 sm:px-6 lg:px-8">
                    Built for customer web, powered by Teamart design tokens and product-led commerce flows.
                </div>
            </footer>
        </div>
    );
}
