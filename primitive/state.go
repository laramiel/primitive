package primitive

import "github.com/laramiel/primitive/primitive/shape"

type State struct {
	Worker      *Worker
	Shape       shape.Shape
	Z           int
	Alpha       int
	MutateAlpha bool
	Score       float64
}

func NewState(worker *Worker, shape shape.Shape, alpha int) *State {
	mutateAlpha := false
	if alpha == 0 {
		alpha = 128
		mutateAlpha = true
	}
	alpha = clampInt(alpha, 1, 255)
	z := worker.RandomZ()
	return &State{worker, shape, z, alpha, mutateAlpha, -1}
}

func (state *State) Energy() float64 {
	if state.Score < 0 {
		state.Score = state.Worker.Energy(state.Shape, state.Alpha)
	}
	return state.Score
}

func (state *State) DoMove() interface{} {
	oldState := state.Copy()
	state.Shape.Mutate(&state.Worker.Plane)
	if state.MutateAlpha {
		rnd := state.Worker.Plane.Rnd
		state.Alpha = clampInt(state.Alpha+rnd.Intn(21)-10, 1, 255)
	}
	state.Score = -1
	return oldState
}

func (state *State) UndoMove(undo interface{}) {
	oldState := undo.(*State)
	state.Shape = oldState.Shape
	state.Alpha = oldState.Alpha
	state.Score = oldState.Score
}

func (state *State) Copy() Annealable {
	return &State{
		state.Worker, state.Shape.Copy(), state.Z, state.Alpha, state.MutateAlpha, state.Score}
}
