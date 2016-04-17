// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

// FromList represents a rooted tree (or forest) where each node is associated
// with a half arc identifying an arc "from" another node.
//
// Other terms for this data structure include "parent list",
// "predecessor list", "in-tree", "inverse arborescence", and
// "spaghetti stack."
//
// The Paths member represents the tree structure.  Leaves and MaxLen are
// not always needed.  Where Leaves is used it serves as a bitmap where
// Leaves.Bit(n) == 1 for each leaf n of the tree.  Where MaxLen is used it is
// provided primarily as a convenience for functions that might want to
// anticipate the maximum path length that would be encountered traversing
// the tree.
//
// Various graph search methods use a FromList to returns search results.
// For a start node of a search, From will be -1 and Len will be 1. For other
// nodes reached by the search, From represents a half arc in a path back to
// start and Len represents the number of nodes in the path.  For nodes not
// reached by the search, From will be -1 and Len will be 0.
//
// A single FromList can also represent a forest.  In this case paths from
// all leaves do not return to a single root node, but multiple root nodes.
//
// While a FromList generally encodes a tree or forest, it is technically
// possible to encode a cyclic graph.  A number of FromList methods require
// the receiver to be acyclic.  Graph methods documented to return a tree or
// forest will never return a cyclic FromList.  In other cases however,
// where a FromList is not known to by cyclic, the Cyclic method can be
// useful to validate the acyclic property.
type FromList struct {
	Paths  []PathEnd // tree representation
	Leaves Bits      // leaves of tree
	MaxLen int       // length of longest path, max of all PathEnd.Len values
}

// PathEnd associates a half arc and a path length.
//
// A PathEnd list is an element type of FromList.
type PathEnd struct {
	From NI  // a "from" half arc, the node the arc comes from
	Len  int // number of nodes in path from start
}

// NewFromList creates a FromList object of given order.
//
// The Paths member is allocated to length n but there is no other
// initialization.
func NewFromList(n int) FromList {
	return FromList{Paths: make([]PathEnd, n)}
}

// BoundsOk validates the "from" values in the list.
//
// Negative values are allowed as they indicate root nodes.
//
// BoundsOk returns true when all from values are less than len(t).
// Otherwise it returns false and a node with a from value >= len(t).
func (f FromList) BoundsOk() (ok bool, n NI) {
	for n, e := range f.Paths {
		if int(e.From) >= len(f.Paths) {
			return false, NI(n)
		}
	}
	return true, -1
}

// CommonStart returns the common start node of minimal paths to a and b.
//
// It returns -1 if a and b cannot be traced back to a common node.
//
// The method relies on populated PathEnd.Len members.  Use RecalcLen if
// the Len members are not known to be present and correct.
func (f FromList) CommonStart(a, b NI) NI {
	p := f.Paths
	if p[a].Len < p[b].Len {
		a, b = b, a
	}
	for bl := p[b].Len; p[a].Len > bl; {
		a = p[a].From
		if a < 0 {
			return -1
		}
	}
	for a != b {
		a = p[a].From
		if a < 0 {
			return -1
		}
		b = p[b].From
	}
	return a
}

// Cyclic determines if f contains a cycle, a non-empty path from a node
// back to itself.
//
// Cyclic returns true if g contains at least one cycle.  It also returns
// an example of a node involved in a cycle.
//
// Cyclic returns (false, -1) in the normal case where f is acyclic.
// Note that the bool is not an "ok" return.  A cyclic FromList is usually
// not okay.
func (f FromList) Cyclic() (cyclic bool, n NI) {
	var vis Bits
	p := f.Paths
	for i := range p {
		var path Bits
		for n := NI(i); vis.Bit(n) == 0; {
			vis.SetBit(n, 1)
			path.SetBit(n, 1)
			if n = p[n].From; n < 0 {
				break
			}
			if path.Bit(n) == 1 {
				return true, n
			}
		}
	}
	return false, -1
}

// IsolatedNodeBits returns a bitmap of isolated nodes in receiver graph f.
//
// An isolated node is one with no arcs going to or from it.
func (f FromList) IsolatedNodes() (iso Bits) {
	p := f.Paths
	iso.SetAll(len(p))
	for n, e := range p {
		if e.From >= 0 {
			iso.SetBit(NI(n), 0)
			iso.SetBit(e.From, 0)
		}
	}
	return
}

