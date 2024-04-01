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
	printText := func(text string) *actions.Function {
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
	typeText := func(text string) *actions.Collection {

		textIndex := 0

		// When this is added to a Block definition, the block will take the actions from the Collection.
		return actions.NewCollection(

			// We use "loop:"+text and "finish:"+text here because multiple slowType functions will define multliple labels
			// for jumping, so we use the text as essentially an identifier of which parts to jump to in the Block.

			actions.NewLabel("loop:"+text),

			actions.NewFunction(func(block *routine.Block) routine.Flow {
				if textIndex >= len(text)+1 {
					fmt.Print("\n")
					block.JumpTo("finish:" + text)
					return routine.FlowNext
				} else {
					fmt.Print(text[:textIndex] + "\r")
				}
				textIndex++
				return routine.FlowNext
			}),

			actions.NewWait(time.Second/20),

			actions.NewJumpTo("loop:"+text),

			actions.NewLabel("finish:"+text),
		)
	}

	myRoutine.Define("first",

		printText("You can easily make your own Actions by using Function Actions or Collection Actions for groups of Actions."),

		actions.NewWait(time.Second*2),

		typeText("For example, here's a typing Action."),

		typeText("This is a function that returns a Collection of multiple Actions that, when combined, type a message out to the console."),

		actions.NewWait(time.Second),

		printText("Done!"),

		actions.NewFinish(),
	)

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
