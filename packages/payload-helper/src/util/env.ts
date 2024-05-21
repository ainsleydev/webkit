function hasKey(key: string): boolean {
	return Object.prototype.hasOwnProperty.call(process.env, key);
}

function envFn<T>(key: string, defaultValue?: T): string | T | undefined {
	return hasKey(key) ? process.env[key] : defaultValue;
}

function getKey(key: string): string {
	return process.env[key] ?? '';
}

const utils = {
	isProduction: getKey('NODE_ENV') === 'production',

	int(key: string, defaultValue?: number): number | undefined {
		if (!hasKey(key)) {
			return defaultValue;
		}

		return parseInt(getKey(key), 10);
	},

	float(key: string, defaultValue?: number): number | undefined {
		if (!hasKey(key)) {
			return defaultValue;
		}

		return parseFloat(getKey(key));
	},

	bool(key: string, defaultValue?: boolean): boolean | undefined {
		if (!hasKey(key)) {
			return defaultValue;
		}

		return getKey(key) === 'true';
	},

	json(key: string, defaultValue?: object) {
		if (!hasKey(key)) {
			return defaultValue;
		}

		try {
			return JSON.parse(getKey(key));
		} catch (error) {
			if (error instanceof Error) {
				throw new Error(`Invalid json environment variable ${key}: ${error.message}`);
			}

			throw error;
		}
	},

	array(key: string, defaultValue?: string[]): string[] | undefined {
		if (!hasKey(key)) {
			return defaultValue;
		}

		let value = getKey(key);

		if (value.startsWith('[') && value.endsWith(']')) {
			value = value.substring(1, value.length - 1);
		}

		return value.split(',').map((v) => {
			return v.trim().replace(/^"(.*)"$/, '$1');
		});
	},

	date(key: string, defaultValue?: Date): Date | undefined {
		if (!hasKey(key)) {
			return defaultValue;
		}

		return new Date(getKey(key));
	},

	oneOf(key: string, expectedValues?: unknown[], defaultValue?: unknown) {
		if (!expectedValues) {
			throw new Error(`env.oneOf requires expectedValues`);
		}

		if (defaultValue && !expectedValues.includes(defaultValue)) {
			throw new Error(`env.oneOf requires defaultValue to be included in expectedValues`);
		}

		const rawValue = envFn(key, defaultValue);
		return expectedValues.includes(rawValue) ? rawValue : defaultValue;
	},
};

const env = Object.assign(envFn, utils);

export default env;
