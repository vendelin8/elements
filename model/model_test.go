package model

import (
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/leaanthony/go-ansi-parser"
)

// TestModel tests calendar calculation.
func TestModel(t *testing.T) {
	cases := []struct {
		tim      time.Time
		wantJing Element
		wantQi   Element
		wantShen Element
	}{
		{
			tim:      time.Date(2019, 7, 29, 0, 0, 0, 0, time.UTC),
			wantJing: Wood,
			wantQi:   Fire,
			wantShen: Fire,
		},
		{
			tim:      time.Date(2019, 7, 29, 23, 59, 59, 1e9-7, time.UTC),
			wantJing: Wood,
			wantQi:   Fire,
			wantShen: Fire,
		},
		{
			tim:      time.Date(2019, 8, 26, 0, 0, 0, 0, time.UTC),
			wantJing: Fire,
			wantQi:   Metal,
			wantShen: Wood,
		},
		{
			tim:      time.Date(2019, 8, 31, 0, 0, 0, 0, time.UTC),
			wantJing: Water,
			wantQi:   Earth,
			wantShen: Metal,
		},
		{
			tim:      time.Date(2019, 9, 1, 0, 0, 0, 0, time.UTC),
			wantJing: Earth,
			wantQi:   Earth,
			wantShen: Metal,
		},
		{
			tim:      time.Date(2019, 9, 28, 0, 0, 0, 0, time.UTC),
			wantJing: Earth,
			wantQi:   Wood,
			wantShen: Earth,
		},
		{
			tim:      time.Date(2019, 11, 1, 0, 0, 0, 0, time.UTC),
			wantJing: Wood,
			wantQi:   Metal,
			wantShen: Water,
		},
		{
			tim:      time.Date(2019, 12, 13, 0, 0, 0, 0, time.UTC),
			wantJing: Metal,
			wantQi:   Water,
			wantShen: Wood,
		},
	}

	for _, c := range cases {
		t.Run(c.tim.String(), func(t *testing.T) {
			resJing, resQi, resShen := GetElements(c.tim)
			if resJing != c.wantJing {
				t.Errorf("jing level result %s doesn't match expected %s", resJing, c.wantJing)
			}
			if resQi != c.wantQi {
				t.Errorf("qi level result %s doesn't match expected %s", resQi, c.wantQi)
			}
			if resShen != c.wantShen {
				t.Errorf("shen level result %s doesn't match expected %s", resShen, c.wantShen)
			}
		})
	}
}

// TestModelOutput tests what the model returns as colored strings.
func TestModelOutput(t *testing.T) {
	cases := []struct {
		el      Element
		wantTxt string
		wantFg  string
		wantBg  string
	}{
		{
			el:      Wood,
			wantTxt: "wood",
			wantFg:  "White",
			wantBg:  "Green",
		},
		{
			el:      Fire,
			wantTxt: "fire",
			wantFg:  "White",
			wantBg:  "Red",
		},
		{
			el:      Earth,
			wantTxt: "earth",
			wantFg:  "Black",
			wantBg:  "Yellow",
		},
		{
			el:      Metal,
			wantTxt: "metal",
			wantFg:  "Black",
			wantBg:  "White",
		},
		{
			el:      Water,
			wantTxt: "water",
			wantFg:  "White",
			wantBg:  "Blue",
		},
	}

	color.NoColor = false
	wantStyle := ansi.Bold + ansi.Bright
	for _, c := range cases {
		t.Run(c.el.String(), func(t *testing.T) {
			col := c.el.Color()
			parsed, err := ansi.Parse(col)
			if err != nil {
				t.Errorf("color result has an error: %v", err)
			}
			if len(parsed) != 1 {
				t.Errorf("parsed color should have 1 style, but it's: %#v", parsed)
			}
			p0 := parsed[0]
			if p0.FgCol.Name != c.wantFg {
				t.Errorf("fg color result %s doesn't match expected %s", p0.FgCol.Name, c.wantFg)
			}
			if p0.Style != wantStyle {
				t.Errorf("color should be %d, but it's %v", wantStyle, p0.Style)
			}
			if p0.BgCol.Name != c.wantBg {
				t.Errorf("bg color result %s doesn't match expected %s", p0.BgCol.Name, c.wantBg)
			}
			if p0.Label != c.wantTxt {
				t.Errorf("result text %s doesn't match expected %s", p0.Label, c.wantTxt)
			}
		})
	}
}
