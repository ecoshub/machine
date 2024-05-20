package machine

import (
	"fmt"
	"os"
)

var (
	core *StateMachine
)

func init() {
	core = createCoreMachine()
}

func createCoreMachine() *StateMachine {
	sm := New()
	sm.DefineStringDefineTransition("start_chars", AlphaUpper+AlphaLower+Underline)
	sm.DefineStringDefineTransition("chars", AlphaNumerical+Underline)

	sm.Define("transition_parser", "_start_", " ", "_start_")
	sm.Define("transition_parser", "_start_", `"`, "value_1_start", TransitionSymbolAny, "val_1", TransitionSymbolAny, "val_1", `"`, "value_end", " ", "value_end", "+", "_start_")
	sm.Define("transition_parser", "_start_", `'`, "value_2_start", TransitionSymbolAny, "val_2", TransitionSymbolAny, "val_2", `'`, "value_end")
	sm.Define("transition_parser", "_start_", "{", "ref_start", "$", "ref_name_start", "$start_chars", "ref_name", "$chars", "ref_name", "}", "value_end")
	sm.Define("transition_parser", "value_end", "\n", "_end_")

	sm.Define("transition_value_parser", "_start_", " ", "_start_")
	sm.Define("transition_value_parser", "_start_", `"`, "value_1_start", TransitionSymbolAny, "val_1", TransitionSymbolAny, "val_1", `"`, "_end_")
	sm.Define("transition_value_parser", "_start_", `'`, "value_2_start", TransitionSymbolAny, "val_2", TransitionSymbolAny, "val_2", `'`, "_end_")
	sm.Define("transition_value_parser", "_start_", "*", "free", TransitionSymbolFree, "_end_")
	sm.Define("transition_value_parser", "_start_", ".", "any", TransitionSymbolFree, "_end_")
	sm.Define("transition_value_parser", "_start_", "{", "ref_start", "$", "ref_name_start", "$start_chars", "ref_name", "$chars", "ref_name", "}", "_end_")

	sm.Define("state_name_parser", "_start_", "$start_chars", "val", "$chars", "val", TransitionSymbolFree, "_end_", " ", "_end_")
	sm.Define("state_name_parser", "_start_", ".", "null", TransitionSymbolFree, "_end_")
	sm.Define("state_name_parser", "val", " ", "_end_")
	sm.Define("state_name_parser", "val", "\n", "_end_")
	sm.Define("state_name_parser", "_start_", "{", "tag_start", "$start_chars", "tag", "$chars", "tag", ":", "ref_start", " ", "ref_start", "@", "ref_name_start", "$start_chars", "ref", "$chars", "ref", "}", "_end_")
	sm.Define("state_name_parser", "tag_start", "@", "ref_name_start")

	sm.Define("_main_", "_start_", "#", "comment", TransitionSymbolAny, "comment", "\n", "_start_", "\n", "_start_")
	sm.Define("_main_", "_start_", "$", "transition_name_start", "$start_chars", "transition_name", "$chars", "transition_name", ":", "transition_value_start", " ", "transition_value_start")
	sm.Define("_main_", "transition_value_start", TransitionSymbolFree, "{main_transition_parser: @transition_parser}", TransitionSymbolFree, "transition_name_end", TransitionSymbolFree, "_start_", TransitionSymbolFree, "_end_", "\n", "_end_")

	sm.Define("_main_", "_start_", "@", "state_name_start", "$start_chars", "state_name", "$chars", "state_name", ":", "state_value_start", " ", "state_value_start", "\n", "state_value_start", "\t", "state_value_start")
	sm.Define("_main_", "state_value_start", "[", "state_1_start", "\n", "state_1_start", "\t", "state_1_start", " ", "state_1_start")
	sm.Define("_main_", "state_1_start", TransitionSymbolFree, "{flow_state_parser_1: @state_name_parser}", TransitionSymbolFree, "state_1_end", ">", "transition_1_start", " ", "transition_1_start")
	sm.Define("_main_", "transition_1_start", TransitionSymbolFree, "{flow_transition_parser: @transition_value_parser}", TransitionSymbolFree, "_space", " ", "_space", ">", "state_2_start", " ", "state_2_start")
	sm.Define("_main_", "state_2_start", TransitionSymbolFree, "{flow_state_parser_2: @state_name_parser}", TransitionSymbolFree, "state_2_end", " ", "state_2_end", "\n", "state_2_end")
	sm.Define("_main_", "state_2_end", ">", "transition_1_start")
	sm.Define("_main_", "state_2_end", "]", "state_end_1", TransitionSymbolFree, "_start_", TransitionSymbolFree, "_end_")
	sm.Define("_main_", "state_2_end", "\n", "state_2_end", "\t", "state_2_end", " ", "state_2_end")
	sm.Define("_main_", "state_2_end", ",", "state_end_2", TransitionSymbolFree, "state_2_start", "\n", "state_2_start", "\t", "state_2_start", " ", "state_2_start")
	return sm
}

