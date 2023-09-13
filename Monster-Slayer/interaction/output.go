package interaction

import (
	"fmt"

	asciiArt "github.com/common-nighthawk/go-figure"
	"github.com/johnmerga/Go-Practice/Monster-Slayer/helper"
)

func StartGame() {
	asciiArt.NewFigure("           MONSTER  SLAYER", "", true).Scroll(5000, 200, "right")
	asciiArt.NewFigure("--------------------------------------", "", true).Print()
	println()
	println()
}
func PrintOption(isRoundThree bool) {
	println("**************************************")
	println("Please Select the options below")
	println("(1)- Attack the Monster")
	println("(2)- Heal Yourself")
	if isRoundThree {
		println("(3)- Special Attack(available)")
	}
	println("**************************************")
}

func PrintHealth(playerHealth, monsterHealth int8) {
	fmt.Printf("Player Health: %v \n", playerHealth)
	fmt.Printf("Monster Health: %v \n", monsterHealth)
}

func EndGame(winner string) {
	switch winner {
	case helper.USER:
		asciiArt.NewColorFigure("CONGRATULATIONS", "", "green", true).Blink(10000, 500, -1)
		println("You Won the Game")
	case helper.MONSTER:
		asciiArt.NewColorFigure("GAME OVER", "", "red", true).Blink(10000, 500, -1)
	}

}
