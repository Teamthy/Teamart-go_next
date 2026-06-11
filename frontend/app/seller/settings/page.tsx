"use client";

import { useState } from "react";
import SellerAppShell from "@/components/seller/SellerAppShell";
import PageHeader from "@/components/seller/PageHeader";
import SellerButton from "@/components/seller/Button";
import { Save } from "lucide-react";

export default function SettingsPage() {
    const [activeTab, setActiveTab] = useState("store");

    const tabs = [
        { id: "store", label: "Store Profile" },
        { id: "business", label: "Business Information" },
        { id: "shipping", label: "Shipping Settings" },
        { id: "payments", label: "Payment Methods" },
        { id: "team", label: "Team Members" },
        { id: "notifications", label: "Notifications" },
    ];

    return (
        <SellerAppShell>
            <PageHeader
                title="Settings"
                description="Manage your store configuration and preferences"
                breadcrumbs={[
                    { label: "Dashboard", href: "/seller/dashboard" },
                    { label: "Settings" },
                ]}
            />

            {/* Settings Tabs */}
            <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
                {/* Sidebar Navigation */}
                <div className="lg:col-span-1">
                    <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] overflow-hidden">
                        {tabs.map((tab) => (
                            <button
                                key={tab.id}
                                onClick={() => setActiveTab(tab.id)}
                                className={`w-full text-left px-4 py-3 text-sm font-medium border-b border-[var(--border-light)] last:border-b-0 transition-colors ${activeTab === tab.id
                                        ? "bg-[var(--primary-lighter)] text-[var(--primary)]"
                                        : "text-[var(--text-secondary)] hover:bg-[var(--bg-lighter)]"
                                    }`}
                            >
                                {tab.label}
                            </button>
                        ))}
                    </div>
                </div>

                {/* Content Area */}
                <div className="lg:col-span-3 space-y-6">
                    {/* Store Profile */}
                    {activeTab === "store" && (
                        <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
                            <h3 className="text-lg font-semibold text-[var(--text-primary)] mb-6">
                                Store Profile
                            </h3>

                            <div className="space-y-6">
                                {/* Store Logo */}
                                <div>
                                    <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
                                        Store Logo
                                    </label>
                                    <div className="flex items-center gap-4">
                                        <div className="w-24 h-24 bg-[var(--bg-lighter)] rounded-lg border border-[var(--border-light)] flex items-center justify-center">
                                            <span className="text-4xl">🏪</span>
                                        </div>
                                        <SellerButton variant="secondary" size="sm">
                                            Change Logo
                                        </SellerButton>
                                    </div>
                                </div>

                                {/* Store Name */}
                                <div>
                                    <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
                                        Store Name
                                    </label>
                                    <input
                                        type="text"
                                        defaultValue="My Awesome Store"
                                        className="w-full px-4 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
                                    />
                                </div>

                                {/* Store URL */}
                                <div>
                                    <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
                                        Store URL
                                    </label>
                                    <input
                                        type="text"
                                        defaultValue="teamart.com/store/my-awesome-store"
                                        disabled
                                        className="w-full px-4 py-2 bg-[var(--bg-lighter)] border border-[var(--border-light)] rounded-lg text-sm text-[var(--text-tertiary)] cursor-not-allowed"
                                    />
                                </div>

                                {/* Store Description */}
                                <div>
                                    <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
                                        Store Description
                                    </label>
                                    <textarea
                                        defaultValue="Welcome to our store! We offer premium electronics and accessories."
                                        rows={4}
                                        className="w-full px-4 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
                                    />
                                </div>

                                {/* Support Email */}
                                <div>
                                    <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
                                        Support Email
                                    </label>
                                    <input
                                        type="email"
                                        defaultValue="support@mystore.com"
                                        className="w-full px-4 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
                                    />
                                </div>

                                <div className="flex justify-end">
                                    <SellerButton variant="primary" icon={<Save className="w-4 h-4" />}>
                                        Save Changes
                                    </SellerButton>
                                </div>
                            </div>
                        </div>
                    )}

                    {/* Business Information */}
                    {activeTab === "business" && (
                        <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
                            <h3 className="text-lg font-semibold text-[var(--text-primary)] mb-6">
                                Business Information
                            </h3>

                            <div className="space-y-6">
                                {/* Business Name */}
                                <div>
                                    <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
                                        Business Name
                                    </label>
                                    <input
                                        type="text"
                                        defaultValue="My Awesome Store LLC"
                                        className="w-full px-4 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
                                    />
                                </div>

                                {/* Business Type */}
                                <div>
                                    <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
                                        Business Type
                                    </label>
                                    <select className="w-full px-4 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]">
                                        <option>Sole Proprietor</option>
                                        <option>Partnership</option>
                                        <option>LLC</option>
                                        <option>Corporation</option>
                                    </select>
                                </div>

                                {/* Tax ID */}
                                <div>
                                    <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
                                        Tax ID
                                    </label>
                                    <input
                                        type="text"
                                        defaultValue="XX-XXXXXXX"
                                        className="w-full px-4 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
                                    />
                                </div>

                                {/* Address */}
                                <div className="grid grid-cols-2 gap-4">
                                    <div>
                                        <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
                                            Address
                                        </label>
                                        <input
                                            type="text"
                                            defaultValue="123 Main Street"
                                            className="w-full px-4 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
                                        />
                                    </div>
                                    <div>
                                        <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
                                            City
                                        </label>
                                        <input
                                            type="text"
                                            defaultValue="New York"
                                            className="w-full px-4 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
                                        />
                                    </div>
                                    <div>
                                        <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
                                            State
                                        </label>
                                        <input
                                            type="text"
                                            defaultValue="NY"
                                            className="w-full px-4 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
                                        />
                                    </div>
                                    <div>
                                        <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
                                            ZIP Code
                                        </label>
                                        <input
                                            type="text"
                                            defaultValue="10001"
                                            className="w-full px-4 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
                                        />
                                    </div>
                                </div>

                                <div className="flex justify-end">
                                    <SellerButton variant="primary" icon={<Save className="w-4 h-4" />}>
                                        Save Changes
                                    </SellerButton>
                                </div>
                            </div>
                        </div>
                    )}

                    {/* Shipping Settings */}
                    {activeTab === "shipping" && (
                        <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
                            <h3 className="text-lg font-semibold text-[var(--text-primary)] mb-6">
                                Shipping Settings
                            </h3>

                            <div className="space-y-6">
                                <div className="bg-[var(--bg-lighter)] p-4 rounded-lg">
                                    <p className="text-sm text-[var(--text-secondary)]">
                                        Configure your shipping methods, rates, and delivery options.
                                    </p>
                                </div>

                                <div className="space-y-4">
                                    <div>
                                        <label className="flex items-center gap-3 cursor-pointer">
                                            <input type="checkbox" defaultChecked className="w-4 h-4" />
                                            <span className="text-sm font-medium text-[var(--text-primary)]">
                                                Standard Shipping (5-7 business days)
                                            </span>
                                        </label>
                                        <input
                                            type="text"
                                            placeholder="$9.99"
                                            className="mt-2 w-32 px-3 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
                                        />
                                    </div>

                                    <div>
                                        <label className="flex items-center gap-3 cursor-pointer">
                                            <input type="checkbox" defaultChecked className="w-4 h-4" />
                                            <span className="text-sm font-medium text-[var(--text-primary)]">
                                                Express Shipping (2-3 business days)
                                            </span>
                                        </label>
                                        <input
                                            type="text"
                                            placeholder="$19.99"
                                            className="mt-2 w-32 px-3 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
                                        />
                                    </div>

                                    <div>
                                        <label className="flex items-center gap-3 cursor-pointer">
                                            <input type="checkbox" className="w-4 h-4" />
                                            <span className="text-sm font-medium text-[var(--text-primary)]">
                                                Overnight Shipping (Next business day)
                                            </span>
                                        </label>
                                        <input
                                            type="text"
                                            placeholder="$49.99"
                                            className="mt-2 w-32 px-3 py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)]"
                                        />
                                    </div>
                                </div>

                                <div className="flex justify-end">
                                    <SellerButton variant="primary" icon={<Save className="w-4 h-4" />}>
                                        Save Changes
                                    </SellerButton>
                                </div>
                            </div>
                        </div>
                    )}

                    {/* Payment Methods */}
                    {activeTab === "payments" && (
                        <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
                            <h3 className="text-lg font-semibold text-[var(--text-primary)] mb-6">
                                Payment Methods
                            </h3>

                            <div className="space-y-4">
                                <div className="p-4 border border-[var(--border-light)] rounded-lg">
                                    <div className="flex items-center justify-between">
                                        <div>
                                            <p className="font-medium text-[var(--text-primary)]">
                                                Stripe
                                            </p>
                                            <p className="text-sm text-[var(--text-tertiary)]">
                                                ••••••• 1234
                                            </p>
                                        </div>
                                        <span className="px-3 py-1 bg-[var(--success-lighter)] text-[var(--success)] rounded-full text-xs font-medium">
                                            Connected
                                        </span>
                                    </div>
                                </div>

                                <div className="p-4 border border-[var(--border-light)] rounded-lg">
                                    <div className="flex items-center justify-between">
                                        <div>
                                            <p className="font-medium text-[var(--text-primary)]">
                                                PayPal
                                            </p>
                                            <p className="text-sm text-[var(--text-tertiary)]">
                                                Not connected
                                            </p>
                                        </div>
                                        <SellerButton variant="secondary" size="sm">
                                            Connect
                                        </SellerButton>
                                    </div>
                                </div>

                                <div className="p-4 border border-[var(--border-light)] rounded-lg">
                                    <div className="flex items-center justify-between">
                                        <div>
                                            <p className="font-medium text-[var(--text-primary)]">
                                                Bank Transfer
                                            </p>
                                            <p className="text-sm text-[var(--text-tertiary)]">
                                                Payout method
                                            </p>
                                        </div>
                                        <SellerButton variant="secondary" size="sm">
                                            Edit
                                        </SellerButton>
                                    </div>
                                </div>
                            </div>
                        </div>
                    )}

                    {/* Team Members */}
                    {activeTab === "team" && (
                        <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
                            <h3 className="text-lg font-semibold text-[var(--text-primary)] mb-6">
                                Team Members
                            </h3>

                            <div className="space-y-4 mb-6">
                                <div className="p-4 border border-[var(--border-light)] rounded-lg">
                                    <div className="flex items-center justify-between">
                                        <div className="flex items-center gap-3">
                                            <div className="w-10 h-10 bg-gradient-to-br from-[var(--primary)] to-[var(--primary-hover)] rounded-full flex items-center justify-center text-[var(--bg-white)] font-bold">
                                                Y
                                            </div>
                                            <div>
                                                <p className="font-medium text-[var(--text-primary)]">
                                                    You
                                                </p>
                                                <p className="text-sm text-[var(--text-tertiary)]">
                                                    owner@store.com • Owner
                                                </p>
                                            </div>
                                        </div>
                                    </div>
                                </div>

                                <div className="p-4 border border-[var(--border-light)] rounded-lg">
                                    <div className="flex items-center justify-between">
                                        <div className="flex items-center gap-3">
                                            <div className="w-10 h-10 bg-gradient-to-br from-[var(--primary-light)] to-[var(--primary)] rounded-full flex items-center justify-center text-[var(--bg-white)] font-bold">
                                                A
                                            </div>
                                            <div>
                                                <p className="font-medium text-[var(--text-primary)]">
                                                    Alice Manager
                                                </p>
                                                <p className="text-sm text-[var(--text-tertiary)]">
                                                    alice@store.com • Manager
                                                </p>
                                            </div>
                                        </div>
                                        <SellerButton variant="secondary" size="sm">
                                            Edit
                                        </SellerButton>
                                    </div>
                                </div>
                            </div>

                            <SellerButton variant="primary">
                                Add Team Member
                            </SellerButton>
                        </div>
                    )}

                    {/* Notifications */}
                    {activeTab === "notifications" && (
                        <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-6">
                            <h3 className="text-lg font-semibold text-[var(--text-primary)] mb-6">
                                Notification Preferences
                            </h3>

                            <div className="space-y-4">
                                {[
                                    { label: "New Orders", checked: true },
                                    { label: "Low Inventory Alerts", checked: true },
                                    { label: "Customer Messages", checked: true },
                                    { label: "Livestream Started", checked: true },
                                    { label: "Creator Joined Campaign", checked: false },
                                    { label: "Refund Requests", checked: true },
                                ].map((notification, index) => (
                                    <label
                                        key={index}
                                        className="flex items-center gap-3 cursor-pointer"
                                    >
                                        <input
                                            type="checkbox"
                                            defaultChecked={notification.checked}
                                            className="w-4 h-4"
                                        />
                                        <span className="text-sm text-[var(--text-primary)]">
                                            {notification.label}
                                        </span>
                                    </label>
                                ))}
                            </div>

                            <div className="flex justify-end mt-6">
                                <SellerButton variant="primary" icon={<Save className="w-4 h-4" />}>
                                    Save Preferences
                                </SellerButton>
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </SellerAppShell>
    );
}
