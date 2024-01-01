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
	// FlowFinishBlock indicates the Block should finish its execution.
	FlowFinishBlock
	// FlowFinishRoutine means that the entire Routine should finish its execution.
	FlowFinishRoutine
)

// Action is an interface that represents an object that can Action and direct the flow of a Routine.
type Action interface {
	Init(block *Block)      // The Init function is called when a Action is switched to.
	Poll(block *Block) Flow // The Poll function is called every frame and can return a Flow, indicating what the Routine should do next.
}

type actionCollectionable interface {
	Actions() []Action
}

type actionIdentifiable interface {
	ID() any
}

// Block represents a block of Actions. Blocks execute Actions in sequence, and have an ID that allows them to be
// activated or deactivated at will by their owning Routine.
type Block struct {
	isActive     bool
	Active       bool
	ID           any
	Actions      []Action
	index        int
	indexChanged bool
	Routine      *Routine
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
		if b.Routine.Running() {
			b.indexChanged = true
		}

	}

}

// JumpTo sets the Block's execution index to the index of a ActionLabel, using the label
// provided.
// If it finds the Label, then it will return that index. Otherwise, it will return -1.
func (b *Block) JumpTo(labelID any) int {
	for i, c := range b.Actions {
		if label, ok := c.(actionIdentifiable); ok {
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

// Restart restarts the block.
func (b *Block) Restart() {
	b.index = -1
	b.SetIndex(0)
}

func (b *Block) update() {

	if !b.isActive {
		return
	}

	b.indexChanged = false

	p := b.Actions[b.index].Poll(b)

	switch p {
	case FlowNext:

		if !b.indexChanged {
			b.index++
		}

		if b.index > len(b.Actions)-1 {
			b.index = 0
			b.Active = false // Restart if we're going to the next Action and we're at the end of the block
			b.isActive = false
		}

		b.Actions[b.index].Init(b)

	case FlowFinishBlock:
		b.index = 0
		b.Active = false // Restart if we're going to the next Action and we're at the end of the block
		b.isActive = false
		b.Actions[b.index].Init(b)

	case FlowFinishRoutine:
		b.Routine.Stop()

	case FlowIdle:

		if b.indexChanged {
			b.Actions[b.index].Init(b)
		}

	}

}

// Routine represents a running function that will execute
// until the Routine is finished.
type Routine struct {
	running           bool
	Blocks            []*Block
	properties        *Properties
	AutomaticallyStop bool // If the Routine should automatically stop if no Blocks are active
}

// New creates a new Routine.
func New() *Routine {
	r := &Routine{
		Blocks:            []*Block{},
		properties:        &Properties{},
		AutomaticallyStop: true,
	}
	return r
}

// DefineBlock defines a Block using the ID given and the list of Actions provided and adds it to the Routine.
// The ID can be of any comparable type.
// DefineBlock returns the new Block as well.
func (r *Routine) DefineBlock(blockID any, Actions ...Action) *Block {

	newActions := []Action{}

	for _, c := range Actions {
		if collection, ok := c.(actionCollectionable); ok {
			newActions = append(newActions, collection.Actions()...)
		} else {
			newActions = append(newActions, c)
		}
	}

	newBlock := &Block{
		ID:      blockID,
		Routine: r,
		Actions: newActions,
	}
	r.Blocks = append(r.Blocks, newBlock)
	if len(r.Blocks) == 1 {
		r.Blocks[0].Active = true
	}
	return newBlock
}

// Run starts the Routine.
func (r *Routine) Run() {
	if !r.running {
		r.running = true

		for _, b := range r.Blocks {
			if b.isActive {
				b.Actions[b.index].Init(b)
			}
		}

	}
}

// Running returns if the Routine is running.
func (r *Routine) Running() bool {
	return r.running
}

// Restart restarts the Routine.
func (r *Routine) Restart() {
	r.running = true
	for _, b := range r.Blocks {
		b.SetIndex(0)
	}
}

// Pause pauses the Routine; it does not alter the currently active blocks, or where those blocks are in terms of execution.
func (r *Routine) Pause() {
	if r.running {
		r.running = false
	}
}

// Stop stops the Routine, restarting it from scratch when it runs again.
func (r *Routine) Stop() {
	r.Restart()
	r.Pause()
}

// Properties returns the Properties object for the Routine.
func (r *Routine) Properties() *Properties {
	return r.properties
}

// Update updates the Routine - this should be called once per frame.
func (r *Routine) Update() {

	if r.running {

		for _, block := range r.Blocks {
			block.update()
		}

		anyBlocksActive := false

		for _, block := range r.Blocks {
			block.isActive = block.Active
			if block.isActive {
				anyBlocksActive = true
			}
		}

		if r.AutomaticallyStop && !anyBlocksActive {
			r.Stop()
		}

	}

}

// ActivateBlock activates blocks with the given IDs.
func (r *Routine) ActivateBlock(blockIDs ...any) {
	for _, block := range r.Blocks {
		for _, label := range blockIDs {
			if block.ID == label {
				block.Active = true
				break
			}
		}
	}
}

// IsBlockActive returns if any of the block IDs given belong to active Blocks.
func (r *Routine) IsBlockActive(blockIDs ...any) bool {
	for _, label := range blockIDs {
		for _, block := range r.Blocks {
			if block.isActive && block.ID == label {
				return true
			}
		}
	}
	return false
}

// Deactivate deactivates blocks with the given IDs.
func (r *Routine) DeactivateBlock(blockIDs ...any) {
	for _, block := range r.Blocks {
		for _, label := range blockIDs {
			if block.ID == label {
				block.Restart()
				block.Active = false
				break
			}
		}
	}
}

// SwitchBlock will only activate blocks with any of the given IDs, and deactivates all others.
func (r *Routine) SwitchBlock(blockIDs ...any) {

	for _, block := range r.Blocks {
		block.Active = false

		for _, label := range blockIDs {
			if block.ID == label {
				block.Active = true
				break
			}
		}

		// Restart inactivated blocks
		if !block.Active {
			block.Restart()
		}

	}

}

// BlockByID returns any Block found with the given ID.
func (r *Routine) BlockByID(idLabel any) *Block {
	for _, block := range r.Blocks {
		if block.ID == idLabel {
			return block
		}
	}
	return nil
}
