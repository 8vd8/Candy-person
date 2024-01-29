package main

import (
	simpleansi "awesomeProject2/github"
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

//我们使用包Open中的函数来打开它，并使用缓冲IO包中的扫描器对象(bufio)
//将其读入到内存（读取到名为maze的全局变量）
//最后我们需要通过调用来释放文件处理程序

var maze []string //迷宫的格式
var player sprite //玩家
// sprite:单色画面
// 2D坐标
type sprite struct {
	row int
	col int
}

func loadMaze(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	/*
		Scan()用于触发扫描下一个标记，将其存储在Text()的缓冲区中。
		Text()用于获取上一次扫描到的标记的内容。
		通常，它们一起使用，以便在扫描的同时获取标记的内容。
		在循环中，Scan()通常作为循环条件，而Text()用于在循环体内获取标记的内容。
	*/
	for scanner.Scan() {
		line := scanner.Text()
		maze = append(maze, line)
	}

	//我们需要在加载迷宫后立即捕获玩家位置
	for row, line := range maze {
		for col, char := range line {
			switch char {
			case 'P':
				player = sprite{row, col}
			}
		}
	}
	return nil
}

// 重构这个函数，让他只打印我们想打印的
func printScreen() {
	simpleansi.ClearScreen()
	for _, line := range maze {
		for _, chr := range line {
			switch chr {
			case '#':
				fmt.Printf("%c", chr)
			default:
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	simpleansi.MoveCursor(player.row, player.col)
	fmt.Print("P")

	//将光标移动到迷宫绘制区域外
	simpleansi.MoveCursor(len(maze)+1, 0)
}

// 启用Cbreak模式
func initialise() {
	cbTerm := exec.Command("stty", "cbreak", "-echo")
	cbTerm.Stdin = os.Stdin

	err := cbTerm.Run()
	if err != nil {
		log.Fatalln("unable to activate cbreak mode:", err)
		//如果出现错误，该函数
	}
}

//恢复熟模式

func cleanup() {
	cookedTerm := exec.Command("stty", "-cbreak", "-echo")
	cookedTerm.Stdin = os.Stdin

	err := cookedTerm.Run()
	if err != nil {
		log.Fatalln("unable to restore cooked mode:", err)
	}
}

// 从标准输入读取
func readInput() (string, error) {
	buffer := make([]byte, 100)

	cnt, err := os.Stdin.Read(buffer)
	if err != nil {
		return "", err
	}

	if cnt == 1 && buffer[0] == 0x1b {
		return "ESC", nil
	} else if cnt >= 3 {
		if buffer[0] == 0x1b && buffer[1] == '[' {
			//箭头键的转义序列长度为 3 个字节，ESC+[从 A 到 D 的字母开始
			switch buffer[2] {
			case 'A':
				return "UP", nil
			case 'B':
				return "DOWN", nil
			case 'C':
				return "RIGHT", nil
			case 'D':
				return "LEFT", nil
			}
		}
	}

	return "", nil

}

func makeMove(oldRow, oldCol int, dir string) (newRow, newCol int) {
	newRow, newCol = oldRow, oldCol

	switch dir {
	case "UP":
		newRow = newRow - 1
		if newRow < 0 {
			newRow = len(maze) - 1 //会回到最远的地方
		}
	case "DOWN":
		newRow = newRow + 1
		if newRow == len(maze) {
			newRow = 0
		}
	case "RIGHT":
		newCol = newCol + 1
		if newCol == len(maze[0]) {
			newCol = 0
		}
	case "LEFT":
		newCol = newCol - 1
		if newCol < 0 {
			newCol = len(maze[0]) - 1
		}
	}
	//撞到障碍物
	if maze[newRow][newCol] == '#' {
		newRow = oldRow
		newCol = oldCol
	}
	return
}

func movePlayer(dir string) {
	player.row, player.col = makeMove(player.row, player.col, dir)
}
func main() {
	//initial game
	initialise()
	defer cleanup()
	//load resources
	err := loadMaze("maze01.txt")
	if err != nil {
		log.Println("failed to load maze:", err)
		return
	}

	//game loop
	for {
		//update screen
		printScreen()

		//process input
		input, err := readInput()
		if err != nil {
			log.Print("error reading input:", err)
			break
		}
		if input == "ESC" {
			break
		}
		//process movement
		movePlayer(input)
		//process collisions

		//check game over
		if input == "ESC" {
			break
		}
		//Temp:break infinite loop(死循环)

		//repeat
	}
}
