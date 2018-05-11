package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alexanderbez/alien-invasion/world"
)

func main() {
	var (
		mapFile   string
		numAliens uint
	)

	flag.StringVar(&mapFile, "map", "", "file containing the map definition")
	flag.UintVar(&numAliens, "n", 0, "number of aliens to use in the simulation")

	flag.Parse()

	if len(mapFile) == 0 {
		log.Fatalln("invalid map definition: no file specified")
	}

	if numAliens == 0 {
		log.Fatalln("invalid number of aliens: must be greater than zero")
	}

	worldMap, err := buildWorldMap(mapFile)
	if err != nil {
		log.Fatalf("failed to build map from file: %v", err)
	}

	// We assume there can be no more than twice the number of aliens as there
	// are cities. In otherwords, upon seeding the map with aliens, at most each
	// can be occupied by two aliens.
	if numAliens > worldMap.NumCities()*2 {
		log.Fatalf("invalid number of aliens: cannot have more than 2x of unique cities")
	}

	worldMap.SeedAliens(numAliens)
	fmt.Println(worldMap)
}

// buildWorldMap builds a map from a give file. The map definition file has one
// city per line. The city name is first, followed by 1-4 directions
// (north, south, east, or west). Each one represents a road to another city
// that lies in that direction. The city and each of the pairs are separated by
// a single space, and the directions are separated from their respective
// cities with an equals (=) sign. An error is returned if reading the file
// fails at any point or if the map definition does not adhere to the given
// schema.
func buildWorldMap(mapFile string) (*world.Map, error) {
	file, err := os.Open(mapFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	worldMap := world.NewMap()

	// Create a scanner to read each map entry line by line
	//
	// Note: We assume the line entry can fit into the scanner's buffer
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, " ")

		if len(tokens) == 0 {
			return nil, errors.New("invalid line in map definition")
		}

		cityName := tokens[0]

		if len(tokens) > 1 {
			for _, link := range tokens[1:] {
				linkTokens := strings.Split(link, "=")

				if len(linkTokens) != 2 {
					return nil, errors.New("invalid line in map definition")
				}

				linkDir := linkTokens[0]
				linkCity := linkTokens[1]

				err := worldMap.AddLink(cityName, linkDir, linkCity)
				if err != nil {
					return nil, fmt.Errorf("failed to insert city link to map: %v", err)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return worldMap, nil
}
