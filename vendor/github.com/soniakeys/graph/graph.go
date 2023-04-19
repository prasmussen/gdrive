// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

// graph.go contains type definitions for all graph types and components.
// Also, go generate directives for source transformations.
//
// For readability, the types are defined in a dependency order:
//
//  NI
//  NodeList
//  AdjacencyList
//  Directed
//  Undirected
//  LI
//  Half
//  LabeledAdjacencyList
//  LabeledDirected
//  LabeledUndirected
//  Edge
//  LabeledEdge
//  WeightFunc
//  WeightedEdgeList

//go:generate cp adj_cg.go adj_RO.go
//go:generate gofmt -r "LabeledAdjacencyList -> AdjacencyList" -w adj_RO.go
//go:generate gofmt -r "n.To -> n" -w adj_RO.go
//go:generate gofmt -r "Half -> NI" -w adj_RO.go

//go:generate cp dir_cg.go dir_RO.go
//go:generate gofmt -r "LabeledDirected -> Directed" -w dir_RO.go
//go:generate gofmt -r "LabeledAdjacencyList -> AdjacencyList" -w dir_RO.go
//go:generate gofmt -r "n.To -> n" -w dir_RO.go
//go:generate gofmt -r "Half -> NI" -w dir_RO.go

//go:generate cp undir_cg.go undir_RO.go
//go:generate gofmt -r "LabeledUndirected -> Undirected" -w undir_RO.go
//go:generate gofmt -r "LabeledAdjacencyList -> AdjacencyList" -w undir_RO.go
//go:generate gofmt -r "n.To -> n" -w undir_RO.go
//go:generate gofmt -r "Half -> NI" -w undir_RO.go

// NI is a "node int"
//
// It is a node number or node ID.  NIs are used extensively as slice indexes.
// NIs typically account for a significant fraction of the memory footprint of
// a graph.
type NI int32

// NodeList satisfies sort.Interface.
type NodeList []NI

func (l NodeList) Len() int           { return len(l) }
func (l NodeList) Less(i, j int) bool { return l[i] < l[j] }
func (l NodeList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

// An AdjacencyList represents a graph as a list of neighbors for each node.
// The "node ID" of a node is simply it's slice index in the AdjacencyList.
// For an AdjacencyList g, g[n] represents arcs going from node n to nodes
// g[n].
//
// Adjacency lists are inherently directed but can be used to represent
// directed or undirected graphs.  See types Directed and Undirected.
type AdjacencyList [][]NI

// Directed represents a directed graph.
//
// Directed methods generally rely on the graph being directed, specifically
// that arcs do not have reciprocals.
type Directed struct {
	AdjacencyList // embedded to include AdjacencyList methods
}

// Undirected represents an undirected graph.
//
// In an undirected graph, for each arc between distinct nodes there is also
// a reciprocal arc, an arc in the opposite direction.  Loops do not have
// reciprocals.
//
// Undirected methods generally rely on the graph being undirected,
// specifically that every arc between distinct nodes has a reciprocal.
type Undirected struct {
	AdjacencyList // embedded to include AdjacencyList methods
}

// LI is a label integer, used for associating labels with arcs.
type LI int32

// Half is a half arc, representing a labeled arc and the "neighbor" node
// that the arc leads to.
//
// Halfs can be composed to form a labeled adjacency list.
type Half struct {
	To    NI // node ID, usable as a slice index
	Label LI // half-arc ID for application data, often a weight
}

// A LabeledAdjacencyList represents a graph as a list of neighbors for each
// node, connected by labeled arcs.
//
// Arc labels are not necessarily unique arc IDs.  Different arcs can have
// the same label.
//
// Arc labels are commonly used to assocate a weight with an arc.  Arc labels
// are general purpose however and can be used to associate arbitrary
// information with an arc.
//
// Methods implementing weighted graph algorithms will commonly take a
// weight function that turns a label int into a float64 weight.
//
// If only a small amount of information -- such as an integer weight or
// a single printable character -- needs to be associated, it can sometimes
// be possible to encode the information directly into the label int.  For
// more generality, some lookup scheme will be needed.
//
// In an undirected labeled graph, reciprocal arcs must have identical labels.
// Note this does not preclude parallel arcs with different labels.
type LabeledAdjacencyList [][]Half

// LabeledDirected represents a directed labeled graph.
//
// This is the labeled version of Directed.  See types LabeledAdjacencyList
// and Directed.
type LabeledDirected struct {
	LabeledAdjacencyList // embedded to include LabeledAdjacencyList methods
}

// LabeledUndirected represents an undirected labeled graph.
//
// This is the labeled version of Undirected.  See types LabeledAdjacencyList
// and Undirected.
type LabeledUndirected struct {
	LabeledAdjacencyList // embedded to include LabeledAdjacencyList methods
}

// Edge is an undirected edge between nodes N1 and N2.
type Edge struct{ N1, N2 NI }

// LabeledEdge is an undirected edge with an associated label.
type LabeledEdge struct {
	Edge
	LI
}

// WeightFunc returns a weight for a given label.
//
// WeightFunc is a parameter type for various search functions.  The intent
// is to return a weight corresponding to an arc label.  The name "weight"
// is an abstract term.  An arc "weight" will typically have some application
// specific meaning other than physical weight.
type WeightFunc func(label LI) (weight float64)

// WeightedEdgeList is a graph representation.
//
// It is a labeled edge list, with an associated weight function to return
// a weight given an edge label.
//
// Also associated is the order, or number of nodes of the graph.
// All nodes occurring in the edge list must be strictly less than Order.
//
// WeigtedEdgeList sorts by weight, obtained by calling the weight function.
// If weight computation is expensive, consider supplying a cached or
// memoized version.
type WeightedEdgeList struct {
	Order int
	WeightFunc
	Edges []LabeledEdge
}

// Len implements sort.Interface.
func (l WeightedEdgeList) Len() int { return len(l.Edges) }

// Less implements sort.Interface.
func (l WeightedEdgeList) Less(i, j int) bool {
	return l.WeightFunc(l.Edges[i].LI) < l.WeightFunc(l.Edges[j].LI)
}

// Swap implements sort.Interface.
func (l WeightedEdgeList) Swap(i, j int) {
	l.Edges[i], l.Edges[j] = l.Edges[j], l.Edges[i]
}
