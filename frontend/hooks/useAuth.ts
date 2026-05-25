"use client";

import { useEffect, useState } from "react";
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

export function useAuth() {
    const [state, setState] = useState<AuthState>({
        user: null,
        isLoading: false,
        error: null,
        isAuthenticated: false,
    });

    // Initialize from localStorage
    useEffect(() => {
        const token = localStorage.getItem("access_token");
        const userStr = localStorage.getItem("user");

        if (token && userStr) {
            try {
                const user = JSON.parse(userStr);
                setState({
                    user,
                    isLoading: false,
                    error: null,
                    isAuthenticated: true,
                });
            } catch (e) {
                localStorage.removeItem("user");
                localStorage.removeItem("access_token");
            }
        }
    }, []);

    const login = async (email: string, password: string) => {
        setState((prev) => ({ ...prev, isLoading: true, error: null }));
        try {
            const response = await api.login(email, password);

            localStorage.setItem("access_token", response.access_token);
            localStorage.setItem("refresh_token", response.refresh_token);
            localStorage.setItem("user", JSON.stringify(response.user));

            setState({
                user: response.user,
                isLoading: false,
                error: null,
                isAuthenticated: true,
            });

            return response;
        } catch (err: any) {
            const error = err.message || "Login failed";
            setState((prev) => ({ ...prev, isLoading: false, error }));
            throw err;
        }
    };

    const signup = async (email: string, password: string) => {
        setState((prev) => ({ ...prev, isLoading: true, error: null }));
        try {
            const response = await api.signup(email, password);

            localStorage.setItem("access_token", response.access_token);
            localStorage.setItem("refresh_token", response.refresh_token);
            localStorage.setItem("user", JSON.stringify(response.user));

            setState({
                user: response.user,
                isLoading: false,
                error: null,
                isAuthenticated: true,
            });

            return response;
        } catch (err: any) {
            const error = err.message || "Signup failed";
            setState((prev) => ({ ...prev, isLoading: false, error }));
            throw err;
        }
    };

    const logout = () => {
        localStorage.removeItem("access_token");
        localStorage.removeItem("refresh_token");
        localStorage.removeItem("user");
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

            localStorage.setItem("access_token", response.access_token);
            localStorage.setItem("refresh_token", response.refresh_token);
            localStorage.setItem("user", JSON.stringify(response.user));

            setState({
                user: response.user,
                isLoading: false,
                error: null,
                isAuthenticated: true,
            });

            return response;
        } catch (err: any) {
            const error = err.message || "OTP verification failed";
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
