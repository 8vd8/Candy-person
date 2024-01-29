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

var maze []string

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
	return nil
}

func printScreen() {
	simpleansi.ClearScreen()
	for _, line := range maze {
		fmt.Println(line)
	}
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
	}
	return "", nil
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

		//process collisions

		//check game over

		//Temp:break infinite loop(死循环)
		break

		//repeat
	}
}
