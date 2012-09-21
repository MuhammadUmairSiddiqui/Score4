package main

import (
	"fmt"
	"flag"
)

const (
	width      = 7
	height     = 6
	orangeWins = 1000000
	yellowWins = -orangeWins
)

type Mycell int

const (
	Barren Mycell = 0
	Orange Mycell = 1
	Yellow Mycell = -1
)

type Board struct {
	slots [height][width]Mycell
}

var negativeSlope = [4][2]int{{0, 0}, {1, 1}, {2, 2}, {3, 3}}
var positiveSlope = [4][2]int{{0, 0}, {-1, 1}, {-2, 2}, {-3, 3}}

func ScoreBoard(board *Board) int {
	var counters [9]int
	scores := board.slots

	// Horizontal spans
	for y := 0; y < height; y++ {
		score := scores[y][0] + scores[y][1] + scores[y][2]
		for x := 3; x < width; x++ {
			score += scores[y][x]
			counters[score+4]++
			score -= scores[y][x-3]
		}
	}
	// Vertical spans
	for x := 0; x < width; x++ {
		score := scores[0][x] + scores[1][x] + scores[2][x]
		for y := 3; y < height; y++ {
			score += scores[y][x]
			counters[score+4]++
			score -= scores[y-3][x]
		}
	}
	// Down-right (and up-left) diagonals
	for y := 0; y < height-3; y++ {
		for x := 0; x < width-3; x++ {
			score := 0
			for idx := 0; idx < 4; idx++ {
				yy := y + negativeSlope[idx][0]
				xx := x + negativeSlope[idx][1]
				score += int(scores[yy][xx])
			}
			counters[score+4]++
		}
	}
	// up-right (and down-left) diagonals
	for y := 3; y < height; y++ {
		for x := 0; x < width-3; x++ {
			score := 0
			for idx := 0; idx < 4; idx++ {
				yy := y + positiveSlope[idx][0]
				xx := x + positiveSlope[idx][1]
				score += int(scores[yy][xx])
			}
			counters[score+4]++
		}
	}
	if counters[0] != 0 {
		return yellowWins
	} else if counters[8] != 0 {
		return orangeWins
	}
	return counters[5] + 2*counters[6] + 5*counters[7] -
		counters[3] - 2*counters[2] - 5*counters[1]
}


func dropDisk(board *Board, column int, color Mycell) int {
	for y := height - 1; y >= 0; y-- {
		if board.slots[y][column] == Barren {
			board.slots[y][column] = color
			return y
		}
	}
	return -1
}

var debug = flag.Bool("debug", false, "")
var maxDepth = flag.Int("level", 7, "")

func loadBoard(argv []string) *Board {
	newBoard := new(Board)
	m := map[byte]Mycell{'o': Orange, 'y': Yellow}
	for _, s := range argv {
		if c, ok := m[s[0]]; ok {
			y, x := s[1]-'0', s[2]-'0'
			newBoard.slots[y][x] = c
		}
	}
	return newBoard
}

func abMinimax(maximize bool, color Mycell, depth int, board *Board) (move, score int) {
	if depth == 0 {
		return -1, ScoreBoard(board)
	}
	bestScore, bestMove := 10000000, 1
	winning := yellowWins
	if maximize {
		bestScore, winning = -bestScore, orangeWins
	}
	for column := 0; column < width; column++ {
		if board.slots[0][column] != Barren {
			continue
		}
		rowFilled := dropDisk(board, column, color)
		if rowFilled == -1 {
			continue
		}
		if s := ScoreBoard(board); s == winning {
			board.slots[rowFilled][column] = Barren
			return column, s
		}
		_, scoreInner := abMinimax(!maximize, -color, depth-1, board)
		board.slots[rowFilled][column] = Barren
		if depth == *maxDepth && *debug {
			fmt.Printf("Depth %d, placing on %d, score:%d\n", depth, column, scoreInner)
		}
		if maximize && scoreInner >= bestScore || !maximize && scoreInner <= bestScore {
			bestScore, bestMove = scoreInner, column
		}
	}
	return bestMove, bestScore
}

func main() {
	flag.Parse()
	board := loadBoard(flag.Args())
	score := ScoreBoard(board)
	if score == orangeWins {
		fmt.Println("I win")
		return
	} else if score == yellowWins {
		fmt.Println("You win")
		return
	} else {
		move, _ := abMinimax(true, Orange, *maxDepth, board)
		if move != -1 {
			fmt.Printf("%d\n", move)
			dropDisk(board, move, Orange)
			switch ScoreBoard(board) {
			case orangeWins:
				fmt.Println("I win")
			case yellowWins:
				fmt.Println("You win")
			}
		} else {
			fmt.Println("No move possible")
			return
		}
	}
}
