package mood

type EditEntryRequest struct {
	Feeling string `json:"feeling" validate:"required"`
	Note    string `json:"note,omitempty"`
}
