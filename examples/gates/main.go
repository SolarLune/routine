package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/solarlune/routine"
	"github.com/solarlune/routine/actions"
)

// This example shows how ActionGates work.
// They're used to create conditional branches that evaluate once.
func defineRoutine(myRoutine *routine.Routine) {

	print := func(text string) *actions.Collection {

		return actions.NewCollection(

			actions.NewFunc(func(block *routine.Block) routine.Flow {
				fmt.Println(text)
				return routine.FlowNext
			}),

			actions.NewWait(time.Second*2),
		)

	}

	var choice int

	myRoutine.DefineBlock("first",

		print("OK, so let's try an ActionGate out."),
		print("Let's see which option we get..."),

		actions.NewFunc(func(block *routine.Block) routine.Flow {
			choice = rand.Intn(3)
			return routine.FlowNext
		}),

		// The way a Gate works, each GateOption is checked to see if its function returns true; if so,
		// then that Option is made active, and the Actions within are executed in seactionsuence. The other
		// Options are no longer checked, until the Block comes back around to executing the Gate.
		actions.NewGate(

			actions.NewGateOption(
				func() bool { return choice == 0 },
				print("Option #1 was chosen."),
			),

			actions.NewGateOption(
				func() bool { return choice == 1 },
				print("The second choice, option #2 was selected."),
			),

			actions.NewGateOption(
				func() bool { return choice == 2 },
				print("The third choice was chosen."),
				print("This one is a loser - game over!"),
				actions.NewFinish(),
			),
		),

		print("Nice!"),
	)

}

func main() {

	// Create a new routine.
	myRoutine := routine.New()

	// Define the routine.
	defineRoutine(myRoutine)

	// Run the routine.
	myRoutine.Run()

	// While it's running...
	for myRoutine.Running() {

		// Update it.
		myRoutine.Update()

	}

}
