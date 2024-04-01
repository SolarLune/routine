package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/solarlune/routine"
	"github.com/solarlune/routine/actions"
)

// This example shows how ActionGates work.

// They're used to create conditional branches that evaluate once, then execute a set of actions in sequence.
// While you could use a Function to perform conditional branches, since Gates evaluate just once, they stay active
// regardless of the value of the variable that specified the execution route.
func defineRoutine(myRoutine *routine.Routine) {

	print := func(text string) *actions.Collection {

		return actions.NewCollection(

			actions.NewFunction(func(block *routine.Block) routine.Flow {
				fmt.Println(text)
				return routine.FlowNext
			}),

			actions.NewWait(time.Second*1),
		)

	}

	// A variable that lies outside of a Block definition can act as memory, as all Actions in a Definition
	// can access it. Of course, be aware of the difference between variables that exist within a function
	// definition and those that exist outside of one.
	var choice int

	myRoutine.Define("first",

		print("OK, so let's try an ActionGate."),

		actions.NewLabel("gate start"),

		print("Let's see which option we randomly get..."),

		actions.NewFunction(func(block *routine.Block) routine.Flow {
			choice = rand.Intn(3)
			return routine.FlowNext
		}),

		// Each GateOption is checked to see if its function returns true; if so,
		// then that Option is made active, and the Actions within are executed in sequence. The other
		// Options are no longer checked, until the Block comes back around to executing the Gate.
		actions.NewGate(

			actions.NewGateOption(
				func() bool { return choice == 0 },
				print("1: Option #1 was chosen."),
			),

			actions.NewGateOption(
				func() bool { return choice == 1 },
				print("2: The second choice, option #2 was selected."),
			),

			actions.NewGateOption(
				nil, // You can return a nil for the function in the last GateOption defined to serve as an "else" case; it is always evaluated as true.
				print("3: The third choice was chosen."),
				print("This one is a loser - game over!"),
				actions.NewFinish(),
			),
		),

		print("Nice! Let's try again."),

		actions.NewJumpTo("gate start"),
	).Run() // A quick way to activate a certain Block.

}

func main() {

	// Create a new routine.
	myRoutine := routine.New()

	// Define the routine.
	defineRoutine(myRoutine)

	// While it's running...
	for myRoutine.Running() {

		// Update it.
		myRoutine.Update()

	}

}
