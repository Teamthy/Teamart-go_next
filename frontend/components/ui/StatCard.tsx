import Card from "@/components/ui/card";

interface StatCardProps {
    label: string;
    value: string;
    helper?: string;
}

export default function StatCard({ label, value, helper }: StatCardProps) {
    return (
        <Card className="p-4 sm:p-5">
            <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">{label}</p>
            <p className="mt-3 text-2xl font-semibold text-zinc-900">{value}</p>
            {helper ? <p className="mt-2 text-sm text-zinc-500">{helper}</p> : null}
        </Card>
    );
}
