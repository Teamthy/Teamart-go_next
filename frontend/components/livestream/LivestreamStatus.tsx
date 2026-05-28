export default function LivestreamStatus() {
    return (
        <section className="space-y-4 rounded-[28px] border border-zinc-200 bg-white p-5 sm:p-6">
            <div className="flex items-center justify-between gap-4">
                <div>
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Livestream status</p>
                    <h3 className="mt-2 text-lg font-semibold text-zinc-900">Creator collaboration drop</h3>
                    <p className="mt-1 text-sm text-zinc-500">Real-time broadcast details and creator updates.</p>
                </div>
                <span className="rounded-full bg-[#FCE4EC] px-3 py-1 text-[11px] font-semibold text-[#E91E63]">Live now</span>
            </div>
            <div className="grid gap-3 sm:grid-cols-2">
                <div className="rounded-[24px] bg-[#FFF8FB] p-4">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Current viewers</p>
                    <p className="mt-2 text-3xl font-semibold text-zinc-900">1.2k</p>
                </div>
                <div className="rounded-[24px] bg-[#FFF8FB] p-4">
                    <p className="text-[11px] uppercase tracking-[0.2em] text-zinc-500">Featured product</p>
                    <p className="mt-2 text-lg font-semibold text-zinc-900">Creator Collaboration Hoodie</p>
                </div>
            </div>
        </section>
    );
}
