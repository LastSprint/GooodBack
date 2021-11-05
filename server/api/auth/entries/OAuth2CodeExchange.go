package entries

type OAuth2CodeExchange struct {
	Code string
	Provider string
	RedirectUrl string `json:"redirect_url"`
}
