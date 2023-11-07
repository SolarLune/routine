package actions

import (
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

// Function is a Action that runs a customizeable function.
type Function struct {
	InitFunc func(block *routine.Block)              // The function to run when the ActionFunc object is initialized (before polling)
	PollFunc func(block *routine.Block) routine.Flow // The function to run when polled
}

// NewFunction creates and returns a ActionFunc object with the polling function set to the
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

// NewGateOption creates a new GateOption object, which represents a choice in an ActionGate. The checkFunc()
// argument signifies a function that determines if the entry is made active. If it is made active,
// then the Actions within (and only those Actions) are executed in sequence until the end.
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
	Entries     []*GateOption
	ActiveEntry *GateOption
	onIdle      func()
	onChoose    func()
}

// NewGate creates a Gate action, which allows you to effectively choose one "route" or "choice"
// option among many.
// Once one entry has been made active, it will stay active until the owning Block revisits the
// ActionGate.
func NewGate(entries ...*GateOption) *Gate {
	return &Gate{
		Entries: entries,
	}
}

func (c *Gate) Init(block *routine.Block) {
	for _, entry := range c.Entries {
		entry.actions[0].Init(block)
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
		for _, entry := range c.Entries {
			if entry.CheckFunc() {
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

// Collection creates a ActionCollection, which is a collection of Actions that gets "absorbed" into a
// Block or GateEntry action.
func NewCollection(Actions ...routine.Action) *Collection {
	return &Collection{
		actions: Actions,
	}
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
func NewSwitchBlock(blockIDs ...any) *Function {
	return NewFunction(
		func(block *routine.Block) routine.Flow {
			block.Routine.SwitchBlock(blockIDs...)
			return routine.FlowNext
		},
	)
}

// NewActivateBlock creates a Function action that activates the specified blocks in the
// currently running Routine. Any other blocks are unaffected.
func NewActivateBlock(blockIDs ...any) *Function {
	return NewFunction(
		func(block *routine.Block) routine.Flow {
			block.Routine.ActivateBlock(blockIDs...)
			return routine.FlowNext
		},
	)
}

// NewDeactivateBlock creates a Function action that deactivates the specified blocks
// in the currently running Routine. Any other blocks are unaffected.
func NewDeactivateBlock(blockIDs ...any) *Function {
	return NewFunction(
		func(block *routine.Block) routine.Flow {
			block.Routine.DeactivateBlock(blockIDs...)
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

// NewFinish creates a ActionFunc that simply returns routine.FlowFinish, indicating
// that the Routine has finished and should stop running.
func NewFinish() *Function {
	return NewFunction(
		func(block *routine.Block) routine.Flow {
			return routine.FlowFinish
		},
	)
}
