import { htmlToLexical } from './lexical';

describe('htmlToLexical', () => {
	it('should convert an HTML string to a Lexical editor state', () => {
		const html = '<p>Hello, world!</p>';
		const editorState = htmlToLexical(html);

		expect(editorState).toEqual({
			nodes: [
				{
					type: 'paragraph',
					children: [
						{
							text: 'Hello, world!',
						},
					],
				},
			],
		});
	});
});
