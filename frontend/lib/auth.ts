export type AuthResponse = {
    accessToken?: string;
    refreshToken?: string;
    requiresMfa?: boolean;
    sessionId?: string;
};

export function getStoredSession() {
    const data = typeof window !== "undefined" ? sessionStorage.getItem("session") : null;
    return data ? (JSON.parse(data) as AuthResponse) : null;
}

export function clearSession() {
    if (typeof window !== "undefined") {
        sessionStorage.removeItem("session");
        sessionStorage.removeItem("pendingSession");
    }
}
