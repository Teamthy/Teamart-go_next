"use client";

import { useState } from "react";

export default function ProductPinning() {
    const [pinned, setPinned] = useState(false);

    return (
        <button
            type="button"
            onClick={() => setPinned((value) => !value)}
            className={`rounded-3xl px-4 py-3 text-sm font-semibold transition ${pinned ? "bg-slate-900 text-white" : "border border-slate-200 bg-white text-slate-900 hover:bg-slate-100"
                }`}
        >
            {pinned ? "Pinned to livestream" : "Pin product to stream"}
        </button>
    );
}
