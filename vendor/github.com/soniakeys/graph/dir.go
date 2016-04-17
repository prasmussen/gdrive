// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

// dir.go has methods specific to directed graphs, types Directed and
// LabeledDirected.
//
// Methods on Directed are first, with exported methods alphabetized.

import "errors"

// DAGMaxLenPath finds a maximum length path in a directed acyclic graph.
//
// Argument ordering must be a topological ordering of g.
func (g Directed) DAGMaxLenPath(ordering []NI) (path []NI) {
	// dynamic programming. visit nodes in reverse order. for each, compute
	// longest path as one plus longest of 'to' nodes.
	// Visits each arc once.  O(m).
	//
	// Similar code in label.go
	var n NI
	mlp := make([][]NI, len(g.AdjacencyList)) // index by node number
	for i := len(ordering) - 1; i >= 0; i-- {
		fr := ordering[i] // node number
		to := g.AdjacencyList[fr]
		if len(to) == 0 {
			continue
		}
		mt := to[0]
		for _, to := range to[1:] {
			if len(mlp[to]) > len(mlp[mt]) {
				mt = to
			}
		}
		p := append([]NI{mt}, mlp[mt]...)
		mlp[fr] = p
		if len(p) > len(path) {
			n = fr
			path = p
		}
	}
	return append([]NI{n}, path...)
}

// EulerianCycle finds an Eulerian cycle in a directed multigraph.
//
// * If g has no nodes, result is nil, nil.
//
// * If g is Eulerian, result is an Eulerian cycle with err = nil.
// The cycle result is a list of nodes, where the first and last
// nodes are the same.
//
// * Otherwise, result is nil, error
//
// Internally, EulerianCycle copies the entire graph g.
// See EulerianCycleD for a more space efficient version.
func (g Directed) EulerianCycle() ([]NI, error) {
	c, m := g.Copy()
	return c.EulerianCycleD(m)
}

// EulerianCycleD finds an Eulerian cycle in a directed multigraph.
//
// EulerianCycleD is destructive on its receiver g.  See EulerianCycle for
// a non-destructive version.
//
// Argument ma must be the correct arc size, or number of arcs in g.
//
// * If g has no nodes, result is nil, nil.
//
// * If g is Eulerian, result is an Eulerian cycle with err = nil.
// The cycle result is a list of nodes, where the first and last
// nodes are the same.
//
// * Otherwise, result is nil, error
func (g Directed) EulerianCycleD(ma int) ([]NI, error) {
	if len(g.AdjacencyList) == 0 {
		return nil, nil
	}
	e := newEulerian(g.AdjacencyList, ma)
	for e.s >= 0 {
		v := e.top() // v is node that starts cycle
		e.push()
		// if Eulerian, we'll always come back to starting node
		if e.top() != v {
			return nil, errors.New("not balanced")
		}
		e.keep()
	}
	if !e.uv.Zero() {
		return nil, errors.New("not strongly connected")
	}
	return e.p, nil
}

// EulerianPath finds an Eulerian path in a directed multigraph.
//
// * If g has no nodes, result is nil, nil.
//
// * If g has an Eulerian path, result is an Eulerian path with err = nil.
// The path result is a list of nodes, where the first node is start.
//
// * Otherwise, result is nil, error
//
// Internally, EulerianPath copies the entire graph g.
// See EulerianPathD for a more space efficient version.
func (g Directed) EulerianPath() ([]NI, error) {
	ind := g.InDegree()
	var start NI
	for n, to := range g.AdjacencyList {
		if len(to) > ind[n] {
			start = NI(n)
			break
		}
	}
	c, m := g.Copy()
	return c.EulerianPathD(m, start)
}

