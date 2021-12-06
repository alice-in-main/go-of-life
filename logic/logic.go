package logic

import (
	"runtime"
)

type Board struct {
	data [][]bool
}

var activeBoard *Board = &Board{}
var workBoard *Board = &Board{}
var sizeY int = 100
var sizeX int = 100

func InitGame() {
	*activeBoard = createBoard(sizeX, sizeY)
	*workBoard = createBoard(sizeX, sizeY)

	activeBoard.SpawnCell(2, 6)
	activeBoard.SpawnCell(3, 6)
	activeBoard.SpawnCell(2, 7)
	activeBoard.SpawnCell(3, 7)

	activeBoard.SpawnCell(12, 6)
	activeBoard.SpawnCell(13, 5)
	activeBoard.SpawnCell(14, 4)
	activeBoard.SpawnCell(15, 4)
	activeBoard.SpawnCell(12, 7)
	activeBoard.SpawnCell(12, 8)
	activeBoard.SpawnCell(13, 9)
	activeBoard.SpawnCell(14, 10)
	activeBoard.SpawnCell(15, 10)

	activeBoard.SpawnCell(16, 7)

	activeBoard.SpawnCell(19, 7)
	activeBoard.SpawnCell(18, 6)
	activeBoard.SpawnCell(18, 7)
	activeBoard.SpawnCell(18, 8)
	activeBoard.SpawnCell(17, 5)
	activeBoard.SpawnCell(17, 9)

	activeBoard.SpawnCell(22, 4)
	activeBoard.SpawnCell(22, 5)
	activeBoard.SpawnCell(22, 6)
	activeBoard.SpawnCell(23, 4)
	activeBoard.SpawnCell(23, 5)
	activeBoard.SpawnCell(23, 6)
	activeBoard.SpawnCell(24, 3)
	activeBoard.SpawnCell(24, 7)

	activeBoard.SpawnCell(26, 2)
	activeBoard.SpawnCell(26, 3)

	activeBoard.SpawnCell(26, 7)
	activeBoard.SpawnCell(26, 8)

	activeBoard.SpawnCell(36, 5)
	activeBoard.SpawnCell(36, 4)
	activeBoard.SpawnCell(37, 5)
	activeBoard.SpawnCell(37, 4)
}

func GameLoop(
	spawn_chan chan func() (int, int),
	kill_chan chan func() (int, int),
	board_curr_chan chan interface{},
	board_next_chan chan interface{},
	board_out_chan chan Board,
	clear_chan chan interface{},
) {
	for {
		select {
		case coords_generator := <-spawn_chan:
			{
				x, y := coords_generator()
				activeBoard.SpawnCell(x, y)
			}
		case coords_generator := <-kill_chan:
			{
				x, y := coords_generator()
				activeBoard.KillCell(x, y)
			}
		case <-clear_chan:
			{
				*activeBoard = createBoard(sizeX, sizeY)
			}
		case <-board_curr_chan:
			{
				board_out_chan <- *activeBoard
			}
		case <-board_next_chan:
			{
				activeBoard.calcNextBoardState(workBoard)
				board_out_chan <- *workBoard
				activeBoard, workBoard = workBoard, activeBoard
			}
		}
	}
}

func (board *Board) Get(x int, y int) bool {
	if x < 0 || y < 0 || x >= board.MaxX() || y >= board.MaxY() {
		return false
	}
	return board.data[y][x]
}

func (board *Board) MaxY() int {
	return cap(board.data)
}

func (board *Board) MaxX() int {
	if cap(board.data) == 0 {
		return 0
	}
	return cap(board.data[0])
}

func (board *Board) size() int {
	return board.MaxX() * board.MaxY()
}

func createBoard(x_size int, y_size int) Board {
	new_board := Board{data: make([][]bool, y_size)}
	for i := 0; i < y_size; i++ {
		new_board.data[i] = make([]bool, x_size)
	}
	return new_board
}

func (board *Board) calcNumLiveNeighbors(x int, y int) int {
	counter := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			new_y := y + j
			new_x := x + i
			if new_x < 0 || new_x >= board.MaxX() || new_y < 0 || new_y >= board.MaxY() || (i == 0 && j == 0) {
				continue
			}
			if board.Get(new_x, new_y) {
				counter++
			}
		}
	}
	return counter
}

func (board *Board) calcNextCellState(next *Board, x int, y int) {
	num_live_neighbors := board.calcNumLiveNeighbors(x, y)
	switch {
	case num_live_neighbors < 2:
		next.setCell(false, x, y)
	case num_live_neighbors == 2:
		next.setCell(board.Get(x, y), x, y)
	case num_live_neighbors == 3:
		next.setCell(true, x, y)
	case num_live_neighbors > 3:
		next.setCell(false, x, y)
	}
}

func (board *Board) calcNextBoardState(next *Board) {
	numCpus := runtime.NumCPU()
	state_channel := make(chan interface{}, numCpus)
	for i := 0; i < numCpus; i++ {
		cpuIndex := i
		go func() {
			for y := cpuIndex; y < board.MaxY(); y += numCpus {
				for x := 0; x < board.MaxX(); x++ {
					board.calcNextCellState(next, x, y)
				}
			}
			state_channel <- 0
		}()
	}
	for x := 0; x < numCpus; x++ {
		<-state_channel
	}
}

func (board *Board) setCell(val bool, x int, y int) {
	if x < 0 || x >= board.MaxX() || y < 0 || y >= board.MaxY() {
		return
	}
	board.data[y][x] = val
}

func (board *Board) SpawnCell(x int, y int) {
	board.setCell(true, x, y)
}

func (board *Board) KillCell(x int, y int) {
	board.setCell(false, x, y)
}
