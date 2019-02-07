// machine
package sm

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ActionFunc func(string, Event, *Context)
type GuardFunc func(string, Event, *Context) bool

type Transition struct {
	OnEvent string
	Guard   string
	ToState string
	Action  string
}

type State struct {
	Name        string
	OnEntry     string //func(ev event, context interface{})
	OnExit      string //func(ev event, context interface{})
	Transitions []Transition
}

func (m Machine) getTransition(s State, ev Event, ctxt *Context) (Transition, bool) {
	n := ev.GetName()
	for _, t := range s.Transitions {
		if t.OnEvent == n {
			if t.Guard == "" || m.GuardFuncs[t.Guard](s.Name, ev, ctxt) {
				return t, true
			}
		}
	}
	return Transition{}, false
}

type Machine struct {
	Name        string
	States      []State
	ActionFuncs map[string]ActionFunc
	GuardFuncs  map[string]GuardFunc
}

func (m Machine) AddActionFuncs(fs ...ActionFunc) Machine {
	for _, f := range fs {
		m.ActionFuncs[getFunctionName(f)] = f
	}
	return m
}

func (m Machine) AddGuardFuncs(fs ...GuardFunc) Machine {
	for _, f := range fs {
		m.GuardFuncs[getFunctionName(f)] = f
	}
	return m
}

func (m Machine) getState(name string) (State, bool) {
	for _, s := range m.States {
		if s.Name == name {
			return s, true
		}
	}
	return State{}, false
}

func (m Machine) HandleEvent(ev Event, ctxt *Context) bool {
	cs := ctxt.currentStateName
	currentState, ok := m.getState(cs)
	if !ok {
		panic("Machine doesn't have current state " + cs)
	} else {
		t, ok := m.getTransition(currentState, ev, ctxt)
		if !ok {
			fmt.Printf("Machine %v in state %v has no transition for event %v\n", m.Name, currentState.Name, ev.GetName())
			return false
		}
		fmt.Printf("Machine %v in state %v received event %v, transition to state %v\n", m.Name, currentState.Name, ev.GetName(), t.ToState)
		if currentState.OnExit != "" {
			m.ActionFuncs[currentState.OnExit](currentState.Name, ev, ctxt)
		}
		if t.Action != "" {
			m.ActionFuncs[t.Action](currentState.Name, ev, ctxt)
		}
		newState, ok := m.getState(t.ToState)
		if !ok {
			panic("Machine doesn't have new state " + t.ToState)
		} else {
			ctxt.currentStateName = t.ToState
			if newState.OnEntry != "" {
				m.ActionFuncs[newState.OnEntry](newState.Name, ev, ctxt)
			}
		}
	}
	return true
}

func (m Machine) Validate() error {
	for _, s := range m.States {
		if s.OnEntry != "" {
			if _, ok := m.ActionFuncs[s.OnEntry]; !ok {
				return errors.New("State " + s.Name + " onEntry function " + s.OnEntry + " not registered")
			}
		}
		if s.OnExit != "" {
			if _, ok := m.ActionFuncs[s.OnExit]; !ok {
				return errors.New("State " + s.Name + " onExit function " + s.OnExit + " not registered")
			}
		}
		for _, t := range s.Transitions {
			if t.Action != "" {
				if _, ok := m.ActionFuncs[t.Action]; !ok {
					return errors.New("state:" + s.Name + " transition:" + t.OnEvent + " action function:" + t.Action + " not registered")
				}
			}
			if t.Guard != "" {
				if _, ok := m.GuardFuncs[t.Guard]; !ok {
					return errors.New("state:" + s.Name + " transition:" + t.OnEvent + " guard function:" + t.Guard + " not registered")
				}
			}
			found := false
			for _, s2 := range m.States {
				if s2.Name == t.ToState {
					found = true
					break
				}
			}
			if !found {
				return errors.New("state:" + s.Name + " transition:" + t.OnEvent + " destination:" + t.ToState + " does not exist")
			}
		}
	}
	return nil
}

type Context struct {
	uuid             string
	machine          string
	currentStateName string
	UsrData          interface{}
}

func NewContext(m Machine) Context {
	c := Context{machine: m.Name, currentStateName: m.States[0].Name}
	uuid, _ := newUUID()
	c.uuid = uuid
	return c
}

func (c Context) GetUUID() string {
	return c.uuid
}

func NewMachine(js string) Machine {
	var mac Machine

	if err := json.Unmarshal([]byte(js), &mac); err != nil {
		panic(err)
	}
	//	fmt.Printf("Mac %v\n", mac)
	mac.ActionFuncs = make(map[string]ActionFunc)
	mac.GuardFuncs = make(map[string]GuardFunc)

	return mac
}

func PrintEv(state string, ev Event, ctxt *Context) {
	fmt.Printf("state: %v uuid: %v event: %v payload: %v\n", state, ctxt.uuid, ev.GetName(), ev.GetPayload())
}
