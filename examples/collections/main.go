package main

import (
	"fmt"
	"time"

	"github.com/solarlune/routine"
	"github.com/solarlune/routine/actions"
)

// The point of this example is to show how you can create your own functions, and how Collections work.

func defineRoutine(myRoutine *routine.Routine) {

	// You can create your own Actions easily by making functions.
	print := func(text string) *actions.Function {
		return actions.NewFunction(
			func(b *routine.Block) routine.Flow {
				fmt.Println(text)
				return routine.FlowNext
			},
		)
	}

	// However, Blocks and ActionGates take a variable number of individual Actions -
	// because of this, you can't supply to them a pre-made slice of multiple Actions, like from
	// a function. To bypass this, you can use Collections. They are groups of Actions
	// that are substituted internally for the Actions they contain.
	slowType := func(text string) *actions.Collection {

		textIndex := 0

		return actions.NewCollection(

			actions.NewLabel("loop:"+text),

			actions.NewFunction(func(block *routine.Block) routine.Flow {
				fmt.Print(text[:textIndex] + "\r")
				textIndex++
				if textIndex >= len(text)+1 {
					fmt.Print("\n")
					block.JumpTo("finish")
				}
				return routine.FlowNext
			}),

			actions.NewWait(time.Second/10),

			actions.NewJumpTo("loop:"+text),

			actions.NewLabel("finish"),
		)
	}

	myRoutine.DefineBlock("first",

		print("You can easily make your own Actions by using functions."),

		actions.NewWait(time.Second*2),

		slowType("For example, here's a slow typing Action."),

		actions.NewWait(time.Second),

		print("Done!"),

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
