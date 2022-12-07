package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Content interface {
	isDirectory() bool
	getName() string
	getSize() int64
}

type File struct {
	name string
	size int64
}

func (f *File) isDirectory() bool {
	return false
}

func (f *File) getName() string {
	return f.name
}

func (f *File) getSize() int64 {
	return f.size
}

type Directory struct {
	name  string
	files []Content
}

func (d *Directory) isDirectory() bool {
	return true
}

func (d *Directory) getName() string {
	return d.name
}

func (d *Directory) getSize() int64 {
	var size int64

	for _, content := range d.files {
		isDirectory := content.isDirectory()
		if !isDirectory {
			// We are dealing with a file
			size += content.getSize()
		} else {
			childDirectory, ok := content.(*Directory)
			if !ok {
				log.Fatalf("failed to convert to directory")
			}
			size += childDirectory.getSize()
		}
	}

	return size
}

func parseFile(input string) (*File, error) {
	values := strings.Split(input, " ")

	fileSize, err := strconv.Atoi(values[0])
	if err != nil {
		return nil, err
	}

	file := &File{
		name: values[1],
		size: int64(fileSize),
	}

	return file, nil
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	var currentPaths []string
	var rootDirectory *Directory
	lookup := map[string]*Directory{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		value := scanner.Text()
		isCommand := strings.HasPrefix(value, "$")

		if isCommand {
			// trim the $ prefix
			command := value[2:]

			// Handle user command
			switch {
			case strings.HasPrefix(command, "cd"):
				// We are changing directory using the cd <path> command
				values := strings.Split(command, " ")
				path := values[1]

				if path == ".." {
					// We are going up a directory
					currentPaths = currentPaths[:len(currentPaths)-1]
					continue
				}

				// Handle going down a new directory
				directory := &Directory{
					name:  path,
					files: nil,
				}

				// Store the new path
				previousKey := strings.Join(currentPaths, "->")
				currentPaths = append(currentPaths, path)

				// Add the directory to the lookup map for quick lookups
				key := strings.Join(currentPaths, "->")
				lookup[key] = directory

				if rootDirectory == nil {
					// Set the root directory
					rootDirectory = directory
					continue
				}

				// Add the new directory to root of the previous path
				previousDirectory := lookup[previousKey]
				previousDirectory.files = append(previousDirectory.files, directory)
			case strings.HasPrefix(command, "ls"):
				// We don't need to do anything for the ls command
				continue
			default:
				log.Fatalf("unknown command: %v", command)
			}
		} else {
			// Handle command line output
			isDirectory := strings.HasPrefix(value, "dir")
			if isDirectory {
				continue
			}

			// We are dealing with a file
			file, err := parseFile(value)
			if err != nil {
				log.Fatalf("failed to parse file: %v", err)
			}

			key := strings.Join(currentPaths, "->")
			currentDirectory, ok := lookup[key]
			if !ok {
				log.Fatalf("failed to find current directory: %v", key)
			}

			currentDirectory.files = append(currentDirectory.files, file)
		}
	}

	total := 0
	for _, v := range lookup {
		size := v.getSize()

		if size >= 100000 {
			continue
		}

		total += int(size)
	}
	fmt.Printf("total size of all directories with a size of at most 100000: %v\n", total)

	totalSpaceAvailable := 70000000
	updateSize := 30000000
	unusedSpace := totalSpaceAvailable - int(rootDirectory.getSize())
	minimumSpaceNeededForUpdate := updateSize - unusedSpace

	var eligibleDirectories []*Directory
	for key, directory := range lookup {
		if key == "/" {
			// Skip the root directory
			continue
		}

		if int(directory.getSize()) < minimumSpaceNeededForUpdate {
			continue
		}

		eligibleDirectories = append(eligibleDirectories, directory)
	}

	// Sort from smallest to largest size
	sort.SliceStable(eligibleDirectories, func(i, j int) bool {
		return eligibleDirectories[i].getSize() < eligibleDirectories[j].getSize()
	})

	selectedDirectory := eligibleDirectories[0]

	fmt.Printf("Deleting directory %v which has size of %v\n", selectedDirectory.name, selectedDirectory.getSize())
}
