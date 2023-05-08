package main

import (
	"fmt"
	"math"
)

type Interval struct {
	From          float64
	FromInclusive bool
	To            float64
	ToInclusive   bool

	Left  *Interval
	Right *Interval
}

func NewInterval(fromInclusive bool, from float64, to float64, toInclusive bool) *Interval {
	if from > to {
		fmt.Printf("%f\n%f\n", from, to)
		panic("From should be less than to")
	}
	return &Interval{
		FromInclusive: fromInclusive,
		From:          from,
		To:            to,
		ToInclusive:   toInclusive,
	}
}

// possibly useful api
func EncompassingSpan(a *Interval, b *Interval) *Interval {
	minFrom := math.Min(a.From, b.From)
	fromInclusive := a.FromInclusive
	if minFrom == b.From {
		fromInclusive = b.FromInclusive
	}
	maxTo := math.Min(a.To, b.To)
	toInclusive := a.ToInclusive
	if maxTo == b.To {
		toInclusive = b.ToInclusive
	}
	return &Interval{
		FromInclusive: fromInclusive,
		From:          minFrom,
		To:            maxTo,
		ToInclusive:   toInclusive,
	}
}

func (a *Interval) Copy() *Interval {
	if a == nil {
		return nil
	}
	return &Interval{
		From:          a.From,
		FromInclusive: a.FromInclusive,
		To:            a.To,
		ToInclusive:   a.ToInclusive,
		Left:          a.Left.Copy(),
		Right:         a.Right.Copy(),
	}
}

func (a *Interval) List() []*Interval {
	results := []*Interval{}
	if a.Left == nil && a.Right == nil {
		return []*Interval{a.Copy()}
	}
	results = append(results, a.Left.List()...)
	results = append(results, a.Right.List()...)
	return results
}

// helper function
func (base *Interval) AddLeftRight(a *Interval, b *Interval) {
	if a.To < b.From {
		base.Left = a.Copy()
		base.Right = b.Copy()
	} else if a.From > b.To {
		base.Right = a.Copy()
		base.Left = b.Copy()
	} else {
		panic("I DON'T KNOW")
	}
}

func Overlaps(a *Interval, b *Interval) bool {
	//fmt.Printf("%t\n", a.From < b.From && b.From < a.To)
	//fmt.Printf("%t\n", b.From < a.From && a.From < b.To)
	//fmt.Printf("%t\n", a.From == b.From && (a.FromInclusive || b.FromInclusive))
	//fmt.Printf("%t\n", a.To == b.To && (a.ToInclusive || b.ToInclusive))
	return (a.From < b.From && b.From < a.To) ||
		(b.From < a.From && a.From < b.To) ||
		(a.From == b.From && (a.FromInclusive || b.FromInclusive)) ||
		(a.To == b.To && (a.ToInclusive || b.ToInclusive))
}

func (a *Interval) Intersect(b *Interval) *Interval {
	if Overlaps(a, b) {
		maxFrom := math.Max(a.From, b.From)
		fromInclusive := a.FromInclusive
		if maxFrom == b.From {
			fromInclusive = b.FromInclusive
		}
		minTo := math.Min(a.To, b.To)
		toInclusive := a.ToInclusive
		if minTo == b.To {
			toInclusive = b.ToInclusive
		}
		return NewInterval(fromInclusive, maxFrom, minTo, toInclusive)
	} else {
		return NewInterval(false, 0, 0, false)
	}
}

func (a *Interval) Union(b *Interval) *Interval {
	if a == nil {
		return b.Copy()
	}

	if b == nil {
		return a.Copy()
	}

	// if overlapping combine the nodes into one
	if Overlaps(a, b) {
		return EncompassingSpan(a, b)
	}
	// else dont
	base := EncompassingSpan(a, b)
	base.AddLeftRight(a, b)
	return base
}

func (a *Interval) Complement() *Interval {
	l := a.List()
	var result *Interval
	for i, item := range l {
		if i == 0 && item.From != math.MaxFloat64*-1 {
			result = NewInterval(true, math.MaxFloat64*-1, item.From, !item.FromInclusive)
		}

		if i > 0 && i <= len(l)-1 {
			result = result.Union(NewInterval(!l[i-1].ToInclusive, l[i-1].To, item.From, !item.FromInclusive))
		}

		if i == len(l)-1 && item.To != math.MaxFloat64 {
			result = result.Union(NewInterval(!item.ToInclusive, item.To, math.MaxFloat64, true))
		}
	}
	return result
}

