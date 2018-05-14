package world

import (
	"reflect"
	"sort"
	"testing"
)

func buildMapFixtureEmpty() *Map {
	return NewMap()
}

func buildMapFixtureSimple() *Map {
	a1 := &Alien{name: "alien1", cityName: "foo"}
	a2 := &Alien{name: "alien2", cityName: "foo"}
	a3 := &Alien{name: "alien3", cityName: "bar"}
	a4 := &Alien{name: "alien4", cityName: "bar"}

	m := &Map{
		aliens: map[string]*Alien{
			a1.name: a1,
			a2.name: a2,
			a3.name: a3,
			a4.name: a4,
		},
		cities: map[string]*City{
			"foo": &City{
				name:     "foo",
				inLinks:  []string{"bar"},
				outLinks: []string{"bar"},
				alienOccupancy: map[string]*Alien{
					a1.name: a1,
					a2.name: a2,
				},
			},
			"bar": &City{
				name:     "bar",
				inLinks:  []string{"foo"},
				outLinks: []string{"foo"},
				alienOccupancy: map[string]*Alien{
					a3.name: a3,
					a4.name: a4,
				},
			},
		},
	}

	return m
}

func TestNewMap(t *testing.T) {
	m := buildMapFixtureEmpty()

	if m.cities == nil {
		t.Error("failed to initialize map; invalid cities collection")
	}

	if m.aliens == nil {
		t.Error("failed to initialize map; invalid aliens collection")
	}
}

func TestNumCities(t *testing.T) {
	testCases := []struct {
		m *Map
		e uint
	}{
		{
			m: buildMapFixtureEmpty(),
			e: 0,
		},
		{
			m: buildMapFixtureSimple(),
			e: 2,
		},
	}

	for _, tc := range testCases {
		r := tc.m.NumCities()

		if r != tc.e {
			t.Errorf("incorrect result: expected: %v, got: %v", tc.e, r)
		}
	}
}

func TestNumAliens(t *testing.T) {
	testCases := []struct {
		m *Map
		e uint
	}{
		{
			m: buildMapFixtureEmpty(),
			e: 0,
		},
		{
			m: buildMapFixtureSimple(),
			e: 4,
		},
	}

	for _, tc := range testCases {
		r := tc.m.NumAliens()

		if r != tc.e {
			t.Errorf("incorrect result: expected: %v, got: %v", tc.e, r)
		}
	}
}

func TestAlienNames(t *testing.T) {
	m1 := buildMapFixtureEmpty()

	if len(m1.AlienNames()) != 0 {
		t.Errorf("incorrect result: expected: %v, got: %v", 0, len(m1.AlienNames()))
	}

	m2 := buildMapFixtureSimple()
	e := make([]string, 0, len(m2.aliens))
	r := m2.AlienNames()

	for a := range m2.aliens {
		e = append(e, a)
	}

	sort.Strings(e)
	sort.Strings(r)

	if !reflect.DeepEqual(r, e) {
		t.Errorf("incorrect result: expected: %v, got: %v", r, e)
	}
}

func TestCityNames(t *testing.T) {
	testCases := []struct {
		m *Map
		e []string
	}{
		{
			m: buildMapFixtureEmpty(),
			e: []string{},
		},
		{
			m: buildMapFixtureSimple(),
			e: []string{"foo", "bar"},
		},
	}

	for _, tc := range testCases {
		r := tc.m.CityNames()

		sort.Strings(r)
		sort.Strings(tc.e)

		if !reflect.DeepEqual(r, tc.e) {
			t.Errorf("incorrect result: expected: %v, got: %v", tc.e, r)
		}
	}
}

