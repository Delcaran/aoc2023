package day3

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	color "github.com/fatih/color"
)

type gear_t struct {
	col     int
	row     int
	numbers []*part_number_t
	zone    [][]bool
}

func (g *gear_t) init(s model, row int, col int) {
	g.zone = make([][]bool, s.rows)
	for i := range g.zone {
		g.zone[i] = make([]bool, s.cols)
		for j := range g.zone[i] {
			g.zone[i][j] = false
		}
	}

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

func (g *gear_t) ratio() int {
	if len(g.numbers) != 2 {
		return 0
	}
	prod := 1
	for _, pn := range g.numbers {
		prod = prod * pn.value
	}
	return prod
}

func (g *gear_t) add_part_number(pn *part_number_t) bool {
	if len(g.numbers) < 2 {
		g.numbers = append(g.numbers, pn)
		return true
	}
	return false
}

func (g *gear_t) pn_in_range(pn *part_number_t) bool {
	for r, row := range g.zone {
		for c := range row {
			if row[c] && pn.row == r && pn.begin <= c && c <= pn.end {
				return true
			}
		}
	}
	return false
}

func (g *gear_t) in_range(nrow int, ncol int) bool {
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
	value int
	row   int
	begin int
	end   int
}

func (pn *part_number_t) in_range(nrow int, ncol int) bool {
	for c := pn.begin; c <= pn.end; c++ {
		if nrow == pn.row && ncol == c {
			return true
		}
	}
	return false
}

type model struct {
	schema       []string
	cols         int
	rows         int
	valid        [][]bool
	gears        []gear_t
	part_numbers []part_number_t
	refresh      int
}

func (m *model) find_symbols() {
	for row, line := range m.schema {
		for col, char := range line {
			if !unicode.IsDigit(char) && char != '.' {
				up := max(0, row-1)
				down := min(m.rows-1, row+1)
				left := max(0, col-1)
				right := min(m.cols-1, col+1)
				m.valid[up][left] = true
				m.valid[up][col] = true
				m.valid[up][right] = true
				m.valid[row][left] = true
				m.valid[row][col] = true
				m.valid[row][right] = true
				m.valid[down][left] = true
				m.valid[down][col] = true
				m.valid[down][right] = true

				if char == '*' {
					var new_g gear_t
					new_g.init(*m, row, col)
					m.gears = append(m.gears, new_g)
				}
			}
		}
	}
}

func (m *model) find_part_numbers() {
	m.part_numbers = m.part_numbers[:0]
	for row, line := range m.schema {
		tmp_num := []rune{}
		valid := false
		for col, char := range line {
			if unicode.IsDigit(char) {
				tmp_num = append(tmp_num, char)
				if m.valid[row][col] {
					valid = true
				}
			} else {
				// finished parsing number
				var pn part_number_t
				pn.row = row
				pn.end = col - 1
				pn.begin = col - len(tmp_num)
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

func (m *model) filter_gears() {
	var not_gears []int
	for idx := range m.gears {
		g := &m.gears[idx]
		for pni := range m.part_numbers {
			pn := &m.part_numbers[pni]
			if g.pn_in_range(pn) {
				if !g.add_part_number(pn) {
					// linked to too many part numbers
					not_gears = append(not_gears, idx)
				}
			}
		}
	}
	// second pass: not enough part numbers
	for idx, g := range m.gears {
		if len(g.numbers) != 2 {
			found := false
			for _, x := range not_gears {
				if x == idx {
					found = true
					break
				}
			}
			if !found {
				not_gears = append(not_gears, idx)
			}
		}
	}
	// third pass: removing not gears
	for _, index := range not_gears {
		if index+1 >= len(m.gears) {
			m.gears = m.gears[:index]
		} else {
			m.gears = append(m.gears[:index], m.gears[index+1:]...)
		}
	}
}

type TickMsg time.Time

func tickEvery() tea.Cmd {
	return tea.Every(time.Millisecond*30, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tickEvery()
}

func initialModel(content string) model {
	var m model
	m.schema = strings.Fields(strings.TrimSpace(content))
	m.rows = len(m.schema)
	m.cols = len(m.schema[0])
	m.valid = make([][]bool, m.rows)
	for i := range m.valid {
		m.valid[i] = make([]bool, m.cols)
		for j := range m.valid[i] {
			m.valid[i][j] = false
		}
	}
	return m
}

func (m *model) Run() {
	m.refresh += 1
	m.find_symbols()
	m.find_part_numbers()
	m.filter_gears()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case TickMsg:
		// Return your Every command again to loop.
		m.Run()
		return m, tickEvery()

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

func (m model) View() string {
	s := ""
	for r, line := range m.schema {
		for c, char := range line {
			relevance := 0 // 0: none, 1: valid, 2:gear range, 3: pn, 4: gear pn
			if unicode.IsDigit(char) {
				in_gear_range := false
				for _, g := range m.gears {
					for _, pn := range g.numbers {
						if pn.in_range(r, c) {
							// this "cell" is from a part number owned by a gear
							in_gear_range = true
							break
						}
					}
					if in_gear_range {
						break
					}
				}
				if in_gear_range {
					relevance = 4
				} else {
					for _, pn := range m.part_numbers {
						if pn.in_range(r, c) {
							// this "cell" is from a part number not owned by a gear
							relevance = 3
							break
						}
					}
				}
			} else {
				in_gear_range := false
				for _, g := range m.gears {
					if g.in_range(r, c) {
						in_gear_range = true
						break
					}
				}
				if in_gear_range {
					relevance = 2
				} else {
					if m.valid[r][c] {
						relevance = 1
					} else {
						relevance = 0
					}
				}
			}

			switch relevance {
			case 1:
				s += color.BlueString("%c", char)
			case 2:
				s += color.YellowString("%c", char)
			case 3:
				s += color.GreenString("%c", char)
			case 4:
				s += color.RedString("%c", char)
			default:
				s += fmt.Sprintf("%c", char)
			}
		}
		s += "\n"
	}
	s += "\n"
	s += fmt.Sprintf("Refresh: %d\n", m.refresh)
	s += fmt.Sprintf("Sum Part Number: %d\n", part1(m))
	s += fmt.Sprintf("Sum Gear Ratio: %d\n\n", part2(m))

	// Send the UI for rendering
	return s
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

func part1(m model) int {
	sum := 0
	for _, pn := range m.part_numbers {
		sum += pn.value
	}
	return sum
}

func part2(m model) int {
	sum := 0
	for _, g := range m.gears {
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
