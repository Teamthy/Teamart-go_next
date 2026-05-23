"use client";

import { useState } from "react";

export default function SearchBar() {
    const [query, setQuery] = useState("");

    return (
        <div className="rounded-3xl border border-slate-200 bg-white p-4 shadow-sm">
            <label className="mb-2 block text-sm font-medium text-slate-700">Search products</label>
            <div className="flex flex-col gap-3 sm:flex-row">
                <input
                    value={query}
                    onChange={(event) => setQuery(event.target.value)}
                    placeholder="Search by product, creator, tag..."
                    className="flex-1 rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm text-slate-900 outline-none transition focus:border-slate-400 focus:ring-2 focus:ring-slate-200"
                />
                <button className="inline-flex items-center justify-center rounded-2xl bg-slate-900 px-4 py-3 text-sm font-semibold text-white transition hover:bg-slate-700">
                    Search
                </button>
            </div>
            {query ? <p className="mt-3 text-sm text-slate-500">Searching for “{query}”...</p> : null}
        </div>
    );
}
