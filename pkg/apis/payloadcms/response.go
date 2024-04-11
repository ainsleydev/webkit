package payloadcms

type (

	// ListResponse represents a list of entities that is sent back
	// from the Payload CMS.
	ListResponse[T any] struct {
		Docs          []T  `json:"docs"`
		Total         int  `json:"total"`
		TotalDocs     int  `json:"totalDocs"`
		Limit         int  `json:"limit"`
		TotalPages    int  `json:"totalPages"`
		Page          int  `json:"page"`
		PagingCounter int  `json:"pagingCounter"`
		HasPrevPage   bool `json:"hasPrevPage"`
		HasNextPage   bool `json:"hasNextPage"`
		PrevPage      any  `json:"prevPage"`
		NextPage      any  `json:"nextPage"`
	}
	// CreateResponse represents a response from the Payload CMS
	// when a new entity is created.
	CreateResponse[T any] struct {
		Doc     T      `json:"doc"`
		Message string `json:"message"`
	}
)
