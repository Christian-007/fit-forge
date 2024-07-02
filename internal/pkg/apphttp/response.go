package apphttp

type CollectionRes[Entity any] struct {
	Results []Entity `json:"results"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
