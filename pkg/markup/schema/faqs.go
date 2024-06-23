package schemaorg

// FAQPage defines a WebPage presenting one or more "Frequently asked questions"
//
// See:
// - https://schema.org/FAQPage
// - https://developers.google.com/search/docs/data-types/faqpage
//
//	{
//	   "@context": "https://schema.org",
//	   "@type": "FAQPage",
//	   "mainEntity": [
//	       {
//	           "@type": "Question",
//	           "name": "How to find an apprenticeship?",
//	           "acceptedAnswer": {
//	               "@type": "Answer",
//	               "text": "<p>We provide an official service to search through available apprenticeships. To get started, create an account here, specify the desired region, and your preferences. You will be able to search through all officially registered open apprenticeships.</p>"
//	           }
//	       },
//	       {
//	           "@type": "Question",
//	           "name": "Whom to contact?",
//	           "acceptedAnswer": {
//	               "@type": "Answer",
//	               "text": "You can contact the apprenticeship office through our official phone hotline above, or with the web-form below. We generally respond to written requests within 7-10 days."
//	           }
//	       }
//	   ]
//	}
type FAQPage []FAQQuestionAnswer

// FAQQuestionAnswer represents a single question and answer within an FAQPage.
type FAQQuestionAnswer struct {
	// The Question type defines a single answered question within the FAQ.
	// Every Question instance must be contained within the mainEntity
	// property array of the schema.org/FAQPage.
	Question string

	// The full answer to the question. The answer may contain HTML content
	// such as links and lists.
	//
	// Google Search displays the following HTML tags; all other tags are ignored:
	// <h1> through <h6>, <br>, <ol>, <ul>, <li>, <a>, <p>, <div>, <b>, <strong>, <i>, and <em>.
	Answer string
}

// Alias the types for JSON-LD.
type (
	faqs struct {
		Context    string          `json:"@context"`
		Type       string          `json:"@type"`
		MainEntity []faqMainEntity `json:"mainEntity"`
	}
	faqMainEntity struct {
		Type           string    `json:"@type"`
		Name           string    `json:"name"`
		AcceptedAnswer faqAnswer `json:"acceptedAnswer"`
	}
	faqAnswer struct {
		Type string `json:"@type"`
		Text string `json:"text"`
	}
)

// MarshalJSON implements the json.Marshaler interface to generate
// the JSON-LD for the FAQPage.
func (s FAQPage) MarshalJSON() ([]byte, error) {
	f := faqs{
		Context:    Context,
		Type:       "FAQPage",
		MainEntity: make([]faqMainEntity, len(s)),
	}
	for i, qa := range s {
		f.MainEntity[i] = faqMainEntity{
			Type: "Question",
			Name: qa.Question,
			AcceptedAnswer: faqAnswer{
				Type: "Answer",
				Text: qa.Answer,
			},
		}
	}
	return marshal(f)
}
