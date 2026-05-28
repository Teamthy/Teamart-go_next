import type { InputHTMLAttributes } from "react";

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
    label?: string;
    helperText?: string;
}

export default function Input({
    label,
    helperText,
    className = "",
    ...props
}: InputProps) {
    return (
        <label className="block text-sm font-medium text-zinc-700">
            {label ? <span className="mb-2 block text-[12px] font-semibold text-zinc-700">{label}</span> : null}
            <input
                {...props}
                className={`w-full rounded-[24px] border border-zinc-200 bg-white px-4 py-3 text-[14px] text-zinc-900 outline-none transition placeholder:text-zinc-400 focus:border-[#E91E63] focus:ring-2 focus:ring-[#E91E63]/10 ${className}`}
            />
            {helperText ? <span className="mt-2 block text-[11px] text-zinc-500">{helperText}</span> : null}
        </label>
    );
}
