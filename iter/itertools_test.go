package iter

import (
	"fmt"
	"io"
	"iter"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestConcat(t *testing.T) {
	type TestCase[T any] struct {
		input    [][]T
		expected []T
	}

	table := []TestCase[any]{
		{
			input:    [][]any{{1, 2, 3}, {7, 8, 9}},
			expected: []any{1, 2, 3, 7, 8, 9},
		},
		{
			input:    [][]any{{"hello", "world"}, {"gophers"}},
			expected: []any{"hello", "world", "gophers"},
		},
		{
			input:    [][]any{{1.2, 10.44}, {5.1}},
			expected: []any{1.2, 10.44, 5.1},
		},
	}

	for _, tt := range table {
		var a []any
		for v := range Concat(tt.input...) {
			a = append(a, v)
		}

		if !slices.Equal(a, tt.expected) {
			t.Errorf("not equal")
			t.Logf("got: %v", a)
			t.Logf("exp: %v", tt.expected)
		}
	}
}

func TestConcatIter(t *testing.T) {
	type TestCase[T iter.Seq[any], V any] struct {
		input    []T
		expected []V
	}

	table := []TestCase[iter.Seq[any], any]{
		{
			input:    []iter.Seq[any]{slices.Values([]any{1, 2, 3}), slices.Values([]any{7, 8, 9})},
			expected: []any{1, 2, 3, 7, 8, 9},
		},
		{
			input:    []iter.Seq[any]{slices.Values([]any{"hello", "world"}), slices.Values([]any{"gophers"})},
			expected: []any{"hello", "world", "gophers"},
		},
		{
			input:    []iter.Seq[any]{slices.Values([]any{1.2, 10.44}), slices.Values([]any{5.1})},
			expected: []any{1.2, 10.44, 5.1},
		},
	}

	for _, tt := range table {
		var a []any
		for v := range ConcatIter(tt.input...) {
			a = append(a, v)
		}

		if !slices.Equal(a, tt.expected) {
			t.Errorf("not equal")
			t.Logf("got: %v", a)
			t.Logf("exp: %v", tt.expected)
		}
	}
}

func TestShuffle(t *testing.T) {
	type TestCase[T any] struct {
		input         []T
		shouldContain []T
	}

	table := []TestCase[any]{
		{
			input:         []any{1, 2, 3, 7, 8, 9},
			shouldContain: []any{9, 8, 7, 3, 2, 1},
		},
		{
			input:         []any{"hello", "world", "gophers"},
			shouldContain: []any{"gophers", "world", "hello"},
		},
		{
			input:         []any{1.2, 10.44, 5.1},
			shouldContain: []any{5.1, 10.44, 1.2},
		},
	}

	for _, tt := range table {
		it := Shuffle(tt.input, func(i, j int) {
			tt.input[i], tt.input[j] = tt.input[j], tt.input[i]
		})
		seq := slices.Collect(it)

		if len(seq) != len(tt.shouldContain) {
			t.Errorf("length not equal")
			t.Logf("got: %d", len(seq))
			t.Logf("exp: %d", len(tt.shouldContain))
		}
		for _, v := range tt.shouldContain {
			if !slices.Contains(seq, v) {
				t.Errorf("element not found")
				t.Logf("exp: %v", v)
			}
		}
	}
}

func TestFilter(t *testing.T) {
	type TestCase[T comparable] struct {
		input    iter.Seq[T]
		expected []T
		fn       func(T) bool
	}

	table := []TestCase[any]{
		{
			input:    slices.Values([]any{1, 2, 3, 7, 8, 9}),
			expected: []any{2, 8},
			fn: func(v any) bool {
				val, _ := v.(int)
				return val%2 == 0
			},
		},
		{
			input:    slices.Values([]any{"hello", "world", "gophers"}),
			expected: []any{"hello", "world"},
			fn: func(v any) bool {
				val, _ := v.(string)
				return val != "gophers"
			},
		},
		{
			input:    slices.Values([]any{1.2, 10.44, 5.1}),
			expected: []any{10.44, 5.1},
			fn: func(v any) bool {
				val, _ := v.(float64)
				return val != 1.2
			},
		},
	}

	for _, tt := range table {
		it := Filter(tt.input, tt.fn)
		seq := slices.Collect(it)

		if !slices.Equal(seq, tt.expected) {
			t.Errorf("not equal")
			t.Logf("got: %v", seq)
			t.Logf("exp: %v", tt.expected)
		}
	}
}

func TestMap(t *testing.T) {
	type TestCase[T, R any] struct {
		input    iter.Seq[T]
		expected []T
		fn       func(T) R
	}

	table := []TestCase[any, any]{
		{
			input:    slices.Values([]any{1, 2, 3, 7, 8, 9}),
			expected: []any{2, 4, 6, 14, 16, 18},
			fn: func(v any) any {
				val, _ := v.(int)
				return val * 2
			},
		},
		{
			input:    slices.Values([]any{1, 2, 3, 7, 8, 9}),
			expected: []any{"1", "2", "3", "7", "8", "9"},
			fn: func(v any) any {
				val, _ := v.(int)
				return strconv.Itoa(val)
			},
		},
		{
			input:    slices.Values([]any{"hello", "world", "gophers"}),
			expected: []any{"hello💣", "world💣", "gophers💣"},
			fn: func(v any) any {
				val, _ := v.(string)
				return val + "💣"
			},
		},
		{
			input:    slices.Values([]any{1.2, 10.44, 5.1}),
			expected: []any{2.4, 20.88, 10.2},
			fn: func(v any) any {
				val, _ := v.(float64)
				return val * 2
			},
		},
	}

	for _, tt := range table {
		it := Map(tt.input, tt.fn)
		seq := slices.Collect(it)

		if !slices.Equal(seq, tt.expected) {
			t.Errorf("not equal")
			t.Logf("got: %v", seq)
			t.Logf("exp: %v", tt.expected)
		}
	}
}

func TestForEach(t *testing.T) {
	type TestCase[T comparable] struct {
		input []T
		exp   string
		fn    func(int, T)
	}

	table := []TestCase[any]{
		{
			input: []any{1, 2, 3, 7, 8, 9},
			exp:   "246141618",
			fn: func(_ int, v any) {
				val, _ := v.(int)
				fmt.Print(val * 2)
			},
		},
		{
			input: []any{"hello", "world", "gophers"},
			exp:   "hello💣world💣gophers💣",
			fn: func(_ int, v any) {
				val, _ := v.(string)
				fmt.Print(val + "💣")
			},
		},
		{
			input: []any{1.2, 10.44, 5.1},
			exp:   "2.420.8810.2",
			fn: func(_ int, v any) {
				val, _ := v.(float64)
				fmt.Print(val * 2)
			},
		},
	}

	for _, tt := range table {
		prev := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		ForEach(slices.Values(tt.input), tt.fn)
		w.Close()
		os.Stdout = prev

		msg, err := io.ReadAll(r)
		if err != nil {
			t.Error(err)
		}
		r.Close()

		if !strings.EqualFold(string(msg), tt.exp) {
			t.Error("not equal")
			t.Logf("got: %v", string(msg))
			t.Logf("exp: %v", tt.exp)
		}

	}
}

func TestReduce(t *testing.T) {
	type TestCase[T, R any] struct {
		input    []T
		expected R
		fn       func(R, T) R
		init     R
	}

	table := []TestCase[any, any]{
		{
			input:    []any{1, 2, 3, 7, 8, 9},
			expected: 30,
			fn: func(acc any, cur any) any {
				a, _ := acc.(int)
				c, _ := cur.(int)
				return a + c
			},
			init: 0,
		},
		{
			input:    []any{1, 2, 3, 7, 8, 9},
			expected: "123789",
			fn: func(acc any, cur any) any {
				a, _ := acc.(string)
				c, _ := cur.(int)
				return fmt.Sprintf("%s%v", a, c)
			},
			init: 0,
		},
		{
			input:    []any{1, 2, 3, 7, 8, 9},
			expected: 10,
			fn: func(acc any, cur any) any {
				a, _ := acc.(int)
				c, _ := cur.(int)
				return a - c
			},
			init: 40,
		},

		{
			input:    []any{"hello", "world", "gophers"},
			expected: ":)helloworldgophers",
			fn: func(acc any, cur any) any {
				a, _ := acc.(string)
				c, _ := cur.(string)
				return a + c
			},
			init: ":)",
		},
		{
			input:    []any{1.2, 10.4, 5.1},
			expected: 18.7,
			fn: func(acc any, cur any) any {
				a, _ := acc.(float64)
				c, _ := cur.(float64)
				sum := a + c
				rounded := math.Round(sum*10) / 10
				return rounded
			},
			init: 2.0,
		},
	}

	for _, tt := range table {

		result := Reduce(slices.Values(tt.input), tt.fn, tt.init)

		if result != tt.expected {
			t.Errorf("not equal")
			t.Logf("got: %v", result)
			t.Logf("exp: %v", tt.expected)
		}
	}
}
