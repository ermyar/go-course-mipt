package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
)

func (gf *gitfame) getResult() (result [][]string) {
	order := []int{0, 1, 2, 3}

	if lookuptable[0][gf.order] == 2 {
		order = []int{0, 2, 1, 3}
	}
	if lookuptable[0][gf.order] == 3 {
		order = []int{0, 3, 1, 2}
	}

	// fmt.Println(order)

	type statSt struct {
		name string
		val  []int
	}

	data := make([]statSt, 0, len(gf.result))
	for key, value := range gf.result {
		data = append(data, statSt{name: key, val: []int{value.lines, len(value.commits), len(value.file)}})
	}
	sort.Slice(data, func(i, j int) bool {
		tmp := data[i].val
		orderedI := []int{tmp[order[1]-1], tmp[order[2]-1], tmp[order[3]-1]}
		tmp = data[j].val
		orderedJ := []int{tmp[order[1]-1], tmp[order[2]-1], tmp[order[3]-1]}
		return orderedI[0] > orderedJ[0] || (orderedI[0] == orderedJ[0] && orderedI[1] > orderedJ[1]) || (orderedI[0] == orderedJ[0] && orderedI[1] == orderedJ[1] && orderedI[2] > orderedJ[2]) ||
			(orderedI[0] == orderedJ[0] && orderedI[1] == orderedJ[1] && orderedI[2] == orderedJ[2] && data[i].name < data[j].name)
	})

	for _, d := range data {
		result = append(result, []string{d.name, fmt.Sprint(d.val[0]), fmt.Sprint(d.val[1]), fmt.Sprint(d.val[2])})
	}

	return result
}

func (gf *gitfame) tabular() {
	w := tabwriter.Writer{}
	w.Init(os.Stdout, 1, 0, 1, ' ', 0)

	fmt.Fprintln(&w, "Name\tLines\tCommits\tFiles")
	result := gf.getResult()
	for _, rec := range result {
		fmt.Fprintln(&w, strings.Join(rec, "\t"))
	}

	w.Flush()
}

func (gf *gitfame) csv() {
	w := csv.NewWriter(os.Stdout)

	result := make([][]string, 0)
	result = append(result, []string{"Name", "Lines", "Commits", "Files"})
	result = append(result, gf.getResult()...)
	for _, rec := range result {
		// fmt.Println(rec)
		if err := w.Write(rec); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

type jsonFormat struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}

func (gf *gitfame) json() {
	result := gf.getResult()
	jsonResult := make([]jsonFormat, 0)
	for _, rec := range result {
		lines, _ := strconv.Atoi(rec[1])
		commits, _ := strconv.Atoi(rec[2])
		files, _ := strconv.Atoi(rec[3])
		jsonResult = append(jsonResult, jsonFormat{rec[0], lines, commits, files})
	}
	b, _ := json.Marshal(jsonResult)
	os.Stdout.Write(b)
}

func (gf *gitfame) jsonLine() {
	result := gf.getResult()
	for _, rec := range result {
		lines, _ := strconv.Atoi(rec[1])
		commits, _ := strconv.Atoi(rec[2])
		files, _ := strconv.Atoi(rec[3])
		jsonResult := jsonFormat{rec[0], lines, commits, files}
		b, _ := json.Marshal(jsonResult)
		os.Stdout.Write(b)
		os.Stdout.Write([]byte("\n"))
	}
}

func (gf *gitfame) output() error {
	// if gf.format == "tabular"
	ff := []func(){nil, gf.tabular, gf.csv, gf.json, gf.jsonLine}
	ff[lookuptable[1][gf.format]]()
	return nil
}
