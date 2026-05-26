type IllustrationProps = {
    variant: "onboard1" | "onboard2" | "onboard3" | "forgot" | "success" | "phone";
    className?: string;
};

const PINK = "#E91E63";
const PINK_SOFT = "#F8BBD0";

export default function Illustration({ variant, className = "w-56 h-56" }: IllustrationProps) {
    if (variant === "onboard1") {
        return (
            <div className={className}>
                <svg viewBox="0 0 200 200" className="h-full w-full">
                    <circle cx="100" cy="100" r="80" fill="#FCE4EC" />
                    <rect x="70" y="50" width="60" height="100" rx="10" fill="#2d2d2d" />
                    <rect x="74" y="54" width="52" height="92" rx="6" fill="#fff" />
                    <circle cx="100" cy="95" r="18" fill={PINK} />
                    <path d="M100 105 L86 91 a7 7 0 0 1 14 0 a7 7 0 0 1 14 0 Z" fill="#fff" />
                    <circle cx="40" cy="60" r="4" fill={PINK} />
                    <circle cx="160" cy="80" r="5" fill={PINK} />
                    <circle cx="50" cy="150" r="3" fill={PINK} />
                    <circle cx="170" cy="150" r="4" fill={PINK} />
                    <circle cx="30" cy="120" r="2" fill={PINK_SOFT} />
                </svg>
            </div>
        );
    }

    if (variant === "onboard2") {
        return (
            <div className={className}>
                <svg viewBox="0 0 200 200" className="h-full w-full">
                    <circle cx="100" cy="100" r="80" fill="#FCE4EC" />
                    <circle cx="70" cy="80" r="12" fill="#2d2d2d" />
                    <rect x="60" y="95" width="20" height="30" rx="4" fill={PINK} />
                    <circle cx="140" cy="90" r="12" fill="#2d2d2d" />
                    <rect x="130" y="105" width="20" height="30" rx="4" fill={PINK} />
                    <path d="M80 90 Q110 70 130 95" stroke={PINK} strokeWidth="2" fill="none" strokeDasharray="3 3" />
                    <circle cx="110" cy="70" r="6" fill={PINK} />
                    <circle cx="60" cy="60" r="4" fill={PINK_SOFT} />
                    <circle cx="150" cy="70" r="5" fill={PINK_SOFT} />
                </svg>
            </div>
        );
    }

    if (variant === "onboard3") {
        return (
            <div className={className}>
                <svg viewBox="0 0 200 200" className="h-full w-full">
                    <circle cx="100" cy="100" r="80" fill="#FCE4EC" />
                    <rect x="65" y="80" width="70" height="60" rx="6" fill={PINK} />
                    <path d="M80 80 v-10 a20 20 0 0 1 40 0 v10" stroke="#2d2d2d" strokeWidth="4" fill="none" />
                    <path d="M45 55 l6 -12 l6 12 l12 6 l-12 6 l-6 12 l-6 -12 l-12 -6 Z" fill={PINK} />
                    <path d="M155 140 l5 -10 l5 10 l10 5 l-10 5 l-5 10 l-5 -10 l-10 -5 Z" fill={PINK} />
                    <circle cx="50" cy="150" r="4" fill={PINK_SOFT} />
                    <circle cx="160" cy="60" r="5" fill={PINK_SOFT} />
                </svg>
            </div>
        );
    }

    if (variant === "forgot") {
        return (
            <div className={className}>
                <svg viewBox="0 0 200 200" className="h-full w-full">
                    <circle cx="100" cy="100" r="80" fill="#FCE4EC" />
                    <circle cx="85" cy="70" r="12" fill="#2d2d2d" />
                    <rect x="75" y="85" width="20" height="35" rx="4" fill={PINK} />
                    <line x1="75" y1="100" x2="60" y2="110" stroke="#2d2d2d" strokeWidth="3" />
                    <line x1="95" y1="100" x2="110" y2="110" stroke="#2d2d2d" strokeWidth="3" />
                    <path d="M130 80 l20 10 v18 c0 14 -10 24 -20 30 c-10 -6 -20 -16 -20 -30 v-18 Z" fill={PINK} />
                    <path d="M130 80 l20 10 v18 c0 14 -10 24 -20 30 c-10 -6 -20 -16 -20 -30 v-18 Z" fill="#2d2d2d" opacity="0.1" />
                    <path d="M140 108 l5 5 l12 -12" stroke="#fff" strokeWidth="3" fill="none" strokeLinecap="round" strokeLinejoin="round" />
                </svg>
            </div>
        );
    }

    if (variant === "success") {
        return (
            <div className={className}>
                <svg viewBox="0 0 200 200" className="h-full w-full">
                    <circle cx="100" cy="100" r="80" fill="#FCE4EC" />
                    <rect x="70" y="50" width="60" height="100" rx="10" fill="#2d2d2d" />
                    <rect x="74" y="54" width="52" height="92" rx="6" fill="#fff" />
                    <circle cx="100" cy="100" r="24" fill={PINK} />
                    <path d="M88 100 l8 8 l16 -16" stroke="#fff" strokeWidth="4" fill="none" strokeLinecap="round" strokeLinejoin="round" />
                    <circle cx="50" cy="70" r="4" fill={PINK} />
                    <circle cx="150" cy="80" r="5" fill={PINK} />
                    <circle cx="45" cy="140" r="3" fill={PINK_SOFT} />
                    <circle cx="155" cy="150" r="4" fill={PINK_SOFT} />
                </svg>
            </div>
        );
    }

    return (
        <div className={className}>
            <svg viewBox="0 0 200 200" className="h-full w-full">
                <circle cx="100" cy="100" r="80" fill="#FCE4EC" />
                <rect x="70" y="50" width="60" height="100" rx="10" fill="#2d2d2d" />
                <rect x="74" y="54" width="52" height="92" rx="6" fill="#fff" />
                <circle cx="100" cy="100" r="18" fill={PINK} />
                <circle cx="100" cy="100" r="10" fill="#fff" />
                <circle cx="100" cy="100" r="5" fill={PINK} />
            </svg>
        </div>
    );
}
