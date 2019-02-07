// account
package main

import (
	"fmt"

	"pacs8.com/sm"
)

type stmtLine struct {
	narrative string
	amount    int
	drCr      string
}

type Account struct {
	number   string
	currency string
	balance  int
	stmt     []stmtLine
}

type openPayload struct {
	accountNumber string
	currency      string
}

func NewEventOpen(uuid string, accountNumber string, currency string) sm.Event {
	return sm.NewEvent("evOpen", uuid, openPayload{accountNumber, currency})
}

type withdrawPayload struct {
	amount    int
	currency  string
	narrative string
}

func NewEventWithdraw(uuid string, amount int, currency string, narrative string) sm.Event {
	return sm.NewEvent("evWithdraw", uuid, withdrawPayload{amount, currency, narrative})
}

type depositPayload struct {
	amount    int
	currency  string
	narrative string
}

func NewEventDeposit(uuid string, amount int, currency string, narrative string) sm.Event {
	return sm.NewEvent("evDeposit", uuid, depositPayload{amount, currency, narrative})
}

type closePayload struct {
	reason string
}

func NewEventClose(uuid string, reason string) sm.Event {
	return sm.NewEvent("evClose", uuid, closePayload{reason})
}

type financial interface {
	getCurrency() string
	getAmount() int
}

func (wdr withdrawPayload) getCurrency() string {
	return wdr.currency
}

func (wdr withdrawPayload) getAmount() int {
	return wdr.amount
}

func (depos depositPayload) getCurrency() string {
	return depos.currency
}

func (depos depositPayload) getAmount() int {
	return depos.amount
}

func doOpen(state string, ev sm.Event, ctxt *sm.Context) {
	fmt.Printf("Into doOpen\n")
	open := ev.GetPayload().(openPayload)
	ctxt.UsrData = Account{number: open.accountNumber, currency: open.currency, balance: 0, stmt: make([]stmtLine, 0)}
	fmt.Printf("Account %v opened\n", open.accountNumber)
}

func doDeposit(state string, ev sm.Event, ctxt *sm.Context) {
	fmt.Printf("Into doDeposit\n")
	depos := ev.GetPayload().(depositPayload)
	acc := ctxt.UsrData.(Account)
	acc.balance += depos.amount
	acc.stmt = append(acc.stmt, stmtLine{depos.narrative, depos.amount, "CR"})
	ctxt.UsrData = acc
	fmt.Printf("Deposited %v to Account %v\n", depos.amount, acc.number)
}

func doDepositWrongCurr(state string, ev sm.Event, ctxt *sm.Context) {
	fmt.Printf("Into doDepositWrongCurr\n")
	depos := ev.GetPayload().(depositPayload)
	acc := ctxt.UsrData.(Account)
	fmt.Printf("Cannot deposit %v into  %v account\n", depos.currency, acc.currency)
}

func doWithdraw(state string, ev sm.Event, ctxt *sm.Context) {
	fmt.Printf("Into doWithdraw\n")
	wdr := ev.GetPayload().(withdrawPayload)
	acc := ctxt.UsrData.(Account)
	acc.balance -= wdr.amount
	acc.stmt = append(acc.stmt, stmtLine{wdr.narrative, wdr.amount, "DR"})
	ctxt.UsrData = acc
	fmt.Printf("Withdrew %v from Account %v\n", wdr.amount, acc.number)
}

func doWithdrawWrongCurr(state string, ev sm.Event, ctxt *sm.Context) {
	fmt.Printf("Into doWithdrawWrongCurr\n")
	wdr := ev.GetPayload().(withdrawPayload)
	acc := ctxt.UsrData.(Account)
	fmt.Printf("Cannotwithdraw %v from  %v account\n", wdr.currency, acc.currency)
}

func isSameCurrency(state string, ev sm.Event, ctxt *sm.Context) bool {
	fmt.Printf("Into isSameCurrency\n")
	fin := ev.GetPayload().(financial)
	acc := ctxt.UsrData.(Account)
	return fin.getCurrency() == acc.currency
}
