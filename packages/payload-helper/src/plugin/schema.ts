import type { JSONSchema4 } from 'json-schema';
import type { Config, Field, SanitizedConfig } from 'payload';

/**
 * General Options for Generating Schema
 */
export interface SchemaOptions {
	useWebKitMedia?: boolean;
	assignRelationships?: boolean;
}

/**
 * This function iterates over properties in JSON schema definitions,
 * passing each property and its key to a callback function.
 */
const loopJSONSchemaProperties = (
	jsonSchema: JSONSchema4,
	callback: (args: { key: string; property: JSONSchema4 }) => void,
): JSONSchema4 => {
	if (!jsonSchema.definitions) {
		return jsonSchema;
	}
	Object.entries(jsonSchema.definitions).forEach(([definitionKey, definition]) => {
		if (definition.properties) {
			Object.entries(definition.properties).forEach(([propertyKey, property]) => {
				callback({ key: propertyKey, property });
			});
		}
	});
	return jsonSchema;
};

/**
 * Adds the necessary GoLang type conversions as a helper func.
 */
export const addGoJSONSchema = (type: string, nillable: boolean): Record<string, unknown> => {
	return {
		goJSONSchema: {
			imports: ['github.com/ainsleydev/webkit/pkg/adapters/payload'],
			nillable: nillable,
			type: type,
		},
	};
};

/**
 * Iterates over all the fields within the config, for both collections
 * and globals, and transforms the JSON schema
 * to include the necessary GoLang schema.
 *
 * @param config
 */
