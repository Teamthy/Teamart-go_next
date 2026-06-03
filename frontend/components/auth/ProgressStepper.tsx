/**
 * ProgressStepper Component
 * Visual progress indicator for multi-step forms (customer onboarding)
 */

"use client";

import React from "react";
import { Check } from "lucide-react";

interface Step {
    id: string;
    label: string;
    description?: string;
}

interface ProgressStepperProps {
    steps: Step[];
    currentStep: number;
    onStepClick?: (step: number) => void;
    variant?: "linear" | "circle";
    showLabels?: boolean;
}

export default function ProgressStepper({
    steps,
    currentStep,
    onStepClick,
    variant = "linear",
    showLabels = true,
}: ProgressStepperProps) {
    const progress = ((currentStep) / steps.length) * 100;

    if (variant === "circle") {
        return (
            <div className="mb-8">
                {/* Progress percentage */}
                <div className="mb-6 flex items-baseline justify-between">
                    <h3 className="text-sm font-semibold text-zinc-900">Progress</h3>
                    <span className="text-2xl font-bold text-pink-600">{currentStep * 25}%</span>
                </div>

                {/* Circular progress */}
                <div className="flex items-center justify-center">
                    <div className="relative h-40 w-40">
                        {/* Background circle */}
                        <svg className="h-40 w-40 -rotate-90" viewBox="0 0 160 160">
                            <circle
                                cx="80"
                                cy="80"
                                r="70"
                                fill="none"
                                stroke="#e5e7eb"
                                strokeWidth="8"
                            />
                            {/* Progress circle */}
                            <circle
                                cx="80"
                                cy="80"
                                r="70"
                                fill="none"
                                stroke="currentColor"
                                strokeWidth="8"
                                strokeDasharray={`${2 * Math.PI * 70}`}
                                strokeDashoffset={`${2 * Math.PI * 70 * (1 - currentStep / steps.length)}`}
                                strokeLinecap="round"
                                className="text-pink-600 transition-all duration-500"
                            />
                        </svg>

                        {/* Center content */}
                        <div className="absolute inset-0 flex flex-col items-center justify-center">
                            <div className="text-center">
                                <p className="text-sm text-zinc-600">Step {currentStep}</p>
                                <p className="text-lg font-bold text-zinc-900">
                                    {steps[currentStep - 1]?.label}
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }

    // Linear variant (default)
    return (
        <div className="mb-8 space-y-4">
            {/* Progress bar */}
            <div className="space-y-2">
                <div className="flex items-center justify-between">
                    <span className="text-sm font-semibold text-zinc-900">
                        Step {currentStep} of {steps.length}
                    </span>
                    <span className="text-sm text-zinc-600">{currentStep * 25}%</span>
                </div>

                <div className="h-2 overflow-hidden rounded-full bg-zinc-200">
                    <div
                        className="h-full bg-gradient-to-r from-pink-500 to-pink-600 transition-all duration-500 ease-out"
                        style={{ width: `${progress}%` }}
                    />
                </div>
            </div>

            {/* Step indicators */}
            <div className="grid gap-3">
                {steps.map((step, index) => {
                    const stepNumber = index + 1;
                    const isComplete = stepNumber < currentStep;
                    const isCurrent = stepNumber === currentStep;
                    const isFuture = stepNumber > currentStep;

                    return (
                        <button
                            key={step.id}
                            onClick={() => onStepClick?.(stepNumber)}
                            disabled={isFuture}
                            className={`flex items-start gap-4 rounded-lg p-3 text-left transition-all ${isCurrent
                                    ? "bg-pink-50 ring-2 ring-pink-200"
                                    : isComplete
                                        ? "bg-green-50 hover:bg-green-100"
                                        : "bg-zinc-50 hover:bg-zinc-100"
                                } ${isFuture ? "cursor-not-allowed opacity-50" : "cursor-pointer"}`}
                        >
                            {/* Step indicator */}
                            <div
                                className={`mt-1 flex h-6 w-6 flex-shrink-0 items-center justify-center rounded-full font-semibold text-sm transition-all ${isComplete
                                        ? "bg-green-500 text-white"
                                        : isCurrent
                                            ? "bg-pink-600 text-white ring-2 ring-pink-300 ring-offset-2"
                                            : "bg-zinc-300 text-zinc-600"
                                    }`}
                            >
                                {isComplete ? <Check className="h-4 w-4" /> : stepNumber}
                            </div>

                            {/* Step content */}
                            <div className="flex-1 min-w-0">
                                <p
                                    className={`font-semibold text-sm ${isCurrent
                                            ? "text-pink-900"
                                            : isComplete
                                                ? "text-green-900"
                                                : "text-zinc-900"
                                        }`}
                                >
                                    {step.label}
                                </p>
                                {step.description && (
                                    <p className="text-xs text-zinc-600 mt-0.5">
                                        {step.description}
                                    </p>
                                )}
                            </div>

                            {/* Arrow indicator for current */}
                            {isCurrent && (
                                <div className="mt-1 h-2 w-2 rounded-full bg-pink-600 flex-shrink-0" />
                            )}
                        </button>
                    );
                })}
            </div>
        </div>
    );
}
