"use client";

import { Heart } from "lucide-react";

interface WishlistButtonProps {
  productId: string;
}

export default function WishlistButton({ productId }: WishlistButtonProps) {
  return (
    <button
      type="button"
      className="flex w-full items-center justify-center gap-2 rounded-3xl border border-slate-200 bg-white px-4 py-3 text-sm font-semibold text-slate-900 transition hover:border-slate-300 hover:bg-slate-50"
      aria-label="Add product to wishlist"
      onClick={() => {
        console.log("Add to wishlist", productId);
      }}
    >
      <Heart className="h-4 w-4" />
      Add to wishlist
    </button>
  );
}
