// Copyright (c) 2006-2019, xiaobo
//
// This is free software, licensed under the GNU General Public License v3.
// See /LICENSE for more information.
//

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Conf struct {
	Dir                string
	Template           string
	Adapter            []string
	Excludes           []string
	UseDefaultExcludes bool
	Properties         map[string]string
}

var (
	conf      = &Conf{}
	rawHeader *RawHeader
	MapFiles  map[string][]string
)

func main() {
	check := flag.Bool("check", false, "check all file for ci")
	flag.Parse()

	conf = setupConfig()
	dir := conf.Dir

	MapFiles = make(map[string][]string)
	if err := tree(dir, 0); err != nil {
		fmt.Println(err)
	}

	for ext, files := range MapFiles {
		handler := GetHandler(ext)
		if nil != handler {
			//fmt.Println(files)
			processHeader(files, handler, *check)
		}
	}
}

func processHeader(files []string, handler HeaderHandler, check bool) {
	for _, file := range files {
		buf, err := ioutil.ReadFile(file)
		checkErr("Reads file error", err)
		originalContent := string(buf)
		header := handler.Execute(rawHeader)

		action := getAction(originalContent, header)
		var content string

		switch action {
		case "no":
			defer log.Printf("Don't need to be updated [%s]", file)
		case "add":
			if !check {
				content = header + "\n" + originalContent
				defer log.Printf("Added header to file [%s]", file)
				wirteFile(file, content)
			} else {
				os.Exit(1)
			}
		case "update":
			if !check {
				headerLines := strings.Split(header, "\n")
				originalContentLines := strings.Split(originalContent, "\n")
				contentLines := originalContentLines[len(headerLines):]
				headerLines = append(headerLines, contentLines...)
				content = strings.Join(headerLines, "\n")
				defer log.Printf("Updated header to file [%s]", file)
				wirteFile(file, content)
			} else {
				os.Exit(1)
			}
		default:
			log.Fatalf("Wrong action [%s]", action)
		}
		//fmt.Printf(content)
	}
}

func wirteFile(file, content string) {
	fout, err := os.Create(file)
	checkErr("Prepare to write file err: ", err)

	fout.WriteString(content)
	err = fout.Close()
	checkErr("Write file err: ", err)
}

func setupConfig() *Conf {
	cfg, err := ioutil.ReadFile(".header.cfg")
	checkErr("Loads configuration error", err)

	err = json.Unmarshal(cfg, conf)
	checkErr("Parses configuration error", err)

	dirPath, err := filepath.Abs(conf.Dir)
	checkErr("Reads dir path error", err)

	t, err := template.ParseFiles(conf.Template)
	checkErr("Can't find header template ["+conf.Template+"]", err)

	var buf bytes.Buffer
	t.Execute(&buf, conf.Properties)
	rawHeader = NewRawHeader(buf.String())

	// Excludes
	if conf.UseDefaultExcludes {
		conf.Excludes = append(conf.Excludes, DefaultExcludes...)
	}
	for _, exclude := range conf.Excludes {
		path := filepath.Join(dirPath, exclude)
		path = filepath.ToSlash(path)
		conf.Excludes = append(conf.Excludes, path)
	}
	return conf
}

func tree(dstPath string, level int) error {
	dstF, err := os.Open(dstPath)
	if err != nil {
		return err
	}
	defer dstF.Close()
	match := match(dstF.Name())
	fileInfo, err := dstF.Stat()
	if err != nil {
		return err
	}
	if !fileInfo.IsDir() { //if dstF is file
		if !match {
			ext := filepath.Ext(dstF.Name())
			ok := strings.HasSuffix(fileInfo.Name(), ext)
			if ok {
				var files []string
				if _, exist := MapFiles[ext]; !exist {
					files = append(files, dstPath)
					MapFiles[ext] = files
				} else {
					files = MapFiles[ext]
					files = append(files, dstPath)
					MapFiles[ext] = files
				}
			}
		}
		return nil
	} else { //if dstF is dir
		if !match {
			dir, err := dstF.Readdir(0) //Gets fileInfo for each file or folder under a folder
			if err != nil {
				return err
			}
			for _, fileInfo = range dir {
				err = tree(dstPath+"/"+fileInfo.Name(), level+1)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
}

// getAction returns a handle action for the specified original content and header.
//
//  1. "add" means need to add the header to the original content
//  2. "update" means need to update (replace) the header of the original content
//  3. "no" means nothing need to do
func getAction(originalContent, header string) string {
	headerLines := strings.Split(header, "\n")
	originalContentLines := strings.Split(originalContent, "\n")

	if len(headerLines) > len(originalContentLines) {
		return "add"
	}

	originalHeaderLines := originalContentLines[:len(headerLines)]

	result := similar(originalHeaderLines, headerLines)
	if 100 <= result {
		return "no"
	}

	if result >= 70 {
		return "update"
	}

	return "add"
}

// [0, 100]
//
//  0: not similar at all
//  100: as the same
func similar(lines1, lines2 []string) int {
	if len(lines1) != len(lines2) {
		return 0
	}

	length := len(lines1)
	same := 0
	for i := 0; i < length; i++ {
		l1 := strings.TrimSpace(lines1[i])
		l2 := strings.TrimSpace(lines2[i])

		if l1 == l2 {
			same++
		}
	}

	return int(math.Floor(float64(same) / float64(length) * 100))
}

func match(path string) bool {
	path = filepath.ToSlash(path)
	for _, exclude := range conf.Excludes {
		match, err := filepath.Match(exclude, path)
		//fmt.Println("match %s == path %s", exclude, path)
		checkErr("Exclude path match error", err)
		if match {
			return true
		}
	}
	return false
}

func checkErr(errMsg string, err error) {
	if nil != err {
		log.Fatal(errMsg+", caused by: ", err)
	}
}
