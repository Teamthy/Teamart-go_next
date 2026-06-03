/**
 * RoleSelectorCard Component
 * Premium role selection card with icon, description, and benefits
 */

"use client";

import React from "react";
import Link from "next/link";
import { ArrowRight, ShoppingBag, Sparkles, Store } from "lucide-react";

interface RoleSelectorCardProps {
    role: "customer" | "creator" | "merchant";
    title: string;
    description: string;
    benefits: string[];
    href: string;
    icon?: React.ReactNode;
}

const defaultIcons = {
    customer: <ShoppingBag className="h-8 w-8" />,
    creator: <Sparkles className="h-8 w-8" />,
    merchant: <Store className="h-8 w-8" />,
};

export default function RoleSelectorCard({
    role,
    title,
    description,
    benefits,
    href,
    icon,
}: RoleSelectorCardProps) {
    const bgColors = {
        customer: "from-blue-50 to-blue-50",
        creator: "from-pink-50 to-pink-50",
        merchant: "from-purple-50 to-purple-50",
    };

    const borderColors = {
        customer: "border-blue-200 hover:border-blue-300",
        creator: "border-pink-200 hover:border-pink-300",
        merchant: "border-purple-200 hover:border-purple-300",
    };

    const iconColors = {
        customer: "text-blue-600",
        creator: "text-pink-600",
        merchant: "text-purple-600",
    };

    const buttonColors = {
        customer: "bg-blue-600 hover:bg-blue-700 focus:ring-blue-500",
        creator: "bg-pink-600 hover:bg-pink-700 focus:ring-pink-500",
        merchant: "bg-purple-600 hover:bg-purple-700 focus:ring-purple-500",
    };

    return (
        <Link href={href}>
            <div
                className={`group relative h-full overflow-hidden rounded-2xl border-2 ${borderColors[role]} bg-gradient-to-br ${bgColors[role]} p-6 shadow-sm transition-all duration-300 hover:-translate-y-1 hover:shadow-lg sm:p-8`}
            >
                {/* Background gradient accent */}
                <div className="absolute inset-0 -z-10 opacity-0 transition-opacity group-hover:opacity-100">
                    <div className={`absolute inset-0 bg-gradient-to-br from-white/40 to-transparent`} />
                </div>

                {/* Icon */}
                <div
                    className={`mb-4 inline-flex rounded-xl bg-white p-3 ${iconColors[role]} shadow-sm transition-all group-hover:scale-110 group-hover:shadow-md`}
                >
                    {icon || defaultIcons[role]}
                </div>

                {/* Title */}
                <h3 className="mb-2 text-xl font-bold text-zinc-900 transition-colors group-hover:text-zinc-900">
                    {title}
                </h3>

                {/* Description */}
                <p className="mb-6 text-sm text-zinc-600 leading-relaxed">
                    {description}
                </p>

                {/* Benefits */}
                <ul className="mb-8 space-y-3">
                    {benefits.map((benefit, index) => (
                        <li key={index} className="flex items-start gap-2 text-sm text-zinc-600">
                            <span className="mt-1 inline-block h-1.5 w-1.5 flex-shrink-0 rounded-full bg-zinc-400" />
                            <span>{benefit}</span>
                        </li>
                    ))}
                </ul>

                {/* CTA Button */}
                <div className="inline-flex items-center gap-2 rounded-lg bg-white px-4 py-2.5 font-semibold text-sm text-zinc-900 shadow-sm transition-all group-hover:shadow-md">
                    Get started
                    <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-1" />
                </div>

                {/* Hover overlay */}
                <div className="absolute inset-0 -z-20 transition-opacity opacity-0 group-hover:opacity-5 bg-gradient-to-br from-current to-transparent" />
            </div>
        </Link>
    );
}
