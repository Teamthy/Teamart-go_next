import { cloneElement, isValidElement, type ButtonHTMLAttributes, type ReactElement, type ReactNode } from "react";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: "primary" | "secondary" | "ghost";
    children: ReactNode;
    asChild?: boolean;
}

const variantStyles = {
    primary: "bg-[#E91E63] text-white hover:bg-[#d81b60]",
    secondary: "border border-zinc-200 bg-white text-zinc-900 hover:bg-zinc-50",
    ghost: "border border-[#E91E63] bg-transparent text-[#E91E63] hover:bg-[#FCE4EC]",
};

export default function Button({
    variant = "primary",
    className = "",
    children,
    asChild = false,
    ...props
}: ButtonProps) {
    const baseClassName = `inline-flex items-center justify-center rounded-[24px] px-4 py-3 text-sm font-semibold transition-all duration-200 active:scale-[0.99] ${variantStyles[variant]} ${className}`;

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
