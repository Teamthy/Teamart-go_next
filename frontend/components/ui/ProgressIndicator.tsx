interface ProgressIndicatorProps {
    steps: string[];
    currentStep: number;
}

export default function ProgressIndicator({ steps, currentStep }: ProgressIndicatorProps) {
    return (
        <div className="rounded-[28px] border border-zinc-200 bg-white p-4 sm:p-5">
            <div className="flex flex-wrap items-center gap-2">
                {steps.map((step, index) => {
                    const active = index < currentStep;
                    const current = index === currentStep - 1;

                    return (
                        <div key={step} className="flex items-center gap-2">
                            <span
                                className={`inline-flex h-8 w-8 items-center justify-center rounded-full text-sm font-semibold ${current || active ? "bg-[#E91E63] text-white" : "bg-[#FFF8FB] text-zinc-500"
                                    }`}
                            >
                                {index + 1}
                            </span>
                            <span className={`text-sm ${current ? "font-semibold text-zinc-900" : "text-zinc-500"}`}>{step}</span>
                            {index < steps.length - 1 ? <span className="text-zinc-300">/</span> : null}
                        </div>
                    );
                })}
            </div>
        </div>
    );
}
