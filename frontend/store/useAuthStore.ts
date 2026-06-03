/**
 * Unified Auth Store (Zustand)
 * Manages authentication state across the application
 */

import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { AuthUser, UserRole, AuthState } from "@/types/auth.types";

interface AuthStore extends AuthState {
    // Actions
    setUser: (user: AuthUser | null) => void;
    setToken: (token: string) => void;
    setLoading: (loading: boolean) => void;
    setError: (error: string | null) => void;
    logout: () => void;
    reset: () => void;

    // Auth flow
    login: (email: string, password: string) => Promise<void>;
    signup: (email: string, password: string, role: UserRole) => Promise<void>;
    verifyOTP: (sessionId: string, code: string) => Promise<void>;
    completeProfile: (firstName: string, lastName: string, username: string) => Promise<void>;

    // Utilities
    isAuthenticated: () => boolean;
    canAccess: (role: UserRole) => boolean;
    getUser: () => AuthUser | null;
}

const initialState: AuthState = {
    user: null,
    isAuthenticated: false,
    isLoading: false,
    error: null,
};

export const useAuthStore = create<AuthStore>()(
    persist(
        (set, get) => ({
            ...initialState,

            // Basic setters
            setUser: (user) => set({ user, isAuthenticated: Boolean(user) }),
            setToken: (token) => set({ token }),
            setLoading: (loading) => set({ isLoading: loading }),
            setError: (error) => set({ error }),

            // Logout & reset
            logout: () => set(initialState),
            reset: () => set(initialState),

            // Auth flow (placeholder - implement API calls)
            login: async (email, password) => {
                set({ isLoading: true, error: null });
                try {
                    // API call would happen here
                    // const user = await api.login(email, password);
                    // set({ user, isAuthenticated: true });
                } catch (error) {
                    set({
                        error: error instanceof Error ? error.message : "Login failed",
                    });
                    throw error;
                } finally {
                    set({ isLoading: false });
                }
            },

            signup: async (email, password, role) => {
                set({ isLoading: true, error: null });
                try {
                    // API call would happen here
                    // const response = await api.signup(email, password, role);
                } catch (error) {
                    set({
                        error: error instanceof Error ? error.message : "Signup failed",
                    });
                    throw error;
                } finally {
                    set({ isLoading: false });
                }
            },

            verifyOTP: async (sessionId, code) => {
                set({ isLoading: true, error: null });
                try {
                    // API call would happen here
                    // const user = await api.verifyOTP(sessionId, code);
                    // set({ user, isAuthenticated: true });
                } catch (error) {
                    set({
                        error: error instanceof Error ? error.message : "OTP verification failed",
                    });
                    throw error;
                } finally {
                    set({ isLoading: false });
                }
            },

            completeProfile: async (firstName, lastName, username) => {
                set({ isLoading: true, error: null });
                try {
                    // API call would happen here
                    // const user = await api.completeProfile({ firstName, lastName, username });
                    // set({ user });
                } catch (error) {
                    set({
                        error: error instanceof Error ? error.message : "Profile update failed",
                    });
                    throw error;
                } finally {
                    set({ isLoading: false });
                }
            },

            // Utilities
            isAuthenticated: () => get().isAuthenticated,
            canAccess: (role) => {
                const user = get().user;
                if (!user) return false;
                return user.role === role || user.role === "admin";
            },
            getUser: () => get().user,
        }),
        {
            name: "auth-store",
            partialize: (state) => ({
                user: state.user,
                token: state.token,
                isAuthenticated: state.isAuthenticated,
            }),
        }
    )
);
