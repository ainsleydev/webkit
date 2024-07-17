import type { JSONSchema4 } from 'json-schema';
import type { Config, Field } from 'payload';

/**
 * Iterates over all the fields within the config, for both collections
 * and globals, and transforms the JSON schema
 * to include the necessary GoLang schema.
 *
 * @param config
 */
export const fieldMapper = (config: Config) => {
	const mapper = (field: Field): Field => {
		switch (field.type) {
			case 'blocks':
				field.typescriptSchema = [
					() => ({
						goJSONSchema: {
							imports: ['github.com/ainsleydev/webkit/pkg/adapters/payload'],
							nillable: false,
							type: 'payload.Blocks',
						},
					}),
				];
				break;
			case 'richText':
				field.typescriptSchema = [
					() => ({
						type: 'string',
						goJSONSchema: {
							imports: ['github.com/ainsleydev/webkit/pkg/adapters/payload'],
							nillable: false,
							type: 'payload.RichText',
						},
					}),
				];
				break;
			case 'tabs': {
				field.tabs.forEach((tab) => {
					tab.fields = tab.fields.map((f) => mapper(f));
				});
				break;
			}
			case 'array':
			case 'row':
			case 'collapsible': {
				field.fields = field.fields.map((f) => mapper(f));
			}
		}

		return field;
	};

	if (config.collections) {
		config.collections.forEach((collection) => {
			collection.fields = collection.fields.map((field) => mapper(field));
		});
	}

	if (config.globals) {
		config.globals.forEach((global) => {
			global.fields = global.fields.map((field) => mapper(field));
		});
	}

	return config;
};

/**
 * Adjusts the JSON schema to include the necessary GoLang schema
 *
 */
export const schemas: Array<(args: { jsonSchema: JSONSchema4 }) => JSONSchema4> = [
	/**
	 * Removes the auth property from the schema
	 */
	({ jsonSchema }) => {
		if (!jsonSchema.properties) {
			jsonSchema.properties = {};
		}
		// biome-ignore lint/performance/noDelete: <explanation>
		delete jsonSchema.properties.auth;
		return jsonSchema;
	},
	/**
	 * Adds the settings and media definitions to the schema
	 */
	({ jsonSchema }) => {
		if (!jsonSchema.definitions) {
			jsonSchema.definitions = {};
		}

		jsonSchema.definitions.settings = {
			type: 'object',
			additionalProperties: false,
			fields: [],
			goJSONSchema: {
				imports: ['github.com/ainsleydev/webkit/pkg/adapters/payload'],
				nillable: false,
				type: 'payload.Settings',
			},
		};

		jsonSchema.definitions.media = {
			type: 'object',
			additionalProperties: false,
			goJSONSchema: {
				imports: ['github.com/ainsleydev/webkit/pkg/adapters/payload'],
				nillable: false,
				type: 'payload.Media',
			},
		};

		return jsonSchema;
	},
];
