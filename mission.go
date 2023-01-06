package main

type missionTypeCode uint8

const (
	MISSION_STEAL_TARGET_ITEMS missionTypeCode = iota
	MISSION_STEAL_MINIMUM_LOOT
)

type Mission struct {
	BriefingText           string
	DebriefingText         string
	AdditionalStartingGold int

	MissionType               missionTypeCode
	DifficultyChoosingStr     string
	DifficultyLevelsNames     []string
	TargetNumber              []int
	TargetItemsNames          []string
	AdditionalGuardsNumber    []int
	NumberOfArchersFromGuards []int
	Rewards                   []int
	TotalLoot                 []int
}

func (m *Mission) readFromFile(filename string) {

}
