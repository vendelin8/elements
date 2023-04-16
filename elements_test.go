package main

import (
	"strings"
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/vendelin8/elements/pkg/cio"
)

// TestParseFlags tests parsing command line arguments to Config.
func TestParseFlags(t *testing.T) {
	cases := []struct {
		args   []string
		want   *Config
		hasErr bool
	}{
		{
			args: []string{"apple"},
		},
		{
			args:   []string{"-l"},
			hasErr: true,
		},
		{
			args:   []string{"-location"},
			hasErr: true,
		},
		{
			args:   []string{"--location"},
			hasErr: true,
		},
		{
			args:   []string{"--loc"},
			hasErr: true,
		},
		{
			args:   []string{"-loc"},
			hasErr: true,
		},
		{
			args: []string{"--loc", "some_location"},
			want: &Config{location: "some_location"},
		},
		{
			args: []string{"-loc", ""},
			want: &Config{location: ""},
		},
		{
			args: []string{"-loc", "-"},
			want: &Config{location: ""},
		},
		{
			args: []string{"-loc", "some_loc"},
			want: &Config{location: "some_loc"},
		},
	}

	for _, c := range cases {
		t.Run(strings.Join(c.args, "_"), func(t *testing.T) {
			args := make([]string, len(c.args)+1)
			copy(args[1:], c.args)
			args[0] = "elements"
			conf, err := parseFlags(args...)
			if err != nil != c.hasErr {
				if c.hasErr {
					t.Error("needs parse error, but it's nil")
				} else {
					t.Errorf("parse error should be nil, but it's %v", err)
				}
			}
			if c.want == nil {
				c.want = &Config{}
			}
			if conf.location != c.want.location {
				t.Errorf("parsed location config %s doesn't match expected %s",
					conf.location, c.want.location)
			}
		})
	}
}

// TestHandleFlags tests using command line arguments.
func TestHandleFlags(t *testing.T) {
	cases := []struct {
		name     string
		day      time.Time
		loc      string
		wantDay  time.Time
		wantCode int
	}{
		{
			name:    "emptyLocation",
			day:     time.Date(2019, 8, 26, 0, 0, 0, 0, time.UTC),
			wantDay: time.Date(2019, 8, 26, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "error",
			day:      time.Date(2019, 8, 26, 0, 0, 0, 0, time.UTC),
			loc:      "invalid zone",
			wantCode: 1,
		},
		{
			name:    "goBack",
			day:     time.Date(2019, 8, 26, 0, 0, 0, 0, time.UTC),
			loc:     "America/Detroit",
			wantDay: time.Date(2019, 8, 25, 0, 0, 0, 0, time.UTC),
		},
		{
			name:    "goForward",
			day:     time.Date(2019, 8, 26, 23, 0, 0, 0, time.UTC),
			loc:     "Europe/Budapest",
			wantDay: time.Date(2019, 8, 27, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			visibleDay = c.day
			code := handleFlags(&Config{location: c.loc}, func(loc *time.Location) time.Time {
				return c.day.In(loc)
			})
			if code != c.wantCode {
				t.Errorf("exit code %d doesn't match expected %d", code, c.wantCode)
			}
			if code == 1 {
				return
			}
			if !visibleDay.Equal(c.wantDay) {
				t.Errorf("calculated day %s doesn't match expected %s", visibleDay, c.wantDay)
			}
		})
	}
}

// TestFillElements tests a line that gets printed as date+daily elements.
func TestFillElements(t *testing.T) {
	color.NoColor = false
	cases := []struct {
		day  time.Time
		want string
	}{
		{
			day:  time.Date(2019, 8, 26, 0, 0, 0, 0, time.UTC),
			want: "2019-08-26: shen: \033[97;42;1mwood\033[0m, qi: \033[30;107;1mmetal\033[0m, jing: \033[97;101;1mfire\033[0m",
		},
		{
			day:  time.Date(2019, 8, 31, 0, 0, 0, 0, time.UTC),
			want: "2019-08-31: shen: \033[30;107;1mmetal\033[0m, qi: \033[30;103;1mearth\033[0m, jing: \033[97;104;1mwater\033[0m",
		},
	}

	for _, c := range cases {
		t.Run(c.day.String(), func(t *testing.T) {
			visibleDay = c.day
			res := fillElements()
			if res != c.want {
				t.Errorf("result text %s doesn't match expected %s", res, c.want)
			}
		})
	}
}

// TestChangeDay tests day changing callback function.
func TestChangeDay(t *testing.T) {
	visibleDay = time.Date(2023, 4, 25, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		name string
		act  int
		diff int
		want []int // year, month, day
	}{
		{
			name: "lvl0Up",
			act:  cio.ActLvl0,
			diff: 1,
			want: []int{2024, 4, 25},
		},
		{
			name: "lvl1Down",
			act:  cio.ActLvl1,
			diff: -1,
			want: []int{2024, 3, 25},
		},
		{
			name: "lvl2Up",
			act:  cio.ActLvl2,
			diff: 1,
			want: []int{2024, 3, 26},
		},
		{
			name: "lvl1Up",
			act:  cio.ActLvl1,
			diff: 1,
			want: []int{2024, 4, 26},
		},
		{
			name: "lvl0Down",
			act:  cio.ActLvl0,
			diff: -1,
			want: []int{2023, 4, 26},
		},
		{
			name: "lvl2Down",
			act:  cio.ActLvl2,
			diff: -1,
			want: []int{2023, 4, 25},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			changeDay(c.act, c.diff)
			y, m, d := visibleDay.Date()
			if y != c.want[0] {
				t.Errorf("visible year %d doesn't match expected %d", y, c.want[0])
			}
			if int(m) != c.want[1] {
				t.Errorf("visible month %d doesn't match expected %d", m, c.want[1])
			}
			if d != c.want[2] {
				t.Errorf("visible day %d doesn't match expected %d", d, c.want[2])
			}
		})
	}
}
