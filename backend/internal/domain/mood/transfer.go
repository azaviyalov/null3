package mood

type EditEntryRequest struct {
	Feeling string `json:"feeling" validate:"required"`
	Emoji   string `json:"emoji,omitempty"`
	Note    string `json:"note,omitempty"`
}