// EulerianPathD finds an Eulerian path in a directed multigraph.
//
// EulerianPathD is destructive on its receiver g.  See EulerianPath for
// a non-destructive version.
//
// Argument ma must be the correct arc size, or number of arcs in g.
// Argument start must be a valid start node for the path.
//
// * If g has no nodes, result is nil, nil.
//
// * If g has an Eulerian path, result is an Eulerian path with err = nil.
// The path result is a list of nodes, where the first node is start.
//
// * Otherwise, result is nil, error
func (g Directed) EulerianPathD(ma int, start NI) ([]NI, error) {
	if len(g.AdjacencyList) == 0 {
		return nil, nil
	}
	e := newEulerian(g.AdjacencyList, ma)
	e.p[0] = start
	// unlike EulerianCycle, the first path doesn't have be a cycle.
	e.push()
	e.keep()
	for e.s >= 0 {
		start = e.top()
		e.push()
		// paths after the first must be cycles though
		// (as long as there are nodes on the stack)
		if e.top() != start {
			return nil, errors.New("no Eulerian path")
		}
		e.keep()
	}
	if !e.uv.Zero() {
		return nil, errors.New("no Eulerian path")
	}
	return e.p, nil
}

// starting at the node on the top of the stack, follow arcs until stuck.
// mark nodes visited, push nodes on stack, remove arcs from g.
func (e *eulerian) push() {
	for u := e.top(); ; {
		e.uv.SetBit(u, 0) // reset unvisited bit
		arcs := e.g[u]
		if len(arcs) == 0 {
			return // stuck
		}
		w := arcs[0] // follow first arc
		e.s++        // push followed node on stack
		e.p[e.s] = w
		e.g[u] = arcs[1:] // consume arc
		u = w
	}
}

// like push, but for for undirected graphs.
func (e *eulerian) pushUndir() {
	for u := e.top(); ; {
		e.uv.SetBit(u, 0)
		arcs := e.g[u]
		if len(arcs) == 0 {
			return
		}
		w := arcs[0]
		e.s++
		e.p[e.s] = w
		e.g[u] = arcs[1:] // consume arc
		// here is the only difference, consume reciprocal arc as well:
		a2 := e.g[w]
		for x, rx := range a2 {
			if rx == u { // here it is
				last := len(a2) - 1
				a2[x] = a2[last]   // someone else gets the seat
				e.g[w] = a2[:last] // and it's gone.
				break
			}
		}
		u = w
	}
}

// starting with the node on top of the stack, move nodes with no arcs.
func (e *eulerian) keep() {
	for e.s >= 0 {
		n := e.top()
		if len(e.g[n]) > 0 {
			break
		}
		e.p[e.m] = n
		e.s--
		e.m--
	}
}

type eulerian struct {
	g  AdjacencyList // working copy of graph, it gets consumed
	m  int           // number of arcs in g, updated as g is consumed
	uv Bits          // unvisited
	// low end of p is stack of unfinished nodes
	// high end is finished path
	p []NI // stack + path
	s int  // stack pointer
}

func (e *eulerian) top() NI {
	return e.p[e.s]
}

func newEulerian(g AdjacencyList, m int) *eulerian {
	e := &eulerian{
		g: g,
		m: m,
		p: make([]NI, m+1),
	}
	e.uv.SetAll(len(g))
	return e
}

// MaximalNonBranchingPaths finds all paths in a directed graph that are
// "maximal" and "non-branching".
//
// A non-branching path is one where path nodes other than the first and last
// have exactly one arc leading to the node and one arc leading from the node,
// thus there is no possibility to branch away to a different path.
//
// A maximal non-branching path cannot be extended to a longer non-branching
// path by including another node at either end.
//
// In the case of a cyclic non-branching path, the first and last elements
// of the path will be the same node, indicating an isolated cycle.
//
// The method calls the emit argument for each path or isolated cycle in g,
// as long as emit returns true.  If emit returns false,
// MaximalNonBranchingPaths returns immediately.
func (g Directed) MaximalNonBranchingPaths(emit func([]NI) bool) {
	ind := g.InDegree()
	var uv Bits
	uv.SetAll(len(g.AdjacencyList))
	for v, vTo := range g.AdjacencyList {
		if !(ind[v] == 1 && len(vTo) == 1) {
			for _, w := range vTo {
				n := []NI{NI(v), w}
				uv.SetBit(NI(v), 0)
				uv.SetBit(w, 0)
				wTo := g.AdjacencyList[w]
				for ind[w] == 1 && len(wTo) == 1 {
					u := wTo[0]
					n = append(n, u)
					uv.SetBit(u, 0)
					w = u
					wTo = g.AdjacencyList[w]
				}
				if !emit(n) { // n is a path
					return
				}
			}
		}
	}
	// use uv.From rather than uv.Iterate.
	// Iterate doesn't work here because we're modifying uv
	for b := uv.From(0); b >= 0; b = uv.From(b + 1) {
		v := NI(b)
		n := []NI{v}
		for w := v; ; {
			w = g.AdjacencyList[w][0]
			uv.SetBit(w, 0)
			n = append(n, w)
			if w == v {
				break
			}
		}
		if !emit(n) { // n is an isolated cycle
			return
		}
	}
}

