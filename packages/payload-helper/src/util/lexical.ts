import { createHeadlessEditor } from '@lexical/headless';
import { $generateHtmlFromNodes, $generateNodesFromDOM } from '@lexical/html';
// import { sqliteAdapter } from '@payloadcms/db-sqlite';
// import {
// 	defaultEditorConfig,
// 	getEnabledNodes,
// 	lexicalEditor,
// 	sanitizeServerEditorConfig,
// } from '@payloadcms/richtext-lexical';
import { JSDOM } from 'jsdom';
import { $getRoot, $getSelection, type LexicalEditor } from 'lexical';
import type { SerializedEditorState } from 'lexical';
// import { buildConfig, getPayload } from 'payload';
// import { importWithoutClientFiles } from 'payload/node';

// const loadEditor = async (): Promise<LexicalEditor> => {
// 	const config = {
// 		secret: 'testing',
// 		editor: lexicalEditor({
// 			admin: {
// 				hideGutter: false,
// 			},
// 		}),
// 		db: sqliteAdapter({
// 			client: {
// 				url: 'file:./local.db',
// 			},
// 		}),
// 	};
//
// 	const instance = await getPayload({
// 		config: buildConfig(config),
// 	});
//
// 	const editorConfig = await sanitizeServerEditorConfig(defaultEditorConfig, instance.config);
//
// 	return createHeadlessEditor({
// 		nodes: getEnabledNodes({
// 			editorConfig,
// 		}),
// 	});
// };

/**
 * Converts an HTML string to a Lexical editor state.
 *
 * @param {string} html - The HTML string to convert.
 * @returns {SerializedEditorState} The serialized editor state.
 */
export const htmlToLexical = (html: string): SerializedEditorState => {
	const editor = createHeadlessEditor({
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

			//if (selection) selection.insertNodes(nodes);
		},
		{ discrete: true },
	);

	return editor.getEditorState().toJSON();

	// let state = {};
	//
	// loadEditor().then((editor) => {
	// 	editor.update(
	// 		() => {
	// 			// In a headless environment you can use a package such as JSDom to parse the HTML string.
	// 			const dom = new JSDOM(`<!DOCTYPE html><body>${html}</body>`);
	//
	// 			// Once you have the DOM instance it's easy to generate LexicalNodes.
	// 			const nodes = $generateNodesFromDOM(editor, dom.window.document);
	//
	// 			// Select the root
	// 			$getRoot().select();
	//
	// 			// Insert them at a selection.
	// 			const selection = $getSelection();
	//
	// 			if (selection) selection.insertNodes(nodes);
	// 		},
	// 		{ discrete: true },
	// 	);
	//
	// 	state = editor.getEditorState().toJSON();
	// });
	//
	// return state as SerializedEditorState;
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
