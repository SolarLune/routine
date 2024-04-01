// routine is a package for creating sequences of events, primarily for game development in Golang.
package routine

// Properties represents a kind of "local memory" for an Execution object.
type Properties map[any]any

// Init will initialize a property by the given name with the given value
// if it doesn't already exist.
func (p *Properties) Init(propName any, toValue any) {
	if p.Has(propName) {
		return
	}
	p.Set(propName, toValue)
}

// Get returns the value associated with the given property identifier.
func (p Properties) Get(propName any) any {
	return p[propName]
}

// Has returns if the Properties object has a property associated with
// the given identifier.
func (p Properties) Has(propName any) bool {
	_, exists := p[propName]
	return exists
}

// Set sets the Properties object with the given property name to the
// value passed.
func (p *Properties) Set(propName any, value any) {
	(*p)[propName] = value
}

// Clear clears the properties map.
func (p *Properties) Clear() {
	for k := range *p {
		delete(*p, k)
	}
}

// Delete deletes a key out of the properties map.
func (p *Properties) Delete(keyName string) {
	delete(*p, keyName)
}

// Flow is simply a uint8, and represents what a Routine should do following a Action's action.
type Flow uint8

const (
	// FlowIdle means that the Routine should cycle and do the same Action again the following Update() cycle.
	FlowIdle Flow = iota
	// FlowNext means that the Routine should move on to the next Action in the Block.
	// If this is returned from the last Action in a Block, the Block will loop.
	FlowNext
	// FlowFinish indicates the Block should finish its execution, deactivating afterwards.
	FlowFinish
)

// Action is an interface that represents an object that can Action and direct the flow of a Routine.
type Action interface {
	Init(block *Block)      // The Init function is called when a Action is switched to.
	Poll(block *Block) Flow // The Poll function is called every frame and can return a Flow, indicating what the Routine should do next.
}

// ActionCollectionable identifies an interface for an Action that allows it to return a slice of Actions to be added to Blocks, Gates, or Collections in definition.
type ActionCollectionable interface {
	Actions() []Action
}

// ActionIdentifiable identifies an interface for an action that allows that Action to be used for jumping (as though it were a label).
type ActionIdentifiable interface {
	ID() any
}

// Block represents a block of Actions. Blocks execute Actions in sequence, and have an ID that allows them to be
// activated or deactivated at will by their owning Routine.
type Block struct {
	currentlyActive bool
	active          bool
	currentFrame    int // The current frame of the Block for the currently running Action.
	ID              any
	Actions         []Action
	index           int
	indexChanged    bool
	routine         *Routine
}

// SetIndex sets the index of the Action sequence of the Block to the value given.
// This effectively "sets the playhead" of the Block to point to the Action in the given
// slot.
func (b *Block) SetIndex(index int) {

	if index < 0 {
		index = 0
	}

	if index > len(b.Actions)-1 {
		index = len(b.Actions) - 1
	}

	if b.index != index {

		b.index = index
		b.Actions[b.index].Init(b)
		b.currentFrame = 0
		if b.currentlyActive {
			b.indexChanged = true
		}

	}

}

// JumpTo sets the Block's execution index to the index of a ActionLabel, using the label
// provided.
// If it finds the Label, then it will jump to and return that index. Otherwise, it will return -1.
func (b *Block) JumpTo(labelID any) int {
	for i, c := range b.Actions {
		if label, ok := c.(ActionIdentifiable); ok {
			if label.ID() == labelID {
				b.SetIndex(i)
				return i
			}
		}
	}
	return -1
}

// Index returns the index of the currently active Action in the Block.
func (b *Block) Index() int {
	return b.index
}

func (b *Block) update() {

	if !b.currentlyActive {
		return
	}

	b.indexChanged = false

	p := b.Actions[b.index].Poll(b)

	b.currentFrame++

	switch p {
	case FlowNext:

		if !b.indexChanged {
			b.index++
		}

		if b.index > len(b.Actions)-1 {
			b.index = 0
			b.active = false
			b.currentlyActive = false
		}

		b.Actions[b.index].Init(b)
		b.currentFrame = 0

		if b.active {
			b.update() // We call update again because it should move on unless it's idling, specifically
		}

	case FlowFinish:
		b.index = 0
		b.active = false // Restart if we're going to the next Action and we're at the end of the block
		b.currentlyActive = false
		b.Actions[b.index].Init(b)
		b.currentFrame = 0

	case FlowIdle:

		if b.indexChanged {
			b.Actions[b.index].Init(b)
			b.currentFrame = 0
		}

	}

}

// Run runs the specified block.
func (b *Block) Run() {
	b.active = true
}

// Running returns if the Block is active.
func (b *Block) Running() bool {
	return b.active
}

