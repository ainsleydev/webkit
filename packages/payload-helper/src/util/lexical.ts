import { createHeadlessEditor } from '@lexical/headless';
import { $generateNodesFromDOM, $generateHtmlFromNodes } from '@lexical/html';
import { $getRoot, $getSelection } from 'lexical';
import { JSDOM } from 'jsdom';
import type { SerializedEditorState } from 'lexical';
import { convertLexicalToHTML } from '@payloadcms/richtext-lexical';

/**
 * Converts an HTML string to a Lexical editor state.
 *
 * @param {string} html - The HTML string to convert.
 * @returns {SerializedEditorState} The serialized editor state.
 */
export const htmlToLexical = (html: string): SerializedEditorState => {
	const editor = createHeadlessEditor({
		nodes: [],
		onError: () => {},
	});

	editor.update(
		() => {
			// In a headless environment you can use a package such as JSDom to parse the HTML string.
			const dom = new JSDOM(html);

			// Once you have the DOM instance it's easy to generate LexicalNodes.
			const nodes = $generateNodesFromDOM(editor, dom.window.document);

			// Select the root
			$getRoot().select();

			// Insert them at a selection.
			const selection = $getSelection();

			if (selection) selection.insertNodes(nodes);
		},
		{ discrete: true },
	);

	return editor.getEditorState().toJSON();
};

/**
 * Converts a Lexical editor state to an HTML string.
 * TODO: This is not working as expected.
 *
 * @param {SerializedEditorState} json - The serialized editor state to convert.
 * @returns {string} The HTML string.
 */
// export const lexicalToHtml = (json: SerializedEditorState): string => {
// 	return convertLexicalToHTML({
// 		data: json,
// 	});
// };