func TestAddLink(t *testing.T) {
	m1 := buildMapFixtureEmpty()
	m1.AddLink("foo", "bar")

	if _, ok := m1.cities["foo"]; !ok {
		t.Errorf("expected %s to exist in map cities", "foo")
	}

	if _, ok := m1.cities["bar"]; !ok {
		t.Errorf("expected %s to exist in map cities", "foo")
	}

	if !reflect.DeepEqual(m1.cities["foo"].outLinks, []string{"bar"}) {
		t.Errorf("expected %s city to have valid out links: %v", "foo", []string{"bar"})
	}

	if !reflect.DeepEqual(m1.cities["bar"].inLinks, []string{"foo"}) {
		t.Errorf("expected %s linked city to have valid in links: %v", "bar", []string{"foo"})
	}
}

func TestMoveAlien(t *testing.T) {
	m1 := buildMapFixtureEmpty()

	if _, err := m1.MoveAlien(); err == nil {
		t.Errorf("expected error: no aliens to move in empty map")
	}

	m2 := buildMapFixtureSimple()
	m2.AddLink("foo", "qu-ux")

	if _, err := m2.MoveAlien(); err != nil {
		t.Errorf("unexpected error: alien should be able to move")
	}
}

func TestDestroyCity(t *testing.T) {
	m := buildMapFixtureSimple()
	c := m.cities["foo"]
	r := m.destroyCity(c)

	e := make([]string, 0, MaxOccupancy)
	for k := range c.alienOccupancy {
		e = append(e, k)
	}

	sort.Strings(e)
	sort.Strings(r)

	if !reflect.DeepEqual(r, e) {
		t.Errorf("incorrect result: expected: %v, got: %v", e, r)
	}

	if _, ok := m.cities["foo"]; ok {
		t.Errorf("expected city %s to be removed from the map", "foo")
	}

	for _, a := range e {
		if _, ok := m.aliens[a]; ok {
			t.Errorf("expected alien %s to be removed from the map", a)
		}
	}

	c = m.cities["bar"]

	if len(c.outLinks) != 0 {
		t.Errorf("expected linked city %s to not have destroyed city %s as an out link", "bar", "foo")
	}

	if len(c.inLinks) != 0 {
		t.Errorf("expected linked city %s to not have destroyed city %s as an in link", "bar", "foo")
	}
}

func TestExecuteFights(t *testing.T) {
	m := buildMapFixtureSimple()
	m.ExecuteFights()

	if len(m.cities) != 0 {
		t.Errorf("expected map to have no remaining cities: cities: %v", m.cities)
	}

	if len(m.aliens) != 0 {
		t.Errorf("expected map to have no remaining aliens: aliens: %v", m.aliens)
	}
}

func TestSeedAliens(t *testing.T) {
	m := buildMapFixtureEmpty()

	m.AddLink("foo", "bar")
	m.AddLink("foo", "qu-ux")
	m.AddLink("foo", "baz")
	m.AddLink("bar", "foo")
	m.AddLink("bar", "bee")

	m.SeedAliens(0)

	if len(m.aliens) != 0 {
		t.Errorf("expected map to have no aliens: got: %d, expected: %d", len(m.aliens), 0)
	}

	m.SeedAliens(10)

	if len(m.aliens) != 10 {
		t.Errorf("expected map to have correct number of aliens: got: %d, expected: %d", len(m.aliens), 10)
	}
}

func TestSeedAliensPriority(t *testing.T) {
	m := buildMapFixtureEmpty()

	m.AddLink("foo", "bar")
	m.AddLink("foo", "qu-ux")
	m.AddLink("foo", "baz")
	m.AddLink("bar", "foo")
	m.AddLink("bar", "bee")

	m.SeedAliens(4)

	if len(m.aliens) != 4 {
		t.Errorf("expected map to have correct number of aliens: got: %d, expected: %d", len(m.aliens), 4)
	}

	if len(m.cities["foo"].alienOccupancy) != 2 {
		t.Errorf("expected map to have correct number of aliens for priority city: got: %d, expected: %d",
			len(m.cities["foo"].alienOccupancy),
			2,
		)
	}

	if len(m.cities["bar"].alienOccupancy) != 2 {
		t.Errorf("expected map to have correct number of aliens for priority city: got: %d, expected: %d",
			len(m.cities["bar"].alienOccupancy),
			2,
		)
	}
}
