# routine 🐱 ➡️ 😺💬

[![Go Reference](https://pkg.go.dev/badge/github.com/SolarLune/routine.svg)](https://pkg.go.dev/github.com/SolarLune/routine)

routine is a package for handling and executing routines, written in Go. The primary reason I made this was for creating cutscenes or running routines for gamedev with Go.

The general idea is that you create a Routine, and then define a Block. Blocks are "modes of execution". They contain groups of Actions, which are objects that are "units of execution" - they perform actions or manipulate the execution of a Block. Multiple Blocks can be active at any given time.

By utilizing Actions and Blocks, you can make up complex behaviors, or sequences of events.

## How do I get it?

`go get github.com/solarlune/routine`

## Example

```go

import "github.com/solarlune/routine"
import "github.com/solarlune/routine/actions"

func main() {

    // Create a new Routine.
    routine := routine.New()

    // And now we begin to define our Blocks, which consist of Action objects that execute in sequence.
    // A Block has an ID (in this case, a string set to "first block"), but the ID can be anything.
    routine.DefineBlock("first block", 
    
        // actions.NewFunction() returns a ActionFunc, which is a function that will run the action provided.
        // Depending on the Flow object returned, the Block will either idle on this Action, or move
        // on to another one.
        actions.NewFunction(func() routine.Flow { 
            fmt.Println("Hi!")
            return routine.FlowNext // FlowNext means to move to the next Action.
        }),

        actions.NewWait(time.Second * 2),

        actions.NewFunction(func() routine.Flow { 
            fmt.Println("OK, that's it. Goodbye!")
            return routine.FlowFinish 
        })
    
    )

    // And now we begin running the Routine. By default, the first defined Block
    // is activated when we run the Routine.
    routine.Run()
    
    for routine.Running() {

        // While the routine runs, we call Routine.Update(). This allows
        // the routine to execute, but also gives Action back to the main
        // thread when it's cycling (done until the next frame / Update() call) 
        // so we can do other stuff, like take input or update a game's screen.
        routine.Update()

    }

    // After Running() is over, we're done with the Routine.

}

```

Check out the `examples` directory for more in-depth examples.

## Anything else?

Not really, that's it. Peace~