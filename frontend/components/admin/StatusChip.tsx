type Tone = "success" | "warning" | "error" | "info" | "default";

export default function StatusChip({ label, tone = "default" }: { label: string; tone?: Tone }) {
    const styles: Record<Tone, string> = {
        success: "bg-emerald-100 text-emerald-800",
        warning: "bg-amber-100 text-amber-800",
        error: "bg-rose-100 text-rose-800",
        info: "bg-sky-100 text-sky-800",
        default: "bg-slate-100 text-slate-700",
    };

    return (
        <span className={`inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold ${styles[tone]}`}>
            {label}
        </span>
    );
}
