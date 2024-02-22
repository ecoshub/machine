package machine

type SpecialState string
type SpecialFlow string

const (
	SpecialFlowMain   SpecialFlow  = "_main_"
	SpecialStateStart SpecialState = "_start_"
	SpecialStateEnd   SpecialState = "_end_"
)

const (
	TransitionSymbolAny  string = "**"
	TransitionSymbolFree string = ">>"
)