// Undirected returns copy of g augmented as needed to make it undirected.
func (g Directed) Undirected() Undirected {
	c, _ := g.AdjacencyList.Copy()                  // start with a copy
	rw := make(AdjacencyList, len(g.AdjacencyList)) // "reciprocals wanted"
	for fr, to := range g.AdjacencyList {
	arc: // for each arc in g
		for _, to := range to {
			if to == NI(fr) {
				continue // loop
			}
			// search wanted arcs
			wf := rw[fr]
			for i, w := range wf {
				if w == to { // found, remove
					last := len(wf) - 1
					wf[i] = wf[last]
					rw[fr] = wf[:last]
					continue arc
				}
			}
			// arc not found, add to reciprocal to wanted list
			rw[to] = append(rw[to], NI(fr))
		}
	}
	// add missing reciprocals
	for fr, to := range rw {
		c[fr] = append(c[fr], to...)
	}
	return Undirected{c}
}

// StronglyConnectedComponents identifies strongly connected components
// in a directed graph.
//
// Algorithm by David J. Pearce, from "An Improved Algorithm for Finding the
// Strongly Connected Components of a Directed Graph".  It is algorithm 3,
// PEA_FIND_SCC2 in
// http://homepages.mcs.vuw.ac.nz/~djp/files/P05.pdf, accessed 22 Feb 2015.
//
// Returned is a list of components, each component is a list of nodes.
/*
func (g Directed) StronglyConnectedComponents() []int {
	rindex := make([]int, len(g))
	S := []int{}
	index := 1
	c := len(g) - 1
	visit := func(v int) {
		root := true
		rindex[v] = index
		index++
		for _, w := range g[v] {
			if rindex[w] == 0 {
				visit(w)
			}
			if rindex[w] < rindex[v] {
				rindex[v] = rindex[w]
				root = false
			}
		}
		if root {
			index--
			for top := len(S) - 1; top >= 0 && rindex[v] <= rindex[top]; top-- {
				w = rindex[top]
				S = S[:top]
				rindex[w] = c
				index--
			}
			rindex[v] = c
			c--
		} else {
			S = append(S, v)
		}
	}
	for v := range g {
		if rindex[v] == 0 {
			visit(v)
		}
	}
	return rindex
}
*/

// Transpose constructs a new adjacency list with all arcs reversed.
//
// For every arc from->to of g, the result will have an arc to->from.
// Transpose also counts arcs as it traverses and returns ma the number of arcs
// in g (equal to the number of arcs in the result.)
func (g Directed) Transpose() (t Directed, ma int) {
	ta := make(AdjacencyList, len(g.AdjacencyList))
	for n, nbs := range g.AdjacencyList {
		for _, nb := range nbs {
			ta[nb] = append(ta[nb], NI(n))
			ma++
		}
	}
	return Directed{ta}, ma
}

// DAGMaxLenPath finds a maximum length path in a directed acyclic graph.
//
// Length here means number of nodes or arcs, not a sum of arc weights.
//
// Argument ordering must be a topological ordering of g.
//
// Returned is a node beginning a maximum length path, and a path of arcs
// starting from that node.
func (g LabeledDirected) DAGMaxLenPath(ordering []NI) (n NI, path []Half) {
	// dynamic programming. visit nodes in reverse order. for each, compute
	// longest path as one plus longest of 'to' nodes.
	// Visits each arc once.  Time complexity O(m).
	//
	// Similar code in dir.go.
	mlp := make([][]Half, len(g.LabeledAdjacencyList)) // index by node number
	for i := len(ordering) - 1; i >= 0; i-- {
		fr := ordering[i] // node number
		to := g.LabeledAdjacencyList[fr]
		if len(to) == 0 {
			continue
		}
		mt := to[0]
		for _, to := range to[1:] {
			if len(mlp[to.To]) > len(mlp[mt.To]) {
				mt = to
			}
		}
		p := append([]Half{mt}, mlp[mt.To]...)
		mlp[fr] = p
		if len(p) > len(path) {
			n = fr
			path = p
		}
	}
	return
}

