package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFileName := flag.String("csv", "problems.csv", "csv file in format for 'question,answer'. ")
	timeLimit := flag.Int("limit", 30, "the time limit for quiz in sec")
	flag.Parse()
	fmt.Println("csvfilename - ", *csvFileName)

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Couldn't open the file - %s ", *csvFileName))
	}
	fmt.Println("file - ", file)
	fileReader := csv.NewReader(file)
	lines, err := fileReader.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Couldn't parse the provided file - %s ", *csvFileName))
	}
	fmt.Println("lines - ", parseLines(lines))

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correctAns := 0
	problems := parseLines(lines)
	problemLoop:
		for i, p := range problems {
			fmt.Printf("Problem #%d: %s = \n", i+1, p.ques)
			answerChan := make(chan string)
			go func() {
				var answer string
				fmt.Scanf("%s\n", &answer)
				answerChan <- answer
			}()
			select {
			case <- timer.C:
				fmt.Println()
				break problemLoop
			
			case answer := <- answerChan:  
				if answer == p.ans {
					correctAns++
				}
			}
		}
		
	
	fmt.Printf("You scored %d out of %d. \n", correctAns, len(problems))
}

type problem struct {
	ques string
	ans string
} 

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			ques: line[0],
			ans:  strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}