import Link from "next/link";
import { Bell, MessageSquare, Search, Settings, User } from "lucide-react";

export default function TopBar() {
  return (
    <header className="fixed top-0 left-0 right-0 h-16 bg-[var(--bg-white)] border-b border-[var(--border-light)] z-40">
      <div className="h-full px-6 flex items-center justify-between">
        {/* Left: Logo */}
        <Link href="/seller/dashboard" className="flex items-center gap-2">
          <div className="w-8 h-8 bg-[var(--primary)] rounded-lg flex items-center justify-center">
            <span className="text-[var(--bg-white)] font-bold text-lg">T</span>
          </div>
          <span className="font-semibold text-[var(--text-primary)] text-lg">
            Teamart Seller Center
          </span>
        </Link>

        {/* Center: Search */}
        <div className="flex-1 max-w-md mx-8">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--text-tertiary)]" />
            <input
              type="text"
              placeholder="Search products, orders, creators..."
              className="w-full pl-10 pr-4 py-2 bg-[var(--bg-lighter)] border border-[var(--border-light)] rounded-lg text-sm text-[var(--text-primary)] placeholder-[var(--text-tertiary)] focus:outline-none focus:border-[var(--primary)]"
            />
          </div>
        </div>

        {/* Right: Notifications, Messages, Settings, Profile */}
        <div className="flex items-center gap-4">
          {/* Notifications */}
          <button className="relative p-2 hover:bg-[var(--bg-lighter)] rounded-lg transition-colors">
            <Bell className="w-5 h-5 text-[var(--text-secondary)]" />
            <span className="absolute top-1 right-1 w-2 h-2 bg-[var(--danger)] rounded-full" />
          </button>

          {/* Messages */}
          <button className="relative p-2 hover:bg-[var(--bg-lighter)] rounded-lg transition-colors">
            <MessageSquare className="w-5 h-5 text-[var(--text-secondary)]" />
            <span className="absolute top-1 right-1 w-2 h-2 bg-[var(--danger)] rounded-full" />
          </button>

          {/* Settings */}
          <button className="p-2 hover:bg-[var(--bg-lighter)] rounded-lg transition-colors">
            <Settings className="w-5 h-5 text-[var(--text-secondary)]" />
          </button>

          {/* Profile Menu */}
          <div className="flex items-center gap-3 pl-4 border-l border-[var(--border-light)]">
            <div className="w-8 h-8 bg-gradient-to-br from-[var(--primary)] to-[var(--primary-hover)] rounded-full flex items-center justify-center">
              <span className="text-[var(--bg-white)] text-sm font-bold">U</span>
            </div>
            <button className="text-[var(--text-secondary)] hover:text-[var(--text-primary)] transition-colors">
              <svg
                className="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M19 9l-7 7-7-7"
                />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </header>
  );
}
