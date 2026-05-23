import SectionHeader from "@/components/ui/SectionHeader";

export default function ForgotPasswordPage() {
    return (
        <div className="mx-auto max-w-xl space-y-8 rounded-3xl border border-slate-200 bg-white p-8 shadow-sm">
            <SectionHeader title="Forgot password" description="Enter your email to receive a password reset link." />
            <form className="space-y-4">
                <label className="block text-sm font-medium text-slate-700">
                    Email address
                    <input
                        type="email"
                        className="mt-2 w-full rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none focus:border-slate-400"
                        placeholder="you@example.com"
                    />
                </label>
                <button className="w-full rounded-3xl bg-slate-900 px-5 py-3 text-sm font-semibold text-white transition hover:bg-slate-700">
                    Send reset link
                </button>
            </form>
        </div>
    );
}
