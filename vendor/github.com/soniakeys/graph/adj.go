// Copyright 2014 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package graph

// adj.go contains methods on AdjacencyList and LabeledAdjacencyList.
//
// AdjacencyList methods are placed first and are alphabetized.
// LabeledAdjacencyList methods follow, also alphabetized.
// Only exported methods need be alphabetized; non-exported methods can
// be left near their use.

import (
	"math"
	"sort"
)

// HasParallelSort identifies if a graph contains parallel arcs, multiple arcs
// that lead from a node to the same node.
//
// If the graph has parallel arcs, the results fr and to represent an example
// where there are parallel arcs from node fr to node to.
//
// If there are no parallel arcs, the method returns false -1 -1.
//
// Multiple loops on a node count as parallel arcs.
//
// "Sort" in the method name indicates that sorting is used to detect parallel
// arcs.  Compared to method HasParallelMap, this may give better performance
// for small or sparse graphs but will have asymtotically worse performance for
// large dense graphs.
func (g AdjacencyList) HasParallelSort() (has bool, fr, to NI) {
	var t NodeList
	for n, to := range g {
		if len(to) == 0 {
			continue
		}
		// different code in the labeled version, so no code gen.
		t = append(t[:0], to...)
		sort.Sort(t)
		t0 := t[0]
		for _, to := range t[1:] {
			if to == t0 {
				return true, NI(n), t0
			}
			t0 = to
		}
	}
	return false, -1, -1
}

// IsUndirected returns true if g represents an undirected graph.
//
// Returns true when all non-loop arcs are paired in reciprocal pairs.
// Otherwise returns false and an example unpaired arc.
func (g AdjacencyList) IsUndirected() (u bool, from, to NI) {
	// similar code in dot/writeUndirected
	unpaired := make(AdjacencyList, len(g))
	for fr, to := range g {
	arc: // for each arc in g
		for _, to := range to {
			if to == NI(fr) {
				continue // loop
			}
			// search unpaired arcs
			ut := unpaired[to]
			for i, u := range ut {
				if u == NI(fr) { // found reciprocal
					last := len(ut) - 1
					ut[i] = ut[last]
					unpaired[to] = ut[:last]
					continue arc
				}
			}
			// reciprocal not found
			unpaired[fr] = append(unpaired[fr], to)
		}
	}
	for fr, to := range unpaired {
		if len(to) > 0 {
			return false, NI(fr), to[0]
		}
	}
	return true, -1, -1
}

// Edgelist constructs the edge list rerpresentation of a graph.
//
// An edge is returned for each arc of the graph.  For undirected graphs
// this includes reciprocal edges.
//
// See also WeightedEdgeList method.
func (g LabeledAdjacencyList) EdgeList() (el []LabeledEdge) {
	for fr, to := range g {
		for _, to := range to {
			el = append(el, LabeledEdge{Edge{NI(fr), to.To}, to.Label})
		}
	}
	return
}

// FloydWarshall finds all pairs shortest distances for a simple weighted
// graph without negative cycles.
//
// In result array d, d[i][j] will be the shortest distance from node i
// to node j.  Any diagonal element < 0 indicates a negative cycle exists.
//
// If g is an undirected graph with no negative edge weights, the result
// array will be a distance matrix, for example as used by package
// github.com/soniakeys/cluster.
func (g LabeledAdjacencyList) FloydWarshall(w WeightFunc) (d [][]float64) {
	d = newFWd(len(g))
	for fr, to := range g {
		for _, to := range to {
			d[fr][to.To] = w(to.Label)
		}
	}
	solveFW(d)
	return
}

// little helper function, makes a blank matrix for FloydWarshall.
func newFWd(n int) [][]float64 {
	d := make([][]float64, n)
	for i := range d {
		di := make([]float64, n)
		for j := range di {
			if j != i {
				di[j] = math.Inf(1)
			}
		}
		d[i] = di
	}
	return d
}

// Floyd Warshall solver, once the matrix d is initialized by arc weights.
func solveFW(d [][]float64) {
	for k, dk := range d {
		for _, di := range d {
			dik := di[k]
			for j := range d {
				if d2 := dik + dk[j]; d2 < di[j] {
					di[j] = d2
				}
			}
		}
	}
}

// HasArcLabel returns true if g has any arc from node fr to node to
// with label l.
//
// Also returned is the index within the slice of arcs from node fr.
// If no arc from fr to to is present, HasArcLabel returns false, -1.
func (g LabeledAdjacencyList) HasArcLabel(fr, to NI, l LI) (bool, int) {
	t := Half{to, l}
	for x, h := range g[fr] {
		if h == t {
			return true, x
		}
	}
	return false, -1
}

