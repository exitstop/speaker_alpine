package google

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Парсим текст
func ParseGoogle6(text string) (fullText string, err error) {
	//fmt.Printf("text: <<<<\n%s\n>>>\n", text)
	text = strings.ReplaceAll(text, ")]}'", "")
	tr := Tr{}

	tr.ToMap(text, 0, utf8.RuneCountInString(text))
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

	//fmt.Println(len(tr.Pair))

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

	//fmt.Println(tr.TextOnly)

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

/*
\"([^\"]*)\"
(?P<one>\[([^\"]*)\])
(?P<one>\\"([^[]*)\\")

		`]}'

1069
[["wrb.fr","MkEWBc","[[null,null,\"en\",[[[0,[[[null,113]],[true]]]],113],[[\"A static page is a page delivered to the user exactly as stored and with no chance on being changed, end of story\",null,null,113]]],[[[null,\"Staticheskaya stranitsa - eto stranitsa, dostavlennaya pol'zovatelyu tochno tak zhe, kak khranitsya, i bez shansov na izmeneniye, konets istorii\",null,null,null,[[\"Статическая страница - это страница, доставленная пользователю точно так же, как хранится, и без шансов на изменение, конец истории\",null,null,null,[[\"1\",[1]],[\"Статическая страница - это страница, доставленная пользователю точно так же, как хранится, и без шансов изменить, конец истории\",[11]]]]]]],\"ru\",1,\"en\",[\"A static page is a page delivered to the user exactly as stored and with no chance on being changed, end of story\",\"auto\",\"ru\",true]],\"en\"]",null,null,null,"generic"],["di",31],["af.httprm",30,"9058989769143286083",8]]
26
[["e",4,null,null,1428]]
*/

var regexpString = `\\"(?P<one>([^,][А-Яа-я- \r\n\v\w'*.:,\d]{4,}))\\"`

func ParseGoogle7(textBeforeTranslate, text string) (fullText string, err error) {

	reg0, err := regexp.Compile(regexpString)
	if err != nil {
		return
	}

	lenTextBeforeTranslate := utf8.RuneCountInString(textBeforeTranslate)
	if lenTextBeforeTranslate < 1 {
		return
	}

	text2 := textBeforeTranslate[:lenTextBeforeTranslate-1]
	// count substring in string
	count := strings.Count(text2, ".")
	count += strings.Count(text2, "?")
	var countSentence = 0

	//if count > 6 {
	//count = 6
	//countSentence = count - 1
	//}
	fmt.Println("lenTextBeforeTranslate:", lenTextBeforeTranslate)
	//if lenTextBeforeTranslate > 1600 {
	//fmt.Println("textBeforeTranslate:", textBeforeTranslate)
	//}

	// iterate all matches
	iter := reg0.FindAllStringSubmatchIndex(text, -1)

	var localText []string
	for index, it := range iter {
		t := text[it[0]:it[1]]
		fmt.Printf("[%d] %s\n", index, t)
		localText = append(localText, t)
	}

	lenText := len(localText)
	if count == 0 {
		count = 3
		countSentence = count - 1
	} else {
		count = int(float64(lenText) * 0.32)
		countSentence = count - 1
	}

	if count > lenText {
		count = 1
		countSentence = 1
	}

	fmt.Println("count: ", count)

	localText = localText[count : lenText-1]
	lenText -= count + 1

	countVariant := lenText / countSentence
	if lenText == 3 || lenText == 2 {
		countVariant = lenText
	}

	fmt.Println(localText)
	fmt.Println("countVariant:", countVariant)
	fmt.Println("lenText:", lenText)
	fmt.Println("countSentence: ", countSentence)

	for i := 0; i < lenText; i += countVariant {
		fullText += localText[i]
	}

	fullText = strings.ReplaceAll(fullText, "\\\"", ".")

	return
}
