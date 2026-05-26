export const BASE = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

type AuthUser = {
    id: number;
    email: string;
    name?: string;
    role?: string;
    created_at?: string;
};

type AuthSessionResponse = {
    user_id?: number;
    session_id?: string;
    email?: string;
    status?: string;
    created_at?: string;
    access_token?: string;
    refresh_token?: string;
    requires_mfa?: boolean;
    requiresMFA?: boolean;
    requires_password_verification?: boolean;
    requiresPassword?: boolean;
    user?: Partial<AuthUser>;
};

function getStoredAuthUser() {
    if (typeof window === "undefined") {
        return null;
    }

    try {
        const raw = localStorage.getItem("user");
        if (!raw) {
            return null;
        }

        const parsed = JSON.parse(raw);
        if (!parsed || typeof parsed !== "object") {
            return null;
        }

        return parsed as AuthUser;
    } catch {
        return null;
    }
}

function normalizeAuthUser(data: AuthSessionResponse): AuthUser {
    const fallbackUser = getStoredAuthUser();
    const email = data.email || fallbackUser?.email || "";

    return {
        id: Number(data.user_id ?? data.user?.id ?? fallbackUser?.id ?? 0),
        email,
        name: data.user?.name || fallbackUser?.name || email.split("@")[0] || "Customer",
        role: data.user?.role || fallbackUser?.role || "customer",
        created_at: data.created_at || data.user?.created_at || fallbackUser?.created_at,
    };
}

function normalizeAuthResponse(data: AuthSessionResponse) {
    return {
        ...data,
        user: normalizeAuthUser(data),
    };
}

// Rate limit tracking
let rateLimitRemaining = 1000;
let rateLimitReset: number | null = null;
let isRefreshingToken = false;
let refreshTokenPromise: Promise<string | null> | null = null;

// Exponential backoff retry logic
async function request(
    path: string,
    options: {
        method?: string;
        body?: any;
        params?: Record<string, any>;
        retries?: number;
        retryDelay?: number;
    } = {}
): Promise<any> {
    const { retries = 3, retryDelay = 1000, ...requestOptions } = options;

    for (let attempt = 0; attempt <= retries; attempt++) {
        try {
            return await rawRequest(path, requestOptions);
        } catch (error: any) {
            // Don't retry on 401 (auth error) or 403 (forbidden)
            if (error.status === 401 || error.status === 403) {
                throw error;
            }

            // Don't retry on 4xx errors except 408, 429, 503, 504
            if (error.status && error.status >= 400 && error.status < 500) {
                if (![408, 429].includes(error.status)) {
                    throw error;
                }
            }

            // Last attempt, throw error
            if (attempt === retries) {
                throw error;
            }

            const delay = retryDelay * Math.pow(2, attempt);
            const backoffDelay =
                error.status === 429 && typeof error.retryAfterMs === "number" && error.retryAfterMs > 0
                    ? error.retryAfterMs
                    : delay;

            console.warn(
                `[API] Request failed (attempt ${attempt + 1}/${retries + 1}), retrying in ${backoffDelay}ms`,
                error.message
            );

            await new Promise((resolve) => setTimeout(resolve, backoffDelay));
        }
    }
}

// Automatic token refresh on 401
async function refreshAccessToken(): Promise<string | null> {
    // Prevent multiple simultaneous refresh attempts
    if (isRefreshingToken && refreshTokenPromise) {
        return refreshTokenPromise;
    }

    isRefreshingToken = true;

    refreshTokenPromise = (async () => {
        try {
            if (typeof window === "undefined") return null;

            const refreshToken = localStorage.getItem("refresh_token");
            if (!refreshToken) {
                localStorage.removeItem("access_token");
                localStorage.removeItem("user");
                return null;
            }

            const res = await fetch(`${BASE}/auth/refresh`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ refresh_token: refreshToken }),
            });

            if (!res.ok) {
                // Refresh failed, clear auth
                localStorage.removeItem("access_token");
                localStorage.removeItem("refresh_token");
                localStorage.removeItem("user");
                return null;
            }

            const data = await res.json();
            localStorage.setItem("access_token", data.access_token);
            if (data.refresh_token) {
                localStorage.setItem("refresh_token", data.refresh_token);
            }

            return data.access_token;
        } catch (error) {
            console.error("[API] Token refresh failed:", error);
            localStorage.removeItem("access_token");
            localStorage.removeItem("refresh_token");
            localStorage.removeItem("user");
            return null;
        } finally {
            isRefreshingToken = false;
            refreshTokenPromise = null;
        }
    })();

    return refreshTokenPromise;
}

