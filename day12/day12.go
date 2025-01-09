package day12

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	SPRING_UNKNOWN = iota
	SPRING_WORKING = iota
	SPRING_DAMAGED = iota
)

type fieldRow struct {
	springsMap     string
	damagedRecords string
}

type field struct {
	rows []fieldRow
}

func parseInput(content string) field {
	var f field
	for _, l := range strings.Split(content, "\n") {
		line := strings.TrimSpace(l)
		if len(line) > 0 {
			var fr fieldRow
			lineData := strings.Split(line, " ")
			fr.springsMap = lineData[0]
			fr.damagedRecords = lineData[1]
			f.rows = append(f.rows, fr)
		}
	}
	return f
}

func minify(input string) (minified string) {
	// remove consecutive dots
	for pos, char := range input {
		if pos == 0 {
			minified = string(char)
		} else {
			if char == '.' && minified[len(minified)-1] == '.' {
				continue
			} else {
				minified = minified + string(char)
			}
		}
	}
	return minified
}

func (f *field) unfold() {
	unfoldingSize := 5
	for rown := range f.rows {
		solvingMap := f.rows[rown].springsMap
		solvingRecords := f.rows[rown].damagedRecords
		unfoldedMap := make([]string, 0)
		unfoldedRecords := make([]string, 0)
		for x := 0; x < unfoldingSize; x++ {
			unfoldedMap = append(unfoldedMap, solvingMap)
			unfoldedRecords = append(unfoldedRecords, solvingRecords)
		}
		f.rows[rown].springsMap = strings.Join(unfoldedMap, "?")
		f.rows[rown].damagedRecords = strings.Join(unfoldedRecords, ",")
	}
}

func buildGroup(size int) string {
	var group string
	for x := 0; x < size; x++ {
		group += "#"
	}
	return group
}

func expandMap(basicMap string) map[string]int {
	expandedMaps := map[string]int{}
	prevch := ' '
	for idx, ch := range basicMap {
		expandedMap := ""
		switch idx {
		case 0:
			expandedMap = "." + basicMap
		case len(basicMap) - 1:
			expandedMap = basicMap + "."
		default:
			if prevch == '#' && ch == '.' {
				expandedMap = basicMap[:idx] + "." + basicMap[idx:]
			}
		}
		expandedMaps[expandedMap] = len(expandedMap)
		prevch = ch
	}
	return expandedMaps
}

func buildMaps(damageGroups []int, length int) map[string]int {
	groups := make([]string, 0)
	for _, size := range damageGroups {
		groups = append(groups, buildGroup(size))
	}
	basicMap := strings.Join(groups, ".")
	maps := map[string]int{basicMap: len(basicMap)}
	count := 0
	for currLen := len(basicMap); currLen != length; {
		count++
		fmt.Printf("Expansion %d\n", count)
		tmp := map[string]int{}
		for currMap := range maps {
			for expandedMap, lenMap := range expandMap(currMap) {
				currLen = lenMap
				tmp[expandedMap] = lenMap
			}
		}
		maps = tmp
	}
	return maps
}

func solveByBuild(rowMap string, damageRecords string) int {
	damageGroups := make([]int, 0)
	for _, g := range strings.Split(damageRecords, ",") {
		if conv, err := strconv.Atoi(g); err == nil {
			damageGroups = append(damageGroups, conv)
		}
	}
	minified := minify(rowMap)
	buildedMaps := buildMaps(damageGroups, len(minified))
	for buildedMap := range buildedMaps {
		if !mapMatch(minified, buildedMap) {
			delete(buildedMaps, buildedMap)
		}
	}
	fmt.Printf("Row %s has %d arrangements\n", rowMap, len(buildedMaps))
	return len(buildedMaps)
}

func mapMatch(original string, solved string) bool {
	if len(original) == len(solved) {
		for idx := 0; idx < len(original); idx++ {
			if original[idx] != '?' {
				if original[idx] != solved[idx] {
					return false
				}
			}
		}
		return true
	}
	return false
}

func Part1(f field) int {
	arrangements := 0
	for rown := range f.rows {
		arrangements += solveByBuild(f.rows[rown].springsMap, f.rows[rown].damagedRecords)
	}
	return arrangements
}

func Part2(f field) int {
	f.unfold()
	return Part1(f)
}

func Run(content string) (int, int) {
	fieldInfo := parseInput(content)
	//return Part1(fieldInfo), 525152
	return 21, Part2(fieldInfo)
	//return Part1(fieldInfo), Part2(fieldInfo)
}
