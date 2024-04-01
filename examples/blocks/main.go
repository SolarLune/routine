package main

import (
	"fmt"
	"time"

	"github.com/solarlune/routine"
	"github.com/solarlune/routine/actions"
)

var drawProgress = false
var progress = 0

// This example shows how you can easily jump from one active block to another.
// A block represents a sequence of actions, tied to an ID that identifies that Block.
// Multiple Blocks can be active at any given time.
func defineRoutine(myRoutine *routine.Routine) {

	// You can easily create your own functional Actions by simply creating a function that returns an ActionFunc
	// with your own custom function defined within. You can then call these functions within a Block definition,
	// like you would any of the actionsuickActions functions.
	printText := func(text string) *actions.Function {

		f := func(block *routine.Block) routine.Flow {
			fmt.Println(block.ID.(string) + " : " + text)
			return routine.FlowNext
		}

		return actions.NewFunction(f)

	}

	first := myRoutine.Define("first",

		printText("In this example, we will switch from one block to another."),

		actions.NewWait(time.Second*2),

		printText("Let's fill up a progress bar, but we'll do this in the 'progress' block."),

		actions.NewWait(time.Second*3),

		printText("Let's switch now!"),

		actions.NewWait(time.Second*2),

		printText("-click-"),

		actions.NewWait(time.Second*1),

		actions.NewRunBlock("progress"),
	)

	first.Run() // Schedule the "first" Block to run first

	myRoutine.Define("progress",

		printText("OK. Now we're in the 'progress' block."),

		actions.NewWait(time.Second*2),

		printText("Filling up progress bar..."),

		actions.NewWait(time.Second*2),

		actions.NewFunction(func(block *routine.Block) routine.Flow {
			drawProgress = true
			progress += 5
			if progress >= 100 {
				fmt.Print("\n") // Skip a line
				drawProgress = false
				return routine.FlowNext
			}
			return routine.FlowIdle
		}),

		printText("Done!"),

		actions.NewFunction(func(block *routine.Block) routine.Flow {
			return routine.FlowNext
		}),

		actions.NewWait(time.Second*2),

		// Blocks will automatically stop at the end of their Action list.
		// If no blocks are running, the Routine will also, by default, stop.

	)

}

func main() {

	// Create a new routine.
	myRoutine := routine.New()

	// Define the routine.
	defineRoutine(myRoutine)

	// While it's running...
	for myRoutine.Running() {

		// Update the Routine.
		myRoutine.Update()

		// Draw the progress bar when it's time to do so
		if drawProgress {
			pro := "["

			for i := 0; i < 100; i += 5 {
				if i > progress {
					pro += "▫"
				} else {
					pro += "▪"
				}
			}
			pro += " ]"

			fmt.Print(pro + "\r")
		}

		// Sleep a bit so we don't spam the console when it's time to fill up the progress bar.
		time.Sleep(time.Millisecond * 100)

	}

}
