package main

import (
	"fmt"
	"github.com/johnmerga/Go-Practice/BMI/functions"
)

func main() {
	weight, height := function.UserInput()
	bmiResult := function.CalculateBmi(weight, height)
	fmt.Println(bmiResult)
}
