import Badge from "@/components/ui/badge";

interface CartItemRowProps {
    name: string;
    subtitle: string;
    price: number;
    quantity: number;
    onChange: (next: number) => void;
    onRemove: () => void;
}

export default function CartItemRow({ name, subtitle, price, quantity, onChange, onRemove }: CartItemRowProps) {
    return (
        <div className="rounded-[24px] border border-slate-200 bg-slate-50 p-4">
            <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
                <div className="space-y-2">
                    <div className="flex items-center gap-2">
                        <h3 className="text-base font-semibold text-slate-900">{name}</h3>
                        <Badge tone="info">Popular</Badge>
                    </div>
                    <p className="text-sm text-slate-500">{subtitle}</p>
                    <p className="text-base font-semibold text-slate-900">${price.toFixed(2)}</p>
                </div>
                <div className="flex items-center gap-3">
                    <div className="flex items-center gap-2 rounded-full bg-white px-2 py-1 shadow-sm">
                        <button
                            type="button"
                            onClick={() => onChange(Math.max(1, quantity - 1))}
                            className="h-8 w-8 rounded-full bg-slate-100 text-slate-800"
                        >
                            −
                        </button>
                        <span className="min-w-8 text-center text-sm font-semibold text-slate-900">{quantity}</span>
                        <button
                            type="button"
                            onClick={() => onChange(quantity + 1)}
                            className="h-8 w-8 rounded-full bg-slate-100 text-slate-800"
                        >
                            +
                        </button>
                    </div>
                    <button
                        type="button"
                        onClick={onRemove}
                        className="rounded-full px-3 py-2 text-sm font-semibold text-rose-600 hover:bg-rose-50"
                    >
                        Remove
                    </button>
                </div>
            </div>
        </div>
    );
}
