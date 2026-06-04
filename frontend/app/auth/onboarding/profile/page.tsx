/**
 * Profile Setup Page
 * Set first name, last name, username, and profile picture
 */

"use client";

import React, { useState, useRef } from "react";
import { useRouter } from "next/navigation";
import { Controller } from "react-hook-form";
import { useFormValidation } from "@/hooks/useFormValidation";
import { profileSetupSchema } from "@/schemas/auth.schema";
import { useAuthStore } from "@/store/useAuthStore";
import PremiumAuthLayout from "@/components/auth/PremiumAuthLayout";
import UsernameChecker from "@/components/auth/UsernameChecker";
import ProgressStepper from "@/components/auth/ProgressStepper";
import { Camera, Upload, AlertCircle, CheckCircle2 } from "lucide-react";
import { getErrorMessage } from "@/lib/form";

export default function ProfileSetupPage() {
    const router = useRouter();
    const { completeProfile } = useAuthStore();

    const [avatar, setAvatar] = useState<string | null>(null);
    const [avatarFile, setAvatarFile] = useState<File | null>(null);
    const [usernameAvailable, setUsernameAvailable] = useState<boolean | null>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);

    const { control, formState: { errors }, handleSubmit, watch, onSubmitHandler, isSubmitting } = useFormValidation({
        schema: profileSetupSchema,
        onSubmit: async (data) => {
            await completeProfile(data.firstName, data.lastName, data.username);
            // Navigate based on role (would be stored in auth context)
            router.push("/auth/onboarding/customer/interests");
        },
    });

    const firstName = watch("firstName");
    const lastName = watch("lastName");
    const username = watch("username");

    const handleAvatarClick = () => {
        fileInputRef.current?.click();
    };

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            // Validate file type and size
            if (!file.type.startsWith("image/")) {
                alert("Please select an image file");
                return;
            }
            if (file.size > 5 * 1024 * 1024) {
                alert("File size must be less than 5MB");
                return;
            }

            // Create preview
            const reader = new FileReader();
            reader.onload = (event) => {
                setAvatar(event.target?.result as string);
            };
            reader.readAsDataURL(file);
            setAvatarFile(file);
        }
    };

    const steps = [
        { id: "profile", label: "Profile Picture", description: "Add a photo (optional)" },
        { id: "details", label: "Personal Details", description: "Your name" },
        { id: "username", label: "Username", description: "Your unique handle" },
    ];

    const isProfileComplete = !!firstName && !!lastName && !!username && usernameAvailable === true;

    return (
        <PremiumAuthLayout
            title="Complete Your Profile"
            description="Help others recognize you. You can update this anytime."
            variant="compact"
            showBackButton
            backHref="/auth/verify-email"
        >
            <div className="space-y-6">
                {/* Progress indicator */}
                <ProgressStepper
                    steps={steps}
                    currentStep={1}
                    variant="linear"
                    showLabels={true}
                />

                <form onSubmit={handleSubmit(onSubmitHandler)} className="space-y-6">
                    {/* Profile Picture Section */}
                    <div className="space-y-3">
                        <label className="block text-sm font-semibold text-zinc-900">
                            Profile Picture
                            <span className="text-zinc-500 font-normal ml-2">(optional)</span>
                        </label>

                        <div className="space-y-3">
                            {/* Avatar display */}
                            <button
                                type="button"
                                onClick={handleAvatarClick}
                                className="relative h-24 w-24 overflow-hidden rounded-full border-2 border-dashed border-zinc-300 bg-zinc-50 transition-all hover:border-pink-500 hover:bg-pink-50"
                            >
                                {avatar ? (
                                    <img src={avatar} alt="Preview" className="h-full w-full object-cover" />
                                ) : (
                                    <div className="flex h-full w-full flex-col items-center justify-center gap-1">
                                        <Camera className="h-6 w-6 text-zinc-400" />
                                        <span className="text-xs text-zinc-500">Add photo</span>
                                    </div>
                                )}
                            </button>

                            <input
                                ref={fileInputRef}
                                type="file"
                                accept="image/*"
                                onChange={handleFileChange}
                                className="hidden"
                                aria-label="Upload profile picture"
                            />

                            <p className="text-xs text-zinc-600">
                                <Upload className="h-3 w-3 inline mr-1" />
                                JPG or PNG, max 5MB
                            </p>
                        </div>
                    </div>

                    {/* First Name */}
                    <div>
                        <label htmlFor="firstName" className="block text-sm font-semibold text-zinc-900 mb-2">
                            First Name
                            <span className="text-red-500">*</span>
                        </label>
                        <Controller
                            name="firstName"
                            control={control}
                            render={({ field }) => (
                                <input
                                    {...field}
                                    type="text"
                                    id="firstName"
                                    placeholder="John"
                                    className={`w-full rounded-lg border-2 px-4 py-3 transition-all focus:outline-none ${errors.firstName
                                        ? "border-red-300 bg-red-50 focus:border-red-400 focus:ring-2 focus:ring-red-100"
                                        : firstName
                                            ? "border-green-300 bg-green-50 focus:border-green-400 focus:ring-2 focus:ring-green-100"
                                            : "border-zinc-300 focus:border-pink-500 focus:ring-2 focus:ring-pink-100"
                                        }`}
                                />
                            )}
                        />
                        {errors.firstName && (
                            (() => {
                                const msg = getErrorMessage(errors.firstName, "Invalid first name");
                                return msg ? (
                                    <p className="mt-2 text-xs text-red-600 flex items-center gap-1">
                                        <AlertCircle className="h-3 w-3" />
                                        {msg}
                                    </p>
                                ) : null;
                            })()
                        )}
                    </div>

                    {/* Last Name */}
                    <div>
                        <label htmlFor="lastName" className="block text-sm font-semibold text-zinc-900 mb-2">
                            Last Name
                            <span className="text-red-500">*</span>
                        </label>
                        <Controller
                            name="lastName"
                            control={control}
                            render={({ field }) => (
                                <input
                                    {...field}
                                    type="text"
                                    id="lastName"
                                    placeholder="Doe"
                                    className={`w-full rounded-lg border-2 px-4 py-3 transition-all focus:outline-none ${errors.lastName
                                        ? "border-red-300 bg-red-50 focus:border-red-400 focus:ring-2 focus:ring-red-100"
                                        : lastName
                                            ? "border-green-300 bg-green-50 focus:border-green-400 focus:ring-2 focus:ring-green-100"
                                            : "border-zinc-300 focus:border-pink-500 focus:ring-2 focus:ring-pink-100"
                                        }`}
                                />
                            )}
                        />
                        {errors.lastName && (
                            (() => {
                                const msg = getErrorMessage(errors.lastName, "Invalid last name");
                                return msg ? (
                                    <p className="mt-2 text-xs text-red-600 flex items-center gap-1">
                                        <AlertCircle className="h-3 w-3" />
                                        {msg}
                                    </p>
                                ) : null;
                            })()
                        )}
                    </div>

                    {/* Username (using UsernameChecker component) */}
                    <div>
                        <Controller
                            name="username"
                            control={control}
                            render={({ field }) => (
                                <UsernameChecker
                                    onUsernameChange={(val) => field.onChange(val)}
                                    onAvailabilityChange={setUsernameAvailable}
                                    autoFocus={false}
                                />
                            )}
                        />
                    </div>

                    {/* Submit button */}
                    <button
                        type="submit"
                        disabled={isSubmitting || !isProfileComplete}
                        className="w-full rounded-lg bg-gradient-to-r from-pink-600 to-pink-500 px-4 py-3 font-semibold text-white transition-all hover:shadow-lg hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed disabled:scale-100 flex items-center justify-center gap-2"
                    >
                        {isSubmitting ? (
                            <>
                                <div className="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
                                Saving profile...
                            </>
                        ) : (
                            <>
                                <CheckCircle2 className="h-5 w-5" />
                                Continue to Onboarding
                            </>
                        )}
                    </button>

                    {/* Info text */}
                    <p className="text-xs text-center text-zinc-600">
                        You can always update your profile later from your account settings.
                    </p>
                </form>
            </div>
        </PremiumAuthLayout>
    );
}
