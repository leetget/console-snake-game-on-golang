package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/nsf/termbox-go"
)

const (
	width      int = 50
	height     int = 50
	max_length int = 50
)

var (
	length    int = 1
	head_x    int = width / 2
	head_y    int = height / 2
	snake_len int = 1

	apple_x      int
	apple_y      int
	array_x      [max_length]int
	array_y      [max_length]int
	sleep_time   int  = 100
	snake_head   rune = '@'
	snake_body   rune = '*'
	apple_object rune = 'o'
	wall         rune = '#'
	empty_obj    rune = ' '
	flag         bool = true
	dir          Direction
	done         chan struct{} // канал для завершения горутины
)

type Direction int

const (
	left Direction = iota
	right
	up
	down
)

func InitializeGame() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			var ch rune
			if x == 0 || x == width-1 || y == 0 || y == height-1 {
				ch = wall
			} else if x == head_x && y == head_y {
				ch = snake_head
			} else if x == apple_x && y == apple_y {
				ch = apple_object
			} else {
				ch = empty_obj
			}
			termbox.SetCell(x, y, ch, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	for i := 0; i < length; i++ {
		if i == 0 {
			continue // Пропускаем голову змеи
		}
		termbox.SetCell(array_x[i], array_y[i], snake_body, termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.Flush()
}

func UserInput(done chan struct{}) {
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		switch key {
		case keyboard.KeyArrowUp:
			dir = up
		case keyboard.KeyArrowDown:
			dir = down
		case keyboard.KeyArrowLeft:
			dir = left
		case keyboard.KeyArrowRight:
			dir = right
		case keyboard.KeyEsc:
			fmt.Println("Exit from game")
			done <- struct{}{} // отправляем сигнал завершения
			return
		case 'q':
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		}
	}
}

func MoveSnake() {
	for i := length; i > 0; i-- {
		array_x[i] = array_x[i-1]
		array_y[i] = array_y[i-1]
	}

	switch dir {
	case up:
		head_y--
	case down:
		head_y++
	case left:
		head_x--
	case right:
		head_x++
	}

	if head_x <= 0 || head_x >= width-1 || head_y <= 0 || head_y >= height-1 {
		fmt.Println("Game Over! You hit the wall.")
		flag = false
	}

	array_x[0] = head_x
	array_y[0] = head_y

	if head_x == apple_x && head_y == apple_y {
		length++
		spawnApple()
	}
}

func spawnApple() {
	for {
		newAppleX := rand.Intn(width-2) + 1
		newAppleY := rand.Intn(height-2) + 1

		if !isAppleOnSnake(newAppleX, newAppleY) {
			apple_x = newAppleX
			apple_y = newAppleY
			break
		}
	}
}

func isAppleOnSnake(x, y int) bool {
	for i := 0; i < length; i++ {
		if array_x[i] == x && array_y[i] == y {
			return true
		}
	}
	return false
}

func main() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	if err := keyboard.Open(); err != nil { // Инициализация клавиатуры
		panic(err)
	}
	defer keyboard.Close() // Закрытие клавиатуры при выходе

	dir = right                // Начальное направление
	done = make(chan struct{}) // инициализация канала

	go UserInput(done)

	spawnApple() // Первоначальное спавн яблока

	for flag {
		MoveSnake()
		InitializeGame()
		time.Sleep(time.Duration(sleep_time) * time.Millisecond)

		select {
		case <-done:
			flag = false // завершаем цикл игры при получении сигнала из канала
			break
		default:
			continue
		}
	}
}
