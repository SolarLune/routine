package main

import (
	"fmt"
	"time"

	"github.com/solarlune/routine"
	"github.com/solarlune/routine/actions"
)

// This example shows how jumping to a label works.
// All you have to do is create a new Label action, and then jump to it, either from within a function with block.JumpTo(), or from
// the block definition with actions.NewJumpTo().
func defineRoutine(myRoutine *routine.Routine) {

	myRoutine.DefineBlock("first",

		actions.NewFunction(func(block *routine.Block) routine.Flow {
			fmt.Println("Let's test jumping to a label.")
			// We can also jump within a function with block.JumpTo(). Note that this, of course, wouldn't end the function early
			// automatically; we would still have to return a routine.Flow.
			return routine.FlowNext
		}),

		actions.NewWait(time.Second*3),

		actions.NewJumpTo("after finish"),

		actions.NewFinishRoutine(), // If we didn't jump to the "after finish" label, then this would have ended the Block's execution.

		actions.NewLabel("after finish"),

		actions.NewFunction(func(block *routine.Block) routine.Flow {
			fmt.Println("This wouldn't have printed unless we jumped.")
			return routine.FlowNext
		}),

		actions.NewWait(time.Second*3),

		actions.NewFunction(func(block *routine.Block) routine.Flow {
			fmt.Println("OK, that's it.")
			return routine.FlowFinishRoutine
		}),
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
