// stateMachine project main.go
package main

import (
	"fmt"
	"os"

	"pacs8.com/sm"
)

// TODO handle wildcard transitions?
// multiple contextx with event correlation

func main() {

	js := `{
		"name": "Account",
		"states" : [
			{
				"name": "New",
				"transitions": [
					{"onEvent": "evOpen",
					 "toState": "Open",
					 "action": "doOpen"}
				]
			},
			{
				"name": "Open",
				"transitions": [
					{"onEvent": "evDeposit",
					 "guard": "isSameCurrency",
					 "toState": "Open",
					 "action": "doDeposit"},
					{"onEvent": "evDeposit",
					 "toState": "Open",
					 "action": "doDepositWrongCurr"},
					{"onEvent": "evWithdraw",
					 "guard": "isSameCurrency",
					 "toState": "Open",
					 "action": "doWithdraw"},
					{"onEvent": "evWithdraw",
					 "toState": "Open",
					 "action": "doWithdrawWrongCurr"},
					{"onEvent": "evClose",
					 "toState": "Closed",
					 "action": "PrintEv"}
				]
			},
			{
				"name": "Closed"
			}
		]
	}`

	myMac := sm.NewMachine(js).
		AddActionFuncs(sm.PrintEv, doOpen, doDeposit, doDepositWrongCurr, doWithdraw, doWithdrawWrongCurr).
		AddGuardFuncs(isSameCurrency)

	err := myMac.Validate()
	if err != nil {
		fmt.Println("Machine validation failed: " + err.Error())
		os.Exit(1)
	}

	ctxt := sm.NewContext(myMac)
	uuid := ctxt.GetUUID()

	events := []sm.Event{
		NewEventOpen(uuid, "12345678", "GBP"),
		NewEventDeposit(uuid, 1000, "GBP", "Initial Deposit"),
		NewEventDeposit(uuid, 200, "USD", "FX Deposit"),
		NewEventWithdraw(uuid, 100, "GBP", "ATM Withdrawal"),
		NewEventWithdraw(uuid, 10000, "JPY", "ATM Withdrawal"),
		NewEventClose(uuid, "CloseReason")}

	for _, ev := range events {
		if myMac.HandleEvent(ev, &ctxt) {
			fmt.Printf("After event %v\n", ctxt)
		}
	}
}
