import CartSummary from "@/components/order/CartSummary";
import SectionHeader from "@/components/ui/SectionHeader";

export default function CartPage() {
    return (
        <div className="space-y-8">
            <SectionHeader title="Your cart" description="Review items in your cart and prepare for checkout." />
            <CartSummary />
        </div>
    );
}
