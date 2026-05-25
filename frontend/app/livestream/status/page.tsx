import ChatPanel from "@/components/ChatPanel";
import LiveVideoPlayer from "@/components/LiveVideoPlayer";
import LivestreamStatus from "@/components/LivestreamStatus";
import ProductPinning from "@/components/ProductPinning";
import ReactionPanel from "@/components/ReactionPanel";

export default function LivestreamStatusPage() {
    return (
        <div className="space-y-8">
            <LivestreamStatus />
            <div className="grid gap-8 xl:grid-cols-[0.65fr_0.35fr]">
                <div className="space-y-6">
                    <LiveVideoPlayer />
                    <div className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
                        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                            <div>
                                <h3 className="text-lg font-semibold text-slate-900">Livestream product spotlight</h3>
                                <p className="text-sm text-slate-600">Pinned product currently featured in the live session.</p>
                            </div>
                            <ProductPinning />
                        </div>
                    </div>
                    <ReactionPanel />
                </div>
                <ChatPanel />
            </div>
        </div>
    );
}
