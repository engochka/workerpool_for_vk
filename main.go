package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"workerpool_ch2/workerpool"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	pool := workerpool.NewPool(nil, 2)

	go pool.RunBackground()

	fmt.Println("Команды:")
	fmt.Println("add - добавить воркера")
	fmt.Println("remove - удалить воркера")
	fmt.Println("task [текст] - добавить задачу")
	fmt.Println("exit - выход")

	for {
		fmt.Println("Укажите команду")
		fmt.Print(">")
		scanner.Scan()
		input := scanner.Text()
		parts := strings.SplitN(input, " ", 2)

		switch parts[0] {
		case "add":
			pool.AddWorker()
		case "remove":
			pool.RemoveWorker()
		case "task":
			if len(parts) < 2 {
				fmt.Println("Укажите текст задачи")
				continue
			}
			taskID := len(pool.Tasks) + 1
			task, err := workerpool.NewTask(parts[1], taskID)
			if err != nil {
				fmt.Println("Ошибка создания задачи:", err)
				continue
			}
			pool.AddTask(task)
			pool.Tasks = append(pool.Tasks, task)
			fmt.Printf("Задача #%d добавлена\n", taskID)
		case "exit":
			pool.Stop()
			return
		default:
			fmt.Println("Неизвестная команда")
		}

	}
}
