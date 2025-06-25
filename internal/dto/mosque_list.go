package dto

type MosqueListItemWithProgress struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	City      string  `json:"city,omitempty"`
	Region    string  `json:"region,omitempty"`
	Need      float64 `json:"need"`
	Collected float64 `json:"collected"`
	Remaining float64 `json:"remaining"`
}

type MosqueResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	City       string `json:"city,omitempty"`
	Region     string `json:"region,omitempty"`
	Requisites string `json:"requisites,omitempty"`
}
