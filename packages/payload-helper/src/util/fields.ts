import { Field } from 'payload';

/**
 * Determines if a Payload field has a name property.
 *
 * @param field
 */
export const fieldHasName = (field: Field): boolean => {
	if (field.type === "ui") {
		return false;
	}
	return field.type !== "tabs" && field.type !== "row" && field.type !== "collapsible";
}
