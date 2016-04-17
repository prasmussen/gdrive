// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

// undir_RO.go is code generated from undir_cg.go by directives in graph.go.
// Editing undir_cg.go is okay.  It is the code generation source.
// DO NOT EDIT undir_RO.go.
// The RO means read only and it is upper case RO to slow you down a bit
// in case you start to edit the file.

// Bipartite determines if a connected component of an undirected graph
// is bipartite, a component where nodes can be partitioned into two sets
// such that every edge in the component goes from one set to the other.
//
// Argument n can be any representative node of the component.
//
// If the component is bipartite, Bipartite returns true and a two-coloring
// of the component.  Each color set is returned as a bitmap.  If the component
// is not bipartite, Bipartite returns false and a representative odd cycle.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledUndirected) Bipartite(n NI) (b bool, c1, c2 Bits, oc []NI) {
	b = true
	var open bool
	var df func(n NI, c1, c2 *Bits)
	df = func(n NI, c1, c2 *Bits) {
		c1.SetBit(n, 1)
		for _, nb := range g.LabeledAdjacencyList[n] {
			if c1.Bit(nb.To) == 1 {
				b = false
				oc = []NI{nb.To, n}
				open = true
				return
			}
			if c2.Bit(nb.To) == 1 {
				continue
			}
			df(nb.To, c2, c1)
			if b {
				continue
			}
			switch {
			case !open:
			case n == oc[0]:
				open = false
			default:
				oc = append(oc, n)
			}
			return
		}
	}
	df(n, &c1, &c2)
	if b {
		return b, c1, c2, nil
	}
	return b, Bits{}, Bits{}, oc
}

// BronKerbosch1 finds maximal cliques in an undirected graph.
//
// The graph must not contain parallel edges or loops.
//
// See https://en.wikipedia.org/wiki/Clique_(graph_theory) and
// https://en.wikipedia.org/wiki/Bron%E2%80%93Kerbosch_algorithm for background.
//
// This method implements the BronKerbosch1 algorithm of WP; that is,
// the original algorithm without improvements.
//
// The method calls the emit argument for each maximal clique in g, as long
// as emit returns true.  If emit returns false, BronKerbosch1 returns
// immediately.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also more sophisticated variants BronKerbosch2 and BronKerbosch3.
func (g LabeledUndirected) BronKerbosch1(emit func([]NI) bool) {
	a := g.LabeledAdjacencyList
	var f func(R, P, X *Bits) bool
	f = func(R, P, X *Bits) bool {
		switch {
		case !P.Zero():
			var r2, p2, x2 Bits
			pf := func(n NI) bool {
				r2.Set(*R)
				r2.SetBit(n, 1)
				p2.Clear()
				x2.Clear()
				for _, to := range a[n] {
					if P.Bit(to.To) == 1 {
						p2.SetBit(to.To, 1)
					}
					if X.Bit(to.To) == 1 {
						x2.SetBit(to.To, 1)
					}
				}
				if !f(&r2, &p2, &x2) {
					return false
				}
				P.SetBit(n, 0)
				X.SetBit(n, 1)
				return true
			}
			if !P.Iterate(pf) {
				return false
			}
		case X.Zero():
			return emit(R.Slice())
		}
		return true
	}
	var R, P, X Bits
	P.SetAll(len(a))
	f(&R, &P, &X)
}

// BKPivotMaxDegree is a strategy for BronKerbosch methods.
//
// To use it, take the method value (see golang.org/ref/spec#Method_values)
// and pass it as the argument to BronKerbosch2 or 3.
//
// The strategy is to pick the node from P or X with the maximum degree
// (number of edges) in g.  Note this is a shortcut from evaluating degrees
// in P.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledUndirected) BKPivotMaxDegree(P, X *Bits) (p NI) {
	// choose pivot u as highest degree node from P or X
	a := g.LabeledAdjacencyList
	maxDeg := -1
	P.Iterate(func(n NI) bool { // scan P
		if d := len(a[n]); d > maxDeg {
			p = n
			maxDeg = d
		}
		return true
	})
	X.Iterate(func(n NI) bool { // scan X
		if d := len(a[n]); d > maxDeg {
			p = n
			maxDeg = d
		}
		return true
	})
	return
}

// BKPivotMinP is a strategy for BronKerbosch methods.
//
// To use it, take the method value (see golang.org/ref/spec#Method_values)
// and pass it as the argument to BronKerbosch2 or 3.
//
// The strategy is to simply pick the first node in P.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledUndirected) BKPivotMinP(P, X *Bits) NI {
	return P.From(0)
}