// Generic request handler with method support
async function rawRequest(
    path: string,
    options: {
        method?: string;
        body?: any;
        params?: Record<string, any>;
    } = {}
) {
    const { method = "POST", body, params } = options;

    let url = `${BASE}${path}`;

    // Add query parameters
    if (params) {
        const query = new URLSearchParams();
        Object.entries(params).forEach(([key, value]) => {
            if (value !== undefined && value !== null) {
                query.append(key, String(value));
            }
        });
        if (query.toString()) {
            url += `?${query.toString()}`;
        }
    }

    const fetchOptions: RequestInit = {
        method,
        headers: {
            "Content-Type": "application/json",
        },
    };

    // Add auth token if available
    if (typeof window !== "undefined") {
        const token = localStorage.getItem("access_token");
        if (token) {
            fetchOptions.headers = {
                ...fetchOptions.headers,
                Authorization: `Bearer ${token}`,
            };
        }
    }

    if (body) {
        fetchOptions.body = JSON.stringify(body);
    }

    let res = await fetch(url, fetchOptions);

    // Handle 401 Unauthorized - try to refresh token
    if (res.status === 401 && typeof window !== "undefined") {
        const newToken = await refreshAccessToken();
        if (newToken) {
            // Retry request with new token
            fetchOptions.headers = {
                ...fetchOptions.headers,
                Authorization: `Bearer ${newToken}`,
            };
            res = await fetch(url, fetchOptions);
        }
    }

    // Extract and track rate limit headers
    const rateLimitHeaders = {
        remaining: res.headers.get("x-ratelimit-remaining"),
        reset: res.headers.get("x-ratelimit-reset"),
        limit: res.headers.get("x-ratelimit-limit"),
    };

    if (rateLimitHeaders.remaining) {
        rateLimitRemaining = parseInt(rateLimitHeaders.remaining);
    }
    if (rateLimitHeaders.reset) {
        rateLimitReset = parseInt(rateLimitHeaders.reset) * 1000; // Convert to ms
    }

    // Warn if approaching rate limit
    if (rateLimitRemaining < 100) {
        console.warn(
            `[API] Approaching rate limit: ${rateLimitRemaining} requests remaining`
        );
    }

    const text = await res.text();

    let json: any = null;
    try {
        json = text ? JSON.parse(text) : null;
    } catch (e) {
        throw new Error(`Invalid JSON response: ${text}`);
    }

    if (!res.ok) {
        const msg = (json && json.message) || res.statusText || "Request failed";
        const err: any = new Error(msg);
        err.status = res.status;
        err.body = json;
        err.code = json?.code;

        const retryAfter = res.headers.get("retry-after");
        if (retryAfter) {
            const seconds = parseInt(retryAfter, 10);
            if (!Number.isNaN(seconds)) {
                err.retryAfterMs = seconds * 1000;
            }
        }

        const rateLimitResetHeader = res.headers.get("x-ratelimit-reset");
        if (res.status === 429 && rateLimitResetHeader) {
            const reset = parseInt(rateLimitResetHeader, 10);
            if (!Number.isNaN(reset)) {
                err.retryAfterMs = Math.max(0, reset * 1000 - Date.now());
            }
        }

        throw err;
    }

    return json;
}

// Export rate limit info
export function getRateLimitInfo() {
    return {
        remaining: rateLimitRemaining,
        reset: rateLimitReset ? new Date(rateLimitReset) : null,
    };
}

// =============== AUTH ENDPOINTS ===============

