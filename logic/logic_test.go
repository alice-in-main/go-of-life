package logic

import (
	"testing"

	"gotest.tools/assert"
)

func TestCreateBoard(t *testing.T) {
	x := 5
	y := 7
	board := createBoard(x, y)
	assert.Equal(t, x, board.MaxX())
	assert.Equal(t, y, board.MaxY())
	assert.Equal(t, x*y, board.size())
}

func TestSpawnAndKill(t *testing.T) {
	x := 5
	y := 7
	board := createBoard(x, y)
	board.SpawnCell(2, 3)
	assert.Equal(t, board.Get(2, 3), true)
	board.KillCell(2, 3)
	assert.Equal(t, board.Get(2, 3), false)
}

func TestCalcSingleLiveNeighbors(t *testing.T) {
	x := 10
	y := 10
	board := createBoard(x, y)
	board.SpawnCell(4, 4)
	assert.Equal(t, board.calcNumLiveNeighbors(4, 4), 0)
	assert.Equal(t, board.calcNumLiveNeighbors(3, 4), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(4, 3), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(5, 4), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(4, 5), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(3, 5), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(5, 3), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(3, 3), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(5, 5), 1)
}

func TestCalcSingleLiveNeighborsOnEdge(t *testing.T) {
	x := 10
	y := 10
	board := createBoard(x, y)
	board.SpawnCell(0, 4)
	assert.Equal(t, board.calcNumLiveNeighbors(0, 4), 0)
	assert.Equal(t, board.calcNumLiveNeighbors(0, 5), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(0, 3), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(1, 3), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(1, 4), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(1, 5), 1)

	board.SpawnCell(9, 4)
	assert.Equal(t, board.calcNumLiveNeighbors(9, 4), 0)
	assert.Equal(t, board.calcNumLiveNeighbors(9, 5), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(9, 3), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(8, 3), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(8, 4), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(8, 5), 1)
}

func TestCalcSingleLiveNeighborsInCorner(t *testing.T) {
	x := 10
	y := 10
	board := createBoard(x, y)
	board.SpawnCell(0, 0)
	assert.Equal(t, board.calcNumLiveNeighbors(0, 0), 0)
	assert.Equal(t, board.calcNumLiveNeighbors(0, 1), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(1, 0), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(1, 1), 1)

	board.SpawnCell(0, 9)
	assert.Equal(t, board.calcNumLiveNeighbors(0, 9), 0)
	assert.Equal(t, board.calcNumLiveNeighbors(0, 8), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(1, 9), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(1, 8), 1)

	board.SpawnCell(9, 9)
	assert.Equal(t, board.calcNumLiveNeighbors(9, 9), 0)
	assert.Equal(t, board.calcNumLiveNeighbors(8, 9), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(9, 8), 1)
	assert.Equal(t, board.calcNumLiveNeighbors(8, 8), 1)
}

func TestCalcLiveNeighborsComplexScenario(t *testing.T) {
	x := 10
	y := 10
	board := createBoard(x, y)
	board.SpawnCell(0, 0)
	board.SpawnCell(1, 2)
	board.SpawnCell(2, 1)
	board.SpawnCell(2, 2)
	assert.Equal(t, board.calcNumLiveNeighbors(1, 1), 4)
}

func TestCalcNextCellStateSpawn(t *testing.T) {
	x := 10
	y := 10
	board := createBoard(x, y)
	other_board := createBoard(x, y)
	board.SpawnCell(0, 0)
	board.SpawnCell(1, 2)
	board.SpawnCell(2, 1)
	assert.Equal(t, board.calcNumLiveNeighbors(1, 1), 3)
	assert.Equal(t, board.Get(1, 1), false)
	board.calcNextCellState(&other_board, 1, 1)
	assert.Equal(t, other_board.Get(1, 1), true)
}

func TestCalcNextCellStateStarve(t *testing.T) {
	x := 10
	y := 10
	board := createBoard(x, y)
	other_board := createBoard(x, y)
	board.SpawnCell(0, 0)
	assert.Equal(t, board.calcNumLiveNeighbors(0, 0), 0)
	assert.Equal(t, board.Get(0, 0), true)
	board.calcNextCellState(&other_board, 0, 0)
	assert.Equal(t, other_board.Get(0, 0), false)
}

func TestCalcNextCellStateOverpopDeath(t *testing.T) {
	x := 10
	y := 10
	board := createBoard(x, y)
	other_board := createBoard(x, y)
	board.SpawnCell(1, 1)
	board.SpawnCell(1, 0)
	board.SpawnCell(0, 1)
	board.SpawnCell(1, 2)
	board.SpawnCell(2, 2)
	assert.Equal(t, board.calcNumLiveNeighbors(1, 1), 4)
	assert.Equal(t, board.Get(1, 1), true)
	board.calcNextCellState(&other_board, 1, 1)
	assert.Equal(t, other_board.Get(1, 1), false)
}

func TestCalcNextBoardStateSingleCell(t *testing.T) {
	x := 10
	y := 10
	board := createBoard(x, y)
	other_board := createBoard(x, y)
	board.SpawnCell(1, 1)
	assert.Equal(t, board.calcNumLiveNeighbors(1, 1), 0)
	assert.Equal(t, board.Get(1, 1), true)
	board.calcNextBoardState(&other_board)
	assert.Equal(t, other_board.Get(1, 1), false)
}

func TestCalcNextBoardStateSomeLiveSomeDie(t *testing.T) {
	x := 10
	y := 10
	board := createBoard(x, y)
	other_board := createBoard(x, y)
	board.SpawnCell(0, 0)
	assert.Equal(t, board.Get(0, 0), true)

	board.SpawnCell(2, 1)
	board.SpawnCell(3, 1)
	board.SpawnCell(3, 0)
	board.SpawnCell(4, 1)
	board.SpawnCell(4, 0)
	assert.Equal(t, board.Get(3, 1), true)

	board.SpawnCell(1, 5)
	board.SpawnCell(0, 5)
	board.SpawnCell(2, 5)
	assert.Equal(t, board.Get(1, 5), true)

	board.SpawnCell(6, 5)
	board.SpawnCell(4, 5)
	board.SpawnCell(6, 6)
	assert.Equal(t, board.Get(5, 5), false)

	assert.Equal(t, board.calcNumLiveNeighbors(0, 0), 0)
	assert.Equal(t, board.calcNumLiveNeighbors(3, 1), 4)
	assert.Equal(t, board.calcNumLiveNeighbors(1, 5), 2)
	assert.Equal(t, board.calcNumLiveNeighbors(5, 5), 3)

	board.calcNextBoardState(&other_board)
	assert.Equal(t, other_board.Get(0, 0), false)
	assert.Equal(t, other_board.Get(3, 1), false)
	assert.Equal(t, other_board.Get(1, 5), true)
	assert.Equal(t, other_board.Get(5, 5), true)

	other_board.calcNextBoardState(&board)
	assert.Equal(t, board.Get(2, 1), true)
	assert.Equal(t, board.Get(3, 2), true)
	assert.Equal(t, board.Get(4, 1), true)
	assert.Equal(t, board.Get(3, 1), false)

	board.calcNextBoardState(&other_board)
	assert.Equal(t, other_board.Get(3, 1), true)
	assert.Equal(t, other_board.Get(3, 2), true)
	assert.Equal(t, other_board.Get(3, 3), false)
}
