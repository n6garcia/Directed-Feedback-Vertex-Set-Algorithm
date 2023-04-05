package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type dictInterface interface {
	Print()
	PrintSize()
	loadData(string)
	AddData(*Graph)
	getDef(string) string
	expandDef([]string, string) string
	verify([]string) bool
}

type Dictionary struct {
	definitions map[string]*Definition
	//definitions []*Definition
	// ^--- old DS
}

type Definition struct {
	name  string
	words []string
}

func (d *Dictionary) Print() {
	for _, v := range d.definitions {
		fmt.Println("name: ", v.name)
		fmt.Println("words: ", v.words)
	}
}

func (d *Dictionary) PrintSize() {
	fmt.Println("\nsize : ", len(d.definitions))
}

func (d *Dictionary) addDef(n string, w []string) {
	defn := &Definition{name: n, words: w}
	d.definitions[n] = defn
}

func (d *Dictionary) loadData(fn string) {
	bytes, err := os.ReadFile("wrangle/cleaned/" + fn) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println("isValid: ", json.Valid(bytes))

	var myData map[string][]interface{}

	json.Unmarshal(bytes, &myData)

	for k, v := range myData {
		var words []string
		for _, u := range v {
			words = append(words, u.(string))
		}
		d.addDef(k, words)
	}
}

// Transfers Data in Dictionary to Graph
func (d *Dictionary) AddData(g *Graph) {
	fmt.Println("adding data to graph...")

	for _, v := range d.definitions {
		g.AddVertex(v.name)
		for _, word := range v.words {
			g.AddVertex(word)
		}
	}

	for _, v := range d.definitions {
		for _, word := range v.words {
			// word defines name
			if word != v.name {
				g.AddEdge(word, v.name)
			}
		}
	}

}

func (d *Dictionary) expandDef(delNodes []string, k string) string {
	wordMap := make(map[string]bool)
	var defn []string

	if k == "" {
		return ""
	}

	for _, val := range d.definitions {
		wordMap[val.name] = false
		for _, word := range val.words {
			wordMap[word] = false
		}
	}

	for _, val := range delNodes {
		wordMap[val] = true
	}

	defn = d.findDef(k)

	var newDefn []string = []string{}
	for _, val := range defn {
		// get rid of self loops
		if k == val {
			newDefn = append(newDefn, val)
			continue
		}
		wmBool, ok := wordMap[val]
		if !ok || wmBool {
			newDefn = append(newDefn, val)
			continue
		}
		expand := d.recursiveSearch(wordMap, val)
		if len(expand) != 0 {
			newDefn = append(newDefn, expand...)
		} else {
			newDefn = append(newDefn, val)
		}
	}

	var str string

	for i, val := range newDefn {
		if i == 0 {
			str = str + val
		} else {
			str = str + " " + val
		}
	}

	return str
}

func (d *Dictionary) recursiveSearch(wordMap map[string]bool, k string) []string {
	val, ok := wordMap[k]
	if !ok || val {
		return []string{}
	} else {
		defn := d.findDef(k)
		var newDefn []string = []string{}
		for _, val := range defn {
			// get rid of self loops
			if k == val {
				newDefn = append(newDefn, val)
				continue
			}
			wmBool, ok := wordMap[val]
			if !ok || wmBool {
				newDefn = append(newDefn, val)
				continue
			}

			expand := d.recursiveSearch(wordMap, val)

			if len(expand) != 0 {
				newDefn = append(newDefn, expand...)
			} else {
				newDefn = append(newDefn, val)
			}

		}
		return newDefn
	}
}

// for use in recursiveSearch
func (d *Dictionary) findDef(k string) []string {
	defn, ok := d.definitions[k]
	if ok {
		return defn.words
	} else {
		return []string{}
	}
}

// for use in main()
func (d *Dictionary) getDef(k string) string {
	if k == "" {
		return ""
	}
	defn, ok := d.definitions[k]
	if ok {

		var str string

		for i, val := range defn.words {
			if i == 0 {
				str = str + val
			} else {
				str = str + " " + val
			}
		}

		return str
	} else {
		return ""
	}
}

// very slow implementation!
// implementation takes about 1hr40m on my computer to run (average hardware)
func (d *Dictionary) verify(delNodes []string) bool {

	fmt.Println("verifying...")

	for _, val := range d.definitions {
		d.expandDef(delNodes, val.name)
	}

	return true

}

type WNdict struct {
	IDMappings  map[string]*WNdef
	definitions map[string][]*WNdef
}

type WNdef struct {
	name       string
	origDef    string
	regexDef   string
	regexWords []string
	mappings   []string
}

