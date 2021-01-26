package main

// SortCardByNumber for sorting cards
type SortCardByNumber []Card

func (a SortCardByNumber) Len() int           { return len(a) }
func (a SortCardByNumber) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortCardByNumber) Less(i, j int) bool { return a[i].Number < a[j].Number }
