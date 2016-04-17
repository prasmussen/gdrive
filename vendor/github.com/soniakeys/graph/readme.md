#Graph

A graph library with goals of speed and simplicity, Graph implements
graph algorithms on graphs of zero-based integer node IDs.

[![GoDoc](https://godoc.org/github.com/soniakeys/graph?status.svg)](https://godoc.org/github.com/soniakeys/graph) [![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/soniakeys/graph) [![GoSearch](http://go-search.org/badge?id=github.com%2Fsoniakeys%2Fgraph)](http://go-search.org/view?id=github.com%2Fsoniakeys%2Fgraph)[![Build Status](https://travis-ci.org/soniakeys/graph.svg?branch=master)](https://travis-ci.org/soniakeys/graph)

Status, 4 Apr 2016:  The repo has benefitted recently from being included
in another package.  In response to users of that package, this repo now
builds for 32 bit Windows and ARM, and for Go versions back to 1.2.1.
Thank you all who have filed issues.

###Non-source files of interest

The directory [tutorials](tutorials) is a work in progress - there are only
a couple of tutorials there yet - but the concept is to provide some topical
walk-throughs to supplement godoc.  The source-based godoc documentation
remains the primary documentation.

* [Dijkstra's algorithm](tutorials/dijkstra.md)
* [AdjacencyList types](tutorials/adjacencylist.md)

The directory [bench](bench) is another work in progress.  The concept is
to present some plots showing benchmark performance approaching some
theoretical asymptote.

[hacking.md](hacking.md) has some information about how the library is
developed, built, and tested.  It might be of interest if for example you
plan to fork or contribute to the the repository.

###Test coverage
8 Apr 2016
```
graph          95.3%
graph/df       20.7%
graph/dot      77.5%
graph/treevis  79.4%
```
