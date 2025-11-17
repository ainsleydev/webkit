<script lang="ts">
import type { PayloadFormField } from '../../utils/forms/schema';

/**
 * Payload form configuration
 */
export let form: {
	id: number | string;
	fields?: PayloadFormField[] | null;
	confirmationMessage?: string | null;
	submitButtonLabel?: string | null;
};

/**
 * API endpoint to submit form data to (defaults to '/api/forms')
 */
export const apiEndpoint = '/api/forms';

/**
 * Custom submit handler (overrides default API submission)
 */
export const onSubmit: ((data: Record<string, unknown>) => Promise<void>) | undefined = undefined;

const formData: Record<string, unknown> = $state({});
let errors: Record<string, string> = $state({});
let success = $state(false);
let error: string | null = $state(null);
let submitting = $state(false);

const typeMap = { text: 'text', email: 'email', number: 'number' };

async function handleSubmit(event: SubmitEvent) {
	event.preventDefault();
	errors = {};
	error = null;
	submitting = true;

	try {
		if (onSubmit) {
			await onSubmit(formData);
		} else {
			const response = await fetch(apiEndpoint, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({
					formID: form.id,
					data: formData,
				}),
			});

			const result = await response.json();
			if (!result.success) {
				error = result.error || 'Failed to submit form';
				return;
			}
		}

		success = true;
	} catch (err) {
		error = err instanceof Error ? err.message : 'Failed to submit form';
		success = false;
	} finally {
		submitting = false;
	}
}
</script>

<!--
	@component

	PayloadForm renders a form dynamically from Payload CMS form builder fields.
	Uses basic HTML elements that can be styled via CSS variables or custom classes.

	@example
	```svelte
	<PayloadForm
		form={payloadFormData}
		apiEndpoint="/api/forms"
	/>
	```

	@example
	```svelte
	<PayloadForm
		form={payloadFormData}
		onSubmit={async (data) => {
			// Custom submission logic
			await customAPI.submit(data)
		}}
	/>
	```
-->
<form method="POST" onsubmit={handleSubmit} class="payload-form">
	{#if form.fields?.length}
		{#each form.fields as field, index (index)}
			{@const required = field.required ?? false}
			{@const name = field.name}
			<div class="payload-form__group" class:payload-form__group--error={$errors[name]}>
				{#if field.blockType === 'text' || field.blockType === 'email' || field.blockType === 'number'}
					<label for={field.id || name} class="payload-form__label">
						{field.label}
						{#if required}
							<span class="payload-form__required">*</span>
						{/if}
					</label>
					<input
						id={field.id || name}
						{name}
						type={typeMap[field.blockType]}
						placeholder={field.placeholder ?? ''}
						bind:value={formData[name]}
						disabled={success}
						class="payload-form__input"
						class:payload-form__input--error={errors[name]}
					/>
				{:else if field.blockType === 'textarea'}
					<label for={field.id || name} class="payload-form__label">
						{field.label}
						{#if required}
							<span class="payload-form__required">*</span>
						{/if}
					</label>
					<textarea
						id={field.id || name}
						{name}
						placeholder={field.placeholder ?? ''}
						rows={8}
						bind:value={formData[name]}
						disabled={success}
						class="payload-form__textarea"
						class:payload-form__textarea--error={errors[name]}
					/>
				{:else if field.blockType === 'checkbox'}
					<div class="payload-form__checkbox-wrapper">
						<input
							id={field.id || name}
							type="checkbox"
							{name}
							bind:checked={formData[name]}
							disabled={success}
							class="payload-form__checkbox"
							class:payload-form__checkbox--error={errors[name]}
						/>
						<label for={field.id || name} class="payload-form__label payload-form__label--checkbox">
							{field.label}
							{#if required}
								<span class="payload-form__required">*</span>
							{/if}
						</label>
					</div>
				{/if}

				{#if errors[name]}
					<span class="payload-form__error">{errors[name]}</span>
				{/if}
			</div>
		{/each}
	{/if}

	{#if success}
		<div class="payload-form__success">
			{form.confirmationMessage ?? 'Form submitted successfully!'}
		</div>
	{/if}

	{#if error}
		<div class="payload-form__alert">
			{error}
		</div>
	{/if}

	<button type="submit" disabled={success || submitting} class="payload-form__submit">
		{#if submitting}
			<span class="payload-form__loader"></span>
		{/if}
		{form.submitButtonLabel ?? 'Send'}
	</button>
</form>

<style lang="scss">
	.payload-form {
		$self: &;
		--form-gap: 1rem;
		--form-input-padding: 0.75rem;
		--form-input-border: 1px solid #d1d5db;
		--form-input-border-radius: 0.375rem;
		--form-input-bg: #ffffff;
		--form-input-text: #111827;
		--form-error-color: #ef4444;
		--form-success-color: #10b981;
		--form-button-bg: #3b82f6;
		--form-button-text: #ffffff;
		--form-button-hover-bg: #2563eb;
		--form-button-disabled-bg: #9ca3af;

		display: flex;
		flex-direction: column;
		gap: var(--form-gap);
	}

		&__group {
			display: flex;
			flex-direction: column;
			gap: 0.5rem;
		}

		&__label {
			font-size: 0.875rem;
			font-weight: 500;
			color: var(--form-input-text);

			&--checkbox {
				font-weight: 400;
				cursor: pointer;
			}
		}

		&__required {
			color: var(--form-error-color);
			margin-left: 0.25rem;
		}

		&__input,
		&__textarea {
			padding: var(--form-input-padding);
			border: var(--form-input-border);
			border-radius: var(--form-input-border-radius);
			background: var(--form-input-bg);
			color: var(--form-input-text);
			font-size: 1rem;
			font-family: inherit;

			&:focus {
				outline: 2px solid var(--form-button-bg);
				outline-offset: 2px;
			}

			&--error {
				border-color: var(--form-error-color);
			}
		}

		&__textarea {
			resize: vertical;
			min-height: 8rem;
		}

		&__checkbox-wrapper {
			display: flex;
			align-items: flex-start;
			gap: 0.5rem;
		}

		&__checkbox {
			margin-top: 0.25rem;
			width: 1rem;
			height: 1rem;
			cursor: pointer;
		}

		&__error {
			font-size: 0.875rem;
			color: var(--form-error-color);
		}

		&__success {
			padding: 1rem;
			background-color: #d1fae5;
			color: var(--form-success-color);
			border-radius: var(--form-input-border-radius);
			font-weight: 500;
		}

		&__alert {
			padding: 1rem;
			background-color: #fee2e2;
			color: var(--form-error-color);
			border-radius: var(--form-input-border-radius);
			font-weight: 500;
		}

		&__submit {
			padding: var(--form-input-padding) 1.5rem;
			background-color: var(--form-button-bg);
			color: var(--form-button-text);
			border: none;
			border-radius: var(--form-input-border-radius);
			font-size: 1rem;
			font-weight: 500;
			cursor: pointer;
			display: flex;
			align-items: center;
			justify-content: center;
			gap: 0.5rem;
			transition: background-color 0.2s;

			&:hover:not(:disabled) {
				background-color: var(--form-button-hover-bg);
			}

			&:disabled {
				background-color: var(--form-button-disabled-bg);
				cursor: not-allowed;
			}
		}

		&__loader {
			width: 1rem;
			height: 1rem;
			border: 2px solid currentColor;
			border-top-color: transparent;
			border-radius: 50%;
			animation: spin 0.6s linear infinite;
		}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
