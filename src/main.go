package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

func Error(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func ReadDir(dir string) (filesNames []string) {
	files, err := ioutil.ReadDir(dir)
	Error(err)
	for _, v := range files {
		filesNames = append(filesNames, v.Name())
	}
	return
}

func GetFileName(path string) (fileName string) {
	slplitPath := strings.Split(path, "/")
	fileName = slplitPath[len(slplitPath)-1]
	return
}

func ReadFile(path string) (line [][]string) {
	file, err := os.Open(path)
	Error(err)
	defer file.Close()
	fileName := GetFileName(path)
	fileDate := fmt.Sprintf("%s-%s-%s", fileName[17:21], fileName[13:15], fileName[15:17])
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		value := strings.Split(scanner.Text(), ",")
		value[0] = fmt.Sprintf("%s %s", fileDate, value[0])
		line = append(line, value)
	}
	return
}

func WriteFile(lines []string) {
	file, err := os.Create("/home/edilberto/Desktop/New.txt")
	Error(err)
	defer file.Close()
	write := bufio.NewWriter(file)
	for _, v := range lines {
		write.WriteString(v)
	}
	write.Flush()
}

func CalcAverage(chunk [][]string) (average float64) {
	sum := 0.0
	lenList := float64(len(chunk))
	for _, v := range chunk {
		converFloat, err := strconv.ParseFloat(v[1], 32)
		Error(err)
		sum += converFloat
	}
	average = math.Floor((sum / lenList * 100)) / 100
	return
}

func ProcessLines(list [][]string) (newList []string) {
	var checkLine []string = list[0]
	var chunkList [][]string
	for _, v := range list {
		if checkLine[0] == v[0] {
			chunkList = append(chunkList, v)
		} else {
			line := fmt.Sprintf("%s, %.2f, %s\n", checkLine[0], CalcAverage(chunkList), checkLine[2])
			newList = append(newList, line)
			checkLine = v
			chunkList = [][]string{}
		}
	}
	return
}

func main() {
	WriteFile(ProcessLines(ReadFile("/home/edilberto/Desktop/INAOE parque-01142019.efm")))
	// fmt.Println(ReadDir("/home/edilberto/Desktop/"))
}
