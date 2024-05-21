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
