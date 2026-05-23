export default function LivestreamStatus() {
    return (
        <section className="space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <div className="flex items-center justify-between">
                <div>
                    <h3 className="text-lg font-semibold text-slate-900">Livestream status</h3>
                    <p className="text-sm text-slate-500">Real-time broadcast details and creator updates.</p>
                </div>
                <span className="rounded-full bg-emerald-100 px-3 py-1 text-sm font-semibold text-emerald-800">Live now</span>
            </div>
            <div className="grid gap-4 sm:grid-cols-2">
                <div className="rounded-3xl bg-slate-50 p-4">
                    <p className="text-sm text-slate-500">Current viewers</p>
                    <p className="mt-2 text-3xl font-semibold text-slate-900">1.2k</p>
                </div>
                <div className="rounded-3xl bg-slate-50 p-4">
                    <p className="text-sm text-slate-500">Featured product</p>
                    <p className="mt-2 text-xl font-semibold text-slate-900">Creator Collaboration Hoodie</p>
                </div>
            </div>
        </section>
    );
}
