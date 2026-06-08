import { ReactNode, InputHTMLAttributes } from "react";

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  helpText?: string;
  icon?: ReactNode;
}

export default function Input({
  label,
  error,
  helpText,
  icon,
  className = "",
  ...props
}: InputProps) {
  return (
    <div className="w-full">
      {label && (
        <label className="block text-sm font-medium text-[var(--text-secondary)] mb-2">
          {label}
        </label>
      )}
      <div className="relative">
        {icon && (
          <div className="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--text-tertiary)]">
            {icon}
          </div>
        )}
        <input
          className={`w-full px-${icon ? "10" : "4"} py-2 border border-[var(--border-light)] rounded-lg text-sm focus:outline-none focus:border-[var(--primary)] focus:ring-1 focus:ring-[var(--primary)] ${
            error
              ? "border-[var(--danger)] focus:border-[var(--danger)] focus:ring-[var(--danger)]"
              : ""
          } ${className}`}
          {...props}
        />
      </div>
      {error && (
        <p className="text-xs text-[var(--danger)] mt-1">{error}</p>
      )}
      {helpText && (
        <p className="text-xs text-[var(--text-tertiary)] mt-1">{helpText}</p>
      )}
    </div>
  );
}
