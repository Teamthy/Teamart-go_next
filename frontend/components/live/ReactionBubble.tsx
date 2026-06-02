import { Sparkles } from "lucide-react";

interface ReactionBubbleProps {
    emoji: string;
    label: string;
}

export default function ReactionBubble({ emoji, label }: ReactionBubbleProps) {
    return (
        <div className="group relative inline-flex items-center gap-2 rounded-full bg-white/95 px-3 py-2 text-sm font-semibold text-slate-900 shadow-lg shadow-black/10 transition duration-300 hover:-translate-y-1">
            <span className="text-lg">{emoji}</span>
            <span>{label}</span>
            <Sparkles className="absolute -right-2 -top-2 h-4 w-4 text-amber-400 opacity-80" />
        </div>
    );
}
