package main

import (
	"fmt"
	"time"

	"github.com/solarlune/routine"
	"github.com/solarlune/routine/actions"
)

// This example just shows how Routines work, in their most simple format.

// Here we'll define our Routine. It doesn't have to be done in a function, of course;
// this is just to make it easier-to-understand by segmenting it out from the main function.
func defineRoutine(myRoutine *routine.Routine) {

	block := myRoutine.Define("first",

		actions.NewFunction(func(block *routine.Block) routine.Flow {
			fmt.Println("This will print some text, all at once, even though it's spread across multiple actions.")
			return routine.FlowNext
		}),

		actions.NewFunction(func(block *routine.Block) routine.Flow {
			fmt.Println("This is because Blocks will only yield to the main outer loop if you return FlowIdle constant or if the Block finishes.")
			return routine.FlowNext
		}),

		actions.NewFunction(func(block *routine.Block) routine.Flow {
			fmt.Println("This allows you to easily compose events out of multiple individual actions without having to worry about time unless you explicitly idle or wait across multiple update calls.")
			return routine.FlowNext
		}),

		actions.NewWait(time.Second*3),

		actions.NewFunction(func(block *routine.Block) routine.Flow {
			fmt.Println("That's it for this one.")
			return routine.FlowFinish
		}),
	)

	block.Run() // We activate the Block when we're done, as by default, Blocks aren't active.

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
