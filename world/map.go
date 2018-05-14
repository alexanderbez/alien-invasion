package world

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/alexanderbez/alien-invasion/queue"
	"github.com/alexanderbez/alien-invasion/utils"
)

const (
	// MaxOccupancy reflects the maximum number of aliens that may occupy any
	// given city.
	MaxOccupancy = 2
	// MaxOutDegree reflects the maximum number of out degree edges from a
	// city. Only north, south, east, and west links can be made.
	MaxOutDegree = 4
)

// Map implements a representation of a world map. It's underlying
// implementation is a directed graph. A list of city names are also tracked as
// to be able to pseudo randomly pick cities.
type Map struct {
	cities map[string]*City
	aliens map[string]*Alien
}

// City implements a city in a world map that contains a name, occupied aliens
// and directional links (directional edges) to other cities by name both in
// and out of the city.
//
// Note: We take the space overhead of storing inLinks in order to reduce time
// overhead for certain operations.
type City struct {
	name           string
	inLinks        []string
	outLinks       []string
	alienOccupancy map[string]*Alien
}

// Priority implements the Heapable interface.
func (c *City) Priority(other interface{}) bool {
	if t, ok := other.(*City); ok {
		return len(c.outLinks) > len(t.outLinks)
	}

	return false
}

// NewMap returns a reference to a new initialized Map.
func NewMap() *Map {
	return &Map{
		cities: make(map[string]*City),
		aliens: make(map[string]*Alien),
	}
}

// AlienNames returns a unique list of all the aliens that exist in the map.
func (m *Map) AlienNames() []string {
	alienNames := make([]string, 0, len(m.aliens))

	for alienName := range m.aliens {
		alienNames = append(alienNames, alienName)
	}

	return alienNames
}

// NumCities returns the total number of unique cities in the Map.
func (m *Map) NumCities() uint {
	return uint(len(m.cities))
}

// NumAliens returns the total number of unique aliens occupying a city in the
// map.
func (m *Map) NumAliens() uint {
	return uint(len(m.aliens))
}

// CityNames returns a list of all the unique city names in the map.
func (m *Map) CityNames() []string {
	cityNames := make([]string, 0, m.NumCities())

	for cityName := range m.cities {
		cityNames = append(cityNames, cityName)
	}

	return cityNames
}

// AddLink adds a link (directional edge) from an origin city to a linked city.
// If the origin city or linked city do not exist in the graph, they are
// initialized and added. Finally, the out link is added to the origin city and
// the in link is added to the linked city.
func (m *Map) AddLink(cityName, linkCityName string) {
	// Add the origin city to the map of cities
	if _, ok := m.cities[cityName]; !ok {
		m.cities[cityName] = &City{
			name:           cityName,
			inLinks:        make([]string, 0, MaxOutDegree),
			outLinks:       make([]string, 0, MaxOutDegree),
			alienOccupancy: make(map[string]*Alien, MaxOccupancy),
		}
	}

	// Add the linked city to the map of cities
	if _, ok := m.cities[linkCityName]; !ok {
		m.cities[linkCityName] = &City{
			name:           linkCityName,
			inLinks:        make([]string, 0, MaxOutDegree),
			outLinks:       make([]string, 0, MaxOutDegree),
			alienOccupancy: make(map[string]*Alien, MaxOccupancy),
		}
	}

	// Add outbound and inbound links (directional edges)
	m.cities[cityName].outLinks = append(m.cities[cityName].outLinks, linkCityName)
	m.cities[linkCityName].inLinks = append(m.cities[linkCityName].inLinks, cityName)
}

// MoveAlien attempts to move an alien on the map from one city to another
// following a valid direction. The algorithm for finding a valid move follows:
//
// 1. Find an alien that occupies a city with at least one valid out link (edge)
// 2. If that link leads to a city that has space for an additional alien, then:
// 2a. Update the alien's city
// 2b. Remove alien from current city
// 2c. Add alien to new city
// 3. Otherwise, continue evaluating other out links. If no links are valid,
// then try another alien.
// 4. If no alien can be moved, an error is returned.
func (m *Map) MoveAlien() error {
	// We will get some pseudo randomness iterating over the city's list of
	// aliens.
	for _, alien := range m.aliens {
		occupiedCity := alien.cityName
		city := m.cities[occupiedCity]

		for _, linkCityName := range utils.ShuffleStrings(city.outLinks) {
			linkCity := m.cities[linkCityName]

			if len(linkCity.alienOccupancy) < MaxOccupancy {
				delete(city.alienOccupancy, alien.name)

				alien.cityName = linkCity.name
				linkCity.alienOccupancy[alien.name] = alien

				return nil
			}
		}
	}

	return errors.New("unable to move any alien")
}

