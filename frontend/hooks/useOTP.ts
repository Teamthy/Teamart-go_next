/**
 * useOTP Hook
 * Manages OTP input logic with auto-advance, paste support, and cleanup
 */

import { useCallback, useRef, useState } from "react";

interface UseOTPProps {
    length?: number;
    onComplete: (code: string) => void;
    onError?: (error: string) => void;
}

interface UseOTPReturn {
    values: string[];
    setValues: (values: string[]) => void;
    handleInputChange: (index: number, value: string) => void;
    handleKeyDown: (index: number, e: React.KeyboardEvent) => void;
    handlePaste: (e: React.ClipboardEvent) => void;
    handleFocus: (index: number) => void;
    isFull: boolean;
    clear: () => void;
    inputRefs: React.MutableRefObject<(HTMLInputElement | null)[]>;
}

export function useOTP({ length = 6, onComplete, onError }: UseOTPProps): UseOTPReturn {
    const [values, setValues] = useState<string[]>(Array(length).fill(""));
    const inputRefs = useRef<(HTMLInputElement | null)[]>(Array(length).fill(null));

    const isFull = values.every((v) => v !== "");

    const handleInputChange = useCallback(
        (index: number, value: string) => {
            // Only allow digits
            const sanitized = value.replace(/[^0-9]/g, "");

            if (sanitized.length > 1) {
                // Handle paste of multiple digits
                const newValues = [...values];
                for (let i = 0; i < sanitized.length && index + i < length; i++) {
                    newValues[index + i] = sanitized[i];
                }
                setValues(newValues);

                // Auto-focus next empty field or last field
                const nextEmptyIndex = newValues.findIndex((v, i) => i >= index && v === "");
                const focusIndex = nextEmptyIndex === -1 ? length - 1 : nextEmptyIndex;
                setTimeout(() => inputRefs.current[focusIndex]?.focus(), 0);

                // Check if complete
                if (newValues.every((v) => v !== "")) {
                    onComplete(newValues.join(""));
                }
            } else {
                const newValues = [...values];
                newValues[index] = sanitized;
                setValues(newValues);

                // Auto-advance to next field
                if (sanitized && index < length - 1) {
                    setTimeout(() => inputRefs.current[index + 1]?.focus(), 50);
                }

                // Check if complete
                if (newValues.every((v) => v !== "")) {
                    onComplete(newValues.join(""));
                }
            }
        },
        [values, length, onComplete]
    );

    const handleKeyDown = useCallback(
        (index: number, e: React.KeyboardEvent) => {
            if (e.key === "Backspace" && !values[index] && index > 0) {
                // Move to previous field on backspace if empty
                setTimeout(() => inputRefs.current[index - 1]?.focus(), 0);
            }

            if (e.key === "ArrowLeft" && index > 0) {
                e.preventDefault();
                inputRefs.current[index - 1]?.focus();
            }

            if (e.key === "ArrowRight" && index < length - 1) {
                e.preventDefault();
                inputRefs.current[index + 1]?.focus();
            }
        },
        [values, length]
    );

    const handlePaste = useCallback(
        (e: React.ClipboardEvent) => {
            e.preventDefault();
            const pastedData = e.clipboardData.getData("text").replace(/[^0-9]/g, "");

            if (pastedData.length === length) {
                const newValues = pastedData.split("");
                setValues(newValues);
                onComplete(pastedData);
            } else if (pastedData.length > 0) {
                onError?.(`Please paste a ${length}-digit code`);
            }
        },
        [length, onComplete, onError]
    );

    const handleFocus = useCallback((index: number) => {
        setTimeout(() => {
            const input = inputRefs.current[index];
            if (input) {
                input.select();
            }
        }, 0);
    }, []);

    const clear = useCallback(() => {
        setValues(Array(length).fill(""));
        inputRefs.current[0]?.focus();
    }, [length]);

    return {
        values,
        setValues,
        handleInputChange,
        handleKeyDown,
        handlePaste,
        handleFocus,
        isFull,
        clear,
        inputRefs,
    };
}
