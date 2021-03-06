package primitive

import (
	"math"
	"math/rand"
)

type Annealable interface {
	Energy() float64
	DoMove(temp float64) interface{}
	UndoMove(interface{})
	Copy() Annealable
}

func HillClimb(state Annealable, maxAge int) Annealable {
	state = state.Copy()
	bestState := state.Copy()
	bestEnergy := state.Energy()
	step := 0
	for age := 0; age < maxAge; age++ {
		undo := state.DoMove(1.0)
		energy := state.Energy()
		if energy >= bestEnergy {
			state.UndoMove(undo)
		} else {
			vvv("step: %d, energy: %.6f\n", step, energy)
			bestEnergy = energy
			bestState = state.Copy()
			age = -1
		}
		step++
	}
	return bestState
}

func PreAnneal(state Annealable, iterations int) float64 {
	state = state.Copy()
	previous := state.Energy()
	var total float64
	for i := 0; i < iterations; i++ {
		state.DoMove(1.0)
		energy := state.Energy()
		total += math.Abs(energy - previous)
		previous = energy
	}
	return total / float64(iterations)
}

func Anneal(state Annealable, maxTemp, minTemp float64, steps int) Annealable {
	factor := -math.Log(maxTemp / minTemp)
	state = state.Copy()
	bestState := state.Copy()
	bestEnergy := state.Energy()
	previousEnergy := bestEnergy
	for step := 0; step < steps; step++ {
		pct := float64(step) / float64(steps-1)
		temp := maxTemp * math.Exp(factor*pct)
		undo := state.DoMove(1.0)
		energy := state.Energy()
		change := energy - previousEnergy
		if change > 0 && math.Exp(-change/temp) < rand.Float64() {
			state.UndoMove(undo)
		} else {
			previousEnergy = energy
			if energy < bestEnergy {
				pct := float64(step*100) / float64(steps)
				vvv("step: %d of %d (%.1f%%), temp: %.3f, energy: %.6f\n",
					step, steps, pct, temp, energy)
				bestEnergy = energy
				bestState = state.Copy()
			}
		}
	}
	return bestState
}
