package main

import (
	"fmt"
	"time"

	"github.com/solarlune/routine"
	"github.com/solarlune/routine/actions"
)

// In this example, we see how we can make blocks loop.
func defineRoutine(myRoutine *routine.Routine) {

	loopCount := 4

	myRoutine.DefineBlock("loop",

		actions.NewFunction(func(block *routine.Block) routine.Flow {

			if loopCount == 0 {
				fmt.Println("Welp, that's it. Routine over~")
				return routine.FlowFinishRoutine
			}

			fmt.Printf("This block will loop %d more times.\n", loopCount-1)

			return routine.FlowNext

		}),

		actions.NewWait(time.Second*2),

		actions.NewFunction(func(block *routine.Block) routine.Flow {

			// You can reference global or even local variables outside of a block definition
			// to act as a kind of temporary memory.
			loopCount--

			return routine.FlowNext

		}),

		actions.NewLoop(),
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
