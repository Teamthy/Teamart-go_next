/**
 * useFormValidation Hook
 * Integrates React Hook Form with Zod validation schemas
 */

"use client";

import { useCallback } from "react";
import { useForm, UseFormProps, FieldValues, UseFormReturn } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { type ZodType, type input } from "zod";

interface UseFormValidationProps<S extends ZodType<any, FieldValues>>
    extends Omit<UseFormProps<input<S>>, "resolver"> {
    schema: S;
    onSubmit: (data: input<S>) => Promise<void> | void;
}

/**
 * Hook that combines React Hook Form with Zod validation
 * @param schema - Zod schema for validation
 * @param onSubmit - Callback function on valid form submission
 * @param options - Additional React Hook Form options
 * @returns Form methods from React Hook Form
 */
export function useFormValidation<S extends ZodType<any, FieldValues>>({
    schema,
    onSubmit,
    ...options
}: UseFormValidationProps<S>): UseFormReturn<input<S>> & {
    onSubmitHandler: (data: input<S>) => Promise<void>;
    isSubmitting: boolean;
} {
    const form = useForm<input<S>>({
        resolver: zodResolver(schema, undefined, { raw: true }),
        mode: "onChange",
        ...options,
    });

    const onSubmitHandler = useCallback(
        async (data: input<S>) => {
            try {
                form.clearErrors();
                await Promise.resolve(onSubmit(data));
            } catch (error) {
                if (error instanceof Error) {
                    form.setError("root", {
                        message: error.message,
                    });
                }
            }
        },
        [onSubmit, form]
    );

    return {
        ...form,
        onSubmitHandler,
        isSubmitting: form.formState.isSubmitting,
    };
}

export { useForm, useFormContext, useWatch } from "react-hook-form";
export type { FieldValues, UseFormReturn, Controller } from "react-hook-form";
