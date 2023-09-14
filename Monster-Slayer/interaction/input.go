package interaction

import (
	"bufio"
	"errors"
	"os"

	"runtime"
	"strconv"
	"strings"
)

var Reader = bufio.NewReader(os.Stdin)

func UserInput(isSpecial bool) (int8, error) {
	input, err := Reader.ReadString('\n')
	if err != nil {
		return -1, err
	}
	endingLine := "\n"
	if runtime.GOOS == "windows" {
		endingLine = "\r\n"
	}
	input = strings.Replace(input, endingLine, "", -1)
	// parse the string into int8
	parsedInput, err := strconv.ParseInt(input, 10, 8)
	if err != nil {
		return -1, errors.New("please Only type numbers")
	}
	selectedNum := int8(parsedInput)
	if selectedNum <= 0 || selectedNum >= 4 {
		return -1, errors.New("the number should be (1),(2) or (3) if it's available")
	}
	if !isSpecial && selectedNum == 3 {
		return -1, errors.New("the Special attack is not Available yet")
	}
	return selectedNum, nil
}
