package base

type HTMXSwapType string

type ElementProps struct {
	ID                string
	Classes           []string
	Style             string
	IsDisabled        bool
	IsVisible         bool
	AriaLabel         string
	AriaDescribedByID string
	HTMX              HTMXProps
}

const (
	OuterSwapType HTMXSwapType = "outerHTML"
	InnerSwapType HTMXSwapType = "innerHTML"
)

type HTMXProps struct {
	IsBoosted bool
	// Verb              VerbType
	PostDestination   string
	Encoding          string
	TargetSelector    string
	IndicatorSelector string
	Swap              HTMXSwapType
	SwapOutOfBand     bool
	Include           string
}
