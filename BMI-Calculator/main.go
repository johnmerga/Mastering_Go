package main

import (
	"bmi/functions"
	"fmt"
)

func main() {
	weight, height := function.UserInput()
	bmiResult := function.CalculateBmi(weight, height)
	fmt.Println(bmiResult)
}
