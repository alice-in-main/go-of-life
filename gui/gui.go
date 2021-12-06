package gui

import (
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"

	"go-of-life/logic"
)

func transform_x(x int) int {
	return 2 * x
}

func inverse_transform_x(x int) int {
	return x / 2
}

func transform_y(y int) int {
	return y
}

func inverse_transform_y(y int) int {
	return y
}

func drawBoard(offset_x int, offset_y int, s tcell.Screen, style tcell.Style, board *logic.Board) {
	x1, y1 := transform_x(0), transform_y(0)
	x2, y2 := transform_x(board.MaxX()), transform_y(board.MaxY())
	s.Clear()

	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
	}

	// Draw borders
	for col := x1; col <= x2; col++ {
		s.SetContent(col+offset_x, y1+offset_y, tcell.RuneHLine, nil, style)
		s.SetContent(col+offset_x, y2+offset_y, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1+offset_x, row+offset_y, tcell.RuneVLine, nil, style)
		s.SetContent(x2+offset_x, row+offset_y, tcell.RuneVLine, nil, style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1+offset_x, y1+offset_y, tcell.RuneULCorner, nil, style)
		s.SetContent(x2+offset_x, y1+offset_y, tcell.RuneURCorner, nil, style)
		s.SetContent(x1+offset_x, y2+offset_y, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2+offset_x, y2+offset_y, tcell.RuneLRCorner, nil, style)
	}

	squareStyle := tcell.StyleDefault.Background(tcell.ColorDeepSkyBlue).Foreground(tcell.Color103)
	for x := x1 + 1; x < x2-1; x++ {
		for y := y1 + 1; y < y2-1; y++ {
			if board.Get(inverse_transform_x(x), inverse_transform_y(y)) {
				s.SetContent(x+offset_x, y+offset_y, ' ', nil, squareStyle)
			}
		}
	}
}

type timedEvent struct{}

func (timedEvent) When() time.Time {
	return time.Now()
}

func Start() {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	s.EnableMouse()

	defStyle := tcell.StyleDefault.Background(tcell.ColorDarkRed).Foreground(tcell.Color103)
	s.SetStyle(defStyle)

	logic.InitGame()
	spawn_chan := make(chan func() (int, int))
	kill_chan := make(chan func() (int, int))
	board_request_curr_chan := make(chan interface{})
	board_request_next_chan := make(chan interface{})
	board_in_chan := make(chan logic.Board)
	clear_chan := make(chan interface{})
	go logic.GameLoop(spawn_chan, kill_chan, board_request_curr_chan, board_request_next_chan, board_in_chan, clear_chan)

	s.Clear()

	go func() {
		for {
			time.Sleep(50 * time.Millisecond)
			var tick timedEvent
			s.PostEvent(tick)
		}
	}()

	quit := func() {
		s.Fini()
		os.Exit(0)
	}
	defer quit()

	var run bool = true
	var offset_x int = 0
	var offset_y int = 0

	board_request_curr_chan <- 0
	initial_board := <-board_in_chan
	drawBoard(offset_x, offset_y, s, defStyle, &initial_board)

	for {
		s.Show()
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventMouse:
			x, y := ev.Position()
			switch ev.Buttons() {
			case tcell.Button1:
				spawn_chan <- func() (int, int) { return inverse_transform_x(x - offset_x), inverse_transform_y(y - offset_y) }
				board_request_curr_chan <- 0
				board := <-board_in_chan
				drawBoard(offset_x, offset_y, s, defStyle, &board)
			case tcell.Button2:
				kill_chan <- func() (int, int) { return inverse_transform_x(x - offset_x), inverse_transform_y(y - offset_y) }
				board_request_curr_chan <- 0
				board := <-board_in_chan
				drawBoard(offset_x, offset_y, s, defStyle, &board)
			}
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				quit()
			case tcell.KeyCtrlC:
				quit()
			case tcell.KeyRune:
				switch ev.Rune() {
				case ' ':
					run = !run
				case 's':
					offset_y--
				case 'w':
					offset_y++
				case 'a':
					offset_x--
				case 'd':
					offset_x++
				case 'c':
					clear_chan <- 0
				}
			}
		case timedEvent:
			if run {
				board_request_next_chan <- 0
				board := <-board_in_chan
				drawBoard(offset_x, offset_y, s, defStyle, &board)
			}
		}
	}
}
