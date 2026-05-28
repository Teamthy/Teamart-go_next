"use client";

import { useEffect, useMemo, useRef, useState } from "react";

interface OTPInputProps {
    value: string;
    onChange: (value: string) => void;
    length?: number;
}

export default function OTPInput({ value, onChange, length = 6 }: OTPInputProps) {
    const ids = useMemo(() => Array.from({ length }, (_, index) => index), [length]);
    const inputRefs = useRef<Array<HTMLInputElement | null>>([]);
    const [activeIndex, setActiveIndex] = useState(0);

    useEffect(() => {
        inputRefs.current[activeIndex]?.focus();
    }, [activeIndex]);

    const handleChange = (nextValue: string, index: number) => {
        const sanitized = nextValue.replace(/\D/g, "").slice(0, 1);
        const next = value.split("");
        next[index] = sanitized;
        onChange(next.join(""));

        if (sanitized && index < length - 1) {
            setActiveIndex(index + 1);
        }
    };

    const handleKeyDown = (event: React.KeyboardEvent<HTMLInputElement>, index: number) => {
        if (event.key === "Backspace" && !value[index] && index > 0) {
            setActiveIndex(index - 1);
        }
    };

    return (
        <div className="flex flex-wrap items-center justify-center gap-3">
            {ids.map((index) => (
                <input
                    key={index}
                    ref={(node) => {
                        inputRefs.current[index] = node;
                    }}
                    inputMode="numeric"
                    autoComplete="one-time-code"
                    maxLength={1}
                    value={value[index] ?? ""}
                    onFocus={() => setActiveIndex(index)}
                    onChange={(event) => handleChange(event.target.value, index)}
                    onKeyDown={(event) => handleKeyDown(event, index)}
                    className="h-14 w-12 rounded-[24px] border border-zinc-200 bg-white text-center text-lg font-semibold text-zinc-900 outline-none transition focus:border-[#E91E63] focus:ring-2 focus:ring-[#E91E63]/10"
                />
            ))}
        </div>
    );
}
