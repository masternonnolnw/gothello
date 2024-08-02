package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// create cell type BLANK | X | O
const (
	BLANK = iota
	X
	O
)

// create struct cell for each cell in the board
type Board struct {
	table [8][8]int
	turn  int
}

func getBoard(board Board) string {
	/* create string for board like
	   . . . . . . . .
	   . . . . . . . .
	   . . . . . . . .
	   . . . x o . . .
	   . . . o x . . .
	   . . . . . . . .
	   . . . . . . . .
	   . . . . . . . .
	*/
	boardStr := ""
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			switch board.table[i][j] {
			case BLANK:
				// return " i_j " to know the position of the cell
				boardStr += " " + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j) + " "
			case X:
				boardStr += " (X) "
			case O:
				boardStr += " (O) "
			}
		}
		boardStr += "\n"
	}
	return boardStr
}

func main() {
	// fiber instance
	app := fiber.New()

	directions := [8][2]int{
		{-1, -1}, // top left
		{-1, 0},  // top
		{-1, 1},  // top right
		{0, -1},  // left
		{0, 1},   // right
		{1, -1},  // bottom left
		{1, 0},   // bottom
		{1, 1},   // bottom right
	}

	board := Board{
		table: [8][8]int{
			{BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK},
			{BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK},
			{BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK},
			{BLANK, BLANK, BLANK, X, O, BLANK, BLANK, BLANK},
			{BLANK, BLANK, BLANK, O, X, BLANK, BLANK, BLANK},
			{BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK},
			{BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK},
			{BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK, BLANK},
		},
		turn: X,
	}

	// routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello world 🌈")
	})

	app.Get("/board", func(c *fiber.Ctx) error {
		return c.SendString(getBoard(board))
	})

	// make a move for player x or o and pass play position from body
	app.Post("/move/:player", func(c *fiber.Ctx) error {
		playerJson := c.Params("player")
		if playerJson != "x" && playerJson != "o" {
			return c.SendString("invalid player")
		}

		var player int
		if playerJson == "x" {
			player = X
		} else {
			player = O
		}

		// check if it's the player's turn
		if board.turn != player {
			return c.SendString("not your turn")
		}

		// get position from body
		position := new(struct {
			X int `json:"x"`
			Y int `json:"y"`
		})
		if err := c.BodyParser(position); err != nil {
			return c.SendString("invalid position")
		}

		// check if the position is valid
		if position.X < 0 || position.X > 7 || position.Y < 0 || position.Y > 7 {
			return c.SendString("out of board")
		}

		// check if the position is empty
		if board.table[position.X][position.Y] != BLANK {
			return c.SendString("cell is not empty")
		}

		// check position is possible to play
		// check each direction
		for _, direction := range directions {
			x, y := position.X, position.Y
			x += direction[0]
			y += direction[1]

			println(x, y, board.table[x][y])

			if board.table[x][y] == X {
				println("X")
			} else if board.table[x][y] == O {
				println("O")
			} else {
				println("Blank")
			}

			// check if the next cell is not in the board and has an opponent piece
			if x < 0 || x >= 8 || y < 0 || y >= 8 || board.table[x][y] == player || board.table[x][y] == BLANK {
				println("continue", board.table[x][y])
				continue
			}

			// continue to the next cell in the same direction
			for x >= 0 && x < 8 && y >= 0 && y < 8 {
				// move to the next cell
				x += direction[0]
				y += direction[1]

				println("- ", x, y)

				// check if the next cell is in the board and not empty
				if x < 0 || x >= 8 || y < 0 || y >= 8 || board.table[x][y] == BLANK {
					break
				}

				if board.table[x][y] == x {
					println("- X")
				}
				if board.table[x][y] == O {
					println("- O")
				}

				// if the next cell has the player's piece
				if board.table[x][y] == player {
					// flip the opponent's pieces
					for {
						x -= direction[0]
						y -= direction[1]

						board.table[x][y] = player

						if x == position.X && y == position.Y {
							break
						}
					}
					// switch the turn
					if player == X {
						board.turn = O
					} else {
						board.turn = X
					}
					return c.SendString(getBoard(board))
				}
			}
		}

		return c.SendString("can't move")

	})

	// app listening at PORT: 3000
	app.Listen(":3000")
}
