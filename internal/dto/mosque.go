package dto

type ProductItem struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Need      int     `json:"need"`
	Purchased int     `json:"purchased"`
}

type MosqueFull struct {
	ID         uint          `json:"id"`
	Name       string        `json:"name"`
	City       string        `json:"city"`
	Region     string        `json:"region"`
	Requisites string        `json:"requisites"`
	Products   []ProductItem `json:"products"`
}