func Build(transition []byte) (*StateMachine, error) {
	generatedSM := New()
	transitions := map[string][]rune{}
	transitionReferenceName := ""
	transitionValues := make([]rune, 0, 8)
	core.Record("_main_.transition_name", "_main_.transition_value_start", func(value string) {
		transitionReferenceName = value
	})
	core.Record("main_transition_parser.transition_parser.val_1", "main_transition_parser.transition_parser.value_end", func(value string) {
		transitionValues = append(transitionValues, rawStringToRunes(value)...)
	})
	core.Record("main_transition_parser.transition_parser.val_2", "main_transition_parser.transition_parser.value_end", func(value string) {
		transitionValues = append(transitionValues, rawStringToRunes(value)...)
	})
	core.Record("main_transition_parser.transition_parser.ref_name", "main_transition_parser.transition_parser.value_end", func(value string) {
		v, ok := transitions[value]
		if !ok {
			v, ok = PredefinedTransitions[value]
			if !ok {
				fmt.Println("reference not exists", value)
				os.Exit(1)
				return
			}
		}
		transitionValues = append(transitionValues, v...)
	})

	flowName := ""
	stateAndTransitions := make([]string, 0, 8)
	core.Record("_main_.state_name", "_main_.state_value_start", func(value string) {
		flowName = value
	})

	// state name parse
	core.Record("flow_state_parser_1.state_name_parser.val", "flow_state_parser_1.state_name_parser._end_", func(value string) {
		stateAndTransitions = append(stateAndTransitions, value)
	})
	core.Record("flow_state_parser_1.state_name_parser.tag_start", "flow_state_parser_1.state_name_parser._end_", func(value string) {
		fmt.Printf("value:>%s<\n", value)
		value += "}"
		stateAndTransitions = append(stateAndTransitions, value)
	})
	core.Record("flow_state_parser_1.state_name_parser.null", "flow_state_parser_1.state_name_parser._end_", func(value string) {
		stateAndTransitions = append(stateAndTransitions, value)
	})
	core.Record("flow_state_parser_2.state_name_parser.val", "flow_state_parser_2.state_name_parser._end_", func(value string) {
		stateAndTransitions = append(stateAndTransitions, value)
	})
	core.Record("flow_state_parser_2.state_name_parser.null", "flow_state_parser_2.state_name_parser._end_", func(value string) {
		stateAndTransitions = append(stateAndTransitions, "")
	})
	core.Record("flow_state_parser_2.state_name_parser.tag_start", "flow_state_parser_2.state_name_parser._end_", func(value string) {
		value += "}"
		stateAndTransitions = append(stateAndTransitions, value)
	})

	// transition parse
	core.Record("flow_transition_parser.transition_value_parser.val_1", "flow_transition_parser.transition_value_parser._end_", func(value string) {
		stateAndTransitions = append(stateAndTransitions, string(rawStringToRunes(value)))
	})
	core.Record("flow_transition_parser.transition_value_parser.val_2", "flow_transition_parser.transition_value_parser._end_", func(value string) {
		stateAndTransitions = append(stateAndTransitions, string(rawStringToRunes(value)))
	})
	core.Record("flow_transition_parser.transition_value_parser.free", "flow_transition_parser.transition_value_parser._end_", func(value string) {
		stateAndTransitions = append(stateAndTransitions, TransitionSymbolFree)
	})
	core.Record("flow_transition_parser.transition_value_parser.any", "flow_transition_parser.transition_value_parser._end_", func(value string) {
		stateAndTransitions = append(stateAndTransitions, TransitionSymbolAny)
	})
	core.Record("flow_transition_parser.transition_value_parser.ref_name", "flow_transition_parser.transition_value_parser._end_", func(value string) {
		values, ok := transitions[value]
		if !ok {
			values, ok = PredefinedTransitions[value]
			if !ok {
				fmt.Println("transition reference not found. ref:", value)
				os.Exit(1)
			}
		}
		stateAndTransitions = append(stateAndTransitions, string(values))
	})

	core.Listen(func(s *State, r rune, changed bool) {
		switch s.GetTag() {
		case "main_transition_parser.transition_parser._end_":
			generatedSM.DefineTransition(transitionReferenceName, transitionValues...)
			transitions[transitionReferenceName] = transitionValues
			transitionValues = make([]rune, 0, 8)
			transitionReferenceName = ""
		case "_main_.state_end_1", "_main_.state_end_2":
			if len(stateAndTransitions) != 0 {
				generatedSM.Define(flowName, stateAndTransitions...)
				stateAndTransitions = make([]string, 0, 8)
			}
		}
	})

	core.debug = true

	err := core.Match(string(transition))
	if err != nil {
		return nil, fmt.Errorf("build error. err: %s", err)
	}

	generatedSM.dirty = true
	return generatedSM, nil
}
