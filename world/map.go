package world

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	// MaxOccupancy reflects the maximum number of aliens that may occupy any
	// given city.
	MaxOccupancy = 2
)

// Map implements a representation of a world map. It's underlying
// implementation is a directed graph. A list of city names are also tracked as
// to be able to pseudo randomly pick cities.
type Map struct {
	rw sync.RWMutex

	cityNames []string
	cities    map[string]*City
	aliens    map[string]*Alien
}

// City implements a city in a world map that contains a name and directional
// links to other cities by name.
type City struct {
	name           string
	links          map[string]string
	alienOccupancy map[string]*Alien
}

// NewMap returns a reference to a new initialized Map.
func NewMap() *Map {
	return &Map{
		cityNames: make([]string, 0, 0),
		cities:    make(map[string]*City),
		aliens:    make(map[string]*Alien),
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
// Finally, the link is added to the origin city.
func (m *Map) AddLink(cityName, linkDir, linkCity string) {
	m.rw.Lock()
	defer m.rw.Unlock()

	// Add the origin city to the map of cities
	if _, ok := m.cities[cityName]; !ok {
		m.cityNames = append(m.cityNames, cityName)
		m.cities[cityName] = &City{
			name:           cityName,
			links:          make(map[string]string),
			alienOccupancy: make(map[string]*Alien, 2),
		}
	}

	// Add the linked city to the map of cities
	if _, ok := m.cities[linkCity]; !ok {
		m.cityNames = append(m.cityNames, linkCity)
		m.cities[linkCity] = &City{
			name:           linkCity,
			links:          make(map[string]string),
			alienOccupancy: make(map[string]*Alien, 2),
		}
	}

	// Add link to city
	city := m.cities[cityName]
	city.links[strings.ToLower(linkDir)] = linkCity
}

// MoveAlien attempts to move an alien on the map from one city to another
// following a valid direction. The algorithm for finding a valid move follows:
//
// 1. Find an alien that occupies a city with at least one valid link
// 2. If that link leads to a city that has space for an additional alien,
// then:
// 2a. Update the alien's city
// 2b. Remove alien from current city
// 2c. Add alien to new city
// 3. Otherwise, continue evaluating other links. If no links are valid, then
// try another alien.
// 4. If no alien can move, an error is returned.
func (m *Map) MoveAlien() error {
	m.rw.Lock()
	defer m.rw.Unlock()

	// We will get some pseudo randomness iterating over the city's list of
	// aliens and the map that contains the linked cities.
	for _, alien := range m.aliens {
		occupiedCity := alien.cityName
		city := m.cities[occupiedCity]

		for _, linkCity := range city.links {
			newOccupiedCity := m.cities[linkCity]

			if len(newOccupiedCity.alienOccupancy) < MaxOccupancy {
				alien.cityName = newOccupiedCity.name
				newOccupiedCity.alienOccupancy[alien.name] = alien

				delete(city.alienOccupancy, alien.name)

				return nil
			}
		}
	}

	return errors.New("unable to move any alien")
}

// SeedAliens adds n aliens to the world map at pseudo random cities. At most
// two aliens can occupy a city at any given time. It is assumed the number of
// aliens to seed is valid and as such each alien will find a valid city to
// occupy.
func (m *Map) SeedAliens(n uint) {
	m.rw.Lock()
	defer m.rw.Unlock()

	// Initialize a PRNG using the current time as a seed
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	for i := uint(0); i < n; i++ {
		var city *City

		// Find a city the alien can occupy
		for city == nil {
			tmpCityName := m.cityNames[r.Intn(len(m.cityNames))]
			tmpCity := m.cities[tmpCityName]

			if tmpCity != nil && len(tmpCity.alienOccupancy) < MaxOccupancy {
				city = tmpCity
			}
		}

		alien := &Alien{
			name:     fmt.Sprintf("alien%d", i+1),
			cityName: city.name,
		}

		city.alienOccupancy[alien.name] = alien
		m.aliens[alien.name] = alien
	}
}

// String implements the stringer interface.
func (m *Map) String() (s string) {
	m.rw.RLock()
	defer m.rw.RUnlock()

	for _, city := range m.cities {
		links := make([]string, 0, len(city.links))
		aliens := make([]string, 0, len(city.alienOccupancy))

		for linkDir, linkCity := range city.links {
			if len(linkDir) != 0 {
				links = append(links, fmt.Sprintf("%s:%s", linkDir, linkCity))
			}
		}

		for alien := range city.alienOccupancy {
			if len(alien) != 0 {
				aliens = append(aliens, alien)
			}
		}

		s += fmt.Sprintf(
			"{city: %s, links: [%s], alienOccupancy: [%s]}\n",
			city.name, strings.Join(links, " "), strings.Join(aliens, " "),
		)
	}

	return
}
