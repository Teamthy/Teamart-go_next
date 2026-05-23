"use client";

import React from "react";

interface ErrorBoundaryState {
    hasError: boolean;
    error: Error | null;
}

export default class ErrorBoundary extends React.Component<
    { children: React.ReactNode },
    ErrorBoundaryState
> {
    constructor(props: { children: React.ReactNode }) {
        super(props);
        this.state = { hasError: false, error: null };
    }

    static getDerivedStateFromError(error: Error) {
        return { hasError: true, error };
    }

    componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
        console.error("[ErrorBoundary] Caught error:", error, errorInfo);
    }

    resetErrorState = () => {
        this.setState({ hasError: false, error: null });
    };

    render() {
        if (this.state.hasError) {
            return (
                <div className="min-h-screen flex items-center justify-center bg-slate-950 px-4 py-16 text-slate-100">
                    <div className="w-full max-w-2xl rounded-3xl border border-white/10 bg-slate-900/90 p-10 shadow-2xl shadow-black/20">
                        <h1 className="text-3xl font-semibold text-white">Something went wrong</h1>
                        <p className="mt-4 text-slate-300 leading-relaxed">
                            We were not able to render this page. Please refresh or try again in a moment.
                        </p>
                        <div className="mt-8 flex flex-col gap-3 sm:flex-row">
                            <button
                                type="button"
                                onClick={this.resetErrorState}
                                className="inline-flex items-center justify-center rounded-full bg-indigo-600 px-5 py-3 text-sm font-semibold text-white transition hover:bg-indigo-500"
                            >
                                Try again
                            </button>
                            <button
                                type="button"
                                onClick={() => window.location.reload()}
                                className="inline-flex items-center justify-center rounded-full border border-slate-700 bg-transparent px-5 py-3 text-sm font-semibold text-slate-200 transition hover:border-slate-500"
                            >
                                Reload page
                            </button>
                        </div>
                        {this.state.error ? (
                            <pre className="mt-6 rounded-2xl bg-slate-950 p-4 text-xs text-slate-400 overflow-auto">
                                {this.state.error.message}
                            </pre>
                        ) : null}
                    </div>
                </div>
            );
        }

        return this.props.children;
    }
}
