package machine

type StateType uint8

const (
	StateNormal StateType = 0
	StateStart  StateType = 1
	StateEnd    StateType = 2
)

type State struct {
	_type  StateType
	tag    string
	output []*Transition
}

func NewState(optionalTag ...string) *State {
	return newState(resolveOptionalString(optionalTag...), StateNormal)
}

func NewStartState(optionalTag ...string) *State {
	return newState(resolveOptionalString(optionalTag...), StateStart)
}

func NewEndState(optionalTag ...string) *State {
	return newState(resolveOptionalString(optionalTag...), StateEnd)
}

func newState(tag string, _type StateType) *State {
	return &State{
		_type:  _type,
		tag:    tag,
		output: make([]*Transition, 0, 2),
	}
}

func (s *State) GetTag() string {
	return s.tag
}

func (s *State) GetType() StateType {
	return s._type
}

func (s *State) Bind(other *State, values ...rune) *State {
	t := NewTransition(other, values...)
	s.output = append(s.output, t)
	return other
}

func (s *State) BindString(other *State, values string) *State {
	return s.Bind(other, []rune(values)...)
}
