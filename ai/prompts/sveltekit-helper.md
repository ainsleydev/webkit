# SvelteKit Helper

The majority of web facing applications built with webkit will be SvelteKit with a PayloadCMS
(Next.JS) backend. I want to create a helper npm package for applications using SvelteKit under
packages/sveltekit-helper.

## Requirements

- We should use the latest version of Svelte with runes.
- Applications should be able to run `pnpm add @ainsleydev/sveltekit-helper` and import commonly
  used functions such as `grid/Container.svelte`
- There should be two ways of importing commonly used code (please review this).
	- Via a simple JS import, this would be useful for things like containers, rows or grid elements
	  that don't really change much. But callers are able to customise grid rows etc somehow.
      - I would like to export this cleanly with now `dist` in the import.
	- Via installation (scaffolding) like shadcn. For commonly used components, but change a lot per
	  design, I think it might be best for a scaffolder of sorts where users can scaffold their own
	  files such as `Button.svelte`.
		- I think it makes sense to use the existing architecture and attach this to the Webkit CLI,
		  perhaps `webkit svelte scaffold (something)` using Go.
		- The only concern I have doing it this way is that we're not able to make a nice SvelteKit
		  demo, i.e how would we create a SvelteKit demo when it's scaffolded with Go, updates are
		  harder.
		- If it's scaffolded with Go, how can we ensure the templates are valid Svelte.

## TODO

- Plan and architect this, answering any of the questions below.
- Do not write any files, just come back with architectural design.
- I don't know how it would work for `PayloadForm` as apps may want to define their own `Form` or
  `FormGroup`. Come up with a viable solution.
- Consider a `peerDependencies` approach for SvelteKit to avoid conflicts within the plan.

## Questions

- Should we use https://github.com/shinokada/svelte-lib-helpers as a utility package? What benefit
  would this provide?
- Should we call it `sveltekit-helper` or `svelte-helper`?

## Examples

## `lib/Grid/Container.svelte`

```sveltehtml

<script lang="ts">
	export let fullWidth = false
</script>

<div class="container" class:container--full-width={fullWidth}>
	<slot/>
</div>

<style lang="scss">
	@use '../../scss/abstracts' as a;

	.container {
		$self: &;
		--wrapper-padding-inline: var(--container-padding);
		--wrapper-max-width: var(--container-max-width);
		--breakout-max-width: 1500px;
		--breakout-size: calc((var(--breakout-max-width) - var(--wrapper-max-width)) / 2);

		display: grid;
		width: 100%;
		position: relative;
		grid-template-columns:
			[full-width-start] minmax(var(--wrapper-padding-inline), 1fr)
			[breakout-start] minmax(0, var(--breakout-size))
			[content-start] min(
				100% - (var(--wrapper-padding-inline) * 2),
				var(--wrapper-max-width)
			)
			[content-end]
			minmax(0, var(--breakout-size)) [breakout-end]
			minmax(var(--wrapper-padding-inline), 1fr) [full-width-end];

		:global(> *) {
			grid-column: content;
		}

		:global(> .breakout) {
			grid-column: breakout;
		}

		:global(> .full-width) {
			display: grid;
			grid-column: full-width;
			grid-template-columns: inherit;
		}

		@include a.mq-max(desk) {
			&--full-width {
				display: block;
				grid-template-columns: none;
			}
		}
	}
</style>
```

### `lib/payload/PayloadForm.svelte`

Below is an example of a Payload form that I want to abstract.

