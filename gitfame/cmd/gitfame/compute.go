package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

type stat struct {
	lines   int
	file    map[string]int
	commits map[string]int
}

func checkListMatch(path string, patterns []string) bool {
	for _, p := range patterns {
		if matched, _ := filepath.Match(p, path); matched {
			return false
		}
	}
	return true
}

func (gf *gitfame) checkPathAssert(path string) bool {
	tmp := checkListMatch(path, gf.exclude)
	if !tmp {
		return false
	}
	if len(gf.restricted) > 0 {
		tmp := checkListMatch(path, gf.restricted)
		if tmp {
			return false
		}
	}
	return true
}

func (gf *gitfame) getListOfFiles() ([]string, error) {
	ls := exec.Command("git", "ls-tree", "--name-only", gf.commit, "-r")
	ls.Dir = gf.path
	result, err := ls.Output()
	if err != nil {
		return nil, err
	}
	sliceFiles := strings.Split(string(result), "\n")
	// fmt.Println(sliceFiles)
	listF := make([]string, 0, len(sliceFiles))
	for i := 0; i+1 < len(sliceFiles); i++ {
		if len(gf.extensions) == 0 || gf.extensions[filepath.Ext(sliceFiles[i])] {
			if gf.checkPathAssert(sliceFiles[i]) {
				listF = append(listF, sliceFiles[i])
			}
		}
	}
	return listF, nil
}

type blameSt struct {
	commit    string
	author    string
	committee string
	cntLines  int
}

func (gf *gitfame) addCnt(path string, blame blameSt) {
	name := blame.author
	if gf.useCommitter {
		name = blame.committee
	}
	val, exist := gf.result[name]

	if !exist {
		val = stat{0, make(map[string]int), make(map[string]int)}
	}

	val.commits[blame.commit] = 1
	val.file[path] = 1
	val.lines += blame.cntLines

	gf.result[name] = val
}

func (gf *gitfame) emptyFile(path string) error {
	cmd := exec.Command("git", "log", gf.commit, "-1", "--pretty=format:%an%n%cn%n%H", "--", path)
	cmd.Dir = gf.path
	result, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("git log error")
	}

	tmp := strings.SplitN(string(result), "\n", 3)

	if gf.useCommitter {
		val, exist := gf.result[tmp[1]]

		if !exist {
			val = stat{0, make(map[string]int), make(map[string]int)}
		}

		val.commits[tmp[2]] = 1
		val.file[path] = 1

		gf.result[tmp[1]] = val
		return nil
	}

	val, exist := gf.result[tmp[0]]

	if !exist {
		val = stat{0, make(map[string]int), make(map[string]int)}
	}

	val.commits[tmp[2]] = 1
	val.file[path] = 1

	gf.result[tmp[0]] = val
	return nil
}

func (gf *gitfame) getFileStat(path string) error {
	log.Println("file stat of " + path)
	blame := exec.Command("git", "blame", gf.commit, "--porcelain", path)
	blame.Dir = gf.path
	result, err := blame.Output()
	if err != nil {
		return err
	}
	if len(result) == 0 {
		err := gf.emptyFile(path)
		return err
	}
	mp := make(map[string]blameSt)
	scan := bufio.NewScanner(bytes.NewReader(result))
	lastHash := ""
	for scan.Scan() {
		line := scan.Text()
		// fmt.Println(line)
		if strings.HasPrefix(line, "\t") {
			val := mp[lastHash]
			val.cntLines++
			mp[lastHash] = val
		} else if len(line) >= 40 && !strings.Contains(line[:40], " ") {
			lastHash = line[:40]
			_, exist := mp[lastHash]
			if !exist {
				mp[lastHash] = blameSt{lastHash, "", "", 0}
			}
		} else {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				switch parts[0] {
				case "author":
					val := mp[lastHash]
					val.author = parts[1]
					mp[lastHash] = val
				case "committer":
					val := mp[lastHash]
					val.committee = parts[1]
					mp[lastHash] = val
				}
			}
		}
	}
	// fmt.Println(mp)
	for _, value := range mp {
		// fmt.Println(path, key, value)
		gf.addCnt(path, value)
	}
	return nil
}

func (gf *gitfame) compute() error {
	listFiles, err := gf.getListOfFiles()
	if err != nil {
		return err
	}
	for _, s := range listFiles {
		// fmt.Println(s)
		err := gf.getFileStat(s)
		if err != nil {
			return err
		}
	}
	return nil
}
