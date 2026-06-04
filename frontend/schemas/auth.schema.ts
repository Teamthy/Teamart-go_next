/**
 * Zod Validation Schemas
 * Auth & Onboarding form validation
 */

import { z } from "zod";

// Email validation schema
export const emailSchema = z.string().email("Invalid email address").max(254);

// Password schema (strong)
export const passwordSchema = z
    .string()
    .min(12, "Password must be at least 12 characters")
    .regex(/[A-Z]/, "Password must contain at least one uppercase letter")
    .regex(/[a-z]/, "Password must contain at least one lowercase letter")
    .regex(/[0-9]/, "Password must contain at least one number")
    .regex(/[^A-Za-z0-9]/, "Password must contain at least one special character");

// OTP code (6 digits)
export const otpSchema = z.string().length(6, "OTP must be 6 digits").regex(/^\d{6}$/, "OTP must be numeric");

// Username (3-30 chars, alphanumeric + underscore)
export const usernameSchema = z
    .string()
    .min(3, "Username must be at least 3 characters")
    .max(30, "Username must be at most 30 characters")
    .regex(/^[a-zA-Z0-9_]+$/, "Username can only contain letters, numbers, and underscores")
    .regex(/^[a-zA-Z]/, "Username must start with a letter");

// Name validation
export const nameSchema = z.string().min(2, "Name must be at least 2 characters").max(50, "Name must be at most 50 characters");

// Signup payload
export const signupSchema = z.object({
    email: emailSchema,
    password: passwordSchema,
    confirmPassword: z.string(),
    role: z.enum(["customer", "creator", "merchant"]),
}).refine((data) => data.password === data.confirmPassword, {
    message: "Passwords do not match",
    path: ["confirmPassword"],
});

// OTP verification
export const otpVerificationSchema = z.object({
    sessionId: z.string().min(1, "Session ID required"),
    code: otpSchema,
});

// Profile setup
export const profileSetupSchema = z.object({
    firstName: nameSchema,
    lastName: nameSchema,
    username: usernameSchema,
    profilePicture: z
        .instanceof(File)
        .optional()
        .refine((file) => !file || file.size <= 5 * 1024 * 1024, "File must be less than 5MB")
        .refine((file) => !file || ["image/jpeg", "image/png", "image/webp"].includes(file.type), "File must be JPEG, PNG, or WebP"),
});

// Customer onboarding
export const customerOnboardingSchema = z.object({
    interests: z.array(z.enum(["fashion", "beauty", "electronics", "gaming", "food", "fitness", "home", "tech"])).min(1, "Select at least one interest"),
    followedCreators: z.array(z.number()).optional(),
    notificationPreferences: z.object({
        orderUpdates: z.boolean(),
        creatorAlerts: z.boolean(),
        liveEvents: z.boolean(),
        promotions: z.boolean(),
    }),
});

// Creator onboarding
export const creatorOnboardingSchema = z.object({
    creatorName: nameSchema,
    bio: z.string().max(250, "Bio must be at most 250 characters"),
    category: z.enum(["beauty", "fashion", "gaming", "fitness", "lifestyle", "tech"]),
    socialAccounts: z
        .object({
            tiktok: z.string().optional(),
            instagram: z.string().optional(),
            youtube: z.string().optional(),
            twitter: z.string().optional(),
        })
        .optional(),
    goals: z.array(z.enum(["affiliate", "livestream", "partnerships", "reviews"])).min(1),
    livestreamPreferences: z.object({
        defaultCategory: z.string(),
        displayName: z.string(),
    }),
    taxId: z.string().optional(),
    bankAccount: z
        .object({
            bankName: z.string(),
            accountNumber: z.string(),
            routingNumber: z.string(),
        })
        .optional(),
});

// Merchant onboarding
export const businessInfoSchema = z.object({
    businessName: z.string().min(2, "Business name required").max(100),
    businessType: z.enum(["sole_proprietor", "llc", "corporation", "partnership"]),
    country: z.string().min(2, "Country required"),
    phone: z.string().regex(/^\+?[0-9\s\-()]{7,}$/, "Invalid phone number"),
});

export const storeSetupSchema = z.object({
    storeName: z.string().min(2, "Store name required").max(100),
    storeHandle: z.string().regex(/^[a-z0-9\-]{3,50}$/, "Invalid store handle"),
    storeLogo: z.instanceof(File).optional().refine((file) => !file || file.size <= 2 * 1024 * 1024, "Logo must be less than 2MB"),
    storeBanner: z.instanceof(File).optional().refine((file) => !file || file.size <= 5 * 1024 * 1024, "Banner must be less than 5MB"),
    storeDescription: z.string().max(500, "Description must be at most 500 characters"),
});

export const bankAccountSchema = z.object({
    bankName: z.string().min(1, "Bank name required"),
    accountHolderName: z.string().min(2),
    accountNumber: z.string().regex(/^\d{8,17}$/, "Invalid account number"),
    routingNumber: z.string().regex(/^\d{9}$/, "Invalid routing number"),
    accountType: z.enum(["checking", "savings"]),
});

export const merchantCategoriesSchema = z.object({
    primary: z.enum(["fashion", "electronics", "beauty", "home", "food", "sports", "toys", "other"]),
    secondary: z.array(z.enum(["fashion", "electronics", "beauty", "home", "food", "sports", "toys", "other"])).max(3),
});

export const firstProductSchema = z.object({
    name: z.string().min(1, "Product name required").max(200),
    description: z.string().max(2000),
    price: z.number().positive("Price must be positive"),
    stock: z.number().int().nonnegative("Stock must be non-negative"),
    images: z.array(z.instanceof(File)).min(1, "At least one image required"),
});

export const merchantOnboardingSchema = z.object({
    businessInfo: businessInfoSchema,
    storeSetup: storeSetupSchema,
    bankAccount: bankAccountSchema,
    categories: merchantCategoriesSchema,
    firstProduct: firstProductSchema.optional(),
});

// Login schema
export const loginSchema = z.object({
    email: emailSchema,
    password: z.string().min(1, "Password required"),
    remember: z.boolean().optional(),
});

// Password reset request
export const passwordResetRequestSchema = z.object({
    email: emailSchema,
});

// Password reset
export const passwordResetSchema = z
    .object({
        code: z.string().min(1, "Reset code required"),
        password: passwordSchema,
        confirmPassword: z.string(),
    })
    .refine((data) => data.password === data.confirmPassword, {
        message: "Passwords do not match",
        path: ["confirmPassword"],
    });

// Type exports for form usage
export type SignupForm = z.infer<typeof signupSchema>;
export type OTPVerificationForm = z.infer<typeof otpVerificationSchema>;
export type ProfileSetupForm = z.infer<typeof profileSetupSchema>;
export type CustomerOnboardingForm = z.infer<typeof customerOnboardingSchema>;
export type CreatorOnboardingForm = z.infer<typeof creatorOnboardingSchema>;
export type MerchantOnboardingForm = z.infer<typeof merchantOnboardingSchema>;
export type LoginForm = z.infer<typeof loginSchema>;
export type PasswordResetRequestForm = z.infer<typeof passwordResetRequestSchema>;
export type PasswordResetForm = z.infer<typeof passwordResetSchema>;
