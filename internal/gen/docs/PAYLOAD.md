## Payload


### Code Style

- Use `camelCase` for all field names.
- Always include `admin.description` for Payload collections and fields.
- Collection slugs should be lowercase with hyphens (e.g. `'media'`, `'form-submissions'`).
- Leverage helper functions from `payload-helper` package for common patterns.

### Collection Configuration

Collections follow Payload's `CollectionConfig` type with specific conventions:

```typescript
export const Listings: CollectionConfig = {
	slug: 'listings',
	timestamps: true,
	trash: true,
	versions: {
		drafts: true,
		maxPerDoc: 10,
	},
	admin: {
		useAsTitle: 'title',
		defaultColumns: ['id', 'title', 'vehicle', 'manufacturer'],
		preview: (doc): string | null => {
			return `${frontendURL()}${doc.url}`
		},
	},
	access: {
		read: adminsAndCreatedBy,
		create: ({ req }) => !!req.user,
		update: adminsAndCreatedBy,
		delete: admins,
		admin: admins,
	},
	fields: [
		{
			name: 'title',
			type: 'text',
			required: true,
			maxLength: 80,
			admin: {
				description: `What part are you selling? Keep it clear and under 80 characters.`,
			},
		},
		// ... more fields
	],
}
```

### Access Control

- Create separate access control functions for reusability.
- Use role-based access control (RBAC).
- Define condition functions for field-level visibility.

**Example:**

```typescript
import { checkRole } from './checkRole'

export const admins = (args: { req: { user?: any } }): boolean => {
	const user = args?.req?.user
	return checkRole(Roles.Admin, user) || checkRole(Roles.SuperAdmin, user)
}

export const adminOnly = {
	read: admins,
	create: admins,
	update: admins,
	delete: admins,
	unlock: admins,
}
```

### Preview URLs

Always configure preview URLs for collections that have corresponding frontend pages:

```typescript
admin: {
	preview: (doc): string | null => {
		return `${frontendURL()}${doc.url}`
	},
}
```



### Reusable Fields

For fields that occur more than once within the codebase, they should be abstract within `src/fields`. Every field that
is configurable should accept overrides so the caller can override particular parts of the fieldd.

#### FAQs Example

```typescript
import { ArrayField, deepMerge, Field } from 'payload'

export type FAQsFieldArgs = {
	overrides?: Partial<Omit<ArrayField, 'type'>>
}

export const FAQsField = (args?: FAQsFieldArgs): Field => {
	return deepMerge<ArrayField, Omit<ArrayField, 'type'>>(
		{
			name: 'faqs',
			label: 'FAQs',
			type: 'array',
			fields: [
				{
					name: 'question',
					label: 'Question',
					type: 'text',
					required: true,
					admin: {
						description: 'Add a question for the FAQ item.',
					},
				},
				{
					name: 'answer',
					type: 'textarea',
					label: 'Answer',
					required: true,
					admin: {
						description: 'Add a content (answer) for the FAQ item.',
					},
				},
			],
		},
		args?.overrides || {},
	)
}
```

#### Slug Example

```typescript
export const SlugField: Slug = (fieldToUse = 'title', overrides = {}) => {
	const checkBoxField = deepMerge<CheckboxField, Partial<CheckboxField>>(
		{
			name: 'slugLock',
			type: 'checkbox',
			defaultValue: true,
			admin: {hidden: true, position: 'sidebar'},
		},
		checkboxOverrides || {},
	)

	const slugField = deepMerge<TextField, Partial<TextField>>(
		{
			name: 'slug',
			type: 'text',
			index: true,
			unique: true,
			required: true,
			hooks: {
				beforeValidate: [formatSlugHook(fieldToUse)],
			},
			admin: {
				position: 'sidebar',
				components: {
					Field: {
						path: '/fields/Slug/Component#Component',
						clientProps: {fieldToUse, checkboxFieldPath: checkBoxField.name},
					},
				},
			},
		},
		slugOverrides || {},
	)

	return [slugField, checkBoxField]
}
```

#### Key Patterns

- **Use deepMerge** for composing field configurations with overrides.
- **Provide default values** but allow customisation via overrides.
- **Co-locate fields** that work together (e.g., slug and slugLock).
- **Custom components** via path references for admin UI customisation.

### Width

Use `admin.width` to control field widths in rows:

```typescript
{
	type: 'row',
	fields: [
		{
			name: 'fieldOne',
			type: 'text',
			admin: { width: '50%' },
		},
		{
			name: 'fieldTwo',
			type: 'text',
			admin: { width: '50%' },
		},
	],
}
```



Use hooks for data transformation and business logic. Hooks should be placed in a separate file under `hooks/{file}.ts`
within the collection folder and have a test along side it.

**Example**:

`src/collections/{collection}/setConnectedAt.ts`

```typescript
export const setConnectedAt: CollectionBeforeChangeHook<Connection> = async (args) => {
	const { data, originalDoc } = args

	// Only process if status field exists in the data.
	if (!data?.status) {
		return data
	}

	// If status is being changed to 'accepted', set the connectedAt timestamp.
	if (data.status === 'accepted' && originalDoc?.status !== 'accepted') {
		data.connectedAt = new Date().toISOString()
	}

	// If status is changed from 'accepted' to something else, clear the connectedAt timestamp.
	if (data.status !== 'accepted' && originalDoc?.status === 'accepted') {
		data.connectedAt = null
	}

	return data
}
```

`src/collections/{collection}/setConnectedAt.test.ts`

```typescript
import { getPayload, Payload } from 'payload'
import { describe, it, beforeAll, expect, beforeEach } from 'vitest'

import config from '@/payload.config'
import { createTestPlayer } from '@/test/fixtures'
import { teardown } from '@/test/util'

let payload: Payload

describe('setConnectedAt hook', () => {
	beforeAll(async () => {
		const payloadConfig = await config
		payload = await getPayload({ config: payloadConfig })
	})

	beforeEach(async () => {
		await teardown(payload)
	})

	it('Sets connectedAt timestamp when connection is accepted', async () => {
		const player1 = await createTestPlayer(payload)
		const player2 = await createTestPlayer(payload)

		const connection = await payload.create({
			collection: 'connections',
			data: {
				requester: player1.id,
				recipient: player2.id,
				status: 'pending',
			},
		})

		// Initially should not have connectedAt.
		expect(connection.connectedAt).toBeFalsy()

		// Update to accepted status.
		const updated = await payload.update({
			collection: 'connections',
			id: connection.id,
			data: {
				status: 'accepted',
			},
		})

		// Should now have connectedAt timestamp.
		expect(updated.connectedAt).toBeDefined()
		expect(typeof updated.connectedAt).toBe('string')
	})
})
```


