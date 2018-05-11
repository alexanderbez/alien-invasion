package world

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	westDir  = "west"
	eastDir  = "east"
	northDir = "north"
	southDir = "south"
)

// Map implements a representation of a world map. It's underlying
// implementation is a directed graph. A list of city names are also tracked as
// to be able to pseudo randomly pick cities.
type Map struct {
	rw sync.RWMutex

	cityNames []string
	cities    map[string]*City
}

// City implements a city in a world map that contains a name and directional
// links to other cities by name.
type City struct {
	name           string
	north          string
	south          string
	east           string
	west           string
	alienOccupancy uint8
}

// NewMap returns a reference to a new initialized Map.
func NewMap() *Map {
	return &Map{
		cityNames: make([]string, 0, 0),
		cities:    make(map[string]*City),
	}
}

// NumCities returns the total number of unique cities in the Map.
func (m *Map) NumCities() uint {
	m.rw.RLock()
	defer m.rw.RUnlock()
	return uint(len(m.cityNames))
}

// AddLink adds a link from an origin city to a linked city for a given
// direction. If the origin city or linked city do not exist in the graph, they
// are added in addition to being added to the list of known city names.
// Finally, the link is added to the origin city. An error is returned if the
// link direction is invalid.
func (m *Map) AddLink(cityName, linkDir, linkCity string) error {
	m.rw.Lock()
	defer m.rw.Unlock()

	// Add the origin city to the map of cities
	if _, ok := m.cities[cityName]; !ok {
		m.cityNames = append(m.cityNames, cityName)
		m.cities[cityName] = &City{name: cityName}
	}

	// Add the linked city to the map of cities
	if _, ok := m.cities[linkCity]; !ok {
		m.cityNames = append(m.cityNames, linkCity)
		m.cities[linkCity] = &City{name: linkCity}
	}

	city := m.cities[cityName]
	switch strings.ToLower(linkDir) {
	case northDir:
		city.north = linkCity
	case southDir:
		city.south = linkCity
	case eastDir:
		city.east = linkCity
	case westDir:
		city.west = linkCity
	default:
		return fmt.Errorf("invalid link direction provided for city %s", cityName)
	}

	return nil
}

// SeedAliens adds n aliens to the world map at pseudo random cities. At most
// two aliens can occupy a city at any given time.
func (m *Map) SeedAliens(n uint) {
	m.rw.Lock()
	defer m.rw.Unlock()

	// Initialize a PRNG
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	for i := uint(0); i < n; i++ {
		var city *City

		// Find a city the ith alien can occupy
		for city == nil {
			tmpCityName := m.cityNames[r.Intn(len(m.cityNames))]
			tmpCity := m.cities[tmpCityName]

			if tmpCity != nil && tmpCity.alienOccupancy < 2 {
				city = tmpCity
			}
		}

		city.alienOccupancy++
	}
}

// String implements the stringer interface.
func (m *Map) String() (s string) {
	m.rw.RLock()
	defer m.rw.RUnlock()

	for cityName, city := range m.cities {
		s += fmt.Sprintf(
			"{city: %s, links: [north: %s, south: %s, east: %s, west: %s], alienOccupancy: %d}\n",
			cityName, city.north, city.south, city.east, city.west, city.alienOccupancy,
		)
	}

	return
}
