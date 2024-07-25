package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return printDir(out, path, printFiles, "")
}

func printDir(out io.Writer, path string, printFiles bool, prefix string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var dirs []os.DirEntry
	var files []os.DirEntry

	for _, entry := range entries {
		if entry.Name() == ".DS_Store" {
			continue
		}
		if entry.IsDir() {
			dirs = append(dirs, entry)
		} else if printFiles {
			files = append(files, entry)
		}
	}

	// Сортируем директории и файлы по имени
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Name() < dirs[j].Name()
	})

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	// Объединяем директории и файлы в один список
	entries = append(dirs, files...)

	for i, entry := range entries {
		isLast := i == len(entries)-1
		var newPrefix string
		if isLast {
			fmt.Fprintf(out, "%s└───%s", prefix, entry.Name())
			newPrefix = prefix + "\t"
		} else {
			fmt.Fprintf(out, "%s├───%s", prefix, entry.Name())
			newPrefix = prefix + "│\t"
		}

		if entry.IsDir() {
			fmt.Fprintln(out)
			err = printDir(out, path+string(os.PathSeparator)+entry.Name(), printFiles, newPrefix)
			if err != nil {
				return err
			}
		} else {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			if info.Size() == 0 {
				fmt.Fprintf(out, " (empty)\n")
			} else {
				fmt.Fprintf(out, " (%db)\n", info.Size())
			}
		}
	}

	return nil
}


