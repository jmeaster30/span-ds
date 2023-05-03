package main

import (
	"fmt"
	"math"
)

type Span struct {
	From          float64
	FromInclusive bool
	To            float64
	ToInclusive   bool

	Left  *Span
	Right *Span
}

func NewSpan(fromInclusive bool, from float64, to float64, toInclusive bool) *Span {
	if from > to {
		panic("From shouldn't be less than to")
	}
	return &Span{
		FromInclusive: fromInclusive,
		From:          from,
		To:            to,
		ToInclusive:   toInclusive,
	}
}

// possibly useful api
func EncompassingSpan(a *Span, b *Span) *Span {
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
	return &Span{
		FromInclusive: fromInclusive,
		From:          minFrom,
		To:            maxTo,
		ToInclusive:   toInclusive,
	}
}

// helper function
func (base *Span) AddLeftRight(a *Span, b *Span) {
	if a.To < b.From {
		base.Left = a
		base.Right = b
	} else if a.From > b.To {
		base.Right = a
		base.Left = b
	} else {
		panic("I DON'T KNOW")
	}
}

func (a *Span) Intersect(b *Span) *Span {
	return nil
}

func (a *Span) Union(b *Span) *Span {
	// if overlapping combine the nodes into one
	if (a.From < b.From && b.From < a.To) ||
		(b.From < a.From && a.From < b.To) ||
		(a.To == b.From && (a.ToInclusive || b.FromInclusive)) ||
		(a.From == b.To && (a.FromInclusive || b.ToInclusive)) {
		return EncompassingSpan(a, b)
	}
	// else dont
	base := EncompassingSpan(a, b)
	base.AddLeftRight(a, b)
	return base
}

func (a *Span) Complement() *Span {
	left := NewSpan(true, math.MaxFloat64*-1, a.From, !a.FromInclusive)
	right := NewSpan(!a.ToInclusive, a.To, math.MaxFloat64, true)
	return left.Union(right)
}

func (a *Span) Difference(b *Span) *Span {
	return nil
}

func (a *Span) IsDisjoint(b *Span) bool {
	return false
}

func (a *Span) Contains(element float64) bool {
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

func (a *Span) ToString() string {
	if a.Left == nil && a.Right == nil {
		start := "("
		if a.FromInclusive {
			start = "["
		}
		end := ")"
		if a.ToInclusive {
			end = "]"
		}
		return fmt.Sprintf("%s%.6e, %.6e%s", start, a.From, a.To, end)
	}
	return fmt.Sprintf("%s %s", a.Left.ToString(), a.Right.ToString())
}

func main() {
	fmt.Println("Hello world")

	a := NewSpan(true, 3, 10, false)
	b := NewSpan(false, -5, 1, true)
	fmt.Printf("a = %s\n", a.ToString())
	fmt.Printf("b = %s\n", b.ToString())

	fmt.Printf("\na complement = %s\n", a.Complement().ToString())
	fmt.Printf("b complement = %s\n", b.Complement().ToString())

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

	d := NewSpan(true, 2, 5, false)
	e := c.Union(d)
	fmt.Printf("\nd = %s\n", d.ToString())
	fmt.Printf("e = %s\n", e.ToString())

	fmt.Printf("\n10 in e = %t\n", e.Contains(10))
	fmt.Printf("2 in e = %t\n", e.Contains(2))
	fmt.Printf("6 in e = %t\n", e.Contains(6))
	fmt.Printf("1 in e = %t\n", e.Contains(1))
	fmt.Printf("-4.9999 in e = %t\n", e.Contains(-4.9999))
	fmt.Printf("700 in e = %t\n", e.Contains(700))

}
