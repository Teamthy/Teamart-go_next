export default function SectionHeader({
    title,
    description,
}: {
    title: string;
    description: string;
}) {
    return (
        <div className="mb-6 space-y-2">
            <h2 className="text-2xl font-semibold text-slate-900">{title}</h2>
            <p className="max-w-2xl text-sm text-slate-600">{description}</p>
        </div>
    );
}
