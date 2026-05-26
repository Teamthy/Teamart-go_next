interface PriceTagProps {
    price: number;
    compareAt?: number;
}

export default function PriceTag({ price, compareAt }: PriceTagProps) {
    return (
        <div className="flex items-end gap-3">
            <span className="text-3xl font-bold text-slate-900">${price.toFixed(2)}</span>
            {compareAt ? (
                <span className="text-sm text-slate-400 line-through">${compareAt.toFixed(2)}</span>
            ) : null}
        </div>
    );
}
