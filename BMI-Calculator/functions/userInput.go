package function

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func UserInput() (float32, float32) {
	// var reader = bufio.NewReader(os.Stdin)
	// fmt.Println("please Enter your weight")
	// var weight float32
	// fmt.Scanln(&weight)

	// fmt.Println("please Enter your height")
	// var height float32
	// fmt.Scanln(&height)
	// checkReader, _ := reader.ReadString('\n')
	// inputFromReader := strings.Replace(checkReader, "\n", "", -1)
	// result, err := strconv.ParseFloat(inputFromReader, 64)
	// // checking if error exist
	// if err != nil {
	// 	fmt.Println("Error while converting string to float64")
	// }
	// fmt.Println(result)

	// return weight, height

	


	// myFullText := fmt.Sprintf("my weight is : %.2f, and my hight: %.2f", weight,hight)
	// println(myFullText)
	weight := userInput()
	hight := userInput()

	return weight,hight

}

func userInput() float32 {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("please type your weight in kilogram(KG). (e.g: 63)")
	weightInput, weightErr := reader.ReadString('\n')
	if weightErr != nil {
		fmt.Println("please only type number")
		userInput()
	}
	cleanedInput := strings.Replace(weightInput,"\n","",-1)
	parsedInput,err := strconv.ParseFloat(cleanedInput,32)
	if err != nil {
		fmt.Println("invalid user input. please only type a number")
		userInput()
	}

	return float32(parsedInput);
}
