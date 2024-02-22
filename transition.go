package machine

type TransitionType = rune

const (
	TransitionNone TransitionType = 0
	TransitionAny  TransitionType = -1
	TransitionFree TransitionType = -2
)

type Transition struct {
	values []rune
	to     *State
}

func NewTransition(to *State, values ...rune) *Transition {
	return &Transition{
		values: values,
		to:     to,
	}
}
