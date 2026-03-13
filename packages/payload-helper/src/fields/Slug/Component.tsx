'use client';

import {
	Button,
	FieldDescription,
	FieldLabel,
	TextInput,
	useField,
	useForm,
	useFormFields,
} from '@payloadcms/ui';
import type { TextFieldClientProps } from 'payload';
import type React from 'react';
import { useCallback, useEffect } from 'react';

import { formatSlug } from './formatSlug.js';
import './index.scss';

type SlugComponentProps = {
	fieldToUse: string;
	checkboxFieldPath: string;
} & TextFieldClientProps;

export const Component: React.FC<SlugComponentProps> = ({
	field,
	fieldToUse,
	checkboxFieldPath: checkboxFieldPathFromProps,
	path,
	readOnly: readOnlyFromProps,
}) => {
	const {
		admin: { description } = {},
		label,
	} = field;

	const checkboxFieldPath = path?.includes('.')
		? `${path}.${checkboxFieldPathFromProps}`
		: checkboxFieldPathFromProps;
	const resolvedPath = path || field.name;

	const { value, setValue } = useField<string>({ path: resolvedPath });
	const { dispatchFields } = useForm();

	const checkboxValue = useFormFields(([fields]) => {
		return fields[checkboxFieldPath]?.value as string;
	});

	const targetFieldValue = useFormFields(([fields]) => {
		return fields[fieldToUse]?.value as string;
	});

	useEffect(() => {
		if (checkboxValue) {
			if (targetFieldValue) {
				const formattedSlug = formatSlug(targetFieldValue);

				if (value !== formattedSlug) {
					setValue(formattedSlug);
				}
			} else if (value !== '') {
				setValue('');
			}
		}
	}, [targetFieldValue, checkboxValue, setValue, value]);

	const handleLock = useCallback(
		(event: React.MouseEvent<Element>) => {
			event.preventDefault();

			dispatchFields({
				type: 'UPDATE',
				path: checkboxFieldPath,
				value: !checkboxValue,
			});
		},
		[checkboxValue, checkboxFieldPath, dispatchFields],
	);

	const readOnly = readOnlyFromProps || checkboxValue;

	return (
		<div className='field-type slug-field-component'>
			<div className='label-wrapper'>
				<FieldLabel htmlFor={`field-${resolvedPath}`} label={label} />
				<Button className='lock-button' buttonStyle='none' onClick={handleLock}>
					{checkboxValue ? 'Unlock' : 'Lock'}
				</Button>
			</div>
			<TextInput
				value={value}
				onChange={setValue}
				path={resolvedPath}
				readOnly={Boolean(readOnly)}
			/>
			<FieldDescription
				className={`field-description-${resolvedPath.replace(/\./g, '__')}`}
				description={description ?? ''}
				path={resolvedPath}
			/>
		</div>
	);
};
