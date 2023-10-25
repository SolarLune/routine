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

	// Actions, to their name, Action a Routine's flow.

	// A Routine runs through its Blocks.
	// You customize Actions that live within Blocks, and can change execution for Blocks freely - that's how
	// you create a Routine that does what you want.

	// Routine.DefineBlock defines a block of Actions to execute in seactionsuence.
	// Whatever block is the first to be defined will be the default block to run when a Routine is run.
	// When a Action has completed its behavior, the Block will move on to the next one, until it's at the end.
	// At that point, the Block will loop.

	// Below, we define a Block with the ID "first". The ID is a string here, but can be any comparable object.
	myRoutine.DefineBlock("first",

		// actions.NewFunc() creates a Func, which is a Action that executes a customizeable function.
		// This function must take the current block and return a RoutineFlow.
		// A RoutineFlow signals to the running Block in the Routine what to do after the Action ends.

		// The Block can:

		// - Stay on the current Action, repeating the Action the next time Routine.Update() is called (routine.FlowIdle)
		// - Move to the next Action (routine.FlowNext), or
		// - End the Routine entirely (routine.FlowFinish).

		// If you want to simply stop a Block after it's done running, you can do so using Routine.DeactivateBlocks().
		actions.NewFunc(func(block *routine.Block) routine.Flow {
			fmt.Println("Here's a simple block that prints some text, and waits three seconds.")
			return routine.FlowNext
		}),

		actions.NewWait(time.Second*3),

		actions.NewFunc(func(block *routine.Block) routine.Flow {

			fmt.Println("Done!")

			// We can return RoutineFlowFinish here, or use Flow.Finish() as a following Action to end the
			// Routine; whichever works.
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
