import Link from "next/link";
import Card, { CardBody, CardHeader } from "./Card";

interface SectionPlaceholderProps {
    title: string;
    subtitle: string;
    description: string;
    actions: { label: string; href: string }[];
}

export default function SectionPlaceholder({ title, subtitle, description, actions }: SectionPlaceholderProps) {
    return (
        <div className="space-y-6">
            <div className="rounded-3xl border border-[var(--border-light)] bg-[var(--bg-white)] p-6 shadow-sm shadow-slate-200/10">
                <div className="mb-4">
                    <p className="text-sm font-semibold uppercase tracking-[0.2em] text-[var(--primary)]">
                        {subtitle}
                    </p>
                    <h1 className="mt-2 text-3xl font-semibold text-[var(--text-primary)]">{title}</h1>
                    <p className="mt-3 text-sm leading-6 text-[var(--text-secondary)] max-w-2xl">
                        {description}
                    </p>
                </div>
                <div className="grid gap-3 sm:grid-cols-3">
                    {actions.map((action) => (
                        <Link
                            key={action.href}
                            href={action.href}
                            className="inline-flex h-11 items-center justify-center rounded-lg border border-[var(--border-light)] bg-[var(--bg-white)] px-4 text-sm font-semibold text-[var(--text-primary)] transition hover:bg-[var(--bg-lighter)]"
                        >
                            {action.label}
                        </Link>
                    ))}
                </div>
            </div>

            <div className="grid gap-6 lg:grid-cols-3">
                <Card>
                    <CardHeader>
                        <p className="text-sm font-semibold text-[var(--text-primary)]">Seller workflows</p>
                    </CardHeader>
                    <CardBody>
                        <p className="text-sm text-[var(--text-secondary)]">
                            Access catalog management, order operations, livestream tools, and creator onboarding from one place.
                        </p>
                    </CardBody>
                </Card>
                <Card>
                    <CardHeader>
                        <p className="text-sm font-semibold text-[var(--text-primary)]">Analytics & campaign insights</p>
                    </CardHeader>
                    <CardBody>
                        <p className="text-sm text-[var(--text-secondary)]">
                            Track revenue performance and livestream conversions while launching promotions in a single workflow.
                        </p>
                    </CardBody>
                </Card>
                <Card>
                    <CardHeader>
                        <p className="text-sm font-semibold text-[var(--text-primary)]">Inventory health</p>
                    </CardHeader>
                    <CardBody>
                        <p className="text-sm text-[var(--text-secondary)]">
                            Keep products in stock, view low inventory alerts, and manage product status across channels.
                        </p>
                    </CardBody>
                </Card>
            </div>

            <div className="rounded-3xl border border-[var(--border-light)] bg-[var(--bg-white)] p-6 shadow-sm shadow-slate-200/10">
                <div className="flex items-center justify-between gap-4 mb-4">
                    <div>
                        <p className="text-sm font-semibold text-[var(--text-primary)]">Need help getting started?</p>
                        <p className="text-sm text-[var(--text-secondary)]">Visit seller onboarding resources or connect with customer success.</p>
                    </div>
                    <Link href="/seller/settings" className="text-sm font-semibold text-[var(--primary)] hover:text-[var(--primary-hover)]">
                        Open seller settings →
                    </Link>
                </div>
            </div>
        </div>
    );
}
