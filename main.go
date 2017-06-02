package main

import (
	"bufio"
	"crypto/sha512"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var path = flag.String("path", "", "path to search for dedupe")
var files = make(map[[sha512.Size]byte]string)
var delete []string
var dir []string

func main() {
	flag.Parse()
	if *path == "" {
		log.Fatal("path must be provided")
	}

	files := scan(*path)
	walk(*path, files)

	if len(dir) > 0 {
		for _, d := range dir {
			files := scan(d)
			walk(d, files)
		}
	}
	if len(delete) > 0 {
		fmt.Println("Can remove the following files:")
		for _, d := range delete {
			fmt.Println(d)
		}

		prompt := bufio.NewReader(os.Stdin)
		fmt.Print("would you like to delete? [Y/N]: ")
		text, _ := prompt.ReadString('\n')
		if strings.Contains(strings.ToUpper(text), "Y") {
			deleteAll()
		}
	}

}

func scan(path string) []os.FileInfo {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalf("reading directory", err)
	}
	return files
}

func walk(root string, infos []os.FileInfo) {
	for _, f := range infos {
		checkDuplicate(root, f)
	}
}

func checkDuplicate(root string, f os.FileInfo) {
	path := filepath.Join(root, f.Name())
	fmt.Printf("scanning %s\n", path)
	if f.IsDir() {
		dir = append(dir, path)
		return
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("reading file ", err)
	}
	hash := sha512.Sum512(data)
	if _, ok := files[hash]; ok {
		delete = append(delete, path)
	} else {
		files[hash] = path
	}
}

func deleteAll() {
	for _, f := range delete {
		fmt.Println("deleting", f)
		err := os.Remove(f)
		if err != nil {
			log.Fatalf("delete", err)
		}
	}
}
