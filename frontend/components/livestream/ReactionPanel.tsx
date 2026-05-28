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
        <div className="rounded-[28px] border border-zinc-200 bg-white p-5 sm:p-6">
            <h3 className="text-lg font-semibold text-zinc-900">Reactions</h3>
            <p className="mt-2 text-sm text-zinc-600">Send an in-stream reaction while viewing the product.</p>
            <div className="mt-4 flex flex-wrap gap-3">
                {reactions.map((reaction) => (
                    <button
                        type="button"
                        key={reaction.label}
                        onClick={() => setLast(reaction.label)}
                        className="rounded-[24px] border border-zinc-200 bg-[#FFF8FB] px-4 py-3 text-sm font-semibold text-zinc-700 transition hover:bg-[#FCE4EC]"
                    >
                        {reaction.emoji} {reaction.label}
                    </button>
                ))}
            </div>
            <p className="mt-4 text-sm text-zinc-500">
                Last reaction: <span className="font-semibold text-zinc-900">{last}</span>
            </p>
        </div>
    );
}
