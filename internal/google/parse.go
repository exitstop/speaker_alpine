package google

import (
	"fmt"
	"strings"
)

// Парсим текст
func ParseGoogle6(text string) (fullText string, err error) {
	//fmt.Printf("text: <<<<\n%s\n>>>\n", text)
	text = strings.ReplaceAll(text, ")]}'", "")
	tr := Tr{}

	tr.ToMap(text, 0, len(text))
	for _, it := range tr.TextOnly {
		fullText += it
	}
	return
}

type Tr struct {
	//mBr []ItemG
	Brakets  []ItemG
	Pair     []Brakets
	TextOnly []string
}
type Brakets struct {
	Level int
	Open  int
	Close int
}
type ItemG struct {
	Index int
	Open  bool
}

func (tr *Tr) ToMap(text string, braketsStart, braketsEnd int) {
	var accumulate = 0
	for {
		indexOpen := strings.Index(text[accumulate:], "[")
		indexClose := strings.Index(text[accumulate:], "]")
		var (
			flag  bool
			index int
		)
		if indexClose > indexOpen {
			flag = true
			index = indexOpen
		} else {
			index = indexClose
		}
		if index < 0 {
			break
		}
		accumulate += index
		itmeG := ItemG{Index: accumulate, Open: flag}
		tr.Brakets = append(tr.Brakets, itmeG)
		if braketsStart < 0 || accumulate > braketsEnd {
			break
		}
		accumulate++
	}
	var (
		openB      = -1
		closeB     = -1
		level      = 0
		lenBrakets = len(tr.Brakets)
		start      = 0
		//fFindedPair  = false
		stack []int
	)
	for {
		if start >= lenBrakets {
			break
		}
		if tr.Brakets[start].Open {
			stack = append(stack, tr.Brakets[start].Index)
			level++
		}
		if !tr.Brakets[start].Open {
			closeB = tr.Brakets[start].Index
			level--
		}

		if closeB != -1 {
			n := len(stack) - 1
			openB = stack[n]
			stack = stack[:n]
			pair := Brakets{Open: openB, Close: closeB, Level: level}
			tr.Pair = append(tr.Pair, pair)
			closeB = -1
		}
		start++
	}

	fmt.Println(len(tr.Pair))

	for _, it := range tr.Pair {
		if it.Level != 8 {
			continue
		}
		//fmt.Printf("(%d) [%s]\n", it.Level, text[it.Open+1:it.Close])
		out, err := ExtratText(text[it.Open+1 : it.Close])
		if err != nil {
			//fmt.Println("ExtratText:", err)
			continue
		}
		tr.TextOnly = append(tr.TextOnly, out)
	}

	fmt.Println(tr.TextOnly)

}

func ExtratText(in string) (out string, err error) {
	start := strings.Index(in, `\"`)
	if start < 0 {
		err = fmt.Errorf("start = -1")
		return
	}
	end := strings.Index(in[start+1:], `\"`)
	if start < 0 {
		err = fmt.Errorf("end = -1")
		return
	}
	if start+2 < end {
		out = in[start+2 : end]
	}
	return
}
