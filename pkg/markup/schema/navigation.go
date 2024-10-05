package schemaorg

// NavItemList defines a structured data representation for a navigational list
// of items on a webpage. This helps search engines understand the website's
// structure and potentially improve search ranking.
//
// See:
// - https://schema.org/WebPage
// - https://developers.google.com/search/docs/appearance/structured-data/carousel
//
//	{
//	   "@context": "https://schema.org",
//	   "@type": "BreadcrumbList",
//	   "itemListElement": [
//	       {
//	           "@type": "ListItem",
//	           "position": 1,
//	           "name": "Books",
//	           "item": "https://example.com/books"
//	       },
//	       {
//	           "@type": "ListItem",
//	           "position": 2,
//	           "name": "Science Fiction",
//	           "item": "https://example.com/books/sciencefiction"
//	       },
//	       {
//	           "@type": "ListItem",
//	           "position": 3,
//	           "name": "Award Winners"
//	       }
//	   ]
//	}

//type Navigation []NavItemListElement
//
//// NavigationItem represents a single item within a navigational
//// list on a webpage.
////
//// See: https://schema.org/ItemList
//type NavigationItem struct {
//	Position    int    `json:"position"`    // I.e 1, 2, 3
//	Name        string `json:"name"`        // I.e "Home"
//	AgencyDescription string `json:"description"` // I.e "The homepage of the website" usually the same as the description tag.
//	URL         string `json:"url"`         // Full URL
//}
//
//type T struct {
//	Context         string `json:"@context"`
//	Type            string `json:"@type"`
//	ItemListElement []struct {
//		Type        string `json:"@type"`
//		Position    int    `json:"position"`
//		Name        string `json:"name"`
//		AgencyDescription string `json:"description"`
//		Url         string `json:"url"`
//	} `json:"itemListElement"`
//}

// SiteNavigationElement
// {"@context":"http://schema.org","@type":"ItemList","itemListElement":[{"@type":"SiteNavigationElement","position":1,"name":"Home","description":"","url":"https://ainsley.dev/"},{"@type":"SiteNavigationElement","position":2,"name":"Who we are","description":"Find out why ainsley.dev is a market-leading software design and development company at the forefront of high-tech advancements and techniques.","url":"https://ainsley.dev/who-we-are/"},{"@type":"SiteNavigationElement","position":3,"name":"Services","description":"ainsley.dev offers a wide range of services, from brand strategy, UI\/UX design, website development, and bespoke software development.","url":"https://ainsley.dev/services/"},{"@type":"SiteNavigationElement","position":4,"name":"Portfolio","description":"Be in awe of the award winning projects ainsley.dev has taken on and learn how we helped businesses increase their revenue by crafting stunning designs.","url":"https://ainsley.dev/portfolio/"},{"@type":"SiteNavigationElement","position":5,"name":"Insights","description":"Stay abreast with recent news and insights in the digital sector from leading experts in the website design, development and software fields.","url":"https://ainsley.dev/insights/"},{"@type":"SiteNavigationElement","position":6,"name":"Say hello","description":"Kick off your next project with ainsley.dev. Get in touch with us so we can work together and start your journey to digital freedom.","url":"https://ainsley.dev/contact/"}]}
