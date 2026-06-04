/**
 * Form error utilities
 */
export function getErrorMessage(err: unknown, fallback?: string): string | null {
    if (!err) return fallback ?? null;
    if (typeof err === "string") return err;
    const e = err as any;
    if (e?.message && typeof e.message === "string") return e.message;
    if (Array.isArray(e?._errors) && e._errors.length) return String(e._errors[0]);
    if (Array.isArray(e) && e.length) return String(e[0]);
    return fallback ?? null;
}
