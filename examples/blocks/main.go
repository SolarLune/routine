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

func defineRoutine(myRoutine *routine.Routine) {

	// You can easily create your own functional Actions by simply creating a function that returns an ActionFunc
	// with your own custom function defined within. You can then call these functions within a Block definition,
	// like you would any of the actionsuickActions functions.
	customPrint := func(text string) *actions.Function {

		f := func(block *routine.Block) routine.Flow {
			fmt.Println(block.ID.(string) + " : " + text)
			return routine.FlowNext
		}

		return actions.NewFunction(f)

	}

	myRoutine.DefineBlock("first",

		customPrint("In this example, we will switch from one block to another."),

		actions.NewWait(time.Second*2),

		customPrint("Let's fill up a progress bar, but we'll do this in the 'progress' block."),

		actions.NewWait(time.Second*3),

		customPrint("Let's switch now!"),

		actions.NewWait(time.Second*2),

		customPrint("-click-"),

		actions.NewWait(time.Second*1),

		// We can switch blocks using actions.NewSwitchBlock() or Routine.SwitchBlock(). Any blocks with the given names will be activated.
		actions.NewSwitchBlock("progress"),
	)

	myRoutine.DefineBlock("progress",

		customPrint("OK. Now we're in the 'progress' block."),

		actions.NewWait(time.Second*2),

		customPrint("Filling up progress bar..."),

		actions.NewWait(time.Second*2),

		actions.NewFunction(func(block *routine.Block) routine.Flow {
			drawProgress = true
			progress += 5
			if progress >= 100 {
				fmt.Println("") // Skip a line
				drawProgress = false
				return routine.FlowNext
			}
			return routine.FlowIdle
		}),

		customPrint("Done!"),

		actions.NewFunction(func(block *routine.Block) routine.Flow {
			return routine.FlowNext
		}),

		actions.NewWait(time.Second*2),

		// If we don't finish explicitly, either with a custom ActionFunc that returns RoutineFlowFinish or
		// Flow.Finish() (which is just a actionsuickie function that creates a ActionFunc returning RoutineFlowFinish),
		// the Routine would loop infinitely.
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
