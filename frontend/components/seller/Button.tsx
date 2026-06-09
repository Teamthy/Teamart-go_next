import { ReactNode, ButtonHTMLAttributes } from "react";

type ButtonVariant =
    | "primary"
    | "secondary"
    | "danger"
    | "ghost"
    | "outline";

interface SellerButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: ButtonVariant;
    size?: "sm" | "md" | "lg";
    children: ReactNode;
    loading?: boolean;
    icon?: ReactNode;
}

const variantStyles: Record<ButtonVariant, string> = {
    primary:
        "bg-[var(--primary)] text-[var(--bg-white)] hover:bg-[var(--primary-hover)] active:bg-[var(--primary-active)]",
    secondary:
        "bg-[var(--bg-lighter)] text-[var(--text-primary)] hover:bg-[var(--bg-muted)] border border-[var(--border-light)]",
    danger:
        "bg-[var(--danger)] text-[var(--bg-white)] hover:bg-[var(--danger-hover)] active:bg-[var(--danger-active)]",
    ghost: "text-[var(--text-primary)] hover:bg-[var(--bg-lighter)]",
    outline:
        "border border-[var(--border-light)] text-[var(--text-primary)] hover:bg-[var(--bg-lighter)]",
};

const sizeStyles: Record<"sm" | "md" | "lg", string> = {
    sm: "px-3 py-2 text-xs",
    md: "px-4 py-2 text-sm",
    lg: "px-6 py-3 text-base",
};

export default function SellerButton({
    variant = "primary",
    size = "md",
    children,
    loading = false,
    icon,
    disabled,
    className = "",
    ...props
}: SellerButtonProps) {
    return (
        <button
            disabled={loading || disabled}
            className={`inline-flex items-center justify-center gap-2 font-medium rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed ${variantStyles[variant]} ${sizeStyles[size]} ${className}`}
            {...props}
        >
            {loading && (
                <div className="animate-spin rounded-full h-4 w-4 border-2 border-current border-t-transparent" />
            )}
            {icon && !loading && icon}
            {children}
        </button>
    );
}
