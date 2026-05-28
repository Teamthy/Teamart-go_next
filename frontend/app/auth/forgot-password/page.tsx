<<<<<<< HEAD
import AuthTemplate from "@/components/auth/AuthTemplate";

export default function ForgotPasswordPage() {
    return <AuthTemplate variant="forgot" />;
=======
import Link from "next/link";

export default function ForgotPasswordPage() {
    return (
        <div className="min-h-screen bg-[radial-gradient(circle_at_top,_rgba(236,72,153,0.16),transparent_28%),linear-gradient(180deg,#050816_0%,#0e1632_100%)] px-4 py-16 text-white sm:px-6 lg:px-8">
            <div className="mx-auto max-w-3xl rounded-[2.5rem] border border-white/10 bg-white/5 p-10 shadow-2xl shadow-fuchsia-500/10 backdrop-blur-xl">
                <div className="mb-10 space-y-4">
                    <p className="text-sm uppercase tracking-[0.35em] text-fuchsia-300">Account recovery</p>
                    <h1 className="text-4xl font-semibold tracking-tight text-white sm:text-5xl">Forgot your password?</h1>
                    <p className="max-w-2xl text-base leading-7 text-slate-300">
                        Enter the email address linked to your Teamart account, and we’ll send a secure reset link to restore access.
                    </p>
                </div>
                <form className="space-y-6">
                    <label className="block text-sm text-slate-300">
                        Email address
                        <input
                            type="email"
                            className="mt-3 w-full rounded-3xl border border-white/10 bg-slate-950/80 px-4 py-3 text-sm text-white outline-none transition focus:border-fuchsia-400 focus:ring-4 focus:ring-fuchsia-500/10"
                            placeholder="you@example.com"
                        />
                    </label>
                    <button className="w-full rounded-full bg-fuchsia-500 px-5 py-3 text-sm font-semibold text-white transition hover:bg-fuchsia-400">
                        Send reset link
                    </button>
                </form>
                <p className="mt-6 text-center text-sm text-slate-400">
                    Remembered your password? <Link href="/auth/login" className="text-fuchsia-300 hover:text-white">Sign in</Link>
                </p>
            </div>
        </div>
    );
>>>>>>> 36e8d4c (feat(auth): production auth flows and onboarding UI)
}
