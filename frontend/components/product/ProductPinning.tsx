"use client";

import { useState } from "react";

export default function ProductPinning() {
    const [pinned, setPinned] = useState(false);

    return (
        <button
            type="button"
            onClick={() => setPinned((value) => !value)}
            className={`rounded-[24px] px-4 py-3 text-sm font-semibold transition ${pinned ? "bg-[#E91E63] text-white" : "border border-zinc-200 bg-white text-zinc-900 hover:bg-[#FFF8FB]"}`}
        >
            {pinned ? "Pinned to livestream" : "Pin product to stream"}
        </button>
    );
}
