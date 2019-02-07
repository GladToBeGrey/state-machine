# state-machine

This repo contains two items:
pacs8.com/sm: state machine package
stateMachine: an example of usage

A state machine is a simple yet powerful way to model processes. It consists of three basic concepts:

State: A step in the process. When the process reaches a state, it will halt there until an event is received.
Event: A message received by the process, normally representing something that happens in the outside world
Transition: The processing path from one state to another. A transition is associated with a specific event.

States can have behaviours: onEntry is executed as the process enters the state, and onExit is executed as the process exits the state (irrespective of which event triggered the change of state)

Transitions can also have behaviours: onEvent is executed as the process follows the transition.

