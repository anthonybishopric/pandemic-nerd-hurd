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

_Features_
* Show panic levels in the UI
* Show player turns, which turns caused epidemics
* Track character traits and powerups
* Remind people on their turn what they can do (special abilities)

_Code Fixes_
* Keep pointers to actual epidemic and funded event cards in players / turns
* BUG: current turn on loading a save file is not the correct pointer to a player.