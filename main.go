package main

import (
	"bufio"
	"fmt"
	"os/exec"
)

func main() {
	// Create channels to receive output from each program.
	ch1 := make(chan string)
	ch2 := make(chan string)
	ch3 := make(chan string)

	// Start each program in a separate Goroutine.
	go runProgram("metadata/cmd/main.go", ch1)
	go runProgram("movie/cmd/main.go", ch2)
	go runProgram("rating/cmd/main.go", ch3)

	for {
		select {
		case msg := <-ch1:
			fmt.Println(msg)
		case msg := <-ch2:
			fmt.Println(msg)
		case msg := <-ch3:
			fmt.Println(msg)
		}
	}
}

// Helper function to run a Go program and send its output to a channel.
func runProgram(filename string, ch chan<- string) {
	cmd := exec.Command("go", "run", filename)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		ch <- fmt.Sprintf("Error creating pipe for %s: %s", filename, err)
		return
	}

	if err := cmd.Start(); err != nil {
		ch <- fmt.Sprintf("Error running %s: %s", filename, err)
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		ch <- fmt.Sprintf("%s: %s", filename, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		ch <- fmt.Sprintf("Error reading output from %s: %s", filename, err)
	}

	if err := cmd.Wait(); err != nil {
		ch <- fmt.Sprintf("Error waiting for %s to finish: %s", filename, err)
	}
}