// PathTo decodes a FromList, recovering a single path.
//
// The path is returned as a list of nodes where the first element will be
// a root node and the last element will be the specified end node.
//
// Only the Paths member of the receiver is used.  Other members of the
// FromList do not need to be valid, however the MaxLen member can be useful
// for allocating argument p.
//
// Argument p can provide the result slice.  If p has capacity for the result
// it will be used, otherwise a new slice is created for the result.
//
// See also function PathTo.
func (f FromList) PathTo(end NI, p []NI) []NI {
	return PathTo(f.Paths, end, p)
}

// PathTo decodes a single path from a PathEnd list.
//
// A PathEnd list is the main data representation in a FromList.  See FromList.
//
// PathTo returns a list of nodes where the first element will be
// a root node and the last element will be the specified end node.
//
// Argument p can provide the result slice.  If p has capacity for the result
// it will be used, otherwise a new slice is created for the result.
//
// See also method FromList.PathTo.
func PathTo(paths []PathEnd, end NI, p []NI) []NI {
	n := paths[end].Len
	if n == 0 {
		return nil
	}
	if cap(p) >= n {
		p = p[:n]
	} else {
		p = make([]NI, n)
	}
	for {
		n--
		p[n] = end
		if n == 0 {
			return p
		}
		end = paths[end].From
	}
}

// Preorder traverses f calling Visitor v in preorder.
//
// Nodes are visited in order such that for any node n with from node fr,
// fr is visited before n.  Where f represents a tree, the visit ordering
// corresponds to a preordering, or depth first traversal of the tree.
// Where f represents a forest, the preorderings of the trees can be
// intermingled.
//
// Leaves must be set correctly first.  Use RecalcLeaves if leaves are not
// known to be set correctly.  FromList f cannot be cyclic.
//
// Traversal continues while v returns true.  It terminates if v returns false.
// Preorder returns true if it completes without v returning false.  Preorder
// returns false if traversal is terminated by v returning false.
func (f FromList) Preorder(v OkNodeVisitor) bool {
	p := f.Paths
	var done Bits
	var df func(NI) bool
	df = func(n NI) bool {
		done.SetBit(n, 1)
		if fr := p[n].From; fr >= 0 && done.Bit(fr) == 0 {
			df(fr)
		}
		return v(n)
	}
	for n := range f.Paths {
		p[n].Len = 0
	}
	return f.Leaves.Iterate(func(n NI) bool {
		return df(n)
	})
}

// RecalcLeaves recomputes the Leaves member of f.
func (f *FromList) RecalcLeaves() {
	p := f.Paths
	lv := &f.Leaves
	lv.SetAll(len(p))
	for n := range f.Paths {
		if fr := p[n].From; fr >= 0 {
			lv.SetBit(fr, 0)
		}
	}
}

// RecalcLen recomputes Len for each path end, and recomputes MaxLen.
//
// RecalcLen relies on the Leaves member being valid.  If it is not known
// to be valid, call RecalcLeaves before calling RecalcLen.
func (f *FromList) RecalcLen() {
	p := f.Paths
	var setLen func(NI) int
	setLen = func(n NI) int {
		switch {
		case p[n].Len > 0:
			return p[n].Len
		case p[n].From < 0:
			p[n].Len = 1
			return 1
		}
		l := 1 + setLen(p[n].From)
		p[n].Len = l
		return l
	}
	for n := range f.Paths {
		p[n].Len = 0
	}
	f.MaxLen = 0
	f.Leaves.Iterate(func(n NI) bool {
		if l := setLen(NI(n)); l > f.MaxLen {
			f.MaxLen = l
		}
		return true
	})
}

// ReRoot reorients the tree containing n to make n the root node.
//
// It keeps the tree connected by "reversing" the path from n to the old root.
//
// After ReRoot, the Leaves and Len members are invalid.
// Call RecalcLeaves or RecalcLen as needed.
func (f *FromList) ReRoot(n NI) {
	p := f.Paths
	fr := p[n].From
	if fr < 0 {
		return
	}
	p[n].From = -1
	for {
		ff := p[fr].From
		p[fr].From = n
		if ff < 0 {
			return
		}
		n = fr
		fr = ff
	}
}

