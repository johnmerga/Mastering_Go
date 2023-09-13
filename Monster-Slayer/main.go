package main

import (
	"github.com/johnmerga/Go-Practice/Monster-Slayer/actions"
	"github.com/johnmerga/Go-Practice/Monster-Slayer/helper"
	"github.com/johnmerga/Go-Practice/Monster-Slayer/interaction"
)

var Round = 1
var UserHealth = helper.INITIAL_UER_HEALTH
var MonsterHealth = helper.INITIAL_MONSTER_HEALTH

var GameRecord helper.GameRecord = helper.GameRecord{}

var GameRecords []helper.GameRecord

var attack_monster_damage int8
var attack_player_damage int8
var special_attack_monster_damage int8
var healed_by int8

func main() {

	interaction.StartGame()
	for {

		isSpecial := Round%3 == 0
		interaction.PrintOption(isSpecial)
		val, err := interaction.UserInput(isSpecial)
		if err != nil {
			println(err.Error())
		} else {
			switch val {
			case 1:
				GameRecord.PlayerAttacks = actions.AttackMonster(&MonsterHealth)
				GameRecord.MonsterAttacks = actions.AttackPlayer(&UserHealth)
			case 2:
				GameRecord.PlayerHeals = actions.HealPlayer(&UserHealth)
				GameRecord.MonsterAttacks = actions.AttackPlayer(&UserHealth)
			case 3:
				GameRecord.SpecialAttack = actions.SpecialAttack(&MonsterHealth)
				GameRecord.MonsterAttacks = actions.AttackPlayer(&UserHealth)
			}
			GameRecord.Round = int8(Round)
			GameRecord.MonsterHealth = MonsterHealth
			GameRecord.PlayerHealth = UserHealth
			interaction.PrintHealth(UserHealth, MonsterHealth)
			Round++
			GameRecords = append(GameRecords, GameRecord)
			//reset GameRecord
			GameRecord = helper.GameRecord{}
		}
		isGameEnded, winner := helper.IsTheGameOver(UserHealth, MonsterHealth)
		if isGameEnded {
			interaction.EndGame(winner)
			GameRecord.Round = int8(Round)
			GameRecord.Winner = winner
			GameRecords = append(GameRecords, GameRecord)
			helper.SaveGameHistory(&GameRecords)
			return
		}

	}

}
