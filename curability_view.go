package main

import (
	"fmt"
	"strings"

	"github.com/anthonybishopric/pandemic-nerd-hurd/pandemic"
)

type maxCurability struct {
	prob   float64
	player *pandemic.Player
}

type byCurability struct {
	dts           []pandemic.DiseaseType
	curability    map[pandemic.DiseaseType]float64
	maxCurability map[pandemic.DiseaseType]maxCurability
}

func (b byCurability) Less(i, j int) bool {
	return strings.Compare(b.dts[i].String(), b.dts[j].String()) < 0
}

func (b byCurability) Len() int {
	return len(b.dts)
}

func (b byCurability) Swap(i, j int) {
	b.dts[i], b.dts[j] = b.dts[j], b.dts[i]
}

func (p *PandemicView) colorProbabilityOfCure(prob float64) string {
	str := fmt.Sprintf("%.2f", prob)
	if prob < 0.2 {
		return p.colorOhFuck(str)
	}
	if prob < 0.8 {
		return p.colorWarning(str)
	}
	return p.colorAllGood(str)
}

// ⚗
// Alembic
// Unicode: U+2697, UTF-8: E2 9A 97
