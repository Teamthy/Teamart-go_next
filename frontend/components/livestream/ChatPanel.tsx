"use client";

import { useState } from "react";

const initialMessages = [
    { id: 1, user: "Maya", text: "Love this drop!" },
    { id: 2, user: "Jordan", text: "Is this available in black?" },
];

export default function ChatPanel() {
    const [messages, setMessages] = useState(initialMessages);
    const [draft, setDraft] = useState("");

    const handleSend = () => {
        if (!draft.trim()) return;
        setMessages((current) => [...current, { id: current.length + 1, user: "You", text: draft.trim() }]);
        setDraft("");
    };

    return (
        <section className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <h3 className="text-lg font-semibold text-slate-900">Realtime chat</h3>
            <div className="mt-4 space-y-3">
                {messages.map((message) => (
                    <div key={message.id} className="rounded-3xl bg-slate-50 p-4">
                        <p className="text-sm font-semibold text-slate-900">{message.user}</p>
                        <p className="mt-1 text-sm text-slate-600">{message.text}</p>
                    </div>
                ))}
            </div>
            <div className="mt-4 flex gap-3">
                <input
                    value={draft}
                    onChange={(event) => setDraft(event.target.value)}
                    placeholder="Write a message"
                    className="flex-1 rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none focus:border-slate-400 focus:ring-2 focus:ring-slate-200"
                />
                <button onClick={handleSend} className="rounded-2xl bg-slate-900 px-4 py-3 text-sm font-semibold text-white transition hover:bg-slate-700">
                    Send
                </button>
            </div>
        </section>
    );
}
