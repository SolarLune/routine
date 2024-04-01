package actions

import (
	"math/rand"
	"time"

	"github.com/solarlune/routine"
)

// Wait is an action that waits a customizeable amount of time before continuing.
type Wait struct {
	Duration   time.Duration
	targetTime time.Time
}

// NewWait creates a new Wait Action.
func NewWait(duration time.Duration) *Wait {
	wait := &Wait{
		Duration: duration,
	}
	return wait
}

func (w *Wait) Init(block *routine.Block) {
	w.targetTime = time.Now().Add(w.Duration)
}

func (w *Wait) Poll(block *routine.Block) routine.Flow {
	if time.Now().After(w.targetTime) {
		return routine.FlowNext
	}
	return routine.FlowIdle
}

// NewWaitTicks creates a new action that waits a certain amount of time before proceeding.
func NewWaitTicks(tickCount int) *Function {
	return NewFunction(func(block *routine.Block) routine.Flow {

		if block.CurrentFrame() >= tickCount {
			return routine.FlowNext
		}

		return routine.FlowIdle

	})
}

// NewWaitTicks creates a new action that waits a random amount of time, ranging between minTime and maxTime, before proceeding.
func NewWaitTicksRandom(minTime, maxTime int) *Function {

	tickCount := 0

	return NewFunction(func(block *routine.Block) routine.Flow {

		if block.CurrentFrame() == 0 {
			tickCount = minTime + int((float64(maxTime-minTime) * rand.Float64()))
		}

		if block.CurrentFrame() >= tickCount {
			return routine.FlowNext
		}

		return routine.FlowIdle

	})
}

// Function is a Action that runs a customizeable function.
type Function struct {
	InitFunc func(block *routine.Block)              // The function to run when the ActionFunc object is initialized (before polling)
	PollFunc func(block *routine.Block) routine.Flow // The function to run when polled
}

// NewFunction creates and returns a Function action object with the polling function set to the
// provided function. The routine.Flow returned from the customizeable function influences
// the Routine does after running the function.
func NewFunction(function func(block *routine.Block) routine.Flow) *Function {
	return &Function{
		PollFunc: function,
	}
}

func (f *Function) Init(block *routine.Block) {
	if f.InitFunc != nil {
		f.InitFunc(block)
	}
}

func (f *Function) Poll(block *routine.Block) routine.Flow { return f.PollFunc(block) }

// TimingPair represents an action to take after a specific duration of time
// has passed.
type TimingPair struct {
	Duration   time.Duration
	Function   func()
	targetTime time.Time
}

// Timing is a timing Action, which executes a provided function when
// some amount of time has elapsed.
type Timing struct {
	pairs []TimingPair
	index int
}

// NewTiming creates a new ActionTiming object. A ActionTiming object works with
// TimingPairs, which indicate a function to execute after a specific duration
// of time has passed.
func NewTiming(timingPairs []TimingPair) *Timing {
	return &Timing{
		pairs: timingPairs,
	}
}

func (t *Timing) Init() {
	t.index = 0
}

func (t *Timing) Poll(block *routine.Block) routine.Flow {

	pair := &t.pairs[t.index]

	if pair.targetTime.IsZero() {
		pair.targetTime = time.Now().Add(pair.Duration)
	}

	if time.Now().After(pair.targetTime) {
		pair.Function()

		t.index++
		if t.index >= len(t.pairs) {
			t.index = 0
			return routine.FlowNext
		}

	}

	return routine.FlowIdle
}

// GateOption represents a choice in a ActionGate Action.
type GateOption struct {
	CheckFunc func() bool
	Active    bool
	actions   []routine.Action
	Index     int
}

