import { ArrowLeft } from "lucide-react";

export default function BackBar({
    title,
    onBack,
}: {
    title?: string;
    onBack: () => void;
}) {
    return (
        <div className="flex items-center gap-3 px-4 pt-4 pb-2">
            <button
                type="button"
                onClick={onBack}
                className="grid h-9 w-9 place-items-center rounded-full border border-slate-200 bg-white text-slate-700 transition hover:bg-slate-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#E91E63]/70"
                aria-label="Go back"
            >
                <ArrowLeft className="h-4 w-4" />
            </button>
            <div className="text-sm font-semibold text-slate-900">{title}</div>
        </div>
    );
}
