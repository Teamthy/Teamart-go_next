"use client";

import { useState } from "react";
import * as api from "@/lib/api";

export interface User {
    id: number;
    email: string;
    name?: string;
    role?: string;
    created_at?: string;
}

export interface AuthState {
    user: User | null;
    isLoading: boolean;
    error: string | null;
    isAuthenticated: boolean;
}

type AuthResponse = {
    user?: User;
    session_id?: string;
    access_token?: string;
    refresh_token?: string;
    requires_mfa?: boolean;
    requiresMFA?: boolean;
};

function getErrorMessage(error: unknown, fallback: string) {
    if (error instanceof Error) return error.message;
    if (typeof error === "string") return error;
    return fallback;
}

function buildInitialAuthState(): AuthState {
    if (typeof window === "undefined") {
        return {
            user: null,
            isLoading: false,
            error: null,
            isAuthenticated: false,
        };
    }

    const userStr = localStorage.getItem("user");
    const sessionId = localStorage.getItem("session_id");

    if (userStr) {
        try {
            const parsed = JSON.parse(userStr);
            if (parsed && typeof parsed === "object") {
                return {
                    user: parsed as User,
                    isLoading: false,
                    error: null,
                    isAuthenticated: true,
                };
            }
        } catch {
            localStorage.removeItem("user");
            localStorage.removeItem("session_id");
            localStorage.removeItem("access_token");
            localStorage.removeItem("refresh_token");
        }
    }

    return {
        user: null,
        isLoading: false,
        error: null,
        isAuthenticated: Boolean(sessionId),
    };
}

function persistAuthState(response: AuthResponse | null) {
    if (typeof window === "undefined" || response === null) {
        return response?.user ?? null;
    }

    const user = response.user ?? null;
    const sessionId = response.session_id;

    if (user) {
        localStorage.setItem("user", JSON.stringify(user));
    }

    if (sessionId) {
        localStorage.setItem("session_id", sessionId);
    }

    if (response.access_token) {
        localStorage.setItem("access_token", response.access_token);
    }

    if (response.refresh_token) {
        localStorage.setItem("refresh_token", response.refresh_token);
    }

    localStorage.setItem("session", JSON.stringify(response));

    return user;
}

export function useAuth() {
    const [state, setState] = useState<AuthState>(buildInitialAuthState());

    const login = async (email: string, password: string) => {
        setState((prev) => ({ ...prev, isLoading: true, error: null }));
        try {
            const response = await api.login(email, password);
            const user = persistAuthState(response);

            setState({
                user,
                isLoading: false,
                error: null,
                isAuthenticated: Boolean(user),
            });

            return response;
        } catch (err: unknown) {
            const error = getErrorMessage(err, "Login failed");
            setState((prev) => ({ ...prev, isLoading: false, error }));
            throw err;
        }
    };

    const signup = async (email: string, password: string) => {
        setState((prev) => ({ ...prev, isLoading: true, error: null }));
        try {
            const response = await api.signup(email, password);
            const user = persistAuthState(response);

            setState({
                user,
                isLoading: false,
                error: null,
                isAuthenticated: Boolean(user),
            });

            return response;
        } catch (err: unknown) {
            const error = getErrorMessage(err, "Signup failed");
            setState((prev) => ({ ...prev, isLoading: false, error }));
            throw err;
        }
    };

    const logout = () => {
        if (typeof window !== "undefined") {
            localStorage.removeItem("access_token");
            localStorage.removeItem("refresh_token");
            localStorage.removeItem("session_id");
            localStorage.removeItem("session");
            localStorage.removeItem("user");
        }

        setState({
            user: null,
            isLoading: false,
            error: null,
            isAuthenticated: false,
        });
    };

    const verifyOTP = async (session_id: string, code: string) => {
        setState((prev) => ({ ...prev, isLoading: true, error: null }));
        try {
            const response = await api.verifyOTP(session_id, code);
            const user = persistAuthState(response);

            setState({
                user,
                isLoading: false,
                error: null,
                isAuthenticated: Boolean(user),
            });

            return response;
        } catch (err: unknown) {
            const error = getErrorMessage(err, "OTP verification failed");
            setState((prev) => ({ ...prev, isLoading: false, error }));
            throw err;
        }
    };

    return {
        ...state,
        login,
        signup,
        logout,
        verifyOTP,
    };
}