// NewGateOption creates a new GateOption object, which represents a choice in an ActionGate. The checkFunc
// argument signifies a function that determines if the entry is made active. If it is made active,
// then the Actions within (and only those Actions) are executed in sequence until the end.
// If checkFunc is nil, then that is equivalent to a checkFunc definition of func() bool { return true }.
// A checkFunc of nil can be used as an "else" statement when used after all other entries in a Gate.
// If no actions are passed for the option, then if the check function returns true (or is nil), then the
// Block will move on to the next action after the Gate.
func NewGateOption(checkFunc func() bool, Actions ...routine.Action) *GateOption {

	newActions := []routine.Action{}

	for _, c := range Actions {
		if collection, ok := c.(*Collection); ok {
			newActions = append(newActions, collection.actions...)
		} else {
			newActions = append(newActions, c)
		}
	}

	return &GateOption{
		CheckFunc: checkFunc,
		actions:   newActions,
	}
}

func (g *GateOption) Init(block *routine.Block) {
	g.actions[0].Init(block)
	g.Index = 0
}

func (g *GateOption) Poll(block *routine.Block) routine.Flow {

	if len(g.actions) == 0 {
		return routine.FlowNext
	}

	result := g.actions[g.Index].Poll(block)

	done := false

	if result == routine.FlowNext {
		g.Index++
		if g.Index < len(g.actions) {
			g.actions[g.Index].Init(block)
		} else {
			g.actions[0].Init(block)
			g.Index = 0
			done = true
		}
	}

	if result == routine.FlowFinish {
		return routine.FlowFinish
	} else if done {
		return routine.FlowNext
	}

	return routine.FlowIdle

}

// Gate represents a gate, which allows for executing logic statements to determine
// an execution path (one of the passed GateOptions). Once the logic statement is executed,
// the gate is set until it is reset by revisiting the Action.
type Gate struct {
	Options     []*GateOption
	ActiveEntry *GateOption
	onIdle      func()
	onChoose    func()
}

// NewGate creates a Gate action, which allows you to effectively choose one "route" or "choice"
// option among many.
// Once one GateOption has been made active, it will stay active until the Gate runs
// through all actions the GateOption might have.
func NewGate(entries ...*GateOption) *Gate {
	return &Gate{
		Options: entries,
	}
}

// AddOption adds an option to the Gate action.
func (c *Gate) AddOption(option *GateOption) *Gate {
	c.Options = append(c.Options, option)
	return c
}

func (c *Gate) Init(block *routine.Block) {
	for _, entry := range c.Options {
		if len(entry.actions) > 0 {
			entry.actions[0].Init(block)
		}
	}
	c.ActiveEntry = nil
}

func (c *Gate) Poll(block *routine.Block) routine.Flow {

	if c.ActiveEntry != nil {
		return c.ActiveEntry.Poll(block)
	} else {
		if c.onIdle != nil {
			c.onIdle()
		}
		for _, entry := range c.Options {
			if entry.CheckFunc == nil || entry.CheckFunc() {
				c.ActiveEntry = entry
				if c.onChoose != nil {
					c.onChoose()
				}
				break
			}
		}
	}

	return routine.FlowIdle

}

// SetOnIdle sets the idling function for the ActionGate - when this is set, this function will run
// as long as a gate option isn't chosen.
func (c *Gate) SetOnIdle(onIdle func()) *Gate {
	c.onIdle = onIdle
	return c
}

// SetIdlingFunction sets the "on choose" function for the ActionGate - when this is set, this function will run
// when a gate option is chosen.
func (c *Gate) SetOnChoose(onChoose func()) *Gate {
	c.onChoose = onChoose
	return c
}

// Collection is not actually an Action to be strictly used; it's a container to pass to a Block or ActionGate.
// When either receives it in the process of construction, it will skip adding the Collection itself and instead
// add its contents. This is primarily so that you can, for example, make a function that returns multiple Actions
// in sequence.
type Collection struct {
	actions []routine.Action
}

// Collection creates a ActionCollection, which is a collection of Actions (naturally).
// A Collection by itself does nothing. Instead, the Actions that it is created with are
// supplied in sequence to other Actions that take individual Actions.
func NewCollection(actions ...routine.Action) *Collection {
	collection := &Collection{}

	newActions := []routine.Action{}
	for _, c := range actions {
		if collection, ok := c.(routine.ActionCollectionable); ok {
			newActions = append(newActions, collection.Actions()...)
		} else {
			newActions = append(newActions, c)
		}
	}
	collection.actions = newActions

	return collection
}

