package simulation

import (
	"github.com/alexanderbez/alien-invasion/world"
)

const (
	minAlienMoves = 10000
)

// Simulation reflects a simulation of an alien invasion on a given world map.
type Simulation struct {
	alienMap   *world.Map
	alienMoves map[string]uint
}

// NewSimulation returns a reference to a new initialized alien invasion
// Simulation. It adds all known alien names to the map of alien moves ahead of
// time so they can be removed efficiently once an alien has reached
// minAlienMoves.
func NewSimulation(alienMap *world.Map) *Simulation {
	s := &Simulation{
		alienMap:   alienMap,
		alienMoves: make(map[string]uint),
	}

	for _, alienName := range alienMap.AlienNames() {
		s.alienMoves[alienName] = 0
	}

	return s
}

// Run executes an alien invasion simulation. It will continue to execute
// random alien moves and attempt to fight them to destroy cities. After each
// single random alien move, it will track the total number of moves for that
// given alien. The simulation will terminate when all the aliens have been
// destroyed or each alien has moved at least 'minAlienMoves' times. An error
// is returned if the simulation fails to move any alien during a run.
func (s *Simulation) Run() error {
	for s.canContinue() {
		alienName, err := s.alienMap.MoveAlien()
		if err != nil {
			return err
		}

		_, ok := s.alienMoves[alienName]
		if ok {
			s.alienMoves[alienName]++

			// Once an alien has moved at least 'minAlienMoves' times, we can
			// avoid having to track/count his moves.
			if s.alienMoves[alienName] >= minAlienMoves {
				delete(s.alienMoves, alienName)
			}
		}

		s.alienMap.ExecuteFights()
	}

	return nil
}

// canContinue return a boolean on whether or not a simulation can continue to
// run. A simulation can continue if not all aliens have been destroyed or not
// all aliens have moved at least 'minAlienMoves' times.
func (s *Simulation) canContinue() bool {
	if s.alienMap.NumAliens() == 0 {
		return false
	}

	for _, totalMoves := range s.alienMoves {
		if totalMoves < minAlienMoves {
			return true
		}
	}

	return false
}
