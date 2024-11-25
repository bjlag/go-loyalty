package get

type Response struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
