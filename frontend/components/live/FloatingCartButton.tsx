import { ShoppingCart } from "lucide-react";

interface FloatingCartButtonProps {
    itemCount: number;
    onOpen: () => void;
}

export default function FloatingCartButton({ itemCount, onOpen }: FloatingCartButtonProps) {
    return (
        <button
            type="button"
            onClick={onOpen}
            className="fixed bottom-6 right-6 z-50 inline-flex items-center gap-3 rounded-full bg-pink-500 px-4 py-3 text-sm font-semibold text-white shadow-2xl shadow-pink-500/30 transition hover:bg-pink-600"
        >
            <ShoppingCart className="h-5 w-5" />
            <span>Cart</span>
            <span className="rounded-full bg-white px-2 py-1 text-xs font-semibold text-slate-900">{itemCount}</span>
        </button>
    );
}
