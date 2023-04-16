package cio

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/gosuri/uilive"
	"github.com/micmonay/keybd_event"
)

// TestDispatchKey tests keyboard handler.
func TestDispatchKey(t *testing.T) {
	cases := []struct {
		name    string
		ev      keyboard.KeyEvent
		wantErr bool
		wantAct int
		wantVal int
	}{
		{
			name:    "error",
			ev:      keyboard.KeyEvent{Err: strconv.ErrSyntax},
			wantErr: true,
		},
		{
			name:    "quitQ",
			ev:      keyboard.KeyEvent{Rune: 'q'},
			wantAct: ActQuit,
		},
		{
			name:    "quitEsc",
			ev:      keyboard.KeyEvent{Key: keyboard.KeyEsc},
			wantAct: ActQuit,
		},
		{
			name:    "lvl0Up",
			ev:      keyboard.KeyEvent{Key: keyboard.KeyPgup},
			wantAct: ActLvl0,
			wantVal: up,
		},
		{
			name:    "lvl0Down",
			ev:      keyboard.KeyEvent{Key: keyboard.KeyPgdn},
			wantAct: ActLvl0,
			wantVal: down,
		},
		{
			name:    "lvl1Up",
			ev:      keyboard.KeyEvent{Key: keyboard.KeyArrowUp},
			wantAct: ActLvl1,
			wantVal: up,
		},
		{
			name:    "lvl1Down",
			ev:      keyboard.KeyEvent{Key: keyboard.KeyArrowDown},
			wantAct: ActLvl1,
			wantVal: down,
		},
		{
			name:    "lvl2Up",
			ev:      keyboard.KeyEvent{Key: keyboard.KeyArrowRight},
			wantAct: ActLvl2,
			wantVal: up,
		},
		{
			name:    "lvl2Down",
			ev:      keyboard.KeyEvent{Key: keyboard.KeyArrowLeft},
			wantAct: ActLvl2,
			wantVal: down,
		},
		{
			name:    "notUsedRune",
			ev:      keyboard.KeyEvent{Rune: 'a'},
			wantAct: NoAct,
		},
		{
			name:    "notUsedKey",
			ev:      keyboard.KeyEvent{Key: keyboard.KeyTab},
			wantAct: NoAct,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r != nil != c.wantErr {
					if c.wantErr {
						t.Error("it should have failed, but it didn't")
					} else {
						t.Errorf("it should NOT have failed, but it did with %v", r)
					}
				}
			}()
			act, val := dispatchKey(c.ev)
			if act != c.wantAct {
				t.Errorf("action result %d doesn't match expected %d", act, c.wantAct)
			}
			if val != c.wantVal {
				t.Errorf("action result %d doesn't match expected %d", val, c.wantVal)
			}
		})
	}
}

// TestMain tests initializing and main loop.
func TestMain(t *testing.T) {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		fmt.Println("initializing keybd_event should NOT have failed, but it did:", err)
		return // probably breaks the whole test, but it should be considered as fail
	}
	re := regexp.MustCompile(`\d+$`)

	cases := []struct {
		name string
		ev   int
		want int
	}{
		{
			name: "notUsedRune",
			ev:   keybd_event.VK_A,
			// want: 0 as from the start
		},
		{
			name: "lvl0Up",
			ev:   keybd_event.VK_PAGEUP,
			want: 4,
		},
		{
			name: "lvl1Up",
			ev:   keybd_event.VK_UP,
			want: 6,
		},
		{
			name: "notUsedKey",
			ev:   keybd_event.VK_TAB,
			want: 6,
		},
		{
			name: "lvl2Up",
			ev:   keybd_event.VK_RIGHT,
			want: 7,
		},
	}

	current := 0
	changer := func(act int, value int) {
		switch act {
		case ActLvl0:
			current += value << 2
		case ActLvl1:
			current += value << 1
		case ActLvl2:
			current += value
		}
	}
	liner := func() string {
		return strconv.Itoa(current)
	}

	var b strings.Builder
	uilive.Out = &b
	var wg sync.WaitGroup
	initiated := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		Main(changer, liner, initiated)
	}()

	// wait as written here: https://github.com/micmonay/keybd_event/blob/master/README.md
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	<-initiated

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			kb.SetKeys(c.ev)
			err = kb.Launching()
			if err != nil {
				t.Errorf("pressing key should NOT have failed, but it did with %v", err)
			}
			time.Sleep(200 * time.Millisecond)
			wantStr := strconv.Itoa(c.want)
			if s := re.FindString(strings.TrimSpace(b.String())); s != wantStr {
				t.Errorf("result '%s' doesn't match expected '%s'", s, wantStr)
			}
		})
	}
	kb.SetKeys(keybd_event.VK_ESC)
	err = kb.Launching()
	if err != nil {
		t.Errorf("pressing key should NOT have failed, but it did with %v", err)
	}
	wg.Wait()
}