// Pause pauses the specified block, so that it isn't active when the Routine is run. When it is run again, it resumes execution at its current action.
func (b *Block) Pause() {
	b.active = false
}

// Restart restarts the block.
func (b *Block) Restart() {
	b.index = -1
	b.SetIndex(0)
}

// Stop stops the Block, so that it restarts when it is run again.
func (b *Block) Stop() {
	b.Pause()
	b.Restart()
}

// Routine returns the currently running routine.
func (b *Block) Routine() *Routine {
	return b.routine
}

// CurrentFrame returns the current frame of the Block's execution of the currently executed Action.
// This increases by 1 every Routine.Update() call until the Block executes another Action.
func (b *Block) CurrentFrame() int {
	return b.currentFrame
}

// Routine represents a container to run Blocks of code.
type Routine struct {
	Blocks     []*Block
	properties *Properties
}

// New creates a new Routine.
func New() *Routine {
	r := &Routine{
		Blocks:     []*Block{},
		properties: &Properties{},
	}
	return r
}

// Define defines a Block using the ID given and the list of Actions provided and adds it to the Routine.
// The ID can be of any comparable type.
// Define returns the new Block as well.
// If a block with the given blockID already exists, Define will remove the previous one.
func (r *Routine) Define(id any, Actions ...Action) *Block {

	newActions := []Action{}

	for _, c := range Actions {
		if collection, ok := c.(ActionCollectionable); ok {
			newActions = append(newActions, collection.Actions()...)
		} else {
			newActions = append(newActions, c)
		}
	}

	newBlock := &Block{
		ID:      id,
		routine: r,
		Actions: newActions,
	}

	for i, b := range r.Blocks {
		if b.ID == id {
			r.Blocks[i] = nil
			r.Blocks = append(r.Blocks[:i], r.Blocks[i+1:]...)
		}
	}

	r.Blocks = append(r.Blocks, newBlock)
	return newBlock
}

// Properties returns the Properties object for the Routine.
func (r *Routine) Properties() *Properties {
	return r.properties
}

// Update updates the Routine - this should be called once per frame.
func (r *Routine) Update() {

	for _, block := range r.Blocks {
		block.currentlyActive = block.active
	}

	for _, block := range r.Blocks {
		block.update()
	}

}

// Run runs Blocks with the given IDs.
// If no block IDs are given, then all blocks contained in the Routine are run.
func (r *Routine) Run(blockIDs ...any) {
	if len(blockIDs) == 0 {
		for _, block := range r.Blocks {
			block.Run()
		}
	} else {

		for _, label := range blockIDs {
			for _, block := range r.Blocks {
				if block.ID == label {
					block.Run()
					break
				}
			}
		}

	}
}

// Pause pauses Blocks with the given IDs.
// If no block IDs are given, then all blocks contained in the Routine are paused.
func (r *Routine) Pause(blockIDs ...any) {
	if len(blockIDs) == 0 {
		for _, block := range r.Blocks {
			block.Pause()
		}
	} else {

		for _, label := range blockIDs {
			for _, block := range r.Blocks {
				if block.ID == label {
					block.Pause()
					break
				}
			}
		}

	}

}

// Stop stops Blocks with the given IDs.
// If no block IDs are given, then all blocks contained in the Routine are stopped.
func (r *Routine) Stop(blockIDs ...any) {
	if len(blockIDs) == 0 {
		for _, block := range r.Blocks {
			block.Stop()
		}
	} else {

		for _, label := range blockIDs {
			for _, block := range r.Blocks {
				if block.ID == label {
					block.Stop()
					break
				}
			}
		}
	}

}

// Restart restarts Blocks with the given IDs.
// If no block IDs are given, then all blocks contained in the Routine are restarted.
func (r *Routine) Restart(blockIDs ...any) {
	if len(blockIDs) == 0 {

		for _, block := range r.Blocks {
			block.Restart()
		}

	} else {

		for _, label := range blockIDs {
			for _, block := range r.Blocks {
				if block.ID == label {
					block.Restart()
					break
				}
			}
		}

	}
}

// Running returns true if at least one Block is running with at least one of the given IDs in the Routine.
// If no IDs are given, then any running Blocks will return.
func (r *Routine) Running(ids ...any) bool {
	if len(ids) == 0 {
		for _, b := range r.Blocks {
			if b.Running() {
				return true
			}
		}
	} else {

		for _, id := range ids {
			for _, b := range r.Blocks {
				if b.ID == id && b.Running() {
					return true
				}
			}
		}

	}
	return false
}

// BlockByID returns any Block found with the given ID.
// If no Block with the given id is found, nil is returned.
func (r *Routine) BlockByID(id any) *Block {
	for _, block := range r.Blocks {
		if block.ID == id {
			return block
		}
	}
	return nil
}
