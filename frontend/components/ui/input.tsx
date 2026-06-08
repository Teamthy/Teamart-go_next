import { forwardRef, type InputHTMLAttributes } from "react";

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
    label?: string;
    helperText?: string;
    helperTextId?: string;
    error?: string | boolean | null;
}

const Input = forwardRef<HTMLInputElement, InputProps>(
    ({ label, helperText, helperTextId, error, className = "", ...props }, ref) => {
        const hasError = Boolean(error);

        return (
            <label className="block text-sm font-medium text-zinc-700">
                {label ? (
                    <span className="mb-2 block text-[12px] font-semibold text-zinc-700">
                        {label}
                    </span>
                ) : null}
                <input
                    {...props}
                    ref={ref}
                    className={`w-full rounded-[24px] border px-4 py-3 text-[14px] text-[var(--foreground)] outline-none transition placeholder:text-[var(--text-muted)] focus:ring-2 ${hasError ? "border-[var(--danger)] bg-[var(--danger-muted)] text-[var(--danger)] focus:border-[var(--danger)] focus:ring-[var(--danger-muted)]" : "border-[var(--surface-border)] bg-[var(--surface)] focus:border-[var(--primary)] focus:ring-[var(--primary-muted)]"} ${className}`}
                />
                {helperText ? (
                    <span id={helperTextId} className={`mt-2 block text-[11px] ${hasError ? "text-red-600" : "text-zinc-500"}`}>
                        {helperText}
                    </span>
                ) : null}
                {typeof error === "string" && (
                    <span className="mt-1 block text-[11px] text-red-600">{error}</span>
                )}
            </label>
        );
    }
);

Input.displayName = "Input";

export default Input;
