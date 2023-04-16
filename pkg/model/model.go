//nolint:gomnd
package model

import (
	"time"

	"github.com/fatih/color"
)

const DailyHours = 24

// Element represents one of the Chinese 5 elements zodiac system.
type Element int

var (
	woc = color.New(color.FgHiWhite, color.BgGreen, color.Bold).SprintFunc()
	fc  = color.New(color.FgHiWhite, color.BgHiRed, color.Bold).SprintFunc()
	ec  = color.New(color.FgBlack, color.BgHiYellow, color.Bold).SprintFunc()
	mc  = color.New(color.FgBlack, color.BgHiWhite, color.Bold).SprintFunc()
	wac = color.New(color.FgHiWhite, color.BgHiBlue, color.Bold).SprintFunc()
)

const (
	Wood Element = iota
	Fire
	Earth
	Metal
	Water
)

// String returns the name of the element as a text.
func (e Element) String() string {
	switch e {
	case Wood:
		return "wood"
	case Fire:
		return "fire"
	case Earth:
		return "earth"
	case Metal:
		return "metal"
	case Water:
		fallthrough
	default:
		return "water"
	}
}

// Color returns the element's name as an ansi colored string.
func (e Element) Color() string {
	var c func(a ...interface{}) string
	switch e {
	case Wood:
		c = woc
	case Fire:
		c = fc
	case Earth:
		c = ec
	case Metal:
		c = mc
	case Water:
		fallthrough
	default:
		c = wac
	}
	return c(e.String())
}

// GetElements returns daily elements for a given local time for all three levels
// in jing, qi, shen order. It's repeated every 60 days.
func GetElements(t time.Time) (jing, qi, shen Element) {
	d := int(t.Sub(time.Date(2016, 2, 26, 0, 0, 0, 0, time.UTC)).Hours()) / DailyHours
	return getElementJing(d), getElementQi(d), getElementShen(d)
}

// getElementJing returns the jing level element for the given day.
// It's repeated every 12 days.
func getElementJing(day int) Element {
	switch day % 12 {
	case 0, 1:
		return Wood
	case 4, 5:
		return Fire
	case 2, 3, 8, 11:
		return Earth
	case 6, 7:
		return Metal
	case 9, 10:
		fallthrough
	default:
		return Water
	}
}

// getElementQi returns the qi level element for the given day.
// It's repeated every 30 days.
func getElementQi(day int) Element {
	switch (day % 30) / 2 {
	case 2, 6, 10:
		return Wood
	case 5, 9, 13:
		return Fire
	case 0, 4, 11:
		return Earth
	case 1, 8, 12:
		return Metal
	case 3, 7, 14:
		fallthrough
	default:
		return Water
	}
}

// getElementShen returns the shen level element for the given day.
// It's repeated every 10 days.
func getElementShen(day int) Element {
	switch (day % 10) / 2 {
	case 3:
		return Wood
	case 4:
		return Fire
	case 0:
		return Earth
	case 1:
		return Metal
	case 2:
		fallthrough
	default:
		return Water
	}
}
