package actions

import "github.com/johnmerga/Go-Practice/Monster-Slayer/helper"

func AttackMonster(monsterHealth *int8) (damage int8) {
	genNun := helper.GenerateRandNum(helper.ATTACK_MONSTER_MIN_DMG, helper.ATTACK_MONSTER_MAX_DMG)
	*monsterHealth -= genNun
	return genNun
}
func SpecialAttack(monsterHealth *int8) (damage int8) {
	genNum := helper.GenerateRandNum(helper.SPECIAL_ATTACK_MONSTER_MIN_DMG, helper.SPECIAL_ATTACK_MONSTER_MAX_DMG)
	*monsterHealth -= genNum
	return genNum
}
func AttackPlayer(playerHealth *int8) (damage int8) {
	genNum := helper.GenerateRandNum(helper.ATTACK_PLAYER_MIN_DMG, helper.ATTACK_PLAYER_MAX_DMG)
	*playerHealth -= genNum
	return genNum
}

// heal player
func HealPlayer(playerHealth *int8) (healed int8) {
	genNum := helper.GenerateRandNum(helper.HEAL_PLAYER_MIN, helper.HEAL_PLAYER_MAX)
	if genNum+*playerHealth >= helper.INITIAL_UER_HEALTH {
		*playerHealth = helper.INITIAL_UER_HEALTH
	} else {
		*playerHealth += genNum

	}
	return genNum
}
