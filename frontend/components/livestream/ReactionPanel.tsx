"use client";

import { useState } from "react";

const reactions = [
    { emoji: "❤️", label: "Love" },
    { emoji: "🔥", label: "Hot" },
    { emoji: "👏", label: "Clap" },
    { emoji: "😲", label: "Wow" },
];

export default function ReactionPanel() {
    const [last, setLast] = useState("None yet");

    return (
        <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <h3 className="text-lg font-semibold text-slate-900">Reactions</h3>
            <p className="mt-2 text-sm text-slate-600">Send an in-stream reaction while viewing the product.</p>
            <div className="mt-4 flex flex-wrap gap-3">
                {reactions.map((reaction) => (
                    <button
                        type="button"
                        key={reaction.label}
                        onClick={() => setLast(reaction.label)}
                        className="rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm font-semibold transition hover:bg-slate-100"
                    >
                        {reaction.emoji} {reaction.label}
                    </button>
                ))}
            </div>
            <p className="mt-4 text-sm text-slate-500">Last reaction: <span className="font-semibold text-slate-900">{last}</span></p>
        </div>
    );
}
