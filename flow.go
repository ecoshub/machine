package machine

import (
	"errors"
	"fmt"
	"strings"
)

const (
	stateSeparator string = "."
)

type Flow struct {
	tag    string
	start  *State
	end    *State
	states map[string]*State
}

func newFlow(tag string) *Flow {
	return &Flow{
		tag:    tag,
		states: make(map[string]*State),
	}
}

func (f *Flow) chainTag(tag string) string {
	return f.tag + stateSeparator + tag
}

func (f *Flow) getState(stateTag string) (*State, bool) {
	s, ok := f.states[stateTag]
	return s, ok
}

func (f *Flow) searchStateTag(tag string) (*State, bool) {
	for k, s := range f.states {
		if k == tag {
			return s, true
		}
		tokens := strings.Split(k, stateSeparator)
		for _, t := range tokens {
			if t == tag {
				return s, true
			}
		}
	}
	return nil, false
}

func (f *Flow) searchEndState() (*State, bool) {
	s, ok := f.searchStateTag(string(SpecialStateEnd))
	if ok {
		if s._type == StateEnd {
			return s, true
		}
	}
	return nil, false
}

func (f *Flow) searchStartState() (*State, bool) {
	s, ok := f.searchStateTag(string(SpecialStateStart))
	if ok {
		if s._type == StateStart {
			return s, true
		}
	}
	return nil, false
}

func (f *Flow) newFlow(flows map[string]*Flow, storedTransition map[string][]rune, elements ...string) {
	if len(elements) < 3 || len(elements)%2 != 1 {
		panic("flow must contain at least 2 state and 1 transition")
	}

	states := make([]*State, 0, 16)
	transitions := make([]*Transition, 0, 16)
	for i, tag := range elements {
		if i%2 == 0 {
			if strings.HasPrefix(tag, "{") {
				start, end, err := f.resolveFlowReference(flows, tag)
				if err != nil {
					panic(err)
				}
				states = append(states, start)
				transitions = append(transitions, &Transition{values: []rune{TransitionNone}})
				states = append(states, end)
			} else {
				s := f.resolveState(flows, tag)
				if tag == string(SpecialStateEnd) {
					s._type = StateEnd
				}
				if tag == string(SpecialStateStart) {
					s._type = StateStart
				}
				states = append(states, s)
			}
		} else {
			t := resolveTransition(storedTransition, tag)
			transitions = append(transitions, t)
		}
	}

	for i := 1; i < len(states); i++ {
		f.bindFlows(states, transitions, i)
	}
}

func (f *Flow) resolveFlowReference(flows map[string]*Flow, tag string) (*State, *State, error) {
	tag = strings.TrimPrefix(tag, "{")
	tag = strings.TrimSuffix(tag, "}")

	atIndex := strings.Index(tag, "@")
	if atIndex == -1 {
		return nil, nil, errors.New("reference error. reference must contain '@'")
	}

	ref := tag[atIndex+1:]
	ref = strings.TrimSpace(ref)

	var newFlowName string
	if strings.Contains(tag, ":") {
		tokens := strings.Split(tag, ":")
		newFlowName = tokens[0]
	} else {
		newFlowName = ref + "_" + randomID()
	}

	pointIndex := strings.Index(ref, stateSeparator)
	if pointIndex != -1 {
		ref = ref[:pointIndex]
	}

	f, ok := flows[ref]
	if !ok {
		panic("flow reference not found. name: " + ref)
	}

	ff, err := f.copy(newFlowName, f.chainTag(string(SpecialStateStart)), f.chainTag(string(SpecialStateEnd)))
	if err != nil {
		panic(err)
	}

	flows[newFlowName] = ff
	end, ok := ff.searchEndState()
	if !ok {
		panic(fmt.Sprintf("reference flow must consist '_end_' states. flow: %s", ff.tag))
	}
	ff.end = end
	end._type = StateEnd
	if ff.start == nil {
		panic(fmt.Sprintf("reference flow must consist '_start_' states. flow: %s", ff.tag))
	}
	return ff.start, ff.end, nil
}

func (f *Flow) resolveState(flows map[string]*Flow, tag string) *State {
	if tag == "" || tag == "." {
		return NewState()
	}

	s, ok := f.states[f.chainTag(tag)]
	if ok {
		return s
	}

	s = NewState(f.chainTag(tag))
	f.states[f.chainTag(tag)] = s
	if tag == string(SpecialStateStart) {
		f.start = s
		f.start._type = StateStart
	}
	if tag == string(SpecialStateEnd) {
		f.end = s
		f.end._type = StateEnd
	}
	return s
}

func resolveTransition(storedTransition map[string][]rune, value string) *Transition {
	if len(value) > 1 && strings.HasPrefix(value, "$") {
		ref := strings.TrimPrefix(value, "$")
		var ok bool
		values, ok := storedTransition[ref]
		if !ok {
			values, ok = PredefinedTransitions[ref]
			if !ok {
				panic("transition reference not found. ref: " + ref)
			}
		}
		return &Transition{values: values}
	}

	switch value {
	case TransitionSymbolAny:
		return &Transition{values: []rune{TransitionAny}}
	case TransitionSymbolFree:
		return &Transition{values: []rune{TransitionFree}}
	default:
		return &Transition{values: []rune(value)}
	}
}

func (f *Flow) bindFlows(states []*State, transitions []*Transition, i int) {
	index := i - 1
	s1 := states[index]
	s, ok := f.getState(s1.tag)
	if ok {
		s1 = s
	}

	s2 := states[index+1]
	s, ok = f.getState(s2.tag)
	if ok {
		s2 = s
	}

	t := transitions[index]
	s1.Bind(s2, t.values...)
}

func (f *Flow) copy(newFlowName, from, to string) (*Flow, error) {
	fromState, ok := f.searchStateTag(from)
	if !ok {
		return nil, fmt.Errorf("state not found. tag: %s\n", from)
	}
	visit := make(map[string]bool)
	ff := newFlow(newFlowName)
	tag := ff.chainTag(fromState.tag)
	start := &State{
		_type:  fromState._type,
		tag:    ff.chainTag(fromState.tag),
		output: make([]*Transition, 0, len(fromState.output)),
	}

	ff.states[tag] = start
	ff.start = start
	f.copyCore(fromState, start, ff, to, visit)
	return ff, nil
}

func (f *Flow) copyCore(s, ss *State, other *Flow, stop string, visit map[string]bool) {
	if visit[s.tag] {
		return
	}
	visit[s.tag] = true

	for _, o := range s.output {
		tag := stripTag(o.to.tag)
		ns, ok := other.states[other.chainTag(o.to.tag)]
		if !ok {
			ns = NewState(other.chainTag(o.to.tag))
		}
		if tag == string(SpecialStateStart) {
			ns._type = StateStart
		}
		if tag == string(SpecialStateEnd) {
			ns._type = StateEnd
		}
		ss.Bind(ns, o.values...)
		other.states[other.chainTag(o.to.tag)] = ns
		f.copyCore(o.to, ns, other, stop, visit)
	}

}

func stripTag(tag string) string {
	var t string
	tokens := strings.Split(tag, stateSeparator)
	if len(tokens) == 0 {
		return tag
	}
	t = tokens[len(tokens)-1]
	return t
}
