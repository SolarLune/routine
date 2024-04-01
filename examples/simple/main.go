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

	// Actions, to their name, perform an action and alter a Routine's flow.

	// A Routine runs through its Blocks.
	// You customize Actions that live within Blocks, and can change execution for Blocks freely - that's how
	// you create a Routine that does what you want.

	// Routine.DefineBlock() defines a block of Actions to execute in sequence.

	// When an Action has completed its behavior, the Block will move on to the next one, until it's at the end,
	// at which point, the Block will loop.

	// Below, we define a Block with the ID "first". The ID is a string here, but can be any comparable object.
	myRoutine.Define("first",

		// actions.NewFunction() creates a Function Action that executes a customizeable function.
		// This function must take the current block and return a routine.Flow value.
		// A routine.Flow signals to the running Block in the Routine what to do after the Action ends.

		// Depending on the RoutineFlow received, the Block can:

		// - Stay on the current Action, repeating the Action the next time Routine.Update() is called (routine.FlowIdle)
		// - Move to the next Action in the Block without waiting (routine.FlowNext), or
		// - End the Block entirely (routine.FlowFinish).

		// If you want to stop all Blocks, you can use Routine.DeactivateAllBlocks().

		actions.NewFunction(func(block *routine.Block) routine.Flow {
			fmt.Println("Here's a simple block that prints some text, and waits three seconds.")
			return routine.FlowNext
		}),

		actions.NewWait(time.Second*3),

		actions.NewFunction(func(block *routine.Block) routine.Flow {

			fmt.Println("Done!")
			// We can return routine.FlowFinish here, use actions.NewFinish() to create a finishing Action to end the
			// Routine, or simply do nothing if it's at the end of the block. Any of these options work.
			return routine.FlowFinish

		}),
	)

	// We activate the Block when we're done, as by default, Blocks aren't active.
	// This can be done by calling Routine.Run() with the Block's ID, or calling Block.Run().
	myRoutine.Run("first")

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
