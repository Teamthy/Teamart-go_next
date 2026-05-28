import type { InputHTMLAttributes } from "react";

export default function TextInput({
    label,
    hint,
    className = "",
    ...props
}: InputHTMLAttributes<HTMLInputElement> & {
    label?: string;
    hint?: string;
}) {
    return (
        <div className={className}>
            {label ? <label className="mb-1 block text-[13px] font-medium text-slate-700">{label}</label> : null}
            <input
                className="w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 outline-none transition focus:border-[#E91E63] focus:ring-2 focus:ring-[#E91E63]/10"
                {...props}
            />
            {hint ? <p className="mt-1 text-[12px] text-slate-500">{hint}</p> : null}
        </div>
    );
}
