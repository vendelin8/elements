## Elements

elements is a CLI tool to show Chinese zodiac five elements as a pageable calendar, written in Golang

The project uses `task` to make your life easier. If you're not familiar with Taskfiles you can take a look at [this quickstart guide](https://taskfile.dev/).

## Overview
Chinese zodiac has five elements: wood, fire, earth, metal and water. You can find unlimited number of documents on these.
Chinese zodiac has three treasures:

-jing for the physical body
-qi for the horizontal world: humans, animals, connections, all kinds of energy: eg. breath and food, communication, etc.
-shen for a connection with somewhat spiritual

Every day has an element that effects our life in many ways, in all three levels. This little app shows which ones.

The rest is upon you. You can find yourself a guru/master/consultant who can calculate your own formula based on your birth date and check out how these effects impact you.

I take no responsibility if your life goes up or down based on this. Use with attention and responsibility.

## Usage
```bash
task run
```

Then you can move in the calendar
-a single day by Right and Left cursor key
-one month by Up and Down cursor key
-one year by PageUp and PageDown key

It uses local time zone by default. If you want to change that, overwrite `TIMEZONE` at the top of `Taskfile.yml` or call it like that.

## Test & lint

Run linting

```bash
task lint
```

Run tests

```bash
task test
```
