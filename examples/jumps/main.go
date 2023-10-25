package main

import (
	"fmt"
	"time"

	"github.com/solarlune/routine"
	"github.com/solarlune/routine/actions"
)

func defineRoutine(myRoutine *routine.Routine) {

	myRoutine.DefineBlock("first",

		actions.NewFunc(func(block *routine.Block) routine.Flow {
			fmt.Println("Let's test jumping to a label.")
			// We can also jump within a function with actions.NewCurrentBlock.JumpTo()
			return routine.FlowNext
		}),

		actions.NewWait(time.Second*3),

		actions.NewJumpTo("after finish"),

		actions.NewFinish(), // If we didn't jump to the "after finish" label, then this would have ended the Block's execution.

		actions.NewLabel("after finish"),

		actions.NewFunc(func(block *routine.Block) routine.Flow {
			fmt.Println("This wouldn't have printed unless we jumped.")
			return routine.FlowNext
		}),

		actions.NewWait(time.Second*3),

		actions.NewFunc(func(block *routine.Block) routine.Flow {
			fmt.Println("OK, that's it.")
			return routine.FlowFinish
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
