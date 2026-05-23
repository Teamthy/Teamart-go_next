"use client";

export default function LiveVideoPlayer() {
    return (
        <div className="rounded-3xl border border-slate-200 bg-black p-4 shadow-sm">
            <div className="aspect-video w-full overflow-hidden rounded-3xl bg-slate-900">
                <div className="flex h-full items-center justify-center text-white text-sm">Live video player placeholder</div>
            </div>
            <div className="mt-4 flex items-center justify-between text-sm text-slate-300">
                <span>Live shopping stream · 18:23 elapsed</span>
                <span className="rounded-full bg-emerald-500/15 px-3 py-1 text-emerald-200">Live</span>
            </div>
        </div>
    );
}
