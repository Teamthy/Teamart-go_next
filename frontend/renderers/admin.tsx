import Link from "next/link";
import Button from "@/components/ui/button";
import Card from "@/components/ui/card";
import {
    adminAlerts,
    adminAnalyticsMetrics,
    adminAuditEntries,
    adminComplianceItems,
    adminMetrics,
    adminOperationalSnapshot,
    adminReports,
    adminSettingsCards,
    adminShortcuts,
    adminTickets,
    adminUserRoles,
    moderationEscalation,
    moderationQueue,
} from "@/lib/mock/admin";
import { renderHero } from "./common";

function renderAdminPage(title: string, description: string) {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title, description, badge: "Admin" })}
            <div className="grid gap-4 md:grid-cols-3">
                {adminMetrics.map((item) => (
                    <Card key={item.label} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{item.label}</p>
                        <p className="mt-3 text-2xl font-semibold text-zinc-900">{item.value}</p>
                    </Card>
                ))}
            </div>
            <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Operational snapshot</p>
                    <div className="mt-4 space-y-3">
                        {adminOperationalSnapshot.map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">{item}</div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Admin shortcuts</p>
                    <div className="mt-4 flex flex-wrap gap-3">
                        {adminShortcuts.map((item) => (
                            <Button key={item.href} asChild variant={item.href === "/admin/users" ? "primary" : "secondary"}>
                                <Link href={item.href}>{item.label}</Link>
                            </Button>
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderAdminAnalyticsPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Admin analytics", description: "Review aggregate performance, growth trends, and health signals across the Teamart platform.", badge: "Analytics" })}
            <div className="grid gap-4 lg:grid-cols-3">
                {adminAnalyticsMetrics.map((item) => (
                    <Card key={item.label} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{item.label}</p>
                        <p className="mt-3 text-2xl font-semibold text-zinc-900">{item.value}</p>
                    </Card>
                ))}
            </div>
            <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Operating signals</p>
                    <div className="mt-4 space-y-3">
                        {adminOperationalSnapshot.map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">{item}</div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">What to watch</p>
                    <div className="mt-4 space-y-3">
                        {adminAlerts.map((item) => (
                            <div key={item} className="rounded-[24px] bg-zinc-50 px-4 py-3 text-sm text-zinc-700">{item}</div>
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderAdminUsersPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Admin users", description: "Keep account visibility, permissions, and creator access aligned with operational needs.", badge: "Users" })}
            <div className="grid gap-4 lg:grid-cols-2">
                {adminUserRoles.map((item) => (
                    <Card key={item.label} className="p-5">
                        <p className="text-lg font-semibold text-zinc-900">{item.label}</p>
                        <p className="mt-2 text-sm text-zinc-600">{item.description}</p>
                    </Card>
                ))}
            </div>
        </div>
    );
}

function renderAdminModerationPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Admin moderation", description: "Review flagged content, creator decisions, and prioritized actions to keep the marketplace healthy.", badge: "Moderation" })}
            <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Priority queue</p>
                    <div className="mt-4 space-y-3">
                        {moderationQueue.map((item) => (
                            <div key={item} className="rounded-[24px] bg-[#FFF8FB] px-4 py-3 text-sm text-zinc-700">{item}</div>
                        ))}
                    </div>
                </Card>
                <Card className="p-5">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Escalation</p>
                    <div className="mt-4 flex flex-wrap gap-3">
                        {moderationEscalation.map((item) => (
                            <Button key={item} asChild variant={item === "Review users" ? "primary" : "secondary"}>
                                <Link href={item === "Review users" ? "/admin/users" : "/admin/analytics"}>{item}</Link>
                            </Button>
                        ))}
                    </div>
                </Card>
            </div>
        </div>
    );
}

function renderAdminSettingsPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Admin settings", description: "Inspect core platform configuration, team access, and release operations from a central admin surface.", badge: "Settings" })}
            <div className="grid gap-4 lg:grid-cols-2">
                {adminSettingsCards.map((item) => (
                    <Card key={item.label} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{item.label}</p>
                        <p className="mt-3 text-sm text-zinc-700">{item.value}</p>
                    </Card>
                ))}
            </div>
        </div>
    );
}

function renderAdminAlertsPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Admin alerts", description: "Track active operational alerts and keep the platform response team aligned on urgent issues.", badge: "Alerts" })}
            <div className="grid gap-4 lg:grid-cols-2">
                {adminAlerts.map((item) => (
                    <Card key={item} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Operational alert</p>
                        <p className="mt-3 text-sm text-zinc-700">{item}</p>
                    </Card>
                ))}
            </div>
        </div>
    );
}

function renderAdminTicketsPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Admin tickets", description: "Review support tickets and keep customer-facing escalations organized by priority.", badge: "Tickets" })}
            <div className="grid gap-4 lg:grid-cols-2">
                {adminTickets.map((ticket) => (
                    <Card key={ticket.id} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{ticket.id}</p>
                        <p className="mt-3 text-lg font-semibold text-zinc-900">{ticket.title}</p>
                        <div className="mt-3 flex flex-wrap gap-2 text-sm text-zinc-700">
                            <span>{ticket.priority}</span>
                            <span>•</span>
                            <span>{ticket.status}</span>
                        </div>
                    </Card>
                ))}
            </div>
        </div>
    );
}

function renderAdminCompliancePage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Compliance", description: "Review the latest policy, audit, and compliance checks in one place.", badge: "Compliance" })}
            <div className="grid gap-4 lg:grid-cols-2">
                {adminComplianceItems.map((item) => (
                    <Card key={item.title} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{item.status}</p>
                        <p className="mt-3 text-lg font-semibold text-zinc-900">{item.title}</p>
                        <p className="mt-2 text-sm text-zinc-600">{item.detail}</p>
                    </Card>
                ))}
            </div>
        </div>
    );
}

function renderAdminAuditPage() {
    return (
        <div className="space-y-8 pb-10">
            {renderHero({ title: "Audit log", description: "Inspect recent admin actions and keep a clear record of operational activity.", badge: "Audit" })}
            <div className="grid gap-4 lg:grid-cols-2">
                {adminAuditEntries.map((entry) => (
                    <Card key={entry.id} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{entry.id}</p>
                        <p className="mt-3 text-lg font-semibold text-zinc-900">{entry.action}</p>
                        <p className="mt-2 text-sm text-zinc-600">{entry.actor}</p>
                        <p className="mt-1 text-sm text-zinc-500">{entry.time}</p>
                    </Card>
                ))}
            </div>
            <div className="grid gap-4 lg:grid-cols-2">
                {adminReports.map((item) => (
                    <Card key={item} className="p-5">
                        <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Report</p>
                        <p className="mt-3 text-sm text-zinc-700">{item}</p>
                    </Card>
                ))}
            </div>
        </div>
    );
}

export function renderAdmin(slug: string[]) {
    const route = slug[0] ?? "admin";
    const second = slug[1];

    if (route !== "admin") {
        return renderAdminPage("Admin", "A secure admin surface for analytics, moderation, users, and settings.");
    }

    if (second === "analytics") return renderAdminAnalyticsPage();
    if (second === "users") return renderAdminUsersPage();
    if (second === "moderation") return renderAdminModerationPage();
    if (second === "settings") return renderAdminSettingsPage();
    if (second === "alerts") return renderAdminAlertsPage();
    if (second === "tickets") return renderAdminTicketsPage();
    if (second === "compliance") return renderAdminCompliancePage();
    if (second === "audit") return renderAdminAuditPage();

    return renderAdminPage("Admin", "A secure admin surface for analytics, moderation, users, and settings.");
}
