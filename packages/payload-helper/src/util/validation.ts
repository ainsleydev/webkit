import { Validate } from 'payload/dist/fields/config/types';

const isValidUrl = (url: string): boolean => {
    try {
        new URL(url);
        return true;
    } catch (error) {
        return false;
    }
}

export const validateURL: Validate<string> = async (value, options) => {
    if (!value) {
        return true;
    }
    if (!isValidUrl(value)) {
        return 'Please enter a valid URL';
    }
    return true;
};


export const validatePostcode: Validate<string> = async (value, options) => {
	if (!value) {
		return true;
	}
	const postcodeRegex = /^[A-Z]{1,2}\d[A-Z\d]? ?\d[A-Z]{2}$/i;
	if (!postcodeRegex.test(value)) {
		return 'Invalid postcode format';
	}
	return true;
};