func (a *Interval) Difference(b *Interval) *Interval {
	return nil
}

func (a *Interval) IsDisjoint(b *Interval) bool {
	return false
}

func (a *Interval) Contains(element float64) bool {
	if ((a.FromInclusive && a.From <= element) || (!a.FromInclusive && a.From < element)) &&
		((a.ToInclusive && a.To >= element) || (!a.ToInclusive && a.To > element)) {
		if a.Left == nil && a.Right == nil {
			return true
		}
		return a.Left.Contains(element) || a.Right.Contains(element)
	} else {
		return false
	}
}

func (a *Interval) ToString() string {
	if a.Left == nil && a.Right == nil {
		start := "("
		if a.FromInclusive {
			start = "["
		}
		end := ")"
		if a.ToInclusive {
			end = "]"
		}
		return fmt.Sprintf("%s%f, %f%s", start, a.From, a.To, end)
	}
	return fmt.Sprintf("%s %s", a.Left.ToString(), a.Right.ToString())
}

func main() {
	fmt.Println("Hello world")

	a := NewInterval(true, 3, 10, false)
	b := NewInterval(false, -5, 1, true)
	fmt.Printf("a = %s\n", a.ToString())
	fmt.Printf("b = %s\n", b.ToString())

	fmt.Printf("\na complement = %s\n", a.Complement().ToString())
	fmt.Printf("b complement = %s\n", b.Complement().ToString())
	fmt.Printf("b complement complement = %s\n", b.Complement().Complement().ToString())

	fmt.Printf("\n10 in a = %t\n", a.Contains(10))
	fmt.Printf("2 in a = %t\n", a.Contains(2))
	fmt.Printf("6 in a = %t\n", a.Contains(6))
	fmt.Printf("1 in b = %t\n", b.Contains(1))
	fmt.Printf("-4.9999 in b = %t\n", b.Contains(-4.9999))
	fmt.Printf("700 in b = %t\n", b.Contains(700))

	c := a.Union(b)
	fmt.Printf("\nc = %s\n", c.ToString())

	fmt.Printf("10 in c = %t\n", a.Contains(10))
	fmt.Printf("2 in c = %t\n", a.Contains(2))
	fmt.Printf("6 in c = %t\n", a.Contains(6))
	fmt.Printf("1 in c = %t\n", b.Contains(1))
	fmt.Printf("-4.9999 in c = %t\n", b.Contains(-4.9999))
	fmt.Printf("700 in c = %t\n", b.Contains(700))

	d := NewInterval(true, 2, 5, false)
	e := c.Union(d)
	fmt.Printf("\nd = %s\n", d.ToString())
	fmt.Printf("e = %s\n", e.ToString())

	fmt.Printf("\n10 in e = %t\n", e.Contains(10))
	fmt.Printf("2 in e = %t\n", e.Contains(2))
	fmt.Printf("6 in e = %t\n", e.Contains(6))
	fmt.Printf("1 in e = %t\n", e.Contains(1))
	fmt.Printf("-4.9999 in e = %t\n", e.Contains(-4.9999))
	fmt.Printf("700 in e = %t\n", e.Contains(700))

	f := NewInterval(true, 5, 10, true)
	g := NewInterval(true, 7, 15, true)
	h := f.Intersect(g)
	fmt.Printf("\nf = %s\n", f.ToString())
	fmt.Printf("g = %s\n", g.ToString())
	fmt.Printf("h = %s\n", h.ToString())

	i := NewInterval(true, 1, 4, true)
	j := NewInterval(true, 5, 10, true)
	k := i.Intersect(j)
	fmt.Printf("\ni = %s\n", i.ToString())
	fmt.Printf("j = %s\n", j.ToString())
	fmt.Printf("k = %s\n", k.ToString())

	l := NewInterval(true, 2, 4, true)
	m := NewInterval(true, 8, 12, true)
	n := l.Union(m)
	o := n.Complement()
	fmt.Printf("\nl = %s\n", l.ToString())
	fmt.Printf("m = %s\n", m.ToString())
	fmt.Printf("n = %s\n", n.ToString())
	fmt.Printf("n complement = %s\n", o.ToString())
}
