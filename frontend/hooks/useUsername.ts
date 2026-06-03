/**
 * useUsername Hook
 * Real-time username validation with debounce and suggestions
 */

import { useCallback, useEffect, useRef, useState } from "react";

interface UseUsernameProps {
    onCheck?: (available: boolean) => void;
    debounceMs?: number;
}

interface UseUsernameReturn {
    username: string;
    setUsername: (username: string) => void;
    isAvailable: boolean | null;
    isChecking: boolean;
    error: string | null;
    suggestions: string[];
    isValid: boolean;
    validate: () => boolean;
}

const VALID_USERNAME_REGEX = /^[a-zA-Z][a-zA-Z0-9_]{2,29}$/;
const MIN_LENGTH = 3;
const MAX_LENGTH = 30;

// Simulated API call - replace with real API
async function checkUsernameAvailability(username: string): Promise<{ available: boolean; suggestions?: string[] }> {
    // In real implementation, call backend:
    // const res = await fetch(`/api/auth/check-username?username=${username}`);
    // return res.json();

    // Simulated response
    return new Promise((resolve) => {
        setTimeout(() => {
            const isAvailable = Math.random() > 0.3; // 70% available
            resolve({
                available: isAvailable,
                suggestions: !isAvailable ? [username + "2", username + "_pro", username + "official"] : undefined,
            });
        }, 300);
    });
}

export function useUsername({ onCheck, debounceMs = 500 }: UseUsernameProps = {}): UseUsernameReturn {
    const [username, setUsername] = useState("");
    const [isAvailable, setIsAvailable] = useState<boolean | null>(null);
    const [isChecking, setIsChecking] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [suggestions, setSuggestions] = useState<string[]>([]);
    const debounceTimer = useRef<NodeJS.Timeout | null>(null);

    // Validate format
    const isValid = VALID_USERNAME_REGEX.test(username);

    // Check availability with debounce
    useEffect(() => {
        if (!username) {
            setIsAvailable(null);
            setError(null);
            setSuggestions([]);
            return;
        }

        // Clear previous timer
        if (debounceTimer.current) {
            clearTimeout(debounceTimer.current);
        }

        // Validate format first
        if (username.length < MIN_LENGTH) {
            setError(`Username must be at least ${MIN_LENGTH} characters`);
            setIsAvailable(null);
            return;
        }

        if (username.length > MAX_LENGTH) {
            setError(`Username must be at most ${MAX_LENGTH} characters`);
            setIsAvailable(null);
            return;
        }

        if (!VALID_USERNAME_REGEX.test(username)) {
            setError("Username can only contain letters, numbers, and underscores. Must start with a letter.");
            setIsAvailable(null);
            return;
        }

        // Valid format, check availability after debounce
        setError(null);
        setIsChecking(true);

        debounceTimer.current = setTimeout(async () => {
            try {
                const result = await checkUsernameAvailability(username);
                setIsAvailable(result.available);
                setSuggestions(result.suggestions || []);
                onCheck?.(result.available);
            } catch (err) {
                setError("Failed to check username availability");
                setIsAvailable(null);
            } finally {
                setIsChecking(false);
            }
        }, debounceMs);

        return () => {
            if (debounceTimer.current) {
                clearTimeout(debounceTimer.current);
            }
        };
    }, [username, debounceMs, onCheck]);

    const validate = useCallback(() => {
        if (!isValid) {
            setError("Invalid username format");
            return false;
        }

        if (isAvailable === false) {
            setError("Username is not available. Try a suggestion.");
            return false;
        }

        if (isAvailable === null) {
            setError("Please check username availability");
            return false;
        }

        return true;
    }, [isValid, isAvailable]);

    return {
        username,
        setUsername,
        isAvailable,
        isChecking,
        error,
        suggestions,
        isValid,
        validate,
    };
}
