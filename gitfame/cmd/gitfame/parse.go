package main

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type jsonLang struct {
	Name string   `json:"name"`
	Type string   `json:"type"`
	Ext  []string `json:"extensions"`
}

func readJSON(path string) (listLang []jsonLang, err error) {
	var jsonFile []byte
	jsonFile, err = os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonFile, &listLang)
	return listLang, err
}

func (gf *gitfame) parse() error {
	pflag.StringVar(&gf.path, "repository", ".", "path to repo")
	pflag.StringVar(&gf.commit, "revision", "HEAD", "pointer to commit")
	pflag.StringVar(&gf.order, "order-by", "lines", "result order by this parametr")
	pflag.BoolVar(&gf.useCommitter, "use-committer", false, "bool flag use committer in calc instead of author")
	pflag.StringVar(&gf.format, "format", "tabular", "format of output")
	pflag.StringSliceVar(&gf.exclude, "exclude", nil, "slice of glob-patterns that program will exclude")
	var (
		languages  []string
		extensions []string
	)
	pflag.StringSliceVar(&languages, "languages", nil, "slice of languages which files program will exclude")
	pflag.StringSliceVar(&extensions, "extensions", nil, "slice of files extensions that program will exclude")
	pflag.StringSliceVar(&gf.restricted, "restrict-to", nil, "slice of glob-patterns ...")
	pflag.Parse()

	gf.extensions = make(map[string]bool)

	for _, e := range extensions {
		gf.extensions[e] = true
	}

	mp := make(map[string]bool)
	for _, lang := range languages {
		mp[strings.ToLower(lang)] = true
	}

	listJSON, err := readJSON("../../configs/language_extensions.json")

	if err != nil {
		return err
	}

	for _, jj := range listJSON {
		if _, exists := mp[strings.ToLower(jj.Name)]; exists {
			for _, ext := range jj.Ext {
				gf.extensions[ext] = true
			}
		}
	}

	gf.result = make(map[string]stat)

	// fmt.Println(gf)

	if lookuptable[0][gf.order] == 0 {
		return errors.New("wrong order arg")
	}

	if lookuptable[1][gf.format] == 0 {
		return errors.New("wrong format arg")
	}

	return nil
}
