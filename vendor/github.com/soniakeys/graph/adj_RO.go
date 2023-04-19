// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

// adj_RO.go is code generated from adj_cg.go by directives in graph.go.
// Editing adj_cg.go is okay.
// DO NOT EDIT adj_RO.go.  The RO is for Read Only.

import (
	"math/rand"
	"time"
)

// ArcSize returns the number of arcs in g.
//
// Note that for an undirected graph without loops, the number of undirected
// edges -- the traditional meaning of graph size -- will be ArcSize()/2.
// On the other hand, if g is an undirected graph that has or may have loops,
// g.ArcSize()/2 is not a meaningful quantity.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g AdjacencyList) ArcSize() int {
	m := 0
	for _, to := range g {
		m += len(to)
	}
	return m
}

// BoundsOk validates that all arcs in g stay within the slice bounds of g.
//
// BoundsOk returns true when no arcs point outside the bounds of g.
// Otherwise it returns false and an example arc that points outside of g.
//
// Most methods of this package assume the BoundsOk condition and may
// panic when they encounter an arc pointing outside of the graph.  This
// function can be used to validate a graph when the BoundsOk condition
// is unknown.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g AdjacencyList) BoundsOk() (ok bool, fr NI, to NI) {
	for fr, to := range g {
		for _, to := range to {
			if to < 0 || to >= NI(len(g)) {
				return false, NI(fr), to
			}
		}
	}
	return true, -1, to
}

// BreadthFirst traverses a directed or undirected graph in breadth first order.
//
// Argument start is the start node for the traversal.  If r is nil, nodes are
// visited in deterministic order.  If a random number generator is supplied,
// nodes at each level are visited in random order.
//
// Argument f can be nil if you have no interest in the FromList path result.
// If FromList f is non-nil, the method populates f.Paths and sets f.MaxLen.
// It does not set f.Leaves.  For convenience argument f can be a zero value
// FromList.  If f.Paths is nil, the FromList is initialized first.  If f.Paths
// is non-nil however, the FromList is  used as is.  The method uses a value of
// PathEnd.Len == 0 to indentify unvisited nodes.  Existing non-zero values
// will limit the traversal.
//
// Traversal calls the visitor function v for each node starting with node
// start.  If v returns true, traversal continues.  If v returns false, the
// traversal terminates immediately.  PathEnd Len and From values are updated
// before calling the visitor function.
//
// On return f.Paths and f.MaxLen are set but not f.Leaves.
//
// Returned is the number of nodes visited and ok = true if the traversal
// ran to completion or ok = false if it was terminated by the visitor
// function returning false.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g AdjacencyList) BreadthFirst(start NI, r *rand.Rand, f *FromList, v OkNodeVisitor) (visited int, ok bool) {
	switch {
	case f == nil:
		e := NewFromList(len(g))
		f = &e
	case f.Paths == nil:
		*f = NewFromList(len(g))
	}
	rp := f.Paths
	// the frontier consists of nodes all at the same level
	frontier := []NI{start}
	level := 1
	// assign path when node is put on frontier,
	rp[start] = PathEnd{Len: level, From: -1}
	for {
		f.MaxLen = level
		level++
		var next []NI
		if r == nil {
			for _, n := range frontier {
				visited++
				if !v(n) { // visit nodes as they come off frontier
					return
				}
				for _, nb := range g[n] {
					if rp[nb].Len == 0 {
						next = append(next, nb)
						rp[nb] = PathEnd{From: n, Len: level}
					}
				}
			}
		} else { // take nodes off frontier at random
			for _, i := range r.Perm(len(frontier)) {
				n := frontier[i]
				// remainder of block same as above
				visited++
				if !v(n) {
					return
				}
				for _, nb := range g[n] {
					if rp[nb].Len == 0 {
						next = append(next, nb)
						rp[nb] = PathEnd{From: n, Len: level}
					}
				}
			}
		}
		if len(next) == 0 {
			break
		}
		frontier = next
	}
	return visited, true
}

// BreadthFirstPath finds a single path from start to end with a minimum
// number of nodes.
//
// Returned is the path as list of nodes.
// The result is nil if no path was found.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g AdjacencyList) BreadthFirstPath(start, end NI) []NI {
	var f FromList
	g.BreadthFirst(start, nil, &f, func(n NI) bool { return n != end })
	return f.PathTo(end, nil)
}

// Copy makes a deep copy of g.
// Copy also computes the arc size ma, the number of arcs.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g AdjacencyList) Copy() (c AdjacencyList, ma int) {
	c = make(AdjacencyList, len(g))
	for n, to := range g {
		c[n] = append([]NI{}, to...)
		ma += len(to)
	}
	return
}

