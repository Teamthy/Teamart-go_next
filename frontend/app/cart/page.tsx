import CartSummary from "@/components/CartSummary";
import SectionHeader from "@/components/SectionHeader";

export default function CartPage() {
    return (
        <div className="space-y-8">
            <SectionHeader title="Your cart" description="Review items in your cart and prepare for checkout." />
            <CartSummary />
        </div>
    );
}