export async function login(email: string, password: string) {
    const response = await request("/auth/login", {
        method: "POST",
        body: { email, password, user_agent: navigator.userAgent, ip_address: "0.0.0.0" },
    });

    return normalizeAuthResponse(response);
}

export async function signup(email: string, password: string) {
    const response = await request("/auth/signup", {
        method: "POST",
        body: { email, password },
    });

    return normalizeAuthResponse(response);
}

export async function verifyOTP(session_id: string, code: string) {
    const response = await request("/sessions/validate", {
        method: "POST",
        body: {
            session_id,
            user_agent: navigator.userAgent,
            ip_address: "0.0.0.0",
        },
    });

    return normalizeAuthResponse({
        ...response,
        user_id: response.user_id,
        session_id: response.session_id || session_id,
        email: getStoredAuthUser()?.email || response.email,
    });
}

export async function refreshToken(refresh_token: string) {
    return request("/auth/refresh", {
        method: "POST",
        body: { refresh_token },
    });
}

// =============== USER ENDPOINTS ===============

export async function getUser(userId: number) {
    return request(`/users/${userId}`, { method: "GET" });
}

export async function listUsers(limit = 20, offset = 0) {
    return request("/users", {
        method: "GET",
        params: { limit, offset },
    });
}

export async function updateUser(userId: number, data: any) {
    return request(`/users/${userId}`, {
        method: "PUT",
        body: data,
    });
}

export async function deleteUser(userId: number) {
    return request(`/users/${userId}`, { method: "DELETE" });
}

// =============== PRODUCT ENDPOINTS ===============

export async function listProducts(limit = 20, offset = 0) {
    return request("/products", {
        method: "GET",
        params: { limit, offset },
    });
}

export async function getProduct(productId: number | string) {
    return request(`/products/${productId}`, { method: "GET" });
}

export async function getProductBySKU(sku: string) {
    return request(`/products/sku/${sku}`, { method: "GET" });
}

export async function searchProducts(query: string, limit = 20, offset = 0) {
    return request("/products/search", {
        method: "GET",
        params: { q: query, limit, offset },
    });
}

export async function createProduct(data: any) {
    return request("/products", {
        method: "POST",
        body: data,
    });
}

export async function updateProduct(productId: number, data: any) {
    return request(`/products/${productId}`, {
        method: "PUT",
        body: data,
    });
}

export async function deleteProduct(productId: number) {
    return request(`/products/${productId}`, { method: "DELETE" });
}

// =============== ORDER ENDPOINTS ===============

export async function createOrder(data: any) {
    const storedUser = getStoredAuthUser();
    const user_id = Number(data?.user_id ?? storedUser?.id ?? 0);

    if (!user_id) {
        throw new Error("Please sign in to place an order.");
    }

    const total_amount = Number(data?.total_amount ?? data?.totalAmount ?? 0);
    if (!Number.isFinite(total_amount) || total_amount <= 0) {
        throw new Error("Your cart is empty. Add products before checkout.");
    }

    return request("/orders", {
        method: "POST",
        body: {
            user_id,
            total_amount,
            status: data?.status || "pending",
        },
    });
}

export async function getOrder(orderId: number | string) {
    return request(`/orders/${orderId}`, { method: "GET" });
}

export async function listOrders(limit = 20, offset = 0) {
    return request("/orders", {
        method: "GET",
        params: { limit, offset },
    });
}

export async function listUserOrders(userId: number, limit = 20, offset = 0) {
    return request(`/users/${userId}/orders`, {
        method: "GET",
        params: { limit, offset },
    });
}

export async function listOrdersByStatus(status: string, limit = 20, offset = 0) {
    return request(`/orders/status/${status}`, {
        method: "GET",
        params: { limit, offset },
    });
}

export async function updateOrderStatus(orderId: number, status: string) {
    return request(`/orders/${orderId}`, {
        method: "PUT",
        body: { status },
    });
}

// =============== MERCHANT ENDPOINTS ===============

export async function createMerchant(data: any) {
    return request("/api/v1/merchants", {
        method: "POST",
        body: data,
    });
}

export async function getMerchant(merchantId: number) {
    return request(`/api/v1/merchants/${merchantId}`, { method: "GET" });
}