// BronKerbosch2 finds maximal cliques in an undirected graph.
//
// The graph must not contain parallel edges or loops.
//
// See https://en.wikipedia.org/wiki/Clique_(graph_theory) and
// https://en.wikipedia.org/wiki/Bron%E2%80%93Kerbosch_algorithm for background.
//
// This method implements the BronKerbosch2 algorithm of WP; that is,
// the original algorithm plus pivoting.
//
// The argument is a pivot function that must return a node of P or X.
// P is guaranteed to contain at least one node.  X is not.
// For example see BKPivotMaxDegree.
//
// The method calls the emit argument for each maximal clique in g, as long
// as emit returns true.  If emit returns false, BronKerbosch1 returns
// immediately.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also simpler variant BronKerbosch1 and more sophisticated variant
// BronKerbosch3.
func (g LabeledUndirected) BronKerbosch2(pivot func(P, X *Bits) NI, emit func([]NI) bool) {
	a := g.LabeledAdjacencyList
	var f func(R, P, X *Bits) bool
	f = func(R, P, X *Bits) bool {
		switch {
		case !P.Zero():
			var r2, p2, x2, pnu Bits
			// compute P \ N(u).  next 5 lines are only difference from BK1
			pnu.Set(*P)
			for _, to := range a[pivot(P, X)] {
				pnu.SetBit(to.To, 0)
			}
			// remaining code like BK1
			pf := func(n NI) bool {
				r2.Set(*R)
				r2.SetBit(n, 1)
				p2.Clear()
				x2.Clear()
				for _, to := range a[n] {
					if P.Bit(to.To) == 1 {
						p2.SetBit(to.To, 1)
					}
					if X.Bit(to.To) == 1 {
						x2.SetBit(to.To, 1)
					}
				}
				if !f(&r2, &p2, &x2) {
					return false
				}
				P.SetBit(n, 0)
				X.SetBit(n, 1)
				return true
			}
			if !pnu.Iterate(pf) {
				return false
			}
		case X.Zero():
			return emit(R.Slice())
		}
		return true
	}
	var R, P, X Bits
	P.SetAll(len(a))
	f(&R, &P, &X)
}

// BronKerbosch3 finds maximal cliques in an undirected graph.
//
// The graph must not contain parallel edges or loops.
//
// See https://en.wikipedia.org/wiki/Clique_(graph_theory) and
// https://en.wikipedia.org/wiki/Bron%E2%80%93Kerbosch_algorithm for background.
//
// This method implements the BronKerbosch3 algorithm of WP; that is,
// the original algorithm with pivoting and degeneracy ordering.
//
// The argument is a pivot function that must return a node of P or X.
// P is guaranteed to contain at least one node.  X is not.
// For example see BKPivotMaxDegree.
//
// The method calls the emit argument for each maximal clique in g, as long
// as emit returns true.  If emit returns false, BronKerbosch1 returns
// immediately.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also simpler variants BronKerbosch1 and BronKerbosch2.
func (g LabeledUndirected) BronKerbosch3(pivot func(P, X *Bits) NI, emit func([]NI) bool) {
	a := g.LabeledAdjacencyList
	var f func(R, P, X *Bits) bool
	f = func(R, P, X *Bits) bool {
		switch {
		case !P.Zero():
			var r2, p2, x2, pnu Bits
			// compute P \ N(u).  next lines are only difference from BK1
			pnu.Set(*P)
			for _, to := range a[pivot(P, X)] {
				pnu.SetBit(to.To, 0)
			}
			// remaining code like BK2
			pf := func(n NI) bool {
				r2.Set(*R)
				r2.SetBit(n, 1)
				p2.Clear()
				x2.Clear()
				for _, to := range a[n] {
					if P.Bit(to.To) == 1 {
						p2.SetBit(to.To, 1)
					}
					if X.Bit(to.To) == 1 {
						x2.SetBit(to.To, 1)
					}
				}
				if !f(&r2, &p2, &x2) {
					return false
				}
				P.SetBit(n, 0)
				X.SetBit(n, 1)
				return true
			}
			if !pnu.Iterate(pf) {
				return false
			}
		case X.Zero():
			return emit(R.Slice())
		}
		return true
	}
	var R, P, X Bits
	P.SetAll(len(a))
	// code above same as BK2
	// code below new to BK3
	_, ord, _ := g.Degeneracy()
	var p2, x2 Bits
	for _, n := range ord {
		R.SetBit(n, 1)
		p2.Clear()
		x2.Clear()
		for _, to := range a[n] {
			if P.Bit(to.To) == 1 {
				p2.SetBit(to.To, 1)
			}
			if X.Bit(to.To) == 1 {
				x2.SetBit(to.To, 1)
			}
		}
		if !f(&R, &p2, &x2) {
			return
		}
		R.SetBit(n, 0)
		P.SetBit(n, 0)
		X.SetBit(n, 1)
	}
}

