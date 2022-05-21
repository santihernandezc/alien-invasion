# Alien Invasion

Given a file containing the names of cities, this simulation generates a graph representing a non-existent world. One city is specified per line, starting with the city name and followed by 1-4 directions (north, south, east, or west), each one representing a road to another city that lies in that direction.

The city and each of the pairs are separated by a single space, and the directions are separated from their respective cities with a `=` sign.

Example:

```
Foo north=Bar west=Baz south=Qu-ux
Bar south=Foo west=Bee
```

`N` aliens are created, where `N` is specified as a command-line argument. These aliens start out at random places on the map and wander around randomly.

Each iteration, the aliens can travel in any of the directions leading out of a city.

When two aliens end up in the same place, they fight, and in the process kill each other and destroy the city. When a city is destroyed, it is removed from the map, and so are any roads that lead into or out of it.

In our example above, if Bar were destroyed the map would now be something like:

```
Foo west=Baz south=Qu-ux
```

Once a city is destroyed, aliens can no longer travel to or through it. This may lead to aliens getting "trapped".

This program reads in the world map, creates `N` aliens, and unleashes them. The program runs until all the aliens have been destroyed, trapped, or have moved a predefined number of times each.

Once the program finishes, it prints out whatever is left of the world in the same format as the input file. It can also write this to a file if the `-o` flag is specified.

### Events

Events are logged to stdout in the following format:

```
🚷 Alien 1 is trapped forever in Foo
👾 Alien 2 moved from Bar to Baz
👀 Alien 2 found Alien 3 in Baz
💥 Baz has been destroyed by Alien 2 and Alien 3
```

### Usage

You can run the simulation with `go run main.go` using the following flags:

```
-path string
    path to the config file (default "config.txt")
-n int
    number of aliens for the simulation (default 5)
-movements int
    how many iterations this simulation is going to run (default 10000)
-directed
    use a directed graph (default true)
-o string
    path to the output file (optional)
```

Note: if the simulation is run using a non-directed graph, conflicting directions may be overwritten to allow for bi-directional roads.

```
Bar south=Baz
Foo north=Bar east=Baz
```

In this example we have two conflicting routes: `Bar south=Baz` and `Foo north=Bar`. The final graph would result in something like this:

```
Bar south=Foo north=Baz
Foo north=Bar east=Baz
Baz west=Foo
```
