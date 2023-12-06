package day3

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	color "github.com/fatih/color"
)

type symbol_t struct {
	char         rune
	col          int
	row          int
	part_numbers []*part_number_t
	zone         [][]bool
}

func (g *symbol_t) is_gear() bool {
	return g.char == '*' && len(g.part_numbers) == 2
}

func (g *symbol_t) init(s *schematic, c rune, row int, col int) {
	g.zone = make([][]bool, s.rows)
	for i := range g.zone {
		g.zone[i] = make([]bool, s.cols)
		for j := range g.zone[i] {
			g.zone[i][j] = false
		}
	}

	g.char = c
	g.col = col
	g.row = row

	up := max(0, row-1)
	down := min(s.rows-1, row+1)
	left := max(0, col-1)
	right := min(s.cols-1, col+1)

	g.zone[up][left] = true
	g.zone[up][col] = true
	g.zone[up][right] = true
	g.zone[row][left] = true
	g.zone[row][col] = true
	g.zone[row][right] = true
	g.zone[down][left] = true
	g.zone[down][col] = true
	g.zone[down][right] = true
}

func (g *symbol_t) ratio() int {
	if g.is_gear() {
		prod := 1
		for _, pn := range g.part_numbers {
			prod = prod * pn.value
		}
		return prod
	}
	return 0
}

func (g *symbol_t) add_part_number(pn *part_number_t) {
	g.part_numbers = append(g.part_numbers, pn)
}

func (g *symbol_t) pn_in_range(pn *part_number_t) bool {
	for r, row := range g.zone {
		for c := range row {
			if row[c] && pn.row == r && pn.begin <= c && c <= pn.end {
				return true
			}
		}
	}
	return false
}

func (g *symbol_t) in_range(nrow int, ncol int) bool {
	for r, row := range g.zone {
		for c := range row {
			if row[c] && nrow == r && ncol == c {
				return true
			}
		}
	}
	return false
}

type part_number_t struct {
	value  int
	row    int
	begin  int
	end    int
	lenght int
}

func (pn *part_number_t) in_range(nrow int, ncol int) bool {
	for c := pn.begin; c <= pn.end; c++ {
		if nrow == pn.row && ncol == c {
			return true
		}
	}
	return false
}

type schematic struct {
	schema       []string
	cols         int
	rows         int
	symbols      []symbol_t
	part_numbers []part_number_t
}

func (m *schematic) get_symbols_on_cell(row int, col int) []*symbol_t {
	var symbols []*symbol_t = nil
	for idx := range m.symbols {
		s := &m.symbols[idx]
		if s.in_range(row, col) {
			symbols = append(symbols, s)
		}
	}
	return symbols
}

func (m *schematic) get_symbols_on_part_number(pn *part_number_t) []*symbol_t {
	var symbols []*symbol_t = nil
	for idx := range m.symbols {
		s := &m.symbols[idx]
		for c := pn.begin; c <= pn.end; c++ {
			if s.in_range(pn.row, c) {
				symbols = append(symbols, s)
			}
		}
	}
	return symbols
}

func (m *schematic) find_symbols() {
	for row, line := range m.schema {
		for col, char := range line {
			if !unicode.IsDigit(char) && char != '.' {
				var s symbol_t
				s.init(m, char, row, col)
				m.symbols = append(m.symbols, s)
			}
		}
	}
}

func (m *schematic) find_part_numbers() {
	for row, line := range m.schema {
		tmp_num := []rune{}
		valid := false
		for col, char := range line {
			if unicode.IsDigit(char) {
				tmp_num = append(tmp_num, char)
				if m.get_symbols_on_cell(row, col) != nil {
					valid = true
				}
			} else {
				// finished parsing number
				var pn part_number_t
				pn.row = row
				pn.end = col - 1
				pn.begin = col - len(tmp_num)
				pn.lenght = len(tmp_num)
				pn.value = check_num(string(tmp_num), valid)
				if pn.value != 0 {
					m.part_numbers = append(m.part_numbers, pn)
				}
				valid = false
				tmp_num = []rune{}
			}
		}
		// check right border numbers
		var pn part_number_t
		pn.value = check_num(string(tmp_num), valid)
		if pn.value != 0 {
			m.part_numbers = append(m.part_numbers, pn)
		}
		valid = false
	}
}

func (m *schematic) filter_gears() {
	for idx := range m.symbols {
		g := &m.symbols[idx]
		for pni := range m.part_numbers {
			pn := &m.part_numbers[pni]
			if g.pn_in_range(pn) {
				g.add_part_number(pn)
			}
		}
	}
}

func (m schematic) Init() tea.Cmd {
	return nil
}

func (m *schematic) Run() {
	m.find_symbols()
	m.find_part_numbers()
	m.filter_gears()
}

func (m schematic) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

const (
	NONE             = iota
	VALID            = iota
	GEAR_RANGE       = iota
	PART_NUMBER      = iota
	GEAR_PART_NUMBER = iota
)

func (m schematic) View() string {
	s := ""
	for r, line := range m.schema {
		for c, char := range line {
			relevance := NONE

			if unicode.IsDigit(char) {
				for _, pn := range m.part_numbers {
					if pn.in_range(r, c) {
						relevance = PART_NUMBER
						for _, sym := range m.get_symbols_on_part_number(&pn) {
							if sym.is_gear() {
								relevance = GEAR_PART_NUMBER
								break
							}
						}
						break
					} else {
						relevance = NONE
					}
				}
			} else {
				symbols := m.get_symbols_on_cell(r, c)
				if symbols == nil {
					relevance = NONE
				} else {
					for _, sym := range symbols {
						if sym.is_gear() {
							relevance = GEAR_RANGE
						} else {
							relevance = VALID
						}
						break
					}
				}
			}

			switch relevance {
			case NONE:
				s += fmt.Sprintf("%c", char)
			case VALID:
				s += color.BlueString("%c", char)
			case GEAR_RANGE:
				s += color.YellowString("%c", char)
			case PART_NUMBER:
				s += color.GreenString("%c", char)
			case GEAR_PART_NUMBER:
				s += color.RedString("%c", char)
			default:
				s += fmt.Sprintf("%c", char)
			}
		}
		s += "\n"
	}
	s += "\n"
	s += fmt.Sprintf("Sum Part Number: %d\n", part1(m))
	s += fmt.Sprintf("Sum Gear Ratio: %d\n\n", part2(m))

	// Send the UI for rendering
	return s
}

func initialModel(content string) schematic {
	var m schematic
	m.schema = strings.Fields(strings.TrimSpace(content))
	m.rows = len(m.schema)
	m.cols = len(m.schema[0])
	m.Run()
	return m
}

func check_num(num_str string, valid bool) int {
	if len(num_str) > 0 {
		//fmt.Printf("%s ", num_str)
		if valid {
			num, err := strconv.Atoi(num_str)
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Printf("part number\n")
			return num
		} else {
			//fmt.Printf("not a part number\n")
			return 0
		}
	}
	return 0
}

func part1(m schematic) int {
	sum := 0
	for _, pn := range m.part_numbers {
		sum += pn.value
	}
	return sum
}

func part2(m schematic) int {
	sum := 0
	for _, g := range m.symbols {
		ratio := g.ratio()
		sum += ratio
	}
	return sum
}

func Run(test string) {
	if len(test) > 0 {
		test += "_"
	}
	buffer, err := os.ReadFile("day3/" + test + "input.txt")
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	content := string(buffer[:])
	p := tea.NewProgram(initialModel(content))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	//fmt.Printf("Part 1: %d\n", part1(&schema))
	//fmt.Printf("Part 1: %d\n", part2(&schema))
}
