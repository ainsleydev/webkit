package payload

// See: https://docs.slatejs.org/concepts/12-typescript

type T struct {
	Content []struct {
		Type     string `json:"type,omitempty"`
		Children []struct {
			Text string `json:"text,omitempty"`

			// Link
			Url      string `json:"url,omitempty"`
			NewTab   bool   `json:"newTab,omitempty"`
			LinkType string `json:"linkType,omitempty"`
		} `json:"children"`
		Value struct {
			Id int `json:"id"`
		} `json:"value,omitempty"`
		RelationTo string `json:"relationTo,omitempty"`
	} `json:"content"`
}

// CustomElement defines the interface for elements in your schema
type Element interface {
	Type() string
	Children() []Element
}

// ParagraphElement represents a paragraph element
type ParagraphElement struct {
	nodeType string
	children []Element
}

func (p *ParagraphElement) Type() string {
	return p.nodeType
}

func (p *ParagraphElement) Children() []Element {
	return p.children
}

// HeadingElement represents a heading element (example)
type HeadingElement struct {
	nodeType string
	level    int
	children []Element
}

func (h *HeadingElement) Type() string {
	return h.nodeType
}

func (h *HeadingElement) Children() []Element {
	return h.children
}

type HelloElement struct {
}

// CustomNode is the generic interface for nodes in your schema
type CustomNode interface {
	ParagraphElement | HeadingElement
}

// Node represents a node in the Slate editor content
type Node[T CustomNode] struct {
	Type     string    `json:"type"`
	Children []Node[T] `json:"children"`
	Values   T
	// Add other relevant fields from your schema here
}

type RichTextTextProps struct {
	Bold          bool `json:"bold,omitempty"`
	Italic        bool `json:"italic,omitempty"`
	Underline     bool `json:"underline,omitempty"`
	Strikethrough bool `json:"strikethrough,omitempty"`
	Code          bool `json:"code,omitempty"`

	//Value      T           `json:"value,omitempty"`
	//Text       string      `json:"text,omitempty"`
	//LinkType   string      `json:"linkType,omitempty"`
	//RelationTo string      `json:"relationTo,omitempty"`
}

// Walk iterates over a Node and its children, calling the provided function for each node
func (n *Node[T]) Walk(fn func(node *Node[T])) {
	fn(n)
	for _, child := range n.Children {
		child.Walk(fn)
	}
}

type RichTextBlock struct {
	Type     RichTextType    `json:"type,omitempty"`
	Children []RichTextBlock `json:"children,omitempty"`
}

type RichTextType string

var (
	RichTextTypeH1     RichTextType = "h1"
	RichTextTypeH2     RichTextType = "h2"
	RichTextTypeH3     RichTextType = "h3"
	RichTextTypeH4     RichTextType = "h4"
	RichTextTypeH5     RichTextType = "h5"
	RichTextTypeH6     RichTextType = "h6"
	RichTextTypeUL     RichTextType = "ul"
	RichTextTypeOL     RichTextType = "ol"
	RichTextTypeLI     RichTextType = "li"
	RichTextTypeLink   RichTextType = "link"
	RichTextTypeCode   RichTextType = "code"
	RichTextTypeUpload RichTextType = "upload"
	RichTextTypeRel    RichTextType = "relationship"
)

var example = `
[
	{
		"type": "h1",
		"children": [
			{
				"text": "Heading 1"
			}
		]
	},
	{
		"type": "h2",
		"children": [
			{
				"text": "Heading 2"
			}
		]
	},
	{
		"type": "h3",
		"children": [
			{
				"text": "Heading 3"
			}
		]
	},
	{
		"type": "h4",
		"children": [
			{
				"text": "Heading 4"
			}
		]
	},
	{
		"type": "h5",
		"children": [
			{
				"text": "Heading 5"
			}
		]
	},
	{
		"type": "h6",
		"children": [
			{
				"text": "Heading 6"
			}
		]
	},
	{
		"type": "ul",
		"children": [
			{
				"type": "li",
				"children": [
					{
						"text": "Unordered List Item 1"
					}
				]
			},
			{
				"type": "li",
				"children": [
					{
						"text": "Unordered List Item 2"
					}
				]
			},
			{
				"type": "li",
				"children": [
					{
						"text": "Unordered List Item 3"
					}
				]
			}
		]
	},
	{
		"type": "ol",
		"children": [
			{
				"type": "li",
				"children": [
					{
						"text": "Ordered List Item 1"
					}
				]
			},
			{
				"type": "ol",
				"children": [
					{
						"type": "li",
						"children": [
							{
								"text": "Nested List Item 1"
							}
						]
					},
					{
						"type": "li",
						"children": [
							{
								"text": "Ordered List Item 2"
							}
						]
					}
				]
			}
		]
	},
	{
		"children": [
			{
				"text": ""
			},
			{
				"url": "https://ainsley.dev",
				"type": "link",
				"newTab": true,
				"children": [
					{
						"text": "Link"
					}
				],
				"linkType": "custom"
			},
			{
				"text": ""
			}
		]
	},
	{
		"children": [
			{
				"bold": true,
				"text": "Bold Text"
			}
		]
	},
	{
		"children": [
			{
				"text": "Italic Text",
				"italic": true
			}
		]
	},
	{
		"children": [
			{
				"text": "Underline Text",
				"underline": true
			}
		]
	},
	{
		"children": [
			{
				"text": "Strikethrough Text",
				"strikethrough": true
			}
		]
	},
	{
		"children": [
			{
				"code": true,
				"text": "Code"
			}
		]
	},
	{
		"type": "relationship",
		"value": {
			"id": 1
		},
		"children": [
			{
				"code": true,
				"text": ""
			}
		],
		"relationTo": "users"
	},
	{
		"children": [
			{
				"text": ""
			}
		]
	},
	{
		"type": "upload",
		"value": {
			"id": 15
		},
		"children": [
			{
				"text": ""
			}
		],
		"relationTo": "media"
	},
	{
		"children": [
			{
				"text": ""
			}
		]
	}
],`
