/**
 * useFormValidation Hook
 * Integrates React Hook Form with Zod validation schemas
 */

"use client";

import { useCallback } from "react";
import { useForm, UseFormProps, FieldValues, UseFormReturn } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { ZodSchema } from "zod";

interface UseFormValidationProps<T extends FieldValues>
    extends Omit<UseFormProps<T>, "resolver"> {
    schema: ZodSchema;
    onSubmit: (data: T) => Promise<void> | void;
}

/**
 * Hook that combines React Hook Form with Zod validation
 * @param schema - Zod schema for validation
 * @param onSubmit - Callback function on valid form submission
 * @param options - Additional React Hook Form options
 * @returns Form methods from React Hook Form
 */
export function useFormValidation<T extends FieldValues>({
    schema,
    onSubmit,
    ...options
}: UseFormValidationProps<T>): UseFormReturn<T> & {
    onSubmitHandler: (data: T) => Promise<void>;
    isSubmitting: boolean;
} {
    const form = useForm<T>({
        resolver: zodResolver(schema),
        mode: "onChange",
        ...options,
    });

    const onSubmitHandler = useCallback(
        async (data: T) => {
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
