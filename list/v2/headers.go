package list_v2

import (
	"sort"
)

type Item interface {
	Len() int
}

type Headers struct {
	total   int
	sorted  bool
	lengths sort.IntSlice
}

func (c *Headers) sort() {
	if !c.sorted {
		c.lengths.Sort()
		c.sorted = true
	}
}

func (c *Headers) Add(item string) {
	ln := len(item)
	c.lengths = append(c.lengths, ln)
	c.total += len(item)
	c.sorted = false
}

func (c *Headers) Avg() int {
	if c.lengths.Len() == 0 {
		return 0
	}
	return c.total / c.lengths.Len()
}

func (c *Headers) Max() int {
	if len(c.lengths) == 0 {
		return 0
	}
	c.sort()
	return c.lengths[len(c.lengths)-1]
}

func (c *Headers) Median() int {
	c.sort()
	ln := c.lengths.Len()
	mid := ln / 2
	if ln == 0 {
		return 0
	} else if ln%2 == 1 {
		return c.lengths[mid]
	}
	return (c.lengths[mid-1] + c.lengths[mid]) / 2
}

func (c *Headers) Remove(item string) {
	idx := c.lengths.Search(len(item))
	if idx >= c.lengths.Len() {
		return
	} else if idx == c.lengths.Len()-1 {
		c.lengths = c.lengths[:idx]
	} else {
		c.lengths = append(c.lengths[:idx], c.lengths[idx+1:]...)
	}
	c.total -= len(item)
	c.sorted = false
}

func (c *Headers) Set(items []string) {
	c.lengths = make(sort.IntSlice, len(items))
	c.sorted = false
	c.total = 0

	for i, item := range items {
		ln := len(item)
		c.lengths[i] = ln
		c.total += c.lengths[i]
	}

	c.sort()
}
