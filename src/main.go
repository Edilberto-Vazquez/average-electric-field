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

func ReadFile(path string, c chan [][]string) {
	file, err := os.Open(path)
	Error(err)
	defer file.Close()
	fileName := GetFileName(path)
	fileDate := fileName[17:21] + "-" + fileName[13:15] + "-" + fileName[15:17]
	scanner := bufio.NewScanner(file)
	var lines [][]string
	for scanner.Scan() {
		value := strings.Split(scanner.Text(), ",")
		value[0] = fileDate + " " + value[0]
		lines = append(lines, value)
	}
	c <- lines
}

func WriteFile(file *os.File, path string, lines []string) {
	// file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	// Error(err)
	// defer file.Close()
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
	average = math.Round((sum / lenList * 100)) / 100
	return
}

func ProcessLines(list [][]string, c chan<- []string) {
	var checkLine []string = list[0]
	var chunkList [][]string
	var newList []string
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
	c <- newList
}

func main() {
	files := ReadDir("/media/edilberto/Nuevo vol/DataSets/Conjuntos-originales/medidor-campo-electrico/")
	file, err := os.Create("/home/edilberto/Desktop/electric-field-measurements.txt")
	Error(err)
	defer file.Close()
	read := make(chan [][]string, len(files))
	process := make(chan []string, len(files))
	for _, v := range files {
		go ReadFile("/media/edilberto/Nuevo vol/DataSets/Conjuntos-originales/medidor-campo-electrico/"+v, read)
		go ProcessLines(<-read, process)
	}
	close(read)
	for i := 0; i < len(files); i++ {
		WriteFile(file, "/home/edilberto/Desktop/electric-field-measurements.txt", <-process)
	}
	close(process)
}