// ConnectedComponentBits returns a function that iterates over connected
// components of g, returning a member bitmap for each.
//
// Each call of the returned function returns the order (number of nodes)
// and bits of a connected component.  The returned function returns zeros
// after returning all connected components.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also ConnectedComponentReps, which has lighter weight return values.
func (g LabeledUndirected) ConnectedComponentBits() func() (order int, bits Bits) {
	a := g.LabeledAdjacencyList
	var vg Bits  // nodes visited in graph
	var vc *Bits // nodes visited in current component
	var nc int
	var df func(NI)
	df = func(n NI) {
		vg.SetBit(n, 1)
		vc.SetBit(n, 1)
		nc++
		for _, nb := range a[n] {
			if vg.Bit(nb.To) == 0 {
				df(nb.To)
			}
		}
		return
	}
	var n NI
	return func() (o int, bits Bits) {
		for ; n < NI(len(a)); n++ {
			if vg.Bit(n) == 0 {
				vc = &bits
				nc = 0
				df(n)
				return nc, bits
			}
		}
		return
	}
}

// ConnectedComponentLists returns a function that iterates over connected
// components of g, returning the member list of each.
//
// Each call of the returned function returns a node list of a connected
// component.  The returned function returns nil after returning all connected
// components.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also ConnectedComponentReps, which has lighter weight return values.
func (g LabeledUndirected) ConnectedComponentLists() func() []NI {
	a := g.LabeledAdjacencyList
	var vg Bits // nodes visited in graph
	var m []NI  // members of current component
	var df func(NI)
	df = func(n NI) {
		vg.SetBit(n, 1)
		m = append(m, n)
		for _, nb := range a[n] {
			if vg.Bit(nb.To) == 0 {
				df(nb.To)
			}
		}
		return
	}
	var n NI
	return func() []NI {
		for ; n < NI(len(a)); n++ {
			if vg.Bit(n) == 0 {
				m = nil
				df(n)
				return m
			}
		}
		return nil
	}
}

// ConnectedComponentReps returns a representative node from each connected
// component of g.
//
// Returned is a slice with a single representative node from each connected
// component and also a parallel slice with the order, or number of nodes,
// in the corresponding component.
//
// This is fairly minimal information describing connected components.
// From a representative node, other nodes in the component can be reached
// by depth first traversal for example.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also ConnectedComponentBits and ConnectedComponentLists which can
// collect component members in a single traversal, and IsConnected which
// is an even simpler boolean test.
func (g LabeledUndirected) ConnectedComponentReps() (reps []NI, orders []int) {
	a := g.LabeledAdjacencyList
	var c Bits
	var o int
	var df func(NI)
	df = func(n NI) {
		c.SetBit(n, 1)
		o++
		for _, nb := range a[n] {
			if c.Bit(nb.To) == 0 {
				df(nb.To)
			}
		}
		return
	}
	for n := range a {
		if c.Bit(NI(n)) == 0 {
			reps = append(reps, NI(n))
			o = 0
			df(NI(n))
			orders = append(orders, o)
		}
	}
	return
}

// Copy makes a deep copy of g.
// Copy also computes the arc size ma, the number of arcs.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledUndirected) Copy() (c LabeledUndirected, ma int) {
	l, s := g.LabeledAdjacencyList.Copy()
	return LabeledUndirected{l}, s
}

// Degeneracy computes k-degeneracy, vertex ordering and k-cores.
//
// See Wikipedia https://en.wikipedia.org/wiki/Degeneracy_(graph_theory)
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledUndirected) Degeneracy() (k int, ordering []NI, cores []int) {
	a := g.LabeledAdjacencyList
	// WP algorithm
	ordering = make([]NI, len(a))
	var L Bits
	d := make([]int, len(a))
	var D [][]NI
	for v, nb := range a {
		dv := len(nb)
		d[v] = dv
		for len(D) <= dv {
			D = append(D, nil)
		}
		D[dv] = append(D[dv], NI(v))
	}
	for ox := range a {
		// find a non-empty D
		i := 0
		for len(D[i]) == 0 {
			i++
		}
		// k is max(i, k)
		if i > k {
			for len(cores) <= i {
				cores = append(cores, 0)
			}
			cores[k] = ox
			k = i
		}
		// select from D[i]
		Di := D[i]
		last := len(Di) - 1
		v := Di[last]
		// Add v to ordering, remove from Di
		ordering[ox] = v
		L.SetBit(v, 1)
		D[i] = Di[:last]
		// move neighbors
		for _, nb := range a[v] {
			if L.Bit(nb.To) == 1 {
				continue
			}
			dn := d[nb.To] // old number of neighbors of nb
			Ddn := D[dn]   // nb is in this list
			// remove it from the list
			for wx, w := range Ddn {
				if w == nb.To {
					last := len(Ddn) - 1
					Ddn[wx], Ddn[last] = Ddn[last], Ddn[wx]
					D[dn] = Ddn[:last]
				}
			}
			dn-- // new number of neighbors
			d[nb.To] = dn
			// re--add it to it's new list
			D[dn] = append(D[dn], nb.To)
		}
	}
	cores[k] = len(ordering)
	return
}

