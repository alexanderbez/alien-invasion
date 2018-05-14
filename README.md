# Alien Invasion

Mad aliens are about to invade the earth!

```
You are given a map containing the names of cities in the non-existent world of
X. The map is in a file, with one city per line. The city name is first,
followed by 1-4 directions (north, south, east, or west). Each one represents a
road to another city that lies in that direction.

Given `n` aliens, these aliens start out at random places on the map, and wander around randomly,
following links. Each iteration, the aliens can travel in any of the directions
leading out of a city.

When two aliens end up in the same place, they fight, and in the process kill
each other and destroy the city. When a city is destroyed, it is removed from
the map, and so are any roads that lead into or out of it.
```

The underlying implementation of the fictitious world is a directed graph with both
in and out degree edges. Each city can be thought of as a graph vertex with at
most four in and out degrees. We track in degrees in order to facilitate easy
traversal and removal of vertices when a city is destroyed.

In addition, we keep track of all the aliens that occupy a city separately. Upon
creation of the map, the aliens are placed in random cities where the initial
preferred cities are those with the most number of out degrees as it'll be more
advantageous for these aliens to move about.

## Preliminary

### Tests

```
$ make test
```

### Build

```
$ make build
```

### Usage

```
$ ./alien-invasion-sim --map=<INPUT_FILE> --out=<OUTPUT_FILE> --n=<NUMBER_OF_ALIENS>
```

## Assumptions

- There are no more than 2x aliens of the number of cities in the map
- No more than two aliens can occupy a city, if a third alien attempts to enter it will be denied.
  - Note, this should never happen, as a fight will be initiated before the next alien move.
- A valid map will be provided such that a valid simulation will execute with a given `n`. In other words, the aliens won't be trapped
before 10,000 moves. If so, an error will be returned.