// DepthFirst traverses a graph depth first.
//
// As it traverses it calls visitor function v for each node.  If v returns
// false at any point, the traversal is terminated immediately and DepthFirst
// returns false.  Otherwise DepthFirst returns true.
//
// DepthFirst uses argument bm is used as a bitmap to guide the traversal.
// For a complete traversal, bm should be 0 initially.  During the
// traversal, bits are set corresponding to each node visited.
// The bit is set before calling the visitor function.
//
// Argument bm can be nil if you have no need for it.
// In this case a bitmap is created internally for one-time use.
//
// Alternatively v can be nil.  In this case traversal still proceeds and
// updates the bitmap, which can be a useful result.
// DepthFirst always returns true in this case.
//
// It makes no sense for both bm and v to be nil.  In this case DepthFirst
// returns false immediately.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g AdjacencyList) DepthFirst(start NI, bm *Bits, v OkNodeVisitor) (ok bool) {
	if bm == nil {
		if v == nil {
			return false
		}
		bm = &Bits{}
	}
	var df func(n NI) bool
	df = func(n NI) bool {
		if bm.Bit(n) == 1 {
			return true
		}
		bm.SetBit(n, 1)
		if v != nil && !v(n) {
			return false
		}
		for _, nb := range g[n] {
			if !df(nb) {
				return false
			}
		}
		return true
	}
	return df(start)
}

// DepthFirstRandom traverses a graph depth first, but following arcs in
// random order among arcs from a single node.
//
// If Rand r is nil, the method creates a new source and generator for
// one-time use.
//
// Usage is otherwise like the DepthFirst method.  See DepthFirst.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g AdjacencyList) DepthFirstRandom(start NI, bm *Bits, v OkNodeVisitor, r *rand.Rand) (ok bool) {
	if bm == nil {
		if v == nil {
			return false
		}
		bm = &Bits{}
	}
	if r == nil {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	var df func(n NI) bool
	df = func(n NI) bool {
		if bm.Bit(n) == 1 {
			return true
		}
		bm.SetBit(n, 1)
		if v != nil && !v(n) {
			return false
		}
		to := g[n]
		for _, i := range r.Perm(len(to)) {
			if !df(to[i]) {
				return false
			}
		}
		return true
	}
	return df(start)
}

// HasArc returns true if g has any arc from node fr to node to.
//
// Also returned is the index within the slice of arcs from node fr.
// If no arc from fr to to is present, HasArc returns false, -1.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g AdjacencyList) HasArc(fr, to NI) (bool, int) {
	for x, h := range g[fr] {
		if h == to {
			return true, x
		}
	}
	return false, -1
}

// HasLoop identifies if a graph contains a loop, an arc that leads from a
// a node back to the same node.
//
// If the graph has a loop, the result is an example node that has a loop.
//
// If g contains a loop, the method returns true and an example of a node
// with a loop.  If there are no loops in g, the method returns false, -1.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g AdjacencyList) HasLoop() (bool, NI) {
	for fr, to := range g {
		for _, to := range to {
			if NI(fr) == to {
				return true, to
			}
		}
	}
	return false, -1
}

// HasParallelMap identifies if a graph contains parallel arcs, multiple arcs
// that lead from a node to the same node.
//
// If the graph has parallel arcs, the method returns true and
// results fr and to represent an example where there are parallel arcs
// from node fr to node to.
//
// If there are no parallel arcs, the method returns false, -1 -1.
//
// Multiple loops on a node count as parallel arcs.
//
// "Map" in the method name indicates that a Go map is used to detect parallel
// arcs.  Compared to method HasParallelSort, this gives better asymtotic
// performance for large dense graphs but may have increased overhead for
// small or sparse graphs.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g AdjacencyList) HasParallelMap() (has bool, fr, to NI) {
	for n, to := range g {
		if len(to) == 0 {
			continue
		}
		m := map[NI]struct{}{}
		for _, to := range to {
			if _, ok := m[to]; ok {
				return true, NI(n), to
			}
			m[to] = struct{}{}
		}
	}
	return false, -1, -1
}

// IsSimple checks for loops and parallel arcs.
//
// A graph is "simple" if it has no loops or parallel arcs.
//
// IsSimple returns true, -1 for simple graphs.  If a loop or parallel arc is
// found, simple returns false and a node that represents a counterexample
// to the graph being simple.
//
// See also separate methods HasLoop and HasParallel.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g AdjacencyList) IsSimple() (ok bool, n NI) {
	if lp, n := g.HasLoop(); lp {
		return false, n
	}
	if pa, n, _ := g.HasParallelSort(); pa {
		return false, n
	}
	return true, -1
}

// IsolatedNodes returns a bitmap of isolated nodes in receiver graph g.
//
// An isolated node is one with no arcs going to or from it.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g AdjacencyList) IsolatedNodes() (i Bits) {
	i.SetAll(len(g))
	for fr, to := range g {
		if len(to) > 0 {
			i.SetBit(NI(fr), 0)
			for _, to := range to {
				i.SetBit(to, 0)
			}
		}
	}
	return
}

/*
MaxmimalClique finds a maximal clique containing the node n.

Not sure this is good for anything.  It produces a single maximal clique
but there can be multiple maximal cliques containing a given node.
This algorithm just returns one of them, not even necessarily the
largest one.

func (g LabeledAdjacencyList) MaximalClique(n int) []int {
	c := []int{n}
	var m bitset.BitSet
	m.Set(uint(n))
	for fr, to := range g {
		if fr == n {
			continue
		}
		if len(to) < len(c) {
			continue
		}
		f := 0
		for _, to := range to {
			if m.Test(uint(to.To)) {
				f++
				if f == len(c) {
					c = append(c, to.To)
					m.Set(uint(to.To))
					break
				}
			}
		}
	}
	return c
}
*/
