import SectionHeader from "@/components/ui/SectionHeader";
import { sessionHistory } from "@/lib/mock/users";

export default function SessionsPage() {
    return (
        <div className="space-y-8">
            <SectionHeader title="Session History" description="Review your past sessions and activities." />
            <div className="space-y-4 rounded-3xl border border-slate-200 bg-white p-6 shadow-sm">
                {sessionHistory.map((session) => (
                    <div key={session.id} className="flex items-center justify-between gap-4 rounded-3xl border border-slate-200 bg-slate-50 px-5 py-4">
                        <div className="flex items-center gap-4">
                            <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-slate-100 text-xl">
                                {session.avatar}
                            </div>
                            <div>
                                <p className="font-semibold text-slate-900">{session.device}</p>
                                <p className="text-sm text-slate-500">{session.location}</p>
                            </div>
                        </div>
                        <span className={`rounded-full px-3 py-1 text-xs font-medium ${session.active ? "bg-emerald-100 text-emerald-700" : "bg-slate-100 text-slate-600"}`}>
                            {session.active ? "Active now" : "Last active"}
                        </span>
                    </div>
                ))}
            </div>
        </div>
    );
}
