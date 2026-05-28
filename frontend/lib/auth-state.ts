export type RoleKey = "customer" | "creator" | "merchant";

export interface UserRoles {
    customer: boolean;
    creator: boolean;
    merchant: boolean;
}

export interface CustomerProfile {
    id: string;
    email: string;
    firstName: string;
    lastName: string;
    verified: boolean;
    roles: UserRoles;
    favoriteCategory?: string;
    provider?: string;
}

export interface CustomerOnboardingState {
    profile: CustomerProfile;
    progress: string[];
}

export interface CustomerDraftState {
    firstName: string;
    lastName: string;
    email: string;
    favoriteCategory: string;
    role: RoleKey;
    verified: boolean;
}

export interface AuthSessionState {
    customer: CustomerProfile | null;
    onboarding: CustomerOnboardingState | null;
    activeRole: RoleKey;
    workspaceRole: RoleKey;
    isAuthenticated: boolean;
}

const STORAGE_KEYS = {
    user: "teamart_user",
    onboarding: "teamart_onboarding",
    role: "auth_role",
    workspace: "workspace_role",
    draft: "teamart_customer_draft",
};

function isBrowser() {
    return typeof window !== "undefined";
}

function readJson<T>(key: string): T | null {
    if (!isBrowser()) {
        return null;
    }

    try {
        const stored = localStorage.getItem(key);
        if (!stored) {
            return null;
        }

        return JSON.parse(stored) as T;
    } catch {
        return null;
    }
}

function getStoredRole(): RoleKey {
    if (!isBrowser()) {
        return "customer";
    }

    const stored = localStorage.getItem(STORAGE_KEYS.role);
    if (stored === "creator" || stored === "merchant" || stored === "customer") {
        return stored;
    }

    const workspace = localStorage.getItem(STORAGE_KEYS.workspace);
    if (workspace === "creator" || workspace === "merchant" || workspace === "customer") {
        return workspace;
    }

    return "customer";
}

export function getStoredCustomer(): CustomerProfile | null {
    return readJson<CustomerProfile>(STORAGE_KEYS.user);
}

export function getStoredOnboarding(): CustomerOnboardingState | null {
    return readJson<CustomerOnboardingState>(STORAGE_KEYS.onboarding);
}

export function getCustomerDraft(): CustomerDraftState | null {
    return readJson<CustomerDraftState>(STORAGE_KEYS.draft);
}

export function saveCustomerDraft(draft: Partial<CustomerDraftState>) {
    if (!isBrowser()) {
        return;
    }

    const current = getCustomerDraft();

    localStorage.setItem(
        STORAGE_KEYS.draft,
        JSON.stringify({
            firstName: current?.firstName ?? "",
            lastName: current?.lastName ?? "",
            email: current?.email ?? "",
            favoriteCategory: current?.favoriteCategory ?? "Fashion",
            role: current?.role ?? "customer",
            verified: current?.verified ?? false,
            ...draft,
        }),
    );
}

export function clearCustomerDraft() {
    if (!isBrowser()) {
        return;
    }

    localStorage.removeItem(STORAGE_KEYS.draft);
}

export function getWorkspaceRole(): RoleKey {
    return getStoredRole();
}

export function getAuthState(): AuthSessionState {
    const customer = getStoredCustomer();
    const onboarding = getStoredOnboarding();
    const workspaceRole = getWorkspaceRole();

    return {
        customer,
        onboarding,
        activeRole: workspaceRole,
        workspaceRole,
        isAuthenticated: Boolean(customer),
    };
}

export function saveCustomer(profile: CustomerProfile, progress: string[]) {
    if (!isBrowser()) {
        return;
    }

    const nextRole: RoleKey = profile.roles.creator
        ? "creator"
        : profile.roles.merchant
            ? "merchant"
            : "customer";

    localStorage.setItem(STORAGE_KEYS.user, JSON.stringify(profile));
    localStorage.setItem(STORAGE_KEYS.onboarding, JSON.stringify({ profile, progress }));
    localStorage.setItem(STORAGE_KEYS.role, nextRole);
    localStorage.setItem(STORAGE_KEYS.workspace, nextRole);
    clearCustomerDraft();
}

export function updateRole(role: RoleKey, enabled: boolean) {
    if (!isBrowser()) {
        return;
    }

    const current = getStoredCustomer();
    if (!current) {
        return;
    }

    const next = {
        ...current,
        roles: {
            ...current.roles,
            [role]: enabled,
        },
    };

    localStorage.setItem(STORAGE_KEYS.user, JSON.stringify(next));

    if (enabled) {
        setWorkspaceRole(role);
    }
}

export function setWorkspaceRole(role: RoleKey) {
    if (!isBrowser()) {
        return;
    }

    const current = getStoredCustomer();

    if (current) {
        localStorage.setItem(
            STORAGE_KEYS.user,
            JSON.stringify({
                ...current,
                roles: {
                    customer: true,
                    creator: role === "creator" || current.roles.creator,
                    merchant: role === "merchant" || current.roles.merchant,
                },
            }),
        );
    }

    localStorage.setItem(STORAGE_KEYS.role, role);
    localStorage.setItem(STORAGE_KEYS.workspace, role);
}

export function hasCustomerAccount() {
    return Boolean(getStoredCustomer());
}

export function hasRole(role: RoleKey) {
    const current = getStoredCustomer();
    if (!current) {
        return false;
    }

    return Boolean(current.roles[role]);
}

export function getRoleBadge() {
    const current = getStoredCustomer();
    if (!current) {
        return "Guest";
    }

    if (current.roles.creator && current.roles.merchant) {
        return "Creator + Merchant";
    }

    if (current.roles.creator) {
        return "Creator";
    }

    if (current.roles.merchant) {
        return "Merchant";
    }

    return "Customer";
}

export function logout() {
    if (!isBrowser()) {
        return;
    }

    localStorage.removeItem(STORAGE_KEYS.user);
    localStorage.removeItem(STORAGE_KEYS.onboarding);
    localStorage.removeItem(STORAGE_KEYS.role);
    localStorage.removeItem(STORAGE_KEYS.workspace);
    localStorage.removeItem(STORAGE_KEYS.draft);
}
