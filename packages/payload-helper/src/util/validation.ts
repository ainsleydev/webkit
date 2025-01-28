import type { TextFieldSingleValidation } from 'payload';

export const validateURL: TextFieldSingleValidation = async (value) => {
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

export const validatePostcode: TextFieldSingleValidation = async (value) => {
	if (!value) {
		return true;
	}
	const postcodeRegex = /^[A-Z]{1,2}\d[A-Z\d]? ?\d[A-Z]{2}$/i;
	if (!postcodeRegex.test(value)) {
		return 'Invalid postcode format';
	}
	return true;
};
