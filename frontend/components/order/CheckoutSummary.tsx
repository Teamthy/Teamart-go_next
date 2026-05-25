import { cartItems } from "@/lib/mock-data";

export default function CheckoutSummary() {
    const subtotal = cartItems.reduce((sum, item) => sum + Number(item.price.replace("$", "")) * item.qty, 0);
    const shipping = 8;
    const tax = Math.round(subtotal * 0.08);
    const total = subtotal + shipping + tax;

    return (
        <section className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <h3 className="text-lg font-semibold text-slate-900">Order summary</h3>
            <div className="mt-4 space-y-3">
                {cartItems.map((item) => (
                    <div key={item.id} className="flex items-center justify-between rounded-3xl bg-slate-50 p-4">
                        <div>
                            <p className="font-medium text-slate-900">{item.name}</p>
                            <p className="text-sm text-slate-500">Qty {item.qty}</p>
                        </div>
                        <p className="font-semibold text-slate-900">{item.price}</p>
                    </div>
                ))}
            </div>
            <div className="mt-6 space-y-3 border-t border-slate-200 pt-4 text-sm text-slate-600">
                <div className="flex justify-between">
                    <span>Subtotal</span>
                    <span>${subtotal}</span>
                </div>
                <div className="flex justify-between">
                    <span>Shipping</span>
                    <span>${shipping}</span>
                </div>
                <div className="flex justify-between">
                    <span>Estimated tax</span>
                    <span>${tax}</span>
                </div>
            </div>
            <div className="mt-4 flex items-center justify-between rounded-3xl bg-slate-100 p-4 text-base font-semibold text-slate-900">
                <span>Total</span>
                <span>${total}</span>
            </div>
        </section>
    );
}
