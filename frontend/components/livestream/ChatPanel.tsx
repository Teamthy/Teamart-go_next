"use client";

import { useState } from "react";

const initialMessages = [
    { id: 1, user: "Maya", text: "Love this drop! so soft 😍", time: "just now" },
    { id: 2, user: "Jordan", text: "Is this available in black?", time: "1m ago" },
    { id: 3, user: "Sage", text: "The fit is perfect for layering.", time: "2m ago" },
];

export default function ChatPanel() {
    const [messages, setMessages] = useState(initialMessages);
    const [draft, setDraft] = useState("");

    const handleSend = () => {
        if (!draft.trim()) return;
        setMessages((current) => [
            ...current,
            { id: current.length + 1, user: "You", text: draft.trim(), time: "now" },
        ]);
        setDraft("");
    };

    return (
        <section className="rounded-[28px] border border-zinc-200 bg-white p-5 sm:p-6">
            <div className="flex items-center justify-between">
                <div>
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Comments</p>
                    <h3 className="mt-2 text-lg font-semibold text-zinc-900">Live chat</h3>
                </div>
                <span className="rounded-full bg-[#FCE4EC] px-3 py-1 text-[11px] font-semibold text-[#E91E63]">3 new</span>
            </div>

            <div className="mt-4 space-y-3">
                {messages.map((message) => (
                    <div key={message.id} className="rounded-[24px] bg-[#FFF8FB] p-4">
                        <div className="flex items-center justify-between gap-3">
                            <p className="text-sm font-semibold text-zinc-900">{message.user}</p>
                            <p className="text-[11px] text-zinc-400">{message.time}</p>
                        </div>
                        <p className="mt-2 text-sm leading-6 text-zinc-600">{message.text}</p>
                    </div>
                ))}
            </div>

            <div className="mt-4 flex flex-col gap-3 sm:flex-row">
                <input
                    value={draft}
                    onChange={(event) => setDraft(event.target.value)}
                    placeholder="Write a message"
                    className="flex-1 rounded-[24px] border border-zinc-200 bg-white px-4 py-3 text-sm outline-none transition placeholder:text-zinc-400 focus:border-[#E91E63] focus:ring-2 focus:ring-[#E91E63]/10"
                />
                <button
                    onClick={handleSend}
                    className="rounded-[24px] bg-[#E91E63] px-4 py-3 text-sm font-semibold text-white"
                >
                    Send
                </button>
            </div>
        </section>
    );
}
