import { createHeadlessEditor } from '@lexical/headless';
import { $generateNodesFromDOM, $generateHtmlFromNodes } from '@lexical/html';
import { $getRoot, $getSelection } from 'lexical';
import { JSDOM } from 'jsdom';
import type { SerializedEditorState } from 'lexical';
import { getEnabledNodes, sanitizeServerEditorConfig, defaultEditorConfig } from '@payloadcms/richtext-lexical'
import type { SanitizedConfig} from "payload";

/**
 * Converts an HTML string to a Lexical editor state.
 *
 * @param {string} html - The HTML string to convert.
 * @param config
 * @returns {SerializedEditorState} The serialized editor state.
 */
export const htmlToLexical = async (html: string, config: SanitizedConfig): Promise<SerializedEditorState> => {
	//const cfg = await sanitizeServerEditorConfig(defaultEditorConfig, config)

	const editor = createHeadlessEditor({
		// nodes: getEnabledNodes({
		// 	editorConfig: cfg,
		// }),
		nodes: [],
		onError: (error) => {
			console.error(error);
		},
	});

	editor.update(
		() => {
			// In a headless environment you can use a package such as JSDom to parse the HTML string.
			const dom = new JSDOM(`<!DOCTYPE html><body>${html}</body>`);

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
 *
 * @param {SerializedEditorState} json - The serialized editor state to convert.
 * @returns {string} The HTML string.
 */
export const lexicalToHtml = (json: SerializedEditorState): string => {
	const editor = createHeadlessEditor({
		nodes: [],
		onError: (error) => {
			console.error(error);
		},
	});

	// Initialize a JSDOM instance
	const dom = new JSDOM('');

	// @ts-ignore
	globalThis.window = dom.window;
	globalThis.document = dom.window.document;

	editor.update(() => {
		const editorState = editor.parseEditorState(json);
		editor.setEditorState(editorState);
	});

	// Convert the editor state to HTML
	let html = '';
	editor.getEditorState().read(() => {
		html = $generateHtmlFromNodes(editor);
	});

	return html;
};