// Degree for undirected graphs, returns the degree of a node.
//
// The degree of a node in an undirected graph is the number of incident
// edges, where loops count twice.
//
// If g is known to be loop-free, the result is simply equivalent to len(g[n]).
// See handshaking lemma example at AdjacencyList.ArcSize.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledUndirected) Degree(n NI) int {
	to := g.LabeledAdjacencyList[n]
	d := len(to) // just "out" degree,
	for _, to := range to {
		if to.To == n {
			d++ // except loops count twice
		}
	}
	return d
}

// FromList constructs a FromList representing the tree reachable from
// the given root.
//
// The connected component containing root should represent a simple graph,
// connected as a tree.
//
// For nodes connected as a tree, the Path member of the returned FromList
// will be populated with both From and Len values.  The MaxLen member will be
// set but Leaves will be left a zero value.  Return value cycle will be -1.
//
// If the connected component containing root is not connected as a tree,
// a cycle will be detected.  The returned FromList will be a zero value and
// return value cycle will be a node involved in the cycle.
//
// Loops and parallel edges will be detected as cycles, however only in the
// connected component containing root.  If g is not fully connected, nodes
// not reachable from root will have PathEnd values of {From: -1, Len: 0}.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledUndirected) FromList(root NI) (f FromList, cycle NI) {
	p := make([]PathEnd, len(g.LabeledAdjacencyList))
	for i := range p {
		p[i].From = -1
	}
	ml := 0
	var df func(NI, NI) bool
	df = func(fr, n NI) bool {
		l := p[n].Len + 1
		for _, to := range g.LabeledAdjacencyList[n] {
			if to.To == fr {
				continue
			}
			if p[to.To].Len > 0 {
				cycle = to.To
				return false
			}
			p[to.To] = PathEnd{From: n, Len: l}
			if l > ml {
				ml = l
			}
			if !df(n, to.To) {
				return false
			}
		}
		return true
	}
	p[root].Len = 1
	if !df(-1, root) {
		return
	}
	return FromList{Paths: p, MaxLen: ml}, -1
}

// IsConnected tests if an undirected graph is a single connected component.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also ConnectedComponentReps for a method returning more information.
func (g LabeledUndirected) IsConnected() bool {
	a := g.LabeledAdjacencyList
	if len(a) == 0 {
		return true
	}
	var b Bits
	b.SetAll(len(a))
	var df func(NI)
	df = func(n NI) {
		b.SetBit(n, 0)
		for _, to := range a[n] {
			if b.Bit(to.To) == 1 {
				df(to.To)
			}
		}
	}
	df(0)
	return b.Zero()
}

// IsTree identifies trees in undirected graphs.
//
// Return value isTree is true if the connected component reachable from root
// is a tree.  Further, return value allTree is true if the entire graph g is
// connected.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g LabeledUndirected) IsTree(root NI) (isTree, allTree bool) {
	a := g.LabeledAdjacencyList
	var v Bits
	v.SetAll(len(a))
	var df func(NI, NI) bool
	df = func(fr, n NI) bool {
		if v.Bit(n) == 0 {
			return false
		}
		v.SetBit(n, 0)
		for _, to := range a[n] {
			if to.To != fr && !df(n, to.To) {
				return false
			}
		}
		return true
	}
	v.SetBit(root, 0)
	for _, to := range a[root] {
		if !df(root, to.To) {
			return false, false
		}
	}
	return true, v.Zero()
}

// Size returns the number of edges in g.
//
// See also ArcSize and HasLoop.
func (g LabeledUndirected) Size() int {
	m2 := 0
	for fr, to := range g.LabeledAdjacencyList {
		m2 += len(to)
		for _, to := range to {
			if to.To == NI(fr) {
				m2++
			}
		}
	}
	return m2 / 2
}