func (wn *WNdict) Print() {
	for _, v := range wn.definitions {
		for _, d := range v {
			fmt.Println("name: ", d.name)
			fmt.Println("origDef: ", d.origDef)
		}
	}
}

func (wn *WNdict) PrintSize() {
	fmt.Println("\nsize : ", len(wn.definitions))
}

func (wn *WNdict) addDef(ID string, def *WNdef) {
	wn.IDMappings[ID] = def
	wn.definitions[def.name] = append(wn.definitions[def.name], def)
}

func (wn *WNdict) loadData(fn string) {
	bytes, err := os.ReadFile("wrangle/wordnet/" + fn) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println("isValid: ", json.Valid(bytes))

	var myData map[string][]interface{}

	json.Unmarshal(bytes, &myData)

	for k, v := range myData {
		ID := k
		name := v[0].(string)
		origDef := v[1].(string)
		regexDef := v[2].(string)

		regexWordsInterface := v[3].([]interface{})
		var regexWords []string
		for _, word := range regexWordsInterface {
			regexWords = append(regexWords, word.(string))
		}

		mappingsInterface := v[4].([]interface{})
		var mappings []string
		for _, word := range mappingsInterface {
			mappings = append(mappings, word.(string))
		}

		def := &WNdef{name: name, origDef: origDef, regexDef: regexDef, regexWords: regexWords, mappings: mappings}

		wn.addDef(ID, def)
	}
}

func (wn *WNdict) AddData(g *Graph) {
	fmt.Println("adding data to graph...")

	// add words
	for _, li := range wn.definitions {
		for _, v := range li {
			g.AddVertex(v.name)
			for _, word := range v.regexWords {
				g.AddVertex(word)
			}
		}
	}

	// add edges (has to happen once all words are in graph!)
	for _, li := range wn.definitions {
		for _, v := range li {
			for _, word := range v.regexWords {
				// word defines name
				if word != v.name {
					g.AddEdge(word, v.name)
				}
			}
		}
	}
}

func (wn *WNdict) expandDef(delNodes []string, k string) string {
	wordMap := make(map[string]bool)

	if k == "" {
		return ""
	}

	for _, val := range wn.IDMappings {
		wordMap[val.name] = false
		for _, word := range val.regexWords {
			wordMap[word] = false
		}
	}

	for _, val := range delNodes {
		wordMap[val] = true
	}

	defnArr := wn.findDefArr(k)

	var out string = ""

	for idx, defn := range defnArr {
		var str string = defn.regexDef
		for i, val := range defn.regexWords {
			if k == val {
				str = strings.Replace(str, "%s", val, 1)
				continue
			}
			expand := wn.recursiveSearch(wordMap, defn.mappings[i], val)
			if len(expand) != 0 {
				str = strings.Replace(str, "%s", expand, 1)
			} else {
				str = strings.Replace(str, "%s", val, 1)
			}
		}
		out = out + strconv.Itoa(idx+1) + ". " + str + "\n"
	}

	return out
}

func (wn *WNdict) recursiveSearch(wordMap map[string]bool, ID string, k string) string {
	val, ok := wordMap[k]
	if !ok || val {
		return ""
	} else {
		defn := wn.findDef(ID)
		var str string = defn.regexDef
		for i, val := range defn.regexWords {
			if k == val {
				str = strings.Replace(str, "%s", val, 1)
				continue
			}
			expand := wn.recursiveSearch(wordMap, defn.mappings[i], val)
			if len(expand) != 0 {
				str = strings.Replace(str, "%s", expand, 1)
			} else {
				str = strings.Replace(str, "%s", val, 1)
			}
		}
		return str
	}
}

func (wn *WNdict) findDef(ID string) *WNdef {
	defn, ok := wn.IDMappings[ID]
	if ok {
		return defn
	} else {
		return &WNdef{}
	}
}

// for use in recursiveSearch
func (wn *WNdict) findDefArr(k string) []*WNdef {
	defn, ok := wn.definitions[k]
	if ok {
		return defn
	} else {
		return []*WNdef{}
	}
}

// for use in main()
func (wn *WNdict) getDef(k string) string {
	var str string = ""

	val, ok := wn.definitions[k]
	if ok {
		for idx, def := range val {
			str = str + strconv.Itoa(idx+1) + ". " + def.origDef + "\n"
		}
	}

	return str
}

// very slow implementation! BUT TRUTHFUL!
func (wn *WNdict) verify(delNodes []string) bool {

	fmt.Println("verifying...")

	for _, defnArr := range wn.definitions {
		// expands all synsets anyway!
		wn.expandDef(delNodes, defnArr[0].name)
	}

	return true

}
