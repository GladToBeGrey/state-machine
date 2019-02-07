// event
package sm

// import (
// 	"fmt"
// )

type Event interface {
	GetName() string
	GetCorrelId() string
	GetPayload() interface{}
}

type event struct {
	name     string
	correlId string
	payload  interface{}
}

func (e event) GetName() string {
	return e.name
}

func (e event) GetCorrelId() string {
	return e.correlId
}

func (e event) GetPayload() interface{} {
	return e.payload
}

func NewEvent(aName string, aCorrelId string, aPayload interface{}) Event {
	return event{name: aName, correlId: aCorrelId, payload: aPayload}
}
