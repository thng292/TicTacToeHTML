package main

import (
	"fmt"
	"math"
	"os"
)

type State int8

// We play as X
const (
	Playing  State = 1
	XMachine State = 'X'
	OPlayer  State = 'O'
	Draw     State = 3
	Empty    State = '_'
)

var boardState [9]State

// var board_template *template.Template = nil

type Move struct {
	index uint8
	score int
}

var result map[string]Move = make(map[string]Move)
var done map[string]bool = make(map[string]bool)

func findBestMove(player State, depth int) (index uint8, score int) {
	serialized := serialize()
	if move, ok := result[serialized]; ok {
		return move.index, move.score
	}
	state := getState()
	if state != Playing {
		if state == XMachine {
			return 0, 10 - depth
		} else if state == OPlayer {
			return 0, depth - 10
		} else {
			return 0, 0
		}
	}
	if player == XMachine { // X (Machine)
		var max_i uint8 = 0
		max_score := math.MinInt
		for i, cell := range boardState {
			if cell == Empty {
				boardState[i] = XMachine
				_, score := findBestMove(OPlayer, depth+1)
				if score > max_score {
					max_i = uint8(i)
					max_score = score
				}
				boardState[i] = Empty
			}
		}
		result[serialized] = Move{index: max_i, score: max_score}
		return max_i, max_score
	} else { // O (Human)
		var min_i uint8 = 0
		min_score := math.MaxInt
		for i, cell := range boardState {
			if cell == Empty {
				boardState[i] = OPlayer
				_, score := findBestMove(XMachine, depth+1)
				if score < min_score {
					min_i = uint8(i)
					min_score = score
				}
				boardState[i] = Empty
			}
		}
		result[serialized] = Move{index: min_i, score: min_score}
		return min_i, min_score
	}
}

func serialize() string {
	result := ""
	for _, cell := range boardState {
		switch cell {
		case OPlayer:
			result += "O"
		case XMachine:
			result += "X"
		default:
			result += "_"
		}
	}
	return result
}

func writeFile() {
	current_board := serialize()
	if _, ok := done[current_board]; ok {
		return
	}
	// fmt.Printf("Doing %v\n", current_board)
	done[current_board] = true
	outf, err := os.Create("moves/" + current_board + ".html")
	if err != nil {
		panic(err)
	}
	// fmt.Printf("Before write: %v\n", serialize())
	outf.WriteString(toHTML())
	// fmt.Printf("After write: %v\n", serialize())
	outf.Close()
}

func produce() {
	// print("Start==================\n")
	// defer print("End====================\n")
	now := serialize()
	fmt.Printf("Now: %v\n", now)
	if getState() != Playing {
		return
	}
	for i, cell := range boardState {
		if cell == Empty {
			boardState[i] = OPlayer
			machineMove, _ := findBestMove(XMachine, 0)
			if boardState[machineMove] == Empty {
				boardState[machineMove] = XMachine
			} else {
				machineMove = math.MaxInt8
			}
			// =================
			writeFile()
			if getState() == Playing {
				produce()
				// fmt.Printf("Backed to: %v\n", now)
			}
			// =================
			if machineMove != math.MaxInt8 {
				boardState[machineMove] = Empty
			}
			boardState[i] = Empty
		}
	}
}

func main() {
	for i := range boardState {
		boardState[i] = Empty
	}

	// boardState[0] = OPlayer
	// boardState[1] = OPlayer
	// boardState[2] = XMachine

	// boardState[4] = XMachine
	// boardState[6] = OPlayer

	// println(findBestMove(XMachine, 0))

	outf, err := os.Create("moves/index.html")
	if err != nil {
		panic(err)
	}
	outf.WriteString(toHTML())
	outf.Close()
	produce()

}

// getState determines the current game state based on the board
func getState() State {
	checkWin := func(board [9]State, player State) bool {
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

	allPlayed := func(board [9]State) bool {
		for _, cell := range board {
			if cell == Empty {
				return false
			}
		}
		return true
	}
	// Check for win
	if checkWin(boardState, XMachine) {
		return XMachine
	}
	if checkWin(boardState, OPlayer) {
		return OPlayer
	}

	// Check for draw
	if allPlayed(boardState) {
		return Draw
	}

	// Game is still playing
	return Playing
}

func toHTML() string {
	result := `<title>Tic-tac-toe HTML</title><meta name="viewport" content="width=device-width, initial-scale=30.0"><table>`
	footer := `</table>`
	// print(serialize())
	current_state := getState()
	for i := 0; i < 3; i += 1 {
		result += `<tr>`
		for j := 0; j < 3; j += 1 {
			result += `<td>`
			if boardState[i*3+j] != Empty {
				result += string(boardState[i*3+j])
			} else if current_state == Playing {
				boardState[i*3+j] = OPlayer
				index, _ := findBestMove(XMachine, 0)
				if boardState[index] == Empty {
					boardState[index] = XMachine
				} else {
					index = math.MaxInt8
				}
				result += fmt.Sprintf(`<a href='%v'>_</a>`, serialize()+".html")
				if index != math.MaxInt8 {
					boardState[index] = Empty
				}
				boardState[i*3+j] = Empty
			}
			result += `</td>`
		}
		result += `</tr>`
	}
	result += footer
	if current_state == XMachine {
		result += "You lose"
		result += `<a href='index.html'>Restart</a>`
	}
	if current_state == Draw {
		result += "Draw"
		result += `<a href='index.html'>Restart</a>`
	}
	if current_state == OPlayer {
		result += "You win"
		result += `<a href='index.html'>Restart</a>`
	}
	return result
}
