import { Eye } from "lucide-react";

interface ViewerCountProps {
    viewers: number;
}

export default function ViewerCount({ viewers }: ViewerCountProps) {
    return (
        <div className="inline-flex items-center gap-2 rounded-full bg-black/60 px-3 py-1 text-xs font-semibold uppercase tracking-[0.22em] text-white shadow-lg shadow-black/30">
            <Eye className="h-4 w-4" />
            {viewers.toLocaleString()} viewers
        </div>
    );
}
