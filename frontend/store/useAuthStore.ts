/**
 * Unified Auth Store (Zustand)
 * Manages authentication state across the application
 */

import { create } from "zustand";
import { persist } from "zustand/middleware";
import * as api from "@/lib/api";
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
    verifyOTP: (code: string) => Promise<boolean>;
    completeProfile: (firstName: string, lastName: string, username: string) => Promise<void>;

    // Utilities
    isLoggedIn: () => boolean;
    canAccess: (role: UserRole) => boolean;
    getUser: () => AuthUser | null;
}

const initialState: AuthState = {
    user: null,
    isAuthenticated: false,
    isLoading: false,
    error: null,
};

function isBrowser() {
    return typeof window !== "undefined";
}

function normalizeUserPayload(payload: Record<string, unknown> | null | undefined): AuthUser | null {
    if (!payload || typeof payload !== "object") {
        return null;
    }

    const id = Number(payload.id ?? payload.user_id ?? 0);
    const email = String(payload.email ?? payload.email_address ?? "");

    if (!id || !email) {
        return null;
    }

    const rawName = String(payload.name ?? payload.fullName ?? payload.full_name ?? "");
    const [firstNameFromName, ...rest] = rawName.trim().split(" ");
    const lastNameFromName = rest.join(" ");

    return {
        id,
        email,
        firstName: String(payload.firstName ?? payload.first_name ?? firstNameFromName ?? ""),
        lastName: String(payload.lastName ?? payload.last_name ?? lastNameFromName ?? ""),
        username: String(payload.username ?? payload.user_name ?? ""),
        avatar: String(payload.avatar ?? payload.avatar_url ?? ""),
        role: (payload.role as UserRole | undefined) ?? "customer",
        emailVerified: Boolean(payload.emailVerified ?? payload.email_verified ?? true),
        createdAt: String(payload.createdAt ?? payload.created_at ?? new Date().toISOString()),
    };
}

function getStoredSessionId(): string | null {
    if (!isBrowser()) {
        return null;
    }

    const pending = sessionStorage.getItem("pendingSession");
    if (pending) {
        try {
            const parsed = JSON.parse(pending);
            return String(parsed?.session_id ?? parsed?.sessionID ?? "") || null;
        } catch {
            return null;
        }
    }

    return localStorage.getItem("session_id");
}

function persistAuthResponse(response: {
    user?: Record<string, unknown> | null;
    session_id?: string;
    access_token?: string;
    refresh_token?: string;
} | null) {
    if (!isBrowser() || response === null) {
        return null;
    }

    const user = normalizeUserPayload(response.user ?? response);
    const { session_id, access_token, refresh_token } = response;

    if (user) {
        localStorage.setItem("user", JSON.stringify(user));
    }

    if (session_id) {
        localStorage.setItem("session_id", session_id);
        sessionStorage.setItem("pendingSession", JSON.stringify(response));
    }

    if (access_token) {
        localStorage.setItem("access_token", access_token);
    }

    if (refresh_token) {
        localStorage.setItem("refresh_token", refresh_token);
    }

    localStorage.setItem("session", JSON.stringify(response));

    return user;
}

function clearPersistedAuth() {
    if (!isBrowser()) {
        return;
    }

    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
    localStorage.removeItem("session_id");
    localStorage.removeItem("session");
    localStorage.removeItem("user");
    sessionStorage.removeItem("pendingSession");
}

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
            logout: () => {
                clearPersistedAuth();
                set(initialState);
            },
            reset: () => {
                clearPersistedAuth();
                set(initialState);
            },

            // Auth flow
            login: async (email, password) => {
                set({ isLoading: true, error: null });
                try {
                    const response = await api.login(email, password);
                    const user = persistAuthResponse(response);
                    const requiresMfa = Boolean(response.requires_mfa ?? response.requiresMFA);

                    if (requiresMfa) {
                        set({ user: null, token: response.access_token, isAuthenticated: false, error: null });
                        return;
                    }

                    set({ user, token: response.access_token, isAuthenticated: Boolean(user), error: null });
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
                    const response = await api.signup(email, password, role);
                    const user = persistAuthResponse(response);
                    set({ user, token: response.access_token, isAuthenticated: false, error: null });
                } catch (error) {
                    set({
                        error: error instanceof Error ? error.message : "Signup failed",
                    });
                    throw error;
                } finally {
                    set({ isLoading: false });
                }
            },

            verifyOTP: async (code) => {
                set({ isLoading: true, error: null });
                try {
                    const sessionId = getStoredSessionId();
                    if (!sessionId) {
                        throw new Error("Verification session not found");
                    }

                    let pendingRole: string | undefined;
                    let pendingEmail: string | undefined;
                    if (isBrowser()) {
                        const raw = sessionStorage.getItem("pendingSession");
                        if (raw) {
                            try {
                                const parsed = JSON.parse(raw) as Record<string, unknown>;
                                pendingRole = String(parsed.role ?? parsed.user?.role ?? "");
                                pendingEmail = String(parsed.email ?? parsed.user?.email ?? "");
                            } catch {
                                pendingRole = undefined;
                                pendingEmail = undefined;
                            }
                        }
                    }

                    const response = await api.verifyOTP(sessionId, code, pendingRole, pendingEmail);
                    const user = persistAuthResponse(response);
                    const authenticated = Boolean(user);
                    set({ user, token: response.access_token, isAuthenticated: authenticated, error: null });

                    if (isBrowser()) {
                        sessionStorage.removeItem("pendingSession");
                    }

                    return authenticated;
                } catch (error) {
                    set({
                        error: error instanceof Error ? error.message : "OTP verification failed",
                    });
                    return false;
                } finally {
                    set({ isLoading: false });
                }
            },

            completeProfile: async (firstName, lastName, username) => {
                set({ isLoading: true, error: null });
                try {
                    const currentUser = get().user;
                    if (!currentUser) {
                        throw new Error("No authenticated user available");
                    }

                    const response = await api.updateUser(currentUser.id, {
                        first_name: firstName,
                        last_name: lastName,
                        username,
                    });

                    const normalized = normalizeUserPayload(response as Record<string, unknown>) ?? {
                        ...currentUser,
                        firstName,
                        lastName,
                        username,
                    };

                    localStorage.setItem("user", JSON.stringify(normalized));
                    set({ user: normalized });
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
            isLoggedIn: () => get().isAuthenticated,
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
