package main

import (
	"fmt"
	"time"

	"github.com/solarlune/routine"
	"github.com/solarlune/routine/actions"
)

// The goal of this example is to show an element of parallel logic, in a sense.
// Multiple blocks can be active at any given time, which allows you to combine
// blocks together for various effects.

// In this example, the various blocks all execute together, in sequence, on different
// timings.

func defineRoutine(myRoutine *routine.Routine) {

	next := 0

	// The first block executes every 2 seconds.
	myRoutine.DefineBlock(0,

		actions.NewFunc(func(block *routine.Block) routine.Flow {
			fmt.Println("First block is just the beginning...")

			// Activate the next block; this happens every 2 seconds.
			myRoutine.ActivateBlock(next)
			next++
			return routine.FlowNext
		}),

		actions.NewWait(time.Second*2),
	)

	// The second block executes every half second.
	myRoutine.DefineBlock(1,

		actions.NewFunc(func(block *routine.Block) routine.Flow {
			fmt.Println("second block is alive and well...")
			return routine.FlowNext
		}),

		actions.NewWait(time.Second/2),
	)

	// The third block executes 10 times a second.
	myRoutine.DefineBlock(2,

		actions.NewFunc(func(block *routine.Block) routine.Flow {
			fmt.Println("third block is going crazy...")
			return routine.FlowNext
		}),

		actions.NewWait(time.Second/10),
	)

	// The fourth block executes 20 times a second.
	myRoutine.DefineBlock(3,

		actions.NewFunc(func(block *routine.Block) routine.Flow {
			fmt.Println("fourth block is kinda insane...!!!")
			return routine.FlowNext
		}),

		actions.NewWait(time.Second/20),
	)

	// The last block ends it.
	myRoutine.DefineBlock(4,

		actions.NewFunc(func(block *routine.Block) routine.Flow {
			// Set only this block to be active
			myRoutine.SwitchBlock(4)
			fmt.Println("OK, I'm done. All tuckered out.")
			return routine.FlowNext
		}),

		actions.NewWait(time.Second),

		actions.NewFinish(),
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
