import Link from "next/link";
import { Heart } from "lucide-react";

export default function WishlistBell() {
    return (
        <Link href="/wishlist" aria-label="Open wishlist" className="rounded-full p-2 text-slate-600 transition hover:bg-slate-100 hover:text-slate-900">
            <Heart className="h-5 w-5" />
        </Link>
    );
}
