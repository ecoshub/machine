package machine

import (
	"errors"
	"fmt"
	"strings"
)

func Walk(start *State, text []rune, debug bool, f func(s *State, r rune)) error {
	if len(text) == 0 {
		return errors.New("null text passed")
	}
	var state *State = start
	var found bool
	var tmp *State
	if f != nil {
		f(state, 0)
	}
	for i, r := range text {
		tmp, found = walk(state, r, debug, f)
		if !found {
			index := i
			char := text[index]
			// breakx.Point()
			return unexpectedCharError(debug, string(text), state.tag, index, char)
		}
		state = tmp
	}
	ok := isEnd(state, f)
	if !ok {
		index := len(text) - 1
		char := text[index]
		// breakx.Point()
		return unexpectedCharError(debug, string(text), state.tag, index, char)
	}
	return nil
}

var (
	debug bool = false
)

func walk(s *State, r rune, debug bool, f func(s *State, r rune)) (*State, bool) {
	for _, o := range s.output {
		for _, v := range o.values {
			if v == r {
				if f != nil {
					f(o.to, r)
				}
				return o.to, true
			}
		}
	}
	for _, o := range s.output {
		for _, v := range o.values {
			if v == TransitionNone {
				continue
			}
			if v == TransitionFree {
				if f != nil {
					f(o.to, r)
				}
				s, ok := walk(o.to, r, debug, f)
				if ok {
					return s, true
				}
			}
			if v == TransitionAny {
				if f != nil {
					f(o.to, r)
				}
				return o.to, true
			}
		}
	}
	return nil, false
}

func isEnd(current *State, f func(s *State, r rune)) bool {
	if current._type == StateEnd {
		if f != nil {
			f(current, 0)
		}
		return true
	}

	for _, o := range current.output {
		for _, v := range o.values {
			if v == TransitionNone {
				continue
			}
			if v == TransitionFree {
				ok := isEnd(o.to, f)
				if ok {
					return true
				}
			}
		}
	}
	return false
}

func Print(s *State, tags ...string) {
	visit := make(map[*State]bool)
	// st := stable.New(s.tag)
	// st.AddField("first")
	// st.AddField("transition")
	// st.AddField("second")
	print(s, visit, tags...)
	// fmt.Println(st)
}

func print(s *State, visit map[*State]bool, tags ...string) {
	if visit[s] {
		return
	}
	visit[s] = true

	for _, o := range s.output {
		if len(tags) == 0 {
			// st.Row(s.tag, fmt.Sprintf("'%s' %v", printCharString(o.values), o.values), o.to.tag)
			fmt.Printf("%s\t->\t%s\t%v\t(%s)\n", s.tag, o.to.tag, o.values, printCharString(o.values))
		} else {
			for _, t := range tags {
				if strings.Contains(s.tag, t) || strings.Contains(o.to.tag, t) {
					// st.Row(s.tag, fmt.Sprintf("'%s' %v", printCharString(o.values), o.values), o.to.tag)
					fmt.Printf("%s\t->\t%s\t%v\t(%s)\n", s.tag, o.to.tag, o.values, printCharString(o.values))
				}
			}
		}
		print(o.to, visit, tags...)
	}
}

func unexpectedCharError(debug bool, text string, tag string, index int, char rune) error {
	stringChar, intChar := printChar(char)
	if debug {
		return fmt.Errorf("unexpected token '%s' at: %d char_code: %d last_state: %s", stringChar, index, intChar, tag)
	}
	return fmt.Errorf("unexpected token '%s' at: %d char_code: %d", stringChar, index, intChar)
}
