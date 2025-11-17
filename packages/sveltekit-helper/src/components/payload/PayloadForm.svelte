<script lang="ts">
import { clientForm } from '../../utils/forms/clientForm';
import { type PayloadFormField, generateFormSchema } from '../../utils/forms/schema';

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
export const onSubmit: ((data: Record<string, never>) => Promise<void>) | undefined = undefined;

let success = false;
let error: string | null = null;
const schema = generateFormSchema(form.fields ?? []);
const typeMap = { text: 'text', email: 'email', number: 'number' };

const { fields, errors, validate, submitting, enhance } = clientForm(
	schema,
	{ submissionDelay: 300 },
	onSubmit || handleDefaultSubmit,
);

async function handleDefaultSubmit() {
	try {
		error = null;
		const response = await fetch(apiEndpoint, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({
				formID: form.id,
				data: $fields,
			}),
		});

		const result = await response.json();
		if (!result.success) {
			error = result.error || 'Failed to submit form';
			return;
		}

		success = true;
	} catch (err) {
		error = err instanceof Error ? err.message : 'Failed to submit form';
		success = false;
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
<form method="POST" use:enhance class="payload-form">
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
						bind:value={$fields[name]}
						on:blur={() => validate({ field: name })}
						disabled={success}
						class="payload-form__input"
						class:payload-form__input--error={$errors[name]}
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
						bind:value={$fields[name]}
						on:blur={() => validate({ field: name })}
						disabled={success}
						class="payload-form__textarea"
						class:payload-form__textarea--error={$errors[name]}
					/>
				{:else if field.blockType === 'checkbox'}
					<div class="payload-form__checkbox-wrapper">
						<input
							id={field.id || name}
							type="checkbox"
							{name}
							bind:checked={$fields[name]}
							on:blur={() => validate({ field: name })}
							disabled={success}
							class="payload-form__checkbox"
							class:payload-form__checkbox--error={$errors[name]}
						/>
						<label for={field.id || name} class="payload-form__label payload-form__label--checkbox">
							{field.label}
							{#if required}
								<span class="payload-form__required">*</span>
							{/if}
						</label>
					</div>
				{/if}

				{#if $errors[name]}
					<span class="payload-form__error">{$errors[name]}</span>
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

	<button type="submit" disabled={success || $submitting} class="payload-form__submit">
		{#if $submitting}
			<span class="payload-form__loader"></span>
		{/if}
		{form.submitButtonLabel ?? 'Send'}
	</button>
</form>

<style>
	.payload-form {
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

	.payload-form__group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.payload-form__label {
		font-size: 0.875rem;
		font-weight: 500;
		color: var(--form-input-text);
	}

	.payload-form__label--checkbox {
		font-weight: 400;
		cursor: pointer;
	}

	.payload-form__required {
		color: var(--form-error-color);
		margin-left: 0.25rem;
	}

	.payload-form__input,
	.payload-form__textarea {
		padding: var(--form-input-padding);
		border: var(--form-input-border);
		border-radius: var(--form-input-border-radius);
		background: var(--form-input-bg);
		color: var(--form-input-text);
		font-size: 1rem;
		font-family: inherit;
	}

	.payload-form__input:focus,
	.payload-form__textarea:focus {
		outline: 2px solid var(--form-button-bg);
		outline-offset: 2px;
	}

	.payload-form__input--error,
	.payload-form__textarea--error {
		border-color: var(--form-error-color);
	}

	.payload-form__textarea {
		resize: vertical;
		min-height: 8rem;
	}

	.payload-form__checkbox-wrapper {
		display: flex;
		align-items: flex-start;
		gap: 0.5rem;
	}

	.payload-form__checkbox {
		margin-top: 0.25rem;
		width: 1rem;
		height: 1rem;
		cursor: pointer;
	}

	.payload-form__error {
		font-size: 0.875rem;
		color: var(--form-error-color);
	}

	.payload-form__success {
		padding: 1rem;
		background-color: #d1fae5;
		color: var(--form-success-color);
		border-radius: var(--form-input-border-radius);
		font-weight: 500;
	}

	.payload-form__alert {
		padding: 1rem;
		background-color: #fee2e2;
		color: var(--form-error-color);
		border-radius: var(--form-input-border-radius);
		font-weight: 500;
	}

	.payload-form__submit {
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
	}

	.payload-form__submit:hover:not(:disabled) {
		background-color: var(--form-button-hover-bg);
	}

	.payload-form__submit:disabled {
		background-color: var(--form-button-disabled-bg);
		cursor: not-allowed;
	}

	.payload-form__loader {
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
