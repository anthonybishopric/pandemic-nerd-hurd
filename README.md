# Pandemic: Legacy stats collector

Simple code that emits JSON files containing information about playthroughs. Can do some basic statistical analysis.

## Running

Clone the repo, make sure you have Go 1.6, then:

```
$ go get ./...
$ go build .
$ ./pandemic-nerd-hurd
```

## TODO

* Show panic levels in the UI
* Show player turns, which turns caused epidemics
* Separate turns into city-deck striations
* Incorporate epidemic probability in city outbreak probability
* Toggle outbreak probability view in main UI
* Track character traits and powerups
* Remind people on their turn what they can do (special abilities)
