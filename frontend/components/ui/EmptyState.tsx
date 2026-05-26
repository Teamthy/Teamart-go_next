interface EmptyStateProps {
    title: string;
    description: string;
}

export default function EmptyState({ title, description }: EmptyStateProps) {
    return (
        <div className="rounded-[28px] border border-dashed border-zinc-200 bg-white p-8 text-center">
            <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">No results</p>
            <h3 className="mt-3 text-lg font-semibold text-zinc-900">{title}</h3>
            <p className="mt-2 text-sm leading-6 text-zinc-600">{description}</p>
        </div>
    );
}
