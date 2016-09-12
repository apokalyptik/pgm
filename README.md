PGM
---
A Pokemon Go spawn mapping and scanning program written in Go

Documentation
-------------

Installation
============
Go developers can simply `go get github.com/apokalyptik/pgm`

Usage
=====
You probably just want to read [the documentation for running the command](docs/pgm.md)

Under Development: TODO
-----------------------
* ~~Geocoding and Elevation queries~~
* ~~Jittering the starting point~~
* ~~create simple hex beehive coordinate pattern~~
* ~~log into Pokémon Go with arbitrary accounts~~
* ~~move the accounts around via a queue to scan the beehive~~
* ~~"see" what's arround the accounts~~
* Impliment a non-debug "feed" (which reads encounters from the accounts)
* Some sort of database (probably MySQL)
* Record spawn points *ESPECIALLY non standard types*
* Record spawn contents (encounters)
* Some sort of UI
* Spawnpoint Scanning
* Sitching from beehive to spawnpoint scanning once all locations have been checked properly
* Webhooks
* Account rotation

Error handling
==============
* Detect disconnects ?
* Reconnect
* Detect soft bans
* Banned account removal
