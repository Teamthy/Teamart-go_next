export default function DiscountBadge({ value }: { value: string }) {
    return (
        <span className="inline-flex items-center rounded-full bg-rose-100 px-3 py-1 text-xs font-semibold text-rose-700">
            {value}
        </span>
    );
}
