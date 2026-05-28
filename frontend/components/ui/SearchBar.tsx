"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

export default function SearchBar({ placeholder = "Search..." }: { placeholder?: string }) {
    const [query, setQuery] = useState("");
    const router = useRouter();

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (!query.trim()) return;
        router.push(`/search?query=${encodeURIComponent(query)}`);
    };

    return (
        <form onSubmit={handleSubmit} className="flex w-full items-center gap-2">
            <input
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                placeholder={placeholder}
                className="w-full rounded-3xl border border-white/6 bg-white/5 px-4 py-3 text-sm text-white placeholder:text-slate-300 focus:outline-none"
            />
            <button type="submit" className="rounded-3xl bg-fuchsia-500 px-4 py-2 text-sm font-semibold text-white">Search</button>
        </form>
    );
}
