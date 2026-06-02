"use client";

import { useMemo, useState } from "react";
import type { LiveChatMessage } from "@/types/commerce";
import { Paperclip } from "lucide-react";

interface LiveChatProps {
    messages: LiveChatMessage[];
    onSendMessage: (message: string) => void;
    isSending: boolean;
}

export default function LiveChat({ messages, onSendMessage, isSending }: LiveChatProps) {
    const [draft, setDraft] = useState("");
    const sortedMessages = useMemo(
        () => [...messages].sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()),
        [messages]
    );

    return (
        <section className="flex h-full flex-col rounded-[32px] border border-slate-200/80 bg-white shadow-lg">
            <div className="rounded-t-[32px] bg-slate-950 px-5 py-4 text-white">
                <p className="text-sm font-semibold">Live chat</p>
                <p className="mt-2 text-xs text-slate-300">Fast reactions, product questions, and creator responses in real time.</p>
            </div>
            <div className="flex-1 space-y-4 overflow-y-auto p-4">
                {sortedMessages.map((message) => (
                    <div key={message.id} className="flex gap-3 rounded-3xl bg-slate-100 p-3">
                        <img src={message.avatar_url} alt={message.user_name} className="h-10 w-10 rounded-full object-cover" />
                        <div className="min-w-0">
                            <div className="flex items-center gap-2 text-sm font-semibold text-slate-900">
                                <span>{message.user_name}</span>
                                {message.emoji ? <span className="text-base">{message.emoji}</span> : null}
                            </div>
                            <p className="mt-1 text-sm text-slate-700">{message.message}</p>
                            <p className="mt-1 text-xs text-slate-500">{new Date(message.timestamp).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}</p>
                        </div>
                    </div>
                ))}
            </div>
            <form
                className="rounded-b-[32px] bg-slate-50 p-4"
                onSubmit={(event) => {
                    event.preventDefault();
                    const trimmed = draft.trim();
                    if (!trimmed) return;
                    onSendMessage(trimmed);
                    setDraft("");
                }}
            >
                <div className="flex gap-3">
                    <input
                        value={draft}
                        onChange={(event) => setDraft(event.target.value)}
                        placeholder="Send a message to the room"
                        className="flex-1 rounded-full border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 outline-none transition focus:border-pink-400"
                    />
                    <button
                        type="submit"
                        disabled={isSending || !draft.trim()}
                        className="inline-flex h-12 items-center justify-center rounded-full bg-pink-500 px-5 text-sm font-semibold text-white transition hover:bg-pink-600 disabled:cursor-not-allowed disabled:bg-pink-300"
                    >
                        <Paperclip className="h-4 w-4" />
                    </button>
                </div>
            </form>
        </section>
    );
}