```sveltehtml

<script lang="ts">
	import { fly } from 'svelte/transition'

	import Alert from '$lib/components/Alert.svelte'
	import Button from '$lib/components/Button.svelte'
	import Form from '$lib/components/Form/Form.svelte'
	import FormCheckbox from '$lib/components/Form/FormCheckbox.svelte'
	import FormGroup from '$lib/components/Form/FormGroup.svelte'
	import FormInput from '$lib/components/Form/FormInput.svelte'
	import FormLabel from '$lib/components/Form/FormLabel.svelte'
	import FormTextarea from '$lib/components/Form/FormTextarea.svelte'
	import { clientForm } from '$lib/forms/form'
	import { generateFormSchema } from '$lib/forms/schema'
	import { logger } from '$lib/logger'

	import type { FormSubmissionResponse } from '../../routes/api/forms/+server'
	import type { Form as FormCollection } from '$lib/payload/types'

	export let form: FormCollection

	let success = false
	let schema = generateFormSchema(form.fields ?? [])
	const typeMap = {text: 'text', email: 'email', number: 'number'}

	/**
	 * Create a helper form for error state & validation.
	 */
	const {fields, errors, validate, submitting, enhance} = clientForm(
		schema,
		{submissionDelay: 300},
		handleSubmit,
	)

	/**
	 * Sends a POST to /api/forms which then calls Payload.
	 * Form validation is done server side.
	 */
	async function handleSubmit() {
		try {
			const response = await fetch('/api/forms', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({
					formID: form.id,
					data: $fields,
				}),
			})

			const result: FormSubmissionResponse = await response.json()
			if (!result.success) {
				return
			}

			success = true
		} catch (err) {
			logger.error(`Submitting Payload form: ${err}`)
			success = false
		}
	}
</script>

<Form method="POST" useEnhance={enhance}>
	{#if form.fields?.length}
		{#each form.fields as field, index (index)}
			{@const required = field.required ?? false}
			{@const name = field.name}
			<FormGroup
				id={field.id ?? undefined}
				error={$errors[name]}
				align={field.blockType !== 'checkbox' ? 'vertical' : 'horizontal-top'}
				disabled={success}
			>
				{#if field.blockType === 'text' || field.blockType === 'email' || field.blockType === 'number'}
					<FormLabel {required}>{field.label}</FormLabel>
					<FormInput
						{name}
						placeholder={field.placeholder ?? ''}
						type={typeMap[field.blockType]}
						bind:value={$fields[name]}
						on:blur={() => validate({ field: name })}
					/>
				{:else if field.blockType === 'textarea'}
					<FormLabel {required}>{field.label}</FormLabel>
					<FormTextarea
						{name}
						placeholder={field.placeholder ?? ''}
						rows={8}
						bind:value={$fields[name]}
						on:blur={() => validate({ field: name })}
					/>
				{:else if field.blockType === 'checkbox'}
					<FormCheckbox
						{name}
						bind:checked={$fields[name]}
						on:blur={() => validate({ field: name })}
					/>
					<FormLabel {required} bodyText multiline>
						{field.label}
					</FormLabel>
				{/if}
			</FormGroup>
		{/each}
	{/if}
	<!-- Success -->
	{#if success}
		<div class="my-12" transition:fly={{ y: -30, duration: 500 }}>
			<Alert type="success" title="Form Submitted" dismiss>
				{form.confirmationMessage ?? 'Form submitted successfully!'}
			</Alert>
		</div>
	{/if}
	<!-- Submit -->
	<Button block disabled={success} loading={$submitting} type="submit">
		{form.submitButtonLabel ?? 'Send'}
	</Button>
</Form>
```

### `lib/payload/PayloadMedia.svelte`

