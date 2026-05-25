import Link from "next/link";

export default function CreatorProfileCard() {
    return (
        <section className="rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
            <div className="flex items-center gap-4">
                <div className="h-16 w-16 rounded-full bg-slate-200" />
                <div>
                    <h3 className="text-xl font-semibold text-slate-900">Mia Rivera</h3>
                    <p className="text-sm text-slate-600">Creator partner focused on lifestyle drops and live commerce.</p>
                </div>
            </div>
            <div className="mt-6 grid gap-4 sm:grid-cols-3">
                <div className="rounded-3xl bg-slate-50 p-4 text-center">
                    <p className="text-sm text-slate-500">Followers</p>
                    <p className="mt-2 text-lg font-semibold text-slate-900">42.1k</p>
                </div>
                <div className="rounded-3xl bg-slate-50 p-4 text-center">
                    <p className="text-sm text-slate-500">Live shows</p>
                    <p className="mt-2 text-lg font-semibold text-slate-900">18</p>
                </div>
                <div className="rounded-3xl bg-slate-50 p-4 text-center">
                    <p className="text-sm text-slate-500">Top product</p>
                    <p className="mt-2 text-lg font-semibold text-slate-900">Collab Hoodie</p>
                </div>
            </div>
            <Link href="/creator/mia" className="mt-6 inline-flex rounded-3xl bg-slate-900 px-4 py-3 text-sm font-semibold text-white transition hover:bg-slate-700">
                View creator profile
            </Link>
        </section>
    );
}
