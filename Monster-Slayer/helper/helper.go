package helper

import (
	"encoding/json"
	"math/rand"
	"os"
)

var INITIAL_UER_HEALTH int8 = 100
var INITIAL_MONSTER_HEALTH int8 = 100
var ATTACK_PLAYER_MIN_DMG int8 = 10
var ATTACK_PLAYER_MAX_DMG int8 = 20
var HEAL_PLAYER_MIN int8 = 10
var HEAL_PLAYER_MAX int8 = 20
var ATTACK_MONSTER_MIN_DMG int8 = 10
var ATTACK_MONSTER_MAX_DMG int8 = 17
var SPECIAL_ATTACK_MONSTER_MIN_DMG int8 = 13
var SPECIAL_ATTACK_MONSTER_MAX_DMG int8 = 24
var MONSTER = "Monster"
var USER = "User"

type GameRecord struct {
	Round          int8
	MonsterHealth  int8
	PlayerHealth   int8
	MonsterAttacks int8
	PlayerAttacks  int8
	SpecialAttack  int8
	PlayerHeals    int8
	Winner         string
}

func GenerateRandNum(min, max int8) int8 {
	if min >= max {
		return min
	}
	diff := int(max) - int(min)
	if diff <= 0 {
		return min
	}
	randNum := rand.Intn(diff + 1)
	return int8(randNum) + min
}

func IsTheGameOver(playerHealth, monsterHealth int8) (bool, string) {
	if playerHealth <= 0 {
		return true, MONSTER
	} else if monsterHealth <= 0 {
		return true, USER
	} else {
		return false, ""
	}
}

func SaveGameHistory(gameRecords *[]GameRecord) {
	file, err := os.OpenFile("game-history.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		println(err)
		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)

	// Write an opening bracket to start a JSON array
	file.WriteString("[")

	// Loop through the records and encode them
	for i, v := range *gameRecords {
		err := encoder.Encode(v)
		if err != nil {
			println(err)
			return
		}
		// If it's not the last record, add a comma separator
		if i < len(*gameRecords)-1 {
			file.WriteString(",")
		}
	}

	// Write a closing bracket to end the JSON array
	file.WriteString("]")
}