// HasParallelSort identifies if a graph contains parallel arcs, multiple arcs
// that lead from a node to the same node.
//
// If the graph has parallel arcs, the results fr and to represent an example
// where there are parallel arcs from node fr to node to.
//
// If there are no parallel arcs, the method returns -1 -1.
//
// Multiple loops on a node count as parallel arcs.
//
// "Sort" in the method name indicates that sorting is used to detect parallel
// arcs.  Compared to method HasParallelMap, this may give better performance
// for small or sparse graphs but will have asymtotically worse performance for
// large dense graphs.
func (g LabeledAdjacencyList) HasParallelSort() (has bool, fr, to NI) {
	var t NodeList
	for n, to := range g {
		if len(to) == 0 {
			continue
		}
		// slightly different code needed here compared to AdjacencyList
		t = t[:0]
		for _, to := range to {
			t = append(t, to.To)
		}
		sort.Sort(t)
		t0 := t[0]
		for _, to := range t[1:] {
			if to == t0 {
				return true, NI(n), t0
			}
			t0 = to
		}
	}
	return false, -1, -1
}

// IsUndirected returns true if g represents an undirected graph.
//
// Returns true when all non-loop arcs are paired in reciprocal pairs with
// matching labels.  Otherwise returns false and an example unpaired arc.
//
// Note the requirement that reciprocal pairs have matching labels is
// an additional test not present in the otherwise equivalent unlabeled version
// of IsUndirected.
func (g LabeledAdjacencyList) IsUndirected() (u bool, from NI, to Half) {
	unpaired := make(LabeledAdjacencyList, len(g))
	for fr, to := range g {
	arc: // for each arc in g
		for _, to := range to {
			if to.To == NI(fr) {
				continue // loop
			}
			// search unpaired arcs
			ut := unpaired[to.To]
			for i, u := range ut {
				if u.To == NI(fr) && u.Label == to.Label { // found reciprocal
					last := len(ut) - 1
					ut[i] = ut[last]
					unpaired[to.To] = ut[:last]
					continue arc
				}
			}
			// reciprocal not found
			unpaired[fr] = append(unpaired[fr], to)
		}
	}
	for fr, to := range unpaired {
		if len(to) > 0 {
			return false, NI(fr), to[0]
		}
	}
	return true, -1, to
}

// NegativeArc returns true if the receiver graph contains a negative arc.
func (g LabeledAdjacencyList) NegativeArc(w WeightFunc) bool {
	for _, nbs := range g {
		for _, nb := range nbs {
			if w(nb.Label) < 0 {
				return true
			}
		}
	}
	return false
}

// Unlabeled constructs the unlabeled graph corresponding to g.
func (g LabeledAdjacencyList) Unlabeled() AdjacencyList {
	a := make(AdjacencyList, len(g))
	for n, nbs := range g {
		to := make([]NI, len(nbs))
		for i, nb := range nbs {
			to[i] = nb.To
		}
		a[n] = to
	}
	return a
}

// WeightedEdgeList constructs a WeightedEdgeList object from a
// LabeledAdjacencyList.
//
// Internally it calls g.EdgeList() to obtain the Edges member.
// See LabeledAdjacencyList.EdgeList().
func (g LabeledAdjacencyList) WeightedEdgeList(w WeightFunc) *WeightedEdgeList {
	return &WeightedEdgeList{
		Order:      len(g),
		WeightFunc: w,
		Edges:      g.EdgeList(),
	}
}

// WeightedInDegree computes the weighted in-degree of each node in g
// for a given weight function w.
//
// The weighted in-degree of a node is the sum of weights of arcs going to
// the node.
//
// A weighted degree of a node is often termed the "strength" of a node.
//
// See note for undirected graphs at LabeledAdjacencyList.WeightedOutDegree.
func (g LabeledAdjacencyList) WeightedInDegree(w WeightFunc) []float64 {
	ind := make([]float64, len(g))
	for _, to := range g {
		for _, to := range to {
			ind[to.To] += w(to.Label)
		}
	}
	return ind
}

// WeightedOutDegree computes the weighted out-degree of the specified node
// for a given weight function w.
//
// The weighted out-degree of a node is the sum of weights of arcs going from
// the node.
//
// A weighted degree of a node is often termed the "strength" of a node.
//
// Note for undirected graphs, the WeightedOutDegree result for a node will
// equal the WeightedInDegree for the node.  You can use WeightedInDegree if
// you have need for the weighted degrees of all nodes or use WeightedOutDegree
// to compute the weighted degrees of individual nodes.  In either case loops
// are counted just once, unlike the (unweighted) UndirectedDegree methods.
func (g LabeledAdjacencyList) WeightedOutDegree(n NI, w WeightFunc) (d float64) {
	for _, to := range g[n] {
		d += w(to.Label)
	}
	return
}

// More about loops and strength:  I didn't see consensus on this especially
// in the case of undirected graphs.  Some sources said to add in-degree and
// out-degree, which would seemingly double both loops and non-loops.
// Some said to double loops.  Some said sum the edge weights and had no
// comment on loops.  R of course makes everything an option.  The meaning
// of "strength" where loops exist is unclear.  So while I could write an
// UndirectedWeighted degree function that doubles loops but not edges,
// I'm going to just leave this for now.
