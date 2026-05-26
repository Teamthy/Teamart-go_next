export default function SectionHeader({
    title,
    description,
}: {
    title: string;
    description: string;
}) {
    return (
        <div className="space-y-2">
            <h2 className="text-[22px] font-semibold tracking-tight text-zinc-900 sm:text-[24px]">{title}</h2>
            <p className="max-w-2xl text-sm leading-6 text-zinc-600 sm:text-[15px]">{description}</p>
        </div>
    );
}
