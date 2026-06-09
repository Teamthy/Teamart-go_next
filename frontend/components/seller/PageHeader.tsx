import { ReactNode } from "react";

interface PageHeaderProps {
    title: string;
    description?: string;
    actions?: ReactNode;
    filters?: ReactNode;
    breadcrumbs?: { label: string; href?: string }[];
}

export default function PageHeader({
    title,
    description,
    actions,
    filters,
    breadcrumbs,
}: PageHeaderProps) {
    return (
        <div className="mb-8">
            {/* Breadcrumbs */}
            {breadcrumbs && breadcrumbs.length > 0 && (
                <div className="mb-4 flex items-center gap-2 text-sm text-[var(--text-tertiary)]">
                    {breadcrumbs.map((crumb, index) => (
                        <div key={index} className="flex items-center gap-2">
                            {crumb.href ? (
                                <a
                                    href={crumb.href}
                                    className="hover:text-[var(--text-primary)] transition-colors"
                                >
                                    {crumb.label}
                                </a>
                            ) : (
                                <span className="text-[var(--text-primary)] font-medium">
                                    {crumb.label}
                                </span>
                            )}
                            {index < breadcrumbs.length - 1 && <span>/</span>}
                        </div>
                    ))}
                </div>
            )}

            {/* Title and Actions */}
            <div className="flex items-start justify-between gap-4 mb-6">
                <div>
                    <h1 className="text-4xl font-bold text-[var(--text-primary)] mb-2">
                        {title}
                    </h1>
                    {description && (
                        <p className="text-[var(--text-secondary)]">{description}</p>
                    )}
                </div>
                {actions && <div className="flex items-center gap-3">{actions}</div>}
            </div>

            {/* Filters */}
            {filters && (
                <div className="bg-[var(--bg-white)] rounded-lg border border-[var(--border-light)] p-4">
                    {filters}
                </div>
            )}
        </div>
    );
}
