// Package main (elements) shows.
//
// It accepts cursor keys to move (optionally) in a 3 dimensional array.
// It also maintains a live changing output.
// For example usage check out the upper elements package.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/vendelin8/elements/pkg/cio"
	"github.com/vendelin8/elements/pkg/model"
)

// visibleDay is the date currently shown.
var visibleDay = time.Now().Truncate(time.Hour * model.DailyHours)

type Config struct {
	location string
}

// parseFlags parses the command-line arguments provided to the program.
// Typically os.Args[0] is provided as 'progname' and os.args[1:] as 'args'.
// Returns the Config.
func parseFlags(args ...string) (*Config, error) {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flag.CommandLine = flags
	var conf Config
	flag.StringVar(&conf.location, "loc", "", "A time zone to calculate, local time if omitted.")
	err := flags.Parse(args[1:])
	if conf.location == "-" {
		conf.location = "" // dash means local timezone
	}
	return &conf, err
}

func main() {
	conf, err := parseFlags(os.Args...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if exitCode := handleFlags(conf, getDayNow); exitCode != 0 {
		os.Exit(exitCode)
	}
	cio.Main(changeDay, fillElements, nil)
}

func getDayNow(loc *time.Location) time.Time {
	return time.Now().In(loc)
}

func handleFlags(conf *Config, getDay func(*time.Location) time.Time) int {
	if len(conf.location) == 0 {
		return 0
	}
	loc, err := time.LoadLocation(conf.location)
	if err != nil {
		fmt.Println("Failed to parse argument as time zone: ", err)
		fmt.Println("For valid values please check https://stackoverflow.com/a/40130882/6155997")
		return 1
	}
	visibleDay = getDay(loc)
	y, m, d := visibleDay.Date()
	visibleDay = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	return 0
}

// changeDay moves currently visible day by one year, month or day based on the
// given level/dimension.
func changeDay(level, diff int) {
	switch level {
	case cio.ActLvl0:
		visibleDay = visibleDay.AddDate(diff, 0, 0)
	case cio.ActLvl1:
		visibleDay = visibleDay.AddDate(0, diff, 0)
	case cio.ActLvl2:
		visibleDay = visibleDay.AddDate(0, 0, diff)
	}
}

// fillElements returns currently visible day as a colored string with the daily elements.
func fillElements() string {
	jing, qi, shen := model.GetElements(visibleDay)
	return fmt.Sprintf("%s: shen: %s, qi: %s, jing: %s", visibleDay.Format("2006-01-02"),
		shen.Color(), qi.Color(), jing.Color())
}
