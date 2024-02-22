package machine

func (sm *StateMachine) Record(entry, escape string, f func(value string)) {
	if f == nil {
		return
	}
	sm.events = append(sm.events, record(entry, escape, f))
}

func record(entry, escape string, f func(value string)) func(s *State, r rune) {
	value := ""
	started := false
	return func(s *State, r rune) {
		tag := s.GetTag()
		switch tag {
		case entry:
			started = true
			value += string(r)
		case escape:
			if started {
				f(value)
				value = ""
				started = false
			}
		default:
			if started {
				value += string(r)
			}
		}
	}
}

func (sm *StateMachine) Listen(f func(s *State, r rune, changed bool), stateTags ...string) {
	if f == nil {
		return
	}
	var lastState *State
	changed := false
	sm.events = append(sm.events, func(s *State, r rune) {
		tag := s.GetTag()
		changed = lastState != s
		if len(stateTags) == 0 {
			f(s, r, changed)
			return
		}
		for _, st := range stateTags {
			if st == tag {
				f(s, r, changed)
			}
		}
	})
}
