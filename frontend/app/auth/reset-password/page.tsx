import SectionHeader from "@/components/ui/SectionHeader";

export default function ResetPasswordPage() {
    return (
        <div className="mx-auto max-w-xl space-y-8 rounded-3xl border border-slate-200 bg-white p-8 shadow-sm">
            <SectionHeader title="Reset password" description="Set a new password to restore access to your account." />
            <form className="space-y-4">
                <label className="block text-sm font-medium text-slate-700">
                    New password
                    <input
                        type="password"
                        className="mt-2 w-full rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none focus:border-slate-400"
                        placeholder="Enter your new password"
                    />
                </label>
                <label className="block text-sm font-medium text-slate-700">
                    Confirm password
                    <input
                        type="password"
                        className="mt-2 w-full rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none focus:border-slate-400"
                        placeholder="Confirm your password"
                    />
                </label>
                <button className="w-full rounded-3xl bg-slate-900 px-5 py-3 text-sm font-semibold text-white transition hover:bg-slate-700">
                    Reset password
                </button>
            </form>
        </div>
    );
}
