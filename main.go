package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alexanderbez/alien-invasion/simulation"
	"github.com/alexanderbez/alien-invasion/world"
)

func main() {
	var (
		mapFile   string
		outFile   string
		numAliens uint
	)

	flag.StringVar(&mapFile, "map", "", "file containing the map definition")
	flag.StringVar(&outFile, "out", "", "output file to write resulting map to")
	flag.UintVar(&numAliens, "n", 0, "number of aliens to use in the simulation")

	flag.Parse()

	if len(mapFile) == 0 {
		log.Fatalln("invalid map definition: no file specified")
	} else if len(outFile) == 0 {
		log.Fatalln("invalid output definition: no file specified")
	} else if numAliens == 0 {
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

	// Seed the map with 'n' aliens scattered randomly throughout the city and
	// invoke an initial series of alien fights where a search of the map
	// (graph) is done looking for city alien occupation equal to MaxOccupancy.
	worldMap.SeedAliens(numAliens)
	worldMap.ExecuteFights()

	sim := simulation.NewSimulation(worldMap)

	if err := sim.Run(); err != nil {
		log.Fatalf("failed to execute alien invasion simulation: %v", err)
	}

	log.Println("simulation complete")

	if err := writeMapToFile(worldMap, outFile); err != nil {
		log.Fatalf("failed to write map to file: %v", err)
	}
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

				worldMap.AddLink(cityName, linkTokens[0], linkTokens[1])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return worldMap, nil
}

// writeMapToFile writes a given world map to the file at path 'outPath'. An
// error is returned if the file cannot be created or written to.
func writeMapToFile(worldMap *world.Map, outPath string) error {
	fileHandle, err := os.Create(outPath)
	if err != nil {
		return err
	}

	defer fileHandle.Close()

	writer := bufio.NewWriter(fileHandle)
	defer fileHandle.Close()

	for _, city := range worldMap.Cities() {
		s := city.String()
		if len(s) != 0 {
			fmt.Fprintln(writer, s)
			writer.Flush()
		}
	}

	return nil
}