// destroyCity removes a given city from the map (directed graph) in addition
// to any links (edges) that lead into or out of it. The aliens that occupy the
// city are also destroyed. The resulting list of destroyed aliens is returned.
func (m *Map) destroyCity(city *City) []string {
	destroyedAliens := make([]string, 0, MaxOccupancy)

	for alienName := range city.alienOccupancy {
		destroyedAliens = append(destroyedAliens, alienName)
		delete(m.aliens, alienName)
	}

	// Remove the destroyed city from all inks (inbound and outbound edges) from
	// any city that can get to the destroyed city.
	for _, inCityLinkName := range city.inLinks {
		inCityLink := m.cities[inCityLinkName]

		inCityLink.outLinks = utils.RemoveStrFromSlice(inCityLink.outLinks, city.name)
		inCityLink.inLinks = utils.RemoveStrFromSlice(inCityLink.inLinks, city.name)
	}

	// Remove the destroyed city from any in links (inbound edges) from any city
	// that can get to the destroyed city.
	for _, outCityLinkName := range city.outLinks {
		outCityLink := m.cities[outCityLinkName]

		outCityLink.inLinks = utils.RemoveStrFromSlice(outCityLink.inLinks, city.name)
	}

	delete(m.cities, city.name)
	return destroyedAliens
}

// ExecuteFights simulates a fight between any two aliens if there are any
// found occupying a city. All the aliens are examined along with the city they
// occupy. If any such city is occupied by MaxOccupancy, a fight is simulated
// and the aliens along with the city are destroyed. In addition, any links
// (edges) that lead into or out of the destroyed city are also removed from
// the map.
func (m *Map) ExecuteFights() {
	for _, alien := range m.aliens {
		occupiedCity := alien.cityName
		city := m.cities[occupiedCity]

		// If maximum occupancy has been reached for a city, the occupying
		// aliens will fight and destroy the city. As a result, the following
		// will happen:
		//
		// 1. Both aliens will be removed from the map's known collection of
		// aliens.
		// 2. The city will be removed from the map and so are any links that
		// lead into or out of it.
		if len(city.alienOccupancy) == MaxOccupancy {
			destroyedAliens := m.destroyCity(city)
			log.Printf("%s has been destroyed by %s!", city.name, strings.Join(destroyedAliens, " and "))
		}
	}
}

// SeedAliens adds n aliens to the world map at pseudo random cities. At most
// 'MaxOccupancy' aliens can occupy a city at any given time. It is assumed the
// number of aliens to seed is valid and as such each alien will find a valid
// city to occupy. Alien occupancy is preferred in cities with out roads
// (out edges).
func (m *Map) SeedAliens(n uint) {
	pq := queue.NewPriorityQueue()

	// Add all the cities pseudo-randomly to a priority queue. Priority is
	// based on the total number of out degree links of a city.
	for _, city := range m.cities {
		pq.Push(city)
	}

	// We assume the invariant that there are enough cities to occupy all 'n'
	// aliens.
	seededAliens := uint(0)
	for seededAliens != n {
		city := pq.Pop().(*City)

		for i := 0; i < MaxOccupancy && seededAliens != n; i++ {
			alien := &Alien{
				name:     fmt.Sprintf("alien%d", seededAliens+1),
				cityName: city.name,
			}

			city.alienOccupancy[alien.name] = alien
			m.aliens[alien.name] = alien

			seededAliens++
		}
	}
}

// String implements the stringer interface.
func (m *Map) String() (s string) {
	for _, city := range m.cities {
		aliens := make([]string, 0, len(city.alienOccupancy))

		for alien := range city.alienOccupancy {
			aliens = append(aliens, alien)
		}

		s += fmt.Sprintf(
			"{city: %s, outLinks: %s, inLinks: %s, alienOccupancy: [%s]}\n",
			city.name, city.outLinks, city.inLinks, strings.Join(aliens, " "),
		)
	}

	return
}