// Root finds the root of a node in a FromList.
func (f FromList) Root(n NI) NI {
	for p := f.Paths; ; {
		fr := p[n].From
		if fr < 0 {
			return n
		}
		n = fr
	}
}

// Transpose constructs the directed graph corresponding to FromList f
// but with arcs in the opposite direction.  That is, from roots toward leaves.
//
// The method relies only on the From member of f.Paths.  Other members of
// the FromList are not used.
//
// See FromList.TransposeRoots for a version that also accumulates and returns
// information about the roots.
func (f FromList) Transpose() Directed {
	g := make(AdjacencyList, len(f.Paths))
	for n, p := range f.Paths {
		if p.From == -1 {
			continue
		}
		g[p.From] = append(g[p.From], NI(n))
	}
	return Directed{g}
}

// TransposeLabeled constructs the directed labeled graph corresponding
// to FromList f but with arcs in the opposite direction.  That is, from
// roots toward leaves.
//
// The argument labels can be nil.  In this case labels are generated matching
// the path indexes.  This corresponds to the "to", or child node.
//
// If labels is non-nil, it must be the same length as f.Paths and is used
// to look up label numbers by the path index.
//
// The method relies only on the From member of f.Paths.  Other members of
// the FromList are not used.
//
// See FromList.TransposeLabeledRoots for a version that also accumulates
// and returns information about the roots.
func (f FromList) TransposeLabeled(labels []LI) LabeledDirected {
	g := make(LabeledAdjacencyList, len(f.Paths))
	for n, p := range f.Paths {
		if p.From == -1 {
			continue
		}
		l := LI(n)
		if labels != nil {
			l = labels[n]
		}
		g[p.From] = append(g[p.From], Half{NI(n), l})
	}
	return LabeledDirected{g}
}

// TransposeLabeledRoots constructs the labeled directed graph corresponding
// to FromList f but with arcs in the opposite direction.  That is, from
// roots toward leaves.
//
// TransposeLabeledRoots also returns a count of roots of the resulting forest
// and a bitmap of the roots.
//
// The argument labels can be nil.  In this case labels are generated matching
// the path indexes.  This corresponds to the "to", or child node.
//
// If labels is non-nil, it must be the same length as t.Paths and is used
// to look up label numbers by the path index.
//
// The method relies only on the From member of f.Paths.  Other members of
// the FromList are not used.
//
// See FromList.TransposeLabeled for a simpler verstion that returns the
// forest only.
func (f FromList) TransposeLabeledRoots(labels []LI) (forest LabeledDirected, nRoots int, roots Bits) {
	p := f.Paths
	nRoots = len(p)
	roots.SetAll(len(p))
	g := make(LabeledAdjacencyList, len(p))
	for i, p := range f.Paths {
		if p.From == -1 {
			continue
		}
		l := LI(i)
		if labels != nil {
			l = labels[i]
		}
		n := NI(i)
		g[p.From] = append(g[p.From], Half{n, l})
		if roots.Bit(n) == 1 {
			roots.SetBit(n, 0)
			nRoots--
		}
	}
	return LabeledDirected{g}, nRoots, roots
}

// TransposeRoots constructs the directed graph corresponding to FromList f
// but with arcs in the opposite direction.  That is, from roots toward leaves.
//
// TransposeRoots also returns a count of roots of the resulting forest and
// a bitmap of the roots.
//
// The method relies only on the From member of f.Paths.  Other members of
// the FromList are not used.
//
// See FromList.Transpose for a simpler verstion that returns the forest only.
func (f FromList) TransposeRoots() (forest Directed, nRoots int, roots Bits) {
	p := f.Paths
	nRoots = len(p)
	roots.SetAll(len(p))
	g := make(AdjacencyList, len(p))
	for i, e := range p {
		if e.From == -1 {
			continue
		}
		n := NI(i)
		g[e.From] = append(g[e.From], n)
		if roots.Bit(n) == 1 {
			roots.SetBit(n, 0)
			nRoots--
		}
	}
	return Directed{g}, nRoots, roots
}
