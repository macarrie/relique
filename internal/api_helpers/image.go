package api_helpers

type ImageSearch struct {
	ModuleName string `json:"module"`
	ClientName string `json:"client"`
	Before     string `json:"before"`
	After      string `json:"after"`
}
