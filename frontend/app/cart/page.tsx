import Link from "next/link";
import CartSummary from "@/components/order/CartSummary";
import SectionHeader from "@/components/ui/SectionHeader";

export default function CartPage() {
    return (
        <div className="space-y-8">
            <SectionHeader title="Your cart" description="Review items in your cart and prepare for checkout." />

            <div className="grid gap-8 xl:grid-cols-[0.7fr_0.3fr]">
                <div>
                    <CartSummary />
                </div>

                <aside className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
                    <h2 className="text-lg font-semibold text-slate-900">Ready to checkout?</h2>
                    <p className="mt-3 text-sm text-slate-600">
                        Confirm your cart and move to secure payment. You can still edit quantities before placing your order.
                    </p>
                    <Link
                        href="/checkout"
                        className="mt-6 inline-flex w-full items-center justify-center rounded-[24px] bg-[#E91E63] px-5 py-3 text-sm font-semibold text-white transition hover:bg-[#d81b60]"
                    >
                        Continue to checkout
                    </Link>
                    <Link
                        href="/wishlist"
                        className="mt-3 inline-flex w-full items-center justify-center rounded-[24px] border border-slate-200 bg-white px-5 py-3 text-sm font-semibold text-slate-900 transition hover:bg-slate-50"
                    >
                        View wishlist
                    </Link>
                </aside>
            </div>
        </div>
    );
}
