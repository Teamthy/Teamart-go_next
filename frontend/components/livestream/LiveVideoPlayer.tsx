"use client";

export default function LiveVideoPlayer() {
    return (
        <div className="rounded-[28px] bg-black p-4">
            <div className="relative aspect-video overflow-hidden rounded-[24px] bg-[linear-gradient(135deg,#111827_0%,#E91E63_55%,#111827_100%)]">
                <div className="absolute inset-0 bg-[radial-gradient(circle_at_top_left,_rgba(255,255,255,0.26),transparent_18%)]" />
                <div className="absolute inset-0 flex items-center justify-center">
                    <div className="text-center">
                        <p className="rounded-full bg-white/10 px-3 py-1 text-[11px] font-semibold uppercase tracking-[0.2em] text-white">
                            Live now
                        </p>
                        <p className="mt-4 text-lg font-semibold text-white">Creator Collaboration Hoodie</p>
                        <p className="mt-2 text-sm text-white/80">18:23 elapsed • 1.2k viewers</p>
                    </div>
                </div>
            </div>
            <div className="mt-4 flex items-center justify-between text-sm text-zinc-100">
                <span>Live shopping stream • pinned product</span>
                <span className="rounded-full bg-[#FCE4EC] px-3 py-1 text-[11px] font-semibold text-[#E91E63]">
                    +120 comments
                </span>
            </div>
        </div>
    );
}
