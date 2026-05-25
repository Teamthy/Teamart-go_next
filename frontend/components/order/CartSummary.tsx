import { cartItems } from "@/lib/mock-data";

export default function CartSummary() {
    const subtotal = cartItems.reduce((sum, item) => sum + Number(item.price.replace("$", "")) * item.qty, 0);
    return (
        <section className="space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <div className="flex items-center justify-between">
                <div>
                    <h3 className="text-lg font-semibold text-slate-900">Cart summary</h3>
                    <p className="text-sm text-slate-500">Review your items before checkout.</p>
                </div>
                <span className="rounded-full bg-slate-100 px-3 py-1 text-sm text-slate-700">{cartItems.length} items</span>
            </div>

            <div className="space-y-4">
                {cartItems.map((item) => (
                    <div key={item.id} className="flex items-center justify-between gap-4 rounded-3xl border border-slate-100 bg-slate-50 p-4">
                        <div>
                            <p className="font-medium text-slate-900">{item.name}</p>
                            <p className="text-sm text-slate-500">Qty {item.qty}</p>
                        </div>
                        <p className="font-semibold text-slate-900">{item.price}</p>
                    </div>
                ))}
            </div>

            <div className="flex items-center justify-between border-t border-slate-200 pt-4">
                <span className="text-sm text-slate-600">Subtotal</span>
                <span className="text-lg font-semibold text-slate-900">${subtotal}</span>
            </div>
        </section>
    );
}
