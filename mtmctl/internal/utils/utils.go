package utils

import (
	"fmt"
	"sort"
)

func PrintSortedMap(mapToPrint map[string]string) {
	names := make([]string, 0, len(mapToPrint))
	for name := range mapToPrint {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		fmt.Printf("%s: %s\n", name, mapToPrint[name])
	}
}
