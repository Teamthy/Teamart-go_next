/**
 * Auth & Onboarding Type Definitions
 * Premium auth system for multi-role platform
 */

export type UserRole = "customer" | "creator" | "merchant" | "admin";
export type InterestCategory =
    | "fashion"
    | "beauty"
    | "electronics"
    | "gaming"
    | "food"
    | "fitness"
    | "home"
    | "tech";

export type CreatorCategory =
    | "beauty"
    | "fashion"
    | "gaming"
    | "fitness"
    | "lifestyle"
    | "tech";

export type MerchantCategory =
    | "fashion"
    | "electronics"
    | "beauty"
    | "home"
    | "food"
    | "sports"
    | "toys"
    | "other";

export interface AuthUser {
    id: number;
    email: string;
    firstName?: string;
    lastName?: string;
    username?: string;
    avatar?: string;
    role?: UserRole;
    emailVerified: boolean;
    createdAt: string;
}

export interface SignupPayload {
    email: string;
    password: string;
    role: UserRole;
}

export interface SignupResponse {
    sessionId: string;
    requiresOTP: boolean;
    expiresAt: string;
}

export interface LoginPayload {
    email: string;
    password: string;
}

export interface LoginResponse {
    success: boolean;
    token: string;
    refreshToken?: string;
    user: AuthUser;
}

export interface OTPVerificationPayload {
    sessionId: string;
    code: string;
}

export interface OTPVerificationResponse {
    success: boolean;
    token: string;
    refreshToken?: string;
    user: AuthUser;
}

export interface ProfileSetupPayload {
    firstName: string;
    lastName: string;
    username: string;
    profilePicture?: File;
    avatarUrl?: string;
}

export interface ProfileSetupResponse {
    success: boolean;
    user: AuthUser;
}

export interface CustomerOnboardingPayload {
    interests: InterestCategory[];
    followedCreators: number[];
    notificationPreferences: {
        orderUpdates: boolean;
        creatorAlerts: boolean;
        liveEvents: boolean;
        promotions: boolean;
    };
}

export interface CreatorOnboardingPayload {
    creatorName: string;
    bio: string;
    category: CreatorCategory;
    socialAccounts?: {
        tiktok?: string;
        instagram?: string;
        youtube?: string;
        twitter?: string;
    };
    goals: ("affiliate" | "livestream" | "partnerships" | "reviews")[];
    livestreamPreferences: {
        defaultCategory: string;
        displayName: string;
    };
    taxId?: string;
    bankAccount?: {
        bankName: string;
        accountNumber: string;
        routingNumber: string;
    };
}

export interface MerchantOnboardingPayload {
    businessName: string;
    businessType: string;
    country: string;
    phone: string;
    storeName: string;
    storeHandle: string;
    storeLogo?: File;
    storeBanner?: File;
    storeDescription: string;
    businessDocuments: {
        cacCertificate: File;
        governmentId: File;
        utilityBill: File;
    };
    bankAccount: {
        bankName: string;
        accountHolderName: string;
        accountNumber: string;
        routingNumber: string;
        accountType: "checking" | "savings";
    };
    categories: {
        primary: MerchantCategory;
        secondary: MerchantCategory[];
    };
    firstProduct?: {
        name: string;
        description: string;
        price: number;
        stock: number;
        images: File[];
    };
}

export interface AuthState {
    user: AuthUser | null;
    isAuthenticated: boolean;
    isLoading: boolean;
    error: string | null;
    token?: string;
}

export interface OnboardingState {
    currentStep: number;
    totalSteps: number;
    role: UserRole;
    isComplete: boolean;
    startedAt: Date;
    completedAt?: Date;
}

// Username validation
export interface UsernameCheckPayload {
    username: string;
}

export interface UsernameCheckResponse {
    available: boolean;
    suggestions?: string[];
}

// Email validation
export interface EmailCheckPayload {
    email: string;
}

export interface EmailCheckResponse {
    available: boolean;
    valid: boolean;
}

// OTP request
export interface OTPRequestPayload {
    sessionId: string;
}

export interface OTPRequestResponse {
    success: boolean;
    expiresAt: string;
}
