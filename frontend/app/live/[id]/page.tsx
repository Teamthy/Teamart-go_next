"use client";

import { useMemo, useState } from "react";
import Link from "next/link";
import { useLiveRoom } from "@/hooks/useLiveRoom";
import LiveRoom from "@/components/live/LiveRoom";
import Button from "@/components/ui/button";
import * as api from "@/lib/api";

export default function LiveRoomDetailPage({ params }: { params: { id: string } }) {
    const roomId = params.id;
    const { room, isLoading, isError, error, sendReaction, sendChat, isChatting } = useLiveRoom(roomId);
    const [muted, setMuted] = useState(true);
    const [drawerOpen, setDrawerOpen] = useState(false);
    const [cartCount, setCartCount] = useState(0);
    const [selectedIndex, setSelectedIndex] = useState(0);
    const [purchaseError, setPurchaseError] = useState<string | null>(null);

    const selectedProduct = useMemo(
        () => room?.pinned_products?.[selectedIndex] ?? room?.pinned_products?.[0] ?? null,
        [room, selectedIndex]
    );

    const handleToggleMute = () => setMuted((value) => !value);

    const handleOpenCart = () => setDrawerOpen(true);
    const handleCloseCart = () => setDrawerOpen(false);

    const handleAddToCart = async (quantity: number) => {
        if (!selectedProduct) return;

        try {
            await api.addCartItem({ product_id: selectedProduct.id, quantity });
            setCartCount((count) => count + quantity);
            setPurchaseError(null);
        } catch (error) {
            setPurchaseError(error instanceof Error ? error.message : "Unable to add to cart.");
        }
    };

    const handleBuyNow = async () => {
        if (!selectedProduct) return;
        try {
            await api.addCartItem({ product_id: selectedProduct.id, quantity: 1 });
            window.location.assign("/checkout");
        } catch (error) {
            setPurchaseError(error instanceof Error ? error.message : "Unable to proceed to checkout.");
        }
    };

    const handleSendMessage = (message: string) => {
        sendChat(message);
    };

    const handleReact = (reaction: string) => {
        sendReaction(reaction);
    };

    if (isLoading) {
        return <div className="rounded-[32px] border border-slate-200 bg-white p-8 text-center text-slate-600 shadow-lg">Loading live room...</div>;
    }

    if (isError || !room) {
        return (
            <div className="rounded-[32px] border border-rose-200 bg-rose-50 p-8 text-center text-rose-700 shadow-lg">
                <p className="font-semibold">Unable to load stream.</p>
                <p className="mt-2">{error ?? "The live room could not be loaded right now."}</p>
                <div className="mt-5">
                    <Button variant="secondary" asChild>
                        <Link href="/live">Back to live list</Link>
                    </Button>
                </div>
            </div>
        );
    }

    return (
        <div className="space-y-8 pb-10">
            <div className="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
                <div>
                    <p className="text-sm uppercase tracking-[0.3em] text-pink-600">Live room</p>
                    <h1 className="mt-3 text-4xl font-semibold text-slate-950">{room.title}</h1>
                    <p className="mt-3 text-base leading-7 text-slate-600">{room.description}</p>
                </div>
                <div className="flex flex-wrap gap-3">
                    <Button variant="secondary" asChild>
                        <Link href="/live">All live rooms</Link>
                    </Button>
                    <Button variant="primary" onClick={handleOpenCart}>Open cart ({cartCount})</Button>
                </div>
            </div>
            <LiveRoom
                room={room}
                muted={muted}
                onToggleMute={handleToggleMute}
                onReact={handleReact}
                onOpenCart={handleOpenCart}
                cartCount={cartCount}
                drawerOpen={drawerOpen}
                onCloseDrawer={handleCloseCart}
                onAddToCart={handleAddToCart}
                onBuyNow={handleBuyNow}
                selectedProduct={selectedProduct}
                onSelectProduct={setSelectedIndex}
                chatProps={{
                    messages: room.messages,
                    onSendMessage: handleSendMessage,
                    isSending: isChatting,
                }}
            />
            {purchaseError ? (
                <div className="rounded-3xl bg-rose-500/10 px-4 py-3 text-sm text-rose-700">{purchaseError}</div>
            ) : null}
        </div>
    );
}
