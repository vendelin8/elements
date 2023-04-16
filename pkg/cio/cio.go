// Package cio implements a minimalistic console interface.
//
// It accepts cursor keys to move (optionally) in a 3 dimensional array.
// It also maintains a live changing output.
// For example usage check out the upper elements package.
package cio

import (
	"fmt"
	"sync"

	"github.com/gosuri/uilive"
	"github.com/vendelin8/keyboard"
)

const (
	bufferSize = 10
	up         = 1
	down       = -1
)

const (
	NoAct   int = iota
	ActLvl0     // dimension 0 actions for PgUp and PgDn
	ActLvl1     // dimension 1 actions for Up and Down
	ActLvl2     // dimension 2 actions for Right and Left
	ActQuit     // quit
)

var textLock sync.Mutex

// Main maintains a live changing output based on the given callback functions.
// The first one updates the data to print with two integers. The second one may
// be 1 or -1 to move in the given dimension, the first one defines the dimension
// with 0, 1 and 2 for pageUp-pageDown, up-down, left-right respectively.
func Main(changer func(int, int), liner func() string, initiated chan<- struct{}) {
	keyCh, err := keyboard.GetKeys(bufferSize)
	must("getting keyboard events", err)
	defer func() {
		must("closing keyboard", keyboard.Close())
	}()
	w := uilive.New()
	w.Start()
	defer w.Stop()
	fmt.Fprintln(w, liner())
	w.Flush()
	if initiated != nil {
		initiated <- struct{}{}
	}
	for {
		act, val := dispatchKey(<-keyCh)
		switch act {
		case ActQuit:
			return
		case ActLvl0, ActLvl1, ActLvl2:
			changer(act, val)
			textLock.Lock()
			fmt.Fprintln(w, liner())
			w.Flush()
			textLock.Unlock()
		}
	}
}

// dispatchKey returns the needed action for a keyboard shortcut.
//
//nolint:exhaustive
func dispatchKey(ev keyboard.KeyEvent) (int, int) {
	must("dispatching keyboard event", ev.Err)
	if ev.Rune == 'q' {
		return ActQuit, 0
	}
	if ev.Rune != 0 {
		return NoAct, 0
	}
	switch ev.Key {
	case keyboard.KeyArrowRight:
		return ActLvl2, up
	case keyboard.KeyArrowLeft:
		return ActLvl2, down
	case keyboard.KeyArrowUp:
		return ActLvl1, up
	case keyboard.KeyArrowDown:
		return ActLvl1, down
	case keyboard.KeyPgup:
		return ActLvl0, up
	case keyboard.KeyPgdn:
		return ActLvl0, down
	case keyboard.KeyEsc:
		return ActQuit, 0
	}
	return NoAct, 0
}

func must(descr string, err error) {
	if err != nil {
		fmt.Println("error happened with", descr)
		panic(err)
	}
}
