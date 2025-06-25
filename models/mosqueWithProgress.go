package models

type MosqueWithProgress struct {
	Mosque

	Need      float64 `json:"need"`      // сумма price*need   (общий план)
	Collected float64 `json:"collected"` // сумма price*purchased
	Remaining float64 `json:"remaining"` // Need - Collected
}