export async function createStore(merchantId: number, data: any) {
    return request(`/api/v1/merchants/${merchantId}/stores`, {
        method: "POST",
        body: data,
    });
}

export async function listStores(merchantId: number) {
    return request(`/api/v1/merchants/${merchantId}/stores`, { method: "GET" });
}

export async function addStaff(merchantId: number, data: any) {
    return request(`/api/v1/merchants/${merchantId}/staff`, {
        method: "POST",
        body: data,
    });
}

export async function listStaff(merchantId: number) {
    return request(`/api/v1/merchants/${merchantId}/staff`, { method: "GET" });
}

// =============== ADMIN ENDPOINTS ===============

export async function getAdminDashboard() {
    return request("/admin/dashboard", { method: "GET" });
}

export async function listDisputes() {
    return request("/admin/disputes", { method: "GET" });
}

export async function createDispute(data: any) {
    return request("/admin/disputes", {
        method: "POST",
        body: data,
    });
}

export async function listFraudAlerts() {
    return request("/admin/fraud/alerts", { method: "GET" });
}

export async function listPayouts() {
    return request("/admin/payouts", { method: "GET" });
}

export async function approvePayout(payoutId: string) {
    return request("/admin/payouts/approve", {
        method: "POST",
        params: { id: payoutId },
    });
}

export async function requestPayoutApproval(payoutId: string, requestedBy: number, notes?: string) {
    return request("/admin/payouts/request", {
        method: "POST",
        params: { id: payoutId, requested_by: requestedBy, notes },
    });
}

export async function verifyCreator(data: any) {
    return request("/admin/creators/verify", {
        method: "POST",
        body: data,
    });
}

export async function refund(disputeId: string) {
    return request("/admin/support/refund", {
        method: "POST",
        body: { dispute_id: disputeId },
    });
}

export async function suspendAccount(userId: string) {
    return request("/admin/support/suspend", {
        method: "POST",
        body: { user_id: userId },
    });
}

export async function listAuditLogs() {
    return request("/admin/audit/logs", { method: "GET" });
}

export async function listNotifications() {
    return request("/admin/notifications", { method: "GET" });
}

// =============== MODERATION ENDPOINTS ===============

export async function getUserModerationStatus(userId: number) {
    return request(`/api/v1/moderation/users/${userId}/status`, { method: "GET" });
}

export async function blockUser(userId: number) {
    return request(`/api/v1/moderation/users/${userId}/block`, {
        method: "POST",
    });
}

export async function shadowbanUser(userId: number) {
    return request(`/api/v1/moderation/users/${userId}/shadowban`, {
        method: "POST",
    });
}

export async function muteUser(userId: number) {
    return request(`/api/v1/moderation/users/${userId}/mute`, {
        method: "POST",
    });
}

// =============== ANALYTICS ENDPOINTS ===============

export async function ingestAnalyticsEvent(data: any) {
    return request("/api/v1/analytics/events", {
        method: "POST",
        body: data,
    });
}

export async function getAnalyticsMetrics() {
    return request("/api/v1/analytics/metrics", { method: "GET" });
}

export async function getCreatorMetrics() {
    return request("/api/v1/analytics/metrics/creator", { method: "GET" });
}

export async function getMarketplaceMetrics() {
    return request("/api/v1/analytics/metrics/marketplace", { method: "GET" });
}

// =============== FEED ENDPOINTS ===============

export async function getFeed(limit = 20) {
    return request("/feed", {
        method: "GET",
        params: { limit },
    });
}

export async function ingestRecommendationCandidate(data: any) {
    return request("/feed/candidates", {
        method: "POST",
        body: data,
    });
}

// =============== SESSION ENDPOINTS ===============

export async function getActiveSessions(userId: number) {
    return request(`/sessions/active/${userId}`, { method: "GET" });
}

export async function revokeSession(sessionId: string) {
    return request(`/sessions/${sessionId}/revoke`, {
        method: "POST",
    });
}

export async function revokeAllSessions(userId: number) {
    return request(`/sessions/user/${userId}/revoke-all`, {
        method: "POST",
    });
}
