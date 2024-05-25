import { Validate } from 'payload/dist/fields/config/types';

export const validateURL: Validate<string> = async (value) => {
	if (!value) {
		return true;
	}
	try {
		new URL(value);
		return true;
	} catch (error) {
		return 'Please enter a valid URL';
	}
};

export const validatePostcode: Validate<string> = async (value) => {
	if (!value) {
		return true;
	}
	const postcodeRegex = /^[A-Z]{1,2}\d[A-Z\d]? ?\d[A-Z]{2}$/i;
	if (!postcodeRegex.test(value)) {
		return 'Invalid postcode format';
	}
	return true;
};
