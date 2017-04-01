package dbs

import (
	"github.com/empirefox/bongine/front"
	"github.com/empirefox/bongine/fsm"
)

var orderRules = fsm.CreateRuleset(
	newTransition(front.TOrderStateNopay, front.TOrderStatePrepaid),
	newTransition(front.TOrderStateNopay, front.TOrderStatePaid),
	newTransition(front.TOrderStateNopay, front.TOrderStateCanceled),

	newTransition(front.TOrderStatePrepaid, front.TOrderStatePrepaid),
	newTransition(front.TOrderStatePrepaid, front.TOrderStatePaid),
	newTransition(front.TOrderStatePrepaid, front.TOrderStateCanceled),

	newTransition(front.TOrderStatePaid, front.TOrderStateEnsuring), // auto, if needed
	newTransition(front.TOrderStatePaid, front.TOrderStateCompleted),

	newTransition(front.TOrderStateEnsuring, front.TOrderStateRefound), // admin
	newTransition(front.TOrderStateEnsuring, front.TOrderStateEnsured), // admin

	newTransition(front.TOrderStateEnsured, front.TOrderStateRefound), // admin
	newTransition(front.TOrderStateEnsured, front.TOrderStateCompleted),

	newTransition(front.TOrderStateCompleted, front.TOrderStateEvaled), // standalone

	newTransition(front.TOrderStateEvaled, front.TOrderStateEvaled),  // standalone
	newTransition(front.TOrderStateEvaled, front.TOrderStateHistory), // system
)

func newTransition(from, to front.OrderState) fsm.T {
	return fsm.T{fsm.State(from), fsm.State(to)}
}

func PermitOrderState(o *front.Order, s front.OrderState) error {
	return orderRules.Permitted(o, fsm.State(s))
}