// AddAction allows you to add an Action to the Collection after creation.
func (q *Collection) AddAction(action routine.Action) {
	q.actions = append(q.actions, action)
}

func (q *Collection) Init(block *routine.Block) {}

func (q *Collection) Poll(block *routine.Block) routine.Flow { return routine.FlowNext }

func (q *Collection) Actions() []routine.Action { return q.actions }

// Label doesn't do anything specifically, but rather simply makes it possible
// for Blocks to jump to specific locations with Block.JumpTo(). This is internally
// the same as calling Block.SetIndex(), but with the index of the Label action.
type Label struct {
	Label any
}

// NewLabel creates a ActionLabel with the specified ID at the given location in the
// Block, enabling jumping to this point.
func NewLabel(id any) *Label {
	return &Label{
		Label: id,
	}
}

func (l *Label) Init(block *routine.Block) {}

func (l *Label) Poll(block *routine.Block) routine.Flow { return routine.FlowNext }

func (l *Label) ID() any { return l.Label }

// NewJumpTo creates a Function action that jumps the Block to the ActionLabel that has
// the specified label ID.
// If no Action with the label given is found, then the action will do nothing.
func NewJumpTo(label any) *Function {
	return NewFunction(
		func(block *routine.Block) routine.Flow {
			block.JumpTo(label)
			return routine.FlowNext
		},
	)
}

// NewSwitchBlock creates a Function action that switches the routine to only activate blocks with
// the specified IDs.
// If no block IDs are specified, all blocks are restarted.
func NewSwitchBlock(blockIDs ...any) *Function {
	return NewFunction(
		func(block *routine.Block) routine.Flow {
			r := block.Routine()
			r.Stop(blockIDs...)
			r.Run(blockIDs...)
			return routine.FlowNext
		},
	)
}

// NewRunBlock creates a Function action that activates the specified blocks in the
// currently running Routine. Any other blocks are unaffected.
// If no block IDs are specified, all blocks are run.
func NewRunBlock(blockIDs ...any) *Function {
	return NewFunction(
		func(block *routine.Block) routine.Flow {
			block.Routine().Run(blockIDs...)
			return routine.FlowNext
		},
	)
}

// NewPauseBlock creates a Function action that deactivates the specified blocks
// in the currently running Routine. Any other blocks are unaffected.
// If no block IDs are specified, all blocks are paused.
func NewPauseBlock(blockIDs ...any) *Function {
	return NewFunction(
		func(block *routine.Block) routine.Flow {
			block.Routine().Pause(blockIDs...)
			return routine.FlowNext
		},
	)
}

// NewStopBlock creates a Function action that deactivates the specified blocks
// in the currently running Routine. Any other blocks are unaffected.
// If no block IDs are specified, all blocks are stopped.
func NewStopBlock(blockIDs ...any) *Function {
	return NewFunction(
		func(block *routine.Block) routine.Flow {
			block.Routine().Stop(blockIDs...)
			return routine.FlowNext
		},
	)
}

// NewSetIndex creates a Function action that sets the index of the current block to the
// specified Action index number.
// (In other words, NewSetIndex(0) restarts the Block.)
func NewSetIndex(index int) *Function {
	return NewFunction(
		func(block *routine.Block) routine.Flow {
			block.SetIndex(index)
			return routine.FlowNext
		},
	)
}

// NewFinish creates a Function action that simply returns routine.FlowFinish, indicating
// that the current Block has finished and should stop running.
func NewFinish() *Function {
	return NewFunction(
		func(block *routine.Block) routine.Flow {
			return routine.FlowFinish
		},
	)
}

// NewLoop creates a Function action that simply loops the current block's execution when it is executed.
func NewLoop() *Function {
	return NewFunction(func(block *routine.Block) routine.Flow {
		block.SetIndex(0)
		return routine.FlowNext
	})
}
