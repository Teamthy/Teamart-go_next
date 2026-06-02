"use client";

import { useMemo, useState } from "react";
import type { FeedProduct } from "@/types/commerce";

export function usePinnedProducts(products: FeedProduct[]) {
    const [activeIndex, setActiveIndex] = useState(0);

    const normalized = useMemo(() => products || [], [products]);

    const selected = useMemo(() => normalized[activeIndex] ?? normalized[0] ?? null, [normalized, activeIndex]);

    const selectNext = () => {
        if (!normalized.length) return;
        setActiveIndex((value) => (value + 1) % normalized.length);
    };

    const selectPrevious = () => {
        if (!normalized.length) return;
        setActiveIndex((value) => (value - 1 + normalized.length) % normalized.length);
    };

    return {
        products: normalized,
        selected,
        activeIndex,
        selectNext,
        selectPrevious,
        setActiveIndex,
    };
}
