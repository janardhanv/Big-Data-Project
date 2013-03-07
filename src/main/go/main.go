package main

import (
    "bufio"
    "os"
    "strings"
    "strconv"
    "fmt"
    "math"
    // "runtime"
)

var eta float64 = 0.01

type DataLine struct {
    clicked, depth, position, userid, gender, age int
    tokens []int
}

func main() {
    file, err := os.Open("../resources/train.txt")
    r := bufio.NewReader(file)
    i := 0
    line := ""

    dl := make([]DataLine, 2335859)
    w := make(map[int]float64)

    for err == nil && i < 2335859 {
        line, err = r.ReadString('\n')
        parseLine(line, dl, i)
        i++
    }

    // runtime.GOMAXPROCS(2)

    for i := 0; i < 2335859; i++ {
    	sgd(&dl[i], w)
    }

    fmt.Println("Intercept", w[0])
    fmt.Println("Depth", w[1])
    fmt.Println("Position", w[2])
    fmt.Println("Gender", w[3])
    fmt.Println("Age", w[4])

    fmt.Println("l2 norm", l2norm(w))

    file.Close()
}


func parseLine(s string, dl []DataLine, i int) {
    split := strings.SplitN(s, "|", 7)

    dl[i].clicked, _ = strconv.Atoi(split[0])
    dl[i].depth, _ = strconv.Atoi(split[1])
    dl[i].position, _ = strconv.Atoi(split[2])
    dl[i].userid, _ = strconv.Atoi(split[3])
    dl[i].gender, _ = strconv.Atoi(split[4])
    dl[i].age, _ = strconv.Atoi(split[5])

    tokenSplit := strings.Split(split[6], ",")

    size := len(tokenSplit)
    dl[i].tokens = make([]int, size)

    for k, v := range tokenSplit {
        dl[i].tokens[k], _ = strconv.Atoi(strings.TrimSpace(v))
    }
}


func predictLabel(dl *DataLine, w map[int]float64) float64 {
	offset := 5
	sum :=  w[0] +
			float64(dl.depth) * w[1] +
			float64(dl.position) * w[2] +
			float64(dl.gender) * w[3] +
			float64(dl.age) * w[4]
	for _, v := range dl.tokens {
		sum += w[v + offset]
	}
	numer := math.Exp(sum)
	return numer / (1.0 + numer)
}


func gradient(y int, yhat float64) float64 {
    return (float64(y) - yhat) * eta
}


func updateWeights(dl *DataLine, w map[int]float64, grad float64) {
	offset := 5
	w[0] += grad
	w[1] += grad * float64(dl.depth)
	w[2] += grad * float64(dl.position)
	w[3] += grad * float64(dl.gender)
	w[4] += grad * float64(dl.age)
	for _, v := range dl.tokens {
		w[v + offset] += grad
	}
}

func sgd(dl *DataLine, w map[int]float64) {
	yhat := predictLabel(dl, w)
	grad := gradient(dl.clicked, yhat)
	updateWeights(dl, w, grad)
}

func l2norm(w map[int]float64) float64 {
	sum := float64(0)
	for _, v := range w {
		sum += v * v
	}
	return math.Sqrt(sum)
}