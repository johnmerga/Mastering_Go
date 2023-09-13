package interaction

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

var Reader = bufio.NewReader(os.Stdin)

func UserInput(isSpecial bool) (int8, error) {
	input, err := Reader.ReadString('\n')
	if err != nil {
		return -1, err
	}
	input = strings.Replace(input, "\n", "", -1)
	// parse the string into int8
	parsedInput, err := strconv.ParseInt(input, 10, 8)
	if err != nil {
		return -1, errors.New("Please Only type numbers")
	}
	selectedNum := int8(parsedInput)
	if selectedNum <= 0 || selectedNum >= 4 {
		return -1, errors.New("The Number Should be (1),(2) or (3) if it's available \n")
	}
	if !isSpecial && selectedNum == 3 {
		return -1, errors.New("The Special Attack is not Available yet.\n")
	}
	return selectedNum, nil
}