export const fieldMapper = (config: SanitizedConfig, opts: SchemaOptions) => {
	const mapper = (field: Field): Field => {
		switch (field.type) {
			case 'blocks':
				field.typescriptSchema = [() => ({ ...addGoJSONSchema('payload.Blocks', false) })];
				field.blocks.forEach((block) => {
					block.fields = block.fields.map((f) => mapper(f));
				});
				break;
			case 'json':
				field.typescriptSchema = [() => ({ ...addGoJSONSchema('payload.JSON', false) })];
				break;
			case 'richText':
				field.typescriptSchema = [
					() => ({
						type: 'string',
						...addGoJSONSchema('payload.RichText', false),
					}),
				];
				break;
			case 'upload':
				if (opts.useWebKitMedia) {
					const isArray = field.hasMany; // Assuming `hasMany` indicates an array of uploads
					field.typescriptSchema = [
						() => ({
							...(isArray
								? {
										type: 'array',
										items: {
											...addGoJSONSchema(
												'payload.Media',
												field.required === true,
											),
										},
									}
								: { ...addGoJSONSchema('payload.Media', field.required === true) }),
						}),
					];
				}
				break;
			case 'point':
				field.typescriptSchema = [
					() => ({
						...addGoJSONSchema('payload.Point', field.required === true),
					}),
				];
				break;
			case 'tabs': {
				field.tabs.forEach((tab) => {
					tab.fields = tab.fields.map((f) => mapper(f));
				});
				break;
			}
			case 'relationship': {
				if (field.relationTo === 'forms') {
					field.typescriptSchema = [
						() => ({ ...addGoJSONSchema('payload.Form', field.required === true) }),
					];
				}
				break;
			}

			case 'group':
			case 'array':
			case 'row':
			case 'collapsible': {
				if (field.type === 'group' && field.name === 'meta') {
					field.typescriptSchema = [
						() => ({ ...addGoJSONSchema('payload.SettingsMeta', true) }),
					];
				}
				field.fields = field.fields.map((f) => mapper(f));
				break;
			}
			// SEE: https://github.com/ainsleydev/webkit/blob/cdfa078605bec4ee92f2424f69271a0bf6b71366/packages/payload-helper/src/gen/schema.ts#L235
		}

		if (field.type !== 'ui' && opts.assignRelationships) {
			if (!Array.isArray(field.typescriptSchema)) {
				field.typescriptSchema = [];
			}

			if (field.type !== 'tabs' && field.type !== 'row' && field.type !== 'collapsible') {
				field.typescriptSchema.push(({ jsonSchema }) => {
					const payload = {
						name: field.name,
						type: field.type,
						label: field.label,
					} as Record<string, unknown>;

					if (field.type === 'relationship') {
						payload.hasMany = field.hasMany;
						payload.relationTo = field.relationTo;
					}

					return {
						...jsonSchema,
						payload,
					};
				});
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
		config.globals.forEach((global, index) => {
			global.fields = global.fields.map((field) => mapper(field));
		});
	}

	return config;
};

/**
 * Adjusts the JSON schema to include the necessary GoLang schema
 *
 */
export const schemas = (
	opts: SchemaOptions,
): Array<(args: { jsonSchema: JSONSchema4 }) => JSONSchema4> => [
	/**
	 * Removes the auth & uneeded definitions from the schema.
	 */
	({ jsonSchema }): JSONSchema4 => {
		if (!jsonSchema.properties) {
			jsonSchema.properties = {};
		}
		if (!jsonSchema.definitions) {
			jsonSchema.definitions = {};
		}

		if (opts.useWebKitMedia) {
			delete jsonSchema.definitions.media;
			delete jsonSchema.properties?.collections?.properties?.media;
		}

		delete jsonSchema.properties.auth;
		delete jsonSchema.definitions['payload-locked-documents'];
		delete jsonSchema.properties?.collections?.properties?.['payload-locked-documents'];
		delete jsonSchema.definitions.redirects;
		delete jsonSchema.properties?.collections?.properties?.redirects;
		return jsonSchema;
	},
	/**
	 * Adds the settings and media definitions to the schema
	 */
	({ jsonSchema }): JSONSchema4 => {
		if (!jsonSchema.definitions) {
			jsonSchema.definitions = {};
		}

		if ('settings' in jsonSchema.definitions) {
			jsonSchema.definitions.settings = {
				type: 'object',
				fields: [],
				...addGoJSONSchema('payload.Settings', false),
			};
		}

		if ('forms' in jsonSchema.definitions) {
			jsonSchema.definitions.forms = {
				type: 'object',
				...addGoJSONSchema('payload.Form', false),
				fields: [],
			};
		}

		if ('form-submissions' in jsonSchema.definitions) {
			jsonSchema.definitions['form-submissions'] = {
				type: 'object',
				...addGoJSONSchema('payload.FormSubmission', false),
				fields: [],
			};
		}

		return jsonSchema;
	},
	/**
	 * Updates the JSON schema so that it doesn't feature oneOf, so Go doesn't
	 * output it as an interface{}.
	 */
	({ jsonSchema }): JSONSchema4 => {
		const updateRelationship = (property: JSONSchema4) => {
			const payload = property.payload;
			if (!payload) {
				return;
			}

			if (payload.type === 'relationship') {
				if (payload.hasMany) {
					property.type = 'array';
					property.items = {
						$ref: `#/definitions/${payload.relationTo}`,
					};
					return;
				}
				delete property.oneOf;
				property.$ref = `#/definitions/${property.payload.relationTo}`;
			}

			const pType = payload.type;
			if (
				pType === 'group' ||
				pType === 'row' ||
				pType === 'collapsible' ||
				pType === 'array'
			) {
				if (property.properties) {
					for (const k in property.properties) {
						updateRelationship(property.properties[k]);
					}
				}
				return;
			}
		};

		loopJSONSchemaProperties(jsonSchema, ({ property }) => {
			updateRelationship(property);
		});
		return jsonSchema;
	},
	/**
	 * Changes blockType to a string so it's not an *interface{} when
	 * comparing block types in Go.
	 */
	({ jsonSchema }): JSONSchema4 => {
		loopJSONSchemaProperties(jsonSchema, ({ property, key }) => {
			if (key === 'blockType') {
				property.type = 'string';
				delete property.const;
			}
		});
		return jsonSchema;
	},
	/**
	 * Changes blockType to a string so it's not an *interface{} when
	 * comparing block types in Go.
	 */
	({ jsonSchema }): JSONSchema4 => {
		loopJSONSchemaProperties(jsonSchema, ({ property, key }) => {
			const payload = property.payload;
			if (payload && payload.type === 'relationship' && payload.name === 'form') {
				delete property.$ref;
			}
		});
		return jsonSchema;
	},
];
