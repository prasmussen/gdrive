// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

// undir.go has methods specific to undirected graphs, Undirected and
// LabeledUndirected.

import "errors"

// AddEdge adds an edge to a graph.
//
// It can be useful for constructing undirected graphs.
//
// When n1 and n2 are distinct, it adds the arc n1->n2 and the reciprocal
// n2->n1.  When n1 and n2 are the same, it adds a single arc loop.
//
// The pointer receiver allows the method to expand the graph as needed
// to include the values n1 and n2.  If n1 or n2 happen to be greater than
// len(*p) the method does not panic, but simply expands the graph.
func (p *Undirected) AddEdge(n1, n2 NI) {
	// Similar code in LabeledAdjacencyList.AddEdge.

	// determine max of the two end points
	max := n1
	if n2 > max {
		max = n2
	}
	// expand graph if needed, to include both
	g := p.AdjacencyList
	if int(max) >= len(g) {
		p.AdjacencyList = make(AdjacencyList, max+1)
		copy(p.AdjacencyList, g)
		g = p.AdjacencyList
	}
	// create one half-arc,
	g[n1] = append(g[n1], n2)
	// and except for loops, create the reciprocal
	if n1 != n2 {
		g[n2] = append(g[n2], n1)
	}
}

// EulerianCycleD for undirected graphs is a bit of an experiment.
//
// It is about the same as the directed version, but modified for an undirected
// multigraph.
//
// Parameter m in this case must be the size of the undirected graph -- the
// number of edges.  Use Undirected.Size if the size is unknown.
//
// It works, but contains an extra loop that I think spoils the time
// complexity.  Probably still pretty fast in practice, but a different
// graph representation might be better.
func (g Undirected) EulerianCycleD(m int) ([]NI, error) {
	if len(g.AdjacencyList) == 0 {
		return nil, nil
	}
	e := newEulerian(g.AdjacencyList, m)
	for e.s >= 0 {
		v := e.top()
		e.pushUndir() // call modified method
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

// TarjanBiconnectedComponents decomposes a graph into maximal biconnected
// components, components for which if any node were removed the component
// would remain connected.
//
// The receiver g must be a simple graph.  The method calls the emit argument
// for each component identified, as long as emit returns true.  If emit
// returns false, TarjanBiconnectedComponents returns immediately.
//
// See also the eqivalent labeled TarjanBiconnectedComponents.
func (g Undirected) TarjanBiconnectedComponents(emit func([]Edge) bool) {
	// Implemented closely to pseudocode in "Depth-first search and linear
	// graph algorithms", Robert Tarjan, SIAM J. Comput. Vol. 1, No. 2,
	// June 1972.
	//
	// Note Tarjan's "adjacency structure" is graph.AdjacencyList,
	// His "adjacency list" is an element of a graph.AdjacencyList, also
	// termed a "to-list", "neighbor list", or "child list."
	number := make([]int, len(g.AdjacencyList))
	lowpt := make([]int, len(g.AdjacencyList))
	var stack []Edge
	var i int
	var biconnect func(NI, NI) bool
	biconnect = func(v, u NI) bool {
		i++
		number[v] = i
		lowpt[v] = i
		for _, w := range g.AdjacencyList[v] {
			if number[w] == 0 {
				stack = append(stack, Edge{v, w})
				if !biconnect(w, v) {
					return false
				}
				if lowpt[w] < lowpt[v] {
					lowpt[v] = lowpt[w]
				}
				if lowpt[w] >= number[v] {
					var bcc []Edge
					top := len(stack) - 1
					for number[stack[top].N1] >= number[w] {
						bcc = append(bcc, stack[top])
						stack = stack[:top]
						top--
					}
					bcc = append(bcc, stack[top])
					stack = stack[:top]
					top--
					if !emit(bcc) {
						return false
					}
				}
			} else if number[w] < number[v] && w != u {
				stack = append(stack, Edge{v, w})
				if number[w] < lowpt[v] {
					lowpt[v] = number[w]
				}
			}
		}
		return true
	}
	for w := range g.AdjacencyList {
		if number[w] == 0 && !biconnect(NI(w), 0) {
			return
		}
	}
}

/* half-baked.  Read the 72 paper.  Maybe revisit at some point.
type BiconnectedComponents struct {
	Graph  AdjacencyList
	Start  int
	Cuts   big.Int // bitmap of node cuts
	From   []int   // from-tree
	Leaves []int   // leaves of from-tree
}

func NewBiconnectedComponents(g Undirected) *BiconnectedComponents {
	return &BiconnectedComponents{
		Graph: g,
		From:  make([]int, len(g)),
	}
}

func (b *BiconnectedComponents) Find(start int) {
	g := b.Graph
	depth := make([]int, len(g))
	low := make([]int, len(g))
	// reset from any previous run
	b.Cuts.SetInt64(0)
	bf := b.From
	for n := range bf {
		bf[n] = -1
	}
	b.Leaves = b.Leaves[:0]
	d := 1 // depth. d > 0 means visited
	depth[start] = d
	low[start] = d
	d++
	var df func(int, int)
	df = func(from, n int) {
		bf[n] = from
		depth[n] = d
		dn := d
		l := d
		d++
		cut := false
		leaf := true
		for _, nb := range g[n] {
			if depth[nb] == 0 {
				leaf = false
				df(n, nb)
				if low[nb] < l {
					l = low[nb]
				}
				if low[nb] >= dn {
					cut = true
				}
			} else if nb != from && depth[nb] < l {
				l = depth[nb]
			}
		}
		low[n] = l
		if cut {
			b.Cuts.SetBit(&b.Cuts, n, 1)
		}
		if leaf {
			b.Leaves = append(b.Leaves, n)
		}
		d--
	}
	nbs := g[start]
	if len(nbs) == 0 {
		return
	}
	df(start, nbs[0])
	var rc uint
	for _, nb := range nbs[1:] {
		if depth[nb] == 0 {
			rc = 1
			df(start, nb)
		}
	}
	b.Cuts.SetBit(&b.Cuts, start, rc)
	return
}
*/

// AddEdge adds an edge to a labeled graph.
//
// It can be useful for constructing undirected graphs.
//
// When n1 and n2 are distinct, it adds the arc n1->n2 and the reciprocal
// n2->n1.  When n1 and n2 are the same, it adds a single arc loop.
//
// If the edge already exists in *p, a parallel edge is added.
//
// The pointer receiver allows the method to expand the graph as needed
// to include the values n1 and n2.  If n1 or n2 happen to be greater than
// len(*p) the method does not panic, but simply expands the graph.
func (p *LabeledUndirected) AddEdge(e Edge, l LI) {
	// Similar code in AdjacencyList.AddEdge.

	// determine max of the two end points
	max := e.N1
	if e.N2 > max {
		max = e.N2
	}
	// expand graph if needed, to include both
	g := p.LabeledAdjacencyList
	if max >= NI(len(g)) {
		p.LabeledAdjacencyList = make(LabeledAdjacencyList, max+1)
		copy(p.LabeledAdjacencyList, g)
		g = p.LabeledAdjacencyList
	}
	// create one half-arc,
	g[e.N1] = append(g[e.N1], Half{To: e.N2, Label: l})
	// and except for loops, create the reciprocal
	if e.N1 != e.N2 {
		g[e.N2] = append(g[e.N2], Half{To: e.N1, Label: l})
	}
}

// TarjanBiconnectedComponents decomposes a graph into maximal biconnected
// components, components for which if any node were removed the component
// would remain connected.
//
// The receiver g must be a simple graph.  The method calls the emit argument
// for each component identified, as long as emit returns true.  If emit
// returns false, TarjanBiconnectedComponents returns immediately.
//
// See also the eqivalent unlabeled TarjanBiconnectedComponents.
func (g LabeledUndirected) TarjanBiconnectedComponents(emit func([]LabeledEdge) bool) {
	// Implemented closely to pseudocode in "Depth-first search and linear
	// graph algorithms", Robert Tarjan, SIAM J. Comput. Vol. 1, No. 2,
	// June 1972.
	//
	// Note Tarjan's "adjacency structure" is graph.AdjacencyList,
	// His "adjacency list" is an element of a graph.AdjacencyList, also
	// termed a "to-list", "neighbor list", or "child list."
	//
	// Nearly identical code in undir.go.
	number := make([]int, len(g.LabeledAdjacencyList))
	lowpt := make([]int, len(g.LabeledAdjacencyList))
	var stack []LabeledEdge
	var i int
	var biconnect func(NI, NI) bool
	biconnect = func(v, u NI) bool {
		i++
		number[v] = i
		lowpt[v] = i
		for _, w := range g.LabeledAdjacencyList[v] {
			if number[w.To] == 0 {
				stack = append(stack, LabeledEdge{Edge{v, w.To}, w.Label})
				if !biconnect(w.To, v) {
					return false
				}
				if lowpt[w.To] < lowpt[v] {
					lowpt[v] = lowpt[w.To]
				}
				if lowpt[w.To] >= number[v] {
					var bcc []LabeledEdge
					top := len(stack) - 1
					for number[stack[top].N1] >= number[w.To] {
						bcc = append(bcc, stack[top])
						stack = stack[:top]
						top--
					}
					bcc = append(bcc, stack[top])
					stack = stack[:top]
					top--
					if !emit(bcc) {
						return false
					}
				}
			} else if number[w.To] < number[v] && w.To != u {
				stack = append(stack, LabeledEdge{Edge{v, w.To}, w.Label})
				if number[w.To] < lowpt[v] {
					lowpt[v] = number[w.To]
				}
			}
		}
		return true
	}
	for w := range g.LabeledAdjacencyList {
		if number[w] == 0 && !biconnect(NI(w), 0) {
			return
		}
	}
}
