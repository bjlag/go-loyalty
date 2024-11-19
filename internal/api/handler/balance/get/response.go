package get

type Response struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}
