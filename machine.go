package machine

import "fmt"

type StateMachine struct {
	start       *State
	end         *State
	main        *Flow
	flows       map[string]*Flow
	transitions map[string][]rune
	events      []func(s *State, r rune)
	dirty       bool
	debug       bool
}

func New() *StateMachine {
	return &StateMachine{
		flows:       make(map[string]*Flow),
		transitions: make(map[string][]rune),
		events:      make([]func(s *State, r rune), 0, 2),
	}
}

func (sm *StateMachine) GetStart() (*State, bool) {
	sm.ResolveIfDirty()
	return sm.start, sm.start != nil
}

func (sm *StateMachine) SetDebug(debug bool) {
	sm.debug = debug
}

func (sm *StateMachine) event(s *State, r rune) {
	for _, f := range sm.events {
		f(s, r)
	}
}

func (sm *StateMachine) DefineTransition(name string, values ...rune) {
	sm.transitions[name] = values
}

func (sm *StateMachine) DefineStringDefineTransition(name string, values string) {
	sm.DefineTransition(name, []rune(values)...)
}

func (sm *StateMachine) Define(flowName string, elements ...string) {
	f, ok := sm.flows[flowName]
	if !ok {
		f = newFlow(flowName)
		sm.flows[flowName] = f
	}
	f.newFlow(sm.flows, sm.transitions, elements...)
	sm.dirty = true
}

func (sm *StateMachine) Match(text string) error {
	sm.ResolveIfDirty()
	if sm.main == nil {
		return fmt.Errorf("'%s' flow not found", SpecialFlowMain)
	}
	if sm.start == nil {
		return fmt.Errorf("'%s' state not found", SpecialStateStart)
	}
	if sm.end == nil {
		return fmt.Errorf("'%s' state not found", SpecialStateEnd)
	}
	return Walk(sm.start, []rune(text), sm.debug, sm.event)
}

func (sm *StateMachine) Print(tags ...string) {
	sm.ResolveIfDirty()
	if sm.start == nil {
		fmt.Printf("'%s' state not found\n", SpecialStateStart)
		return
	}
	Print(sm.start, tags...)
}

func (sm *StateMachine) ResolveIfDirty() {
	if !sm.dirty {
		return
	}
	sm.Resolve()
	sm.dirty = false
}
func (sm *StateMachine) Resolve() {
	for _, f := range sm.flows {
		if f.tag != string(SpecialFlowMain) {
			continue
		}
		sm.main = f
		sm.resolveStart()
		sm.resolveEnd()
		break
	}
}

func (sm *StateMachine) resolveStart() bool {
	start, ok := sm.main.searchStartState()
	if ok {
		start._type = StateStart
		sm.start = start
		return false
	}
	return true
}

func (sm *StateMachine) resolveEnd() bool {
	end, ok := sm.main.searchEndState()
	if ok {
		end._type = StateEnd
		sm.end = end
		return false
	}
	return true
}