// FromListLabels transposes a labeled graph into a FromList and associated
// list of labels.
//
// Receiver g should be connected as a tree or forest.  Specifically no node
// can have multiple incoming arcs.  If any node n in g has multiple incoming
// arcs, the method returns (nil, nil, n) where n is a node with multiple
// incoming arcs.
//
// Otherwise (normally) the method populates the From members in a
// FromList.Path, populates a slice of labels, and returns the FromList,
// labels, and -1.
//
// Other members of the FromList are left as zero values.
// Use FromList.RecalcLen and FromList.RecalcLeaves as needed.
func (g LabeledDirected) FromListLabels() (*FromList, []LI, NI) {
	labels := make([]LI, len(g.LabeledAdjacencyList))
	paths := make([]PathEnd, len(g.LabeledAdjacencyList))
	for i := range paths {
		paths[i].From = -1
	}
	for fr, to := range g.LabeledAdjacencyList {
		for _, to := range to {
			if paths[to.To].From >= 0 {
				return nil, nil, to.To
			}
			paths[to.To].From = NI(fr)
			labels[to.To] = to.Label
		}
	}
	return &FromList{Paths: paths}, labels, -1
}

// Transpose constructs a new adjacency list that is the transpose of g.
//
// For every arc from->to of g, the result will have an arc to->from.
// Transpose also counts arcs as it traverses and returns ma the number of
// arcs in g (equal to the number of arcs in the result.)
func (g LabeledDirected) Transpose() (t LabeledDirected, ma int) {
	ta := make(LabeledAdjacencyList, len(g.LabeledAdjacencyList))
	for n, nbs := range g.LabeledAdjacencyList {
		for _, nb := range nbs {
			ta[nb.To] = append(ta[nb.To], Half{To: NI(n), Label: nb.Label})
			ma++
		}
	}
	return LabeledDirected{ta}, ma
}

// Undirected returns a new undirected graph derived from g, augmented as
// needed to make it undirected, with reciprocal arcs having matching labels.
func (g LabeledDirected) Undirected() LabeledUndirected {
	c, _ := g.LabeledAdjacencyList.Copy() // start with a copy
	// "reciprocals wanted"
	rw := make(LabeledAdjacencyList, len(g.LabeledAdjacencyList))
	for fr, to := range g.LabeledAdjacencyList {
	arc: // for each arc in g
		for _, to := range to {
			if to.To == NI(fr) {
				continue // arc is a loop
			}
			// search wanted arcs
			wf := rw[fr]
			for i, w := range wf {
				if w == to { // found, remove
					last := len(wf) - 1
					wf[i] = wf[last]
					rw[fr] = wf[:last]
					continue arc
				}
			}
			// arc not found, add to reciprocal to wanted list
			rw[to.To] = append(rw[to.To], Half{To: NI(fr), Label: to.Label})
		}
	}
	// add missing reciprocals
	for fr, to := range rw {
		c[fr] = append(c[fr], to...)
	}
	return LabeledUndirected{c}
}

// Unlabeled constructs the unlabeled directed graph corresponding to g.
func (g LabeledDirected) Unlabeled() Directed {
	return Directed{g.LabeledAdjacencyList.Unlabeled()}
}

// UnlabeledTranspose constructs a new adjacency list that is the unlabeled
// transpose of g.
//
// For every arc from->to of g, the result will have an arc to->from.
// Transpose also counts arcs as it traverses and returns ma, the number of
// arcs in g (equal to the number of arcs in the result.)
//
// It is equivalent to g.Unlabeled().Transpose() but constructs the result
// directly.
func (g LabeledDirected) UnlabeledTranspose() (t Directed, ma int) {
	ta := make(AdjacencyList, len(g.LabeledAdjacencyList))
	for n, nbs := range g.LabeledAdjacencyList {
		for _, nb := range nbs {
			ta[nb.To] = append(ta[nb.To], NI(n))
			ma++
		}
	}
	return Directed{ta}, ma
}
