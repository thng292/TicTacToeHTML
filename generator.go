package main

import (
	"fmt"
	"html/template"
	"math"
)

var boardState [9]int8

type State int8

// We play as X
const (
	Playing State = 1
	XWon    State = 'X'
	OWon    State = 'O'
	Draw    State = 3
	Empty   State = '_'
)

var board_template *template.Template = nil

type Move struct {
	index uint8
	score int
}

var result map[string]Move = make(map[string]Move)

func findBestMove(player int8, depth int) (index uint8, score int) {
	// serialized := serialize()
	// if move, ok := result[serialized]; ok {
	// return move.index, move.score
	// }
	state := getState()
	if state != Playing {
		if state == XWon {
			return 0, 10 - depth
		} else if state == OWon {
			return 0, depth - 10
		} else {
			return 0, 0
		}
	}
	if player == int8(XWon) { // X (Machine)
		var max_i uint8 = 0
		max_score := math.MinInt
		for i, cell := range boardState {
			if cell == int8(Empty) {
				boardState[i] = int8(XWon)
				_, score := findBestMove(int8(OWon), depth+1)
				if score > max_score {
					max_i = uint8(i)
					max_score = score
				}
				boardState[i] = int8(Empty)
			}
		}
		// result[serialized] = Move{index: max_i, score: max_score}
		return max_i, max_score
	} else { // O (Human)
		var min_i uint8 = 0
		min_score := math.MaxInt
		for i, cell := range boardState {
			if cell == int8(Empty) {
				boardState[i] = int8(OWon)
				_, score := findBestMove(int8(XWon), depth+1)
				if score < min_score {
					min_i = uint8(i)
					min_score = score
				}
				boardState[i] = int8(Empty)
			}
		}
		// result[serialized] = Move{index: min_i, score: min_score}
		return min_i, min_score
	}
}

func serialize() string {
	result := ""
	for _, cell := range boardState {
		switch cell {
		case int8(OWon):
			result += "O"
		case int8(XWon):
			result += "X"
		default:
			result += "_"
		}
	}
	return result
}

func produce() {
	for i, cell := range boardState {
		if cell == int8(Empty) {
			boardState[i] = int8(OWon)
			// current_board := serialize()
			// outf, err := os.Create("moves/" + current_board + ".html")
			// if err != nil {
			// 	panic(err)
			// }
			machineMove, _ := findBestMove(int8(XWon), 0)
			boardState[machineMove] = int8(XWon)
			// outf.WriteString(toHTML())
			// outf.Close()
			produce()
			boardState[i] = int8(Empty)
			boardState[machineMove] = int8(Empty)
		}
	}
}

func main() {
	for i := range boardState {
		boardState[i] = int8(Empty)
	}
	// boardState[0] = int8(OWon)
	// boardState[3] = int8(XWon)
	// boardState[4] = int8(OWon)
	// boardState[6] = int8(XWon)
	// boardState[7] = int8(OWon)
	// boardState[8] = int8(XWon)
	// print(findBestMove(int8(XWon), 0))
	boardState[4] = int8(OWon)
	current := XWon
	for getState() == Playing {
		move, _ := findBestMove(int8(current), 0)
		println(move)
		boardState[move] = int8(current)
		if current == XWon {
			current = OWon
		} else {
			current = XWon
		}
	}

	// produce()
	// count := 100
	// for k, v := range result {
	// 	fmt.Printf("%v: %v\n", k, v)
	// 	count -= 1
	// 	if count == 0 {
	// 		break
	// 	}
	// }
}

// getState determines the current game state based on the board
func getState() State {
	checkWin := func(board [9]int8, player int8) bool {
		// Check rows
		for i := 0; i < 3; i++ {
			if board[i*3] == player && board[i*3+1] == player && board[i*3+2] == player {
				return true
			}
		}

		// Check columns
		for i := 0; i < 3; i++ {
			if board[i] == player && board[i+3] == player && board[i+6] == player {
				return true
			}
		}

		// Check diagonals
		if board[0] == player && board[4] == player && board[8] == player {
			return true
		}
		if board[2] == player && board[4] == player && board[6] == player {
			return true
		}
		return false
	}

	allEmpty := func(board [9]int8) bool {
		for _, cell := range board {
			if cell != int8(Empty) {
				return false
			}
		}
		return true
	}
	// Check for win
	if checkWin(boardState, int8(XWon)) {
		return XWon
	}
	if checkWin(boardState, int8(OWon)) {
		return OWon
	}

	// Check for draw
	if allEmpty(boardState) {
		return Draw
	}

	// Game is still playing
	return Playing
}

func toHTML() string {
	result := `<title>Tic tac toe</title><meta name="viewport" content="width=device-width, initial-scale=3.0"><table>`
	footer := `</table>`
	// print(serialize())
	for i := 0; i < 3; i += 1 {
		result += `<tr>`
		for j := 0; j < 3; j += 1 {
			result += `<td>`
			if boardState[i*3+j] != int8(Empty) {
				result += string(boardState[i*3+j])
			} else if getState() == Playing {
				boardState[i*3+j] = int8(OWon)
				index, _ := findBestMove(int8(XWon), 0)
				boardState[index] = int8(XWon)
				result += fmt.Sprintf(`<a href='%v'>_</a>`, serialize()+".html")
				boardState[index] = int8(Empty)
				boardState[i*3+j] = int8(Empty)
			}
			result += `</td>`
		}
		result += `</tr>`
	}
	return result + footer
}
