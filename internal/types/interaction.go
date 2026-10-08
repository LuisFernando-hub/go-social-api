package types


type InteractionType string

const (
	InteractionLike    InteractionType = "LIKE"
	InteractionDeslike InteractionType = "DESLIKE"
)

func (t InteractionType) IsValid() bool {
	switch t {
	case InteractionLike, InteractionDeslike:
		return true
	default:
		return false
	}
}