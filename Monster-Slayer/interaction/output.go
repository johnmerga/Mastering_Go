package interaction

import (
	"fmt"

	"github.com/johnmerga/Go-Practice/Monster-Slayer/helper"
)

func StartGame() {
	println("--------------------------------------")
	println("           MONSTER SLAYER")
	println("--------------------------------------")
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
		println("CONGRATULATIONS")
		println("You Won the Game")
	case helper.MONSTER:
		println("GAME OVER!!!")
	}

}