```sveltehtml

<script lang="ts" module>
	import type { FolderInterface } from '$lib/payload/types'
	import type { FileSizeImproved } from 'payload'
	import type { HTMLAttributes } from 'svelte/elements'

	export type MediaSizes = Record<string, Partial<FileSizeImproved> | undefined>

	export type PayloadMediaProps = HTMLAttributes<HTMLElement> & {
		data: Media
		loading?: 'lazy' | 'eager' | undefined
		className?: string
		breakpointBuffer?: number
		maxWidth?: number | undefined
		onload?: (event: Event) => void
	}

	export type Media = {
		id: number
		alt?: string
		folder?: (number | null) | FolderInterface
		updatedAt: string
		createdAt: string
		deletedAt?: string | null
		url?: string | null
		thumbnailURL?: string | null
		filename?: string | null
		mimeType?: string | null
		filesize?: number | null
		width?: number | null
		height?: number | null
		focalX?: number | null
		focalY?: number | null
		sizes?: MediaSizes
	}
</script>

<script lang="ts">
	let {
		data,
		loading = undefined,
		className = '',
		breakpointBuffer = 50,
		maxWidth = undefined,
		onload = () => null,
		...restProps
	}: PayloadMediaProps = $props()

	/**
	 * Returns sources grouped by breakpoint/width, then sorted by mime priority within each group:
	 * For each width: avif → webp → jpeg/png → others
	 */
	const getSources = (sizesMap?: MediaSizes) => {
		if (!sizesMap) return []

		const arr = Object.values(sizesMap)
			.filter((s) => s?.url)
			.map((s) => ({
				url: s!.url!,
				width: s!.width!,
				mimeType: s!.mimeType ?? '',
			}))

		// Filter by maxWidth if provided.
		const filtered = maxWidth ? arr.filter((s) => s.width <= maxWidth) : arr

		// Mime type priority order.
		const mimePriority = (mime: string, url: string) => {
			if (/avif/i.test(mime) || /\.avif$/i.test(url)) return 1
			if (/webp/i.test(mime) || /\.webp$/i.test(url)) return 2
			if (/jpe?g|png/i.test(mime) || /\.(jpe?g|png)$/i.test(url)) return 3
			return 4 // others
		}

		// Sort by width first, then mime priority.
		return filtered.sort((a, b) => {
			if (a.width !== b.width) return a.width - b.width
			return mimePriority(a.mimeType, a.url) - mimePriority(b.mimeType, b.url)
		})
	}

	/**
	 * Mime Detectors - now reactive with $derived
	 */
	let isImage = $derived((data?.mimeType ?? '').startsWith('image'))
	let isVideo = $derived((data?.mimeType ?? '').startsWith('video'))
	let isSVG = $derived(!!(data?.filename && /\.svg$/i.test(data.filename)))

	/**
	 * Constants for media - now reactive with $derived
	 */
	let sources = $derived(isImage && !isSVG ? getSources(data.sizes) : [])
	let fallbackSource = $derived(sources[sources.length - 1])
	let imgAlt = $derived(data?.alt ?? '')
	let imgLoading = $derived(loading ?? undefined)
	let imgWidth = $derived(fallbackSource?.width ?? data?.width ?? undefined)
	let imgHeight = $derived(
		fallbackSource && data?.width && data?.height
			? Math.round((fallbackSource.width / data.width) * data.height)
			: (data?.height ?? undefined),
	)
	let fallbackUrl = $derived(sources[sources.length - 1]?.url ?? data?.url ?? '')
</script>

{#if isImage}
	{#if isSVG}
		<img
			src={data.url}
			alt={imgAlt}
			width={imgWidth}
			height={imgHeight}
			loading={imgLoading}
			class={className}
			{onload}
			{...restProps}
		/>
	{:else}
		<picture class={className} {...restProps}>
			{#each sources as source, index (index)}
				<source
					media={`(max-width: ${source.width + breakpointBuffer}px)`}
					srcset={source.url}
					type={source.mimeType}
				/>
			{/each}
			<img
				src={fallbackUrl}
				alt={imgAlt}
				width={imgWidth}
				height={imgHeight}
				loading={imgLoading}
				{onload}
			/>
		</picture>
	{/if}
{:else if isVideo}
	<video
		controls
		width={imgWidth}
		height={imgHeight}
		preload={loading === 'lazy' ? 'metadata' : 'auto'}
		poster={data.thumbnailURL}
		class={className}
		{onload}
		{...restProps}
	>
		<source src={data.url} type={data.mimeType}/>
		<track kind="captions"/>
	</video>
{/if}

<style lang="scss">
	img,
	video {
		max-width: 100%;
		height: auto;
		display: block;
	}
</style>
```
