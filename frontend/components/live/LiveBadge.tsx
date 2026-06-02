import { Radio } from "lucide-react";

interface LiveBadgeProps {
    status?: string;
}

export default function LiveBadge({ status = "LIVE" }: LiveBadgeProps) {
    return (
        <div className="inline-flex items-center gap-2 rounded-full bg-pink-500 px-3 py-1 text-xs font-semibold uppercase tracking-[0.22em] text-white shadow-lg shadow-pink-500/20">
            <Radio className="h-3.5 w-3.5 text-white" />
            {status}
        </div>
    );
}
