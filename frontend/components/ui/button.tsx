import { cloneElement, isValidElement, type ButtonHTMLAttributes, type ReactElement, type ReactNode } from "react";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: "primary" | "secondary" | "ghost";
    children: ReactNode;
    asChild?: boolean;
}

const variantStyles = {
    primary: "bg-[var(--primary)] text-[var(--surface)] hover:bg-[var(--primary-darker)]",
    secondary: "border border-[var(--surface-border)] bg-[var(--surface)] text-[var(--foreground)] hover:bg-[var(--surface-muted)]",
    ghost: "border border-[var(--primary)] bg-transparent text-[var(--primary)] hover:bg-[var(--primary-muted)]",
};

export default function Button({
    variant = "primary",
    className = "",
    children,
    asChild = false,
    ...props
}: ButtonProps) {
    const baseClassName = `inline-flex items-center justify-center rounded-[var(--button-radius)] px-4 py-3 text-sm font-semibold transition-all duration-200 active:scale-[0.99] ${variantStyles[variant]} ${className}`;

    if (asChild && isValidElement(children)) {
        const child = children as ReactElement<{ className?: string }>;
        return cloneElement(child, {
            className: `${baseClassName} ${child.props.className ?? ""}`,
        });
    }

    return (
        <button
            {...props}
            className={baseClassName}
        >
            {children}
        </button>
    );
}
