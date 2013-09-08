/*
 *
 * input.go - Functions for listening to all evdev sources
 *   Copyright Brian Starkey 2013-2014 <stark3y@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of version 2 of the GNU General Public License as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package input

import (
	"fmt"
	"os"
	"github.com/gvalkov/golang-evdev"
)

const (
	KEY_HOME =       102 | 0x100
	KEY_UP =         103 | 0x100
	KEY_PAGEUP =     104 | 0x100
	KEY_LEFT =       105 | 0x100
	KEY_RIGHT =      106 | 0x100
	KEY_END =        107 | 0x100
	KEY_DOWN =       108 | 0x100
	KEY_PAGEDOWN =   109 | 0x100
	KEY_INSERT =     110 | 0x100
	KEY_DELETE =     111 | 0x100
	KEY_MUTE =       113 | 0x100
	KEY_VOLUMEDOWN = 114 | 0x100
	KEY_VOLUMEUP =   115 | 0x100
	KEY_POWER =      116 | 0x100
	KEY_PAUSE =      119 | 0x100
	KEY_STOP =       128 | 0x100
	KEY_MENU =       139 | 0x100
	KEY_BACK =       158 | 0x100
	KEY_FORWARD =    159 | 0x100
	KEY_NEXTSONG =   163 | 0x100
	KEY_PLAYPAUSE =  164 | 0x100
	KEY_PREVIOUSSONG = 165 | 0x100
	KEY_REWIND =     168 | 0x100
	KEY_SCROLLUP =   177 | 0x100
	KEY_SCROLLDOWN = 178 | 0x100
	KEY_PLAY =       207 | 0x100
	KEY_FASTFORWARD = 208 | 0x100
	KEY_BASSBOOST =  209 | 0x100
	KEY_SEARCH =     217 | 0x100
	KEY_ENTER =      13
	KEY_ESC =        27
	KEY_BACKSPACE =  8
)

func FindInputDevices() ([]string) {
	allDevs, error := evdev.ListInputDevices("/dev/input/event*")
	if (error != nil) {
		fmt.Println(error.Error())
	}

	availableDevs := make([]string, 0, len(allDevs))
	for _, dev := range allDevs {
		fp, error := os.Open(dev)
		if (error == nil) {
			fp.Close()
			availableDevs = append(availableDevs, dev)
		}
	}

	return availableDevs
}

func StartListening(out chan *evdev.InputEvent) {

	devs := FindInputDevices()

	for _, dev := range devs {
		id, error := evdev.Open(dev)
		if (error == nil) {
			go func() {
				for {
					ev, err := id.ReadOne()
					if (err != nil) {
						fmt.Println(err.Error())
						return;
					}
					out <- ev
				}
			}()
		}
	}
}

func ProcessInputEvents(in chan *evdev.InputEvent, out chan rune, quit chan int) {
	var shift bool = false
	var shiftMap = getShiftMap()
	var keyMap = getKeyMap()

	for {
		select {
		case ev := <-in:
			if (ev.Type == evdev.EV_KEY) {
				if (ev.Code == evdev.KEY_LEFTSHIFT) ||
					(ev.Code == evdev.KEY_RIGHTSHIFT) {
					if (ev.Value == 0) || (ev.Value == 1) {
						shift = !shift
					}
				} else if (ev.Code == evdev.KEY_CAPSLOCK) {
					if (ev.Value == 1) {
						shift = !shift
					}
				}

				if (ev.Value > 0) {
					if (shift) {
						out <- shiftMap[ev.Code]
					} else {
						out <- keyMap[ev.Code]
					}
				}
			} else if (ev.Type == evdev.EV_REL) && (ev.Code == evdev.REL_WHEEL) {
				if (ev.Value > 0) {
					out <- keyMap[evdev.KEY_SCROLLUP]
				} else if (ev.Value < 0) {
					out <- keyMap[evdev.KEY_SCROLLDOWN]
				}
			}

		case <-quit:
			fmt.Println("Quitting")
			return
		}
	}

}

func getKeyMap() map[uint16]rune {

	return map[uint16]rune {
		evdev.KEY_ESC: 27,
		evdev.KEY_1: '1',
		evdev.KEY_2: '2',
		evdev.KEY_3: '3',
		evdev.KEY_4: '4',
		evdev.KEY_5: '5',
		evdev.KEY_6: '6',
		evdev.KEY_7: '7',
		evdev.KEY_8: '8',
		evdev.KEY_9: '9',
		evdev.KEY_0: '0',
		evdev.KEY_MINUS: '-',
		evdev.KEY_EQUAL: '=',
		evdev.KEY_BACKSPACE: 8,
		evdev.KEY_TAB: '\t',
		evdev.KEY_Q: 'q',
		evdev.KEY_W: 'w',
		evdev.KEY_E: 'e',
		evdev.KEY_R: 'r',
		evdev.KEY_T: 't',
		evdev.KEY_Y: 'y',
		evdev.KEY_U: 'u',
		evdev.KEY_I: 'i',
		evdev.KEY_O: 'o',
		evdev.KEY_P: 'p',
		evdev.KEY_LEFTBRACE: '[',
		evdev.KEY_RIGHTBRACE: ']',
		evdev.KEY_ENTER: 13,
		evdev.KEY_A: 'a',
		evdev.KEY_S: 's',
		evdev.KEY_D: 'd',
		evdev.KEY_F: 'f',
		evdev.KEY_G: 'g',
		evdev.KEY_H: 'h',
		evdev.KEY_J: 'j',
		evdev.KEY_K: 'k',
		evdev.KEY_L: 'l',
		evdev.KEY_SEMICOLON: ';',
		evdev.KEY_APOSTROPHE: '\'',
		evdev.KEY_BACKSLASH: '#',
		evdev.KEY_Z: 'z',
		evdev.KEY_X: 'x',
		evdev.KEY_C: 'c',
		evdev.KEY_V: 'v',
		evdev.KEY_B: 'b',
		evdev.KEY_N: 'n',
		evdev.KEY_M: 'm',
		evdev.KEY_COMMA: ',',
		evdev.KEY_DOT: '.',
		evdev.KEY_SLASH: '/',
		evdev.KEY_KPASTERISK: '*',
		evdev.KEY_SPACE: ' ',
		evdev.KEY_KP7: '7',
		evdev.KEY_KP8: '8',
		evdev.KEY_KP9: '9',
		evdev.KEY_KPMINUS: '-',
		evdev.KEY_KP4: '4',
		evdev.KEY_KP5: '5',
		evdev.KEY_KP6: '6',
		evdev.KEY_KPPLUS: '+',
		evdev.KEY_KP1: '1',
		evdev.KEY_KP2: '2',
		evdev.KEY_KP3: '3',
		evdev.KEY_KP0: '0',
		evdev.KEY_KPDOT: '.',
		evdev.KEY_KPENTER: 13,
		evdev.KEY_KPSLASH: '/',
		evdev.KEY_HOME:       102 | 0x100,
		evdev.KEY_UP:         103 | 0x100,
		evdev.KEY_PAGEUP:     104 | 0x100,
		evdev.KEY_LEFT:       105 | 0x100,
		evdev.KEY_RIGHT:      106 | 0x100,
		evdev.KEY_END:        107 | 0x100,
		evdev.KEY_DOWN:       108 | 0x100,
		evdev.KEY_PAGEDOWN:   109 | 0x100,
		evdev.KEY_INSERT:     110 | 0x100,
		evdev.KEY_DELETE:     111 | 0x100,
		evdev.KEY_MUTE:       113 | 0x100,
		evdev.KEY_VOLUMEDOWN: 114 | 0x100,
		evdev.KEY_VOLUMEUP:   115 | 0x100,
		evdev.KEY_POWER:      116 | 0x100,
		evdev.KEY_KPEQUAL: '=',
		evdev.KEY_PAUSE:      119 | 0x100,
		evdev.KEY_KPCOMMA: ',',
		evdev.KEY_STOP:       128 | 0x100,
		evdev.KEY_MENU:       139 | 0x100,
		evdev.KEY_BACK:       158 | 0x100,
		evdev.KEY_FORWARD:    159 | 0x100,
		evdev.KEY_NEXTSONG:   163 | 0x100,
		evdev.KEY_PLAYPAUSE:  164 | 0x100,
		evdev.KEY_PREVIOUSSONG: 165 | 0x100,
		evdev.KEY_REWIND:     168 | 0x100,
		evdev.KEY_SCROLLUP:   177 | 0x100,
		evdev.KEY_SCROLLDOWN: 178 | 0x100,
		evdev.KEY_KPLEFTPAREN: '(',
		evdev.KEY_KPRIGHTPAREN: ')',
		evdev.KEY_PLAY:       207 | 0x100,
		evdev.KEY_FASTFORWARD: 208 | 0x100,
		evdev.KEY_BASSBOOST: 209 | 0x100,
		evdev.KEY_SEARCH:     217 | 0x100,
		evdev.KEY_OK: 13,
		evdev.KEY_SELECT: 13,
		evdev.KEY_NUMERIC_0: '0',
		evdev.KEY_NUMERIC_1: '1',
		evdev.KEY_NUMERIC_2: '2',
		evdev.KEY_NUMERIC_3: '3',
		evdev.KEY_NUMERIC_4: '4',
		evdev.KEY_NUMERIC_5: '5',
		evdev.KEY_NUMERIC_6: '6',
		evdev.KEY_NUMERIC_7: '7',
		evdev.KEY_NUMERIC_8: '8',
		evdev.KEY_NUMERIC_9: '9',
		evdev.KEY_NUMERIC_STAR: '*',
		evdev.KEY_NUMERIC_POUND: '#',
	}
}

func getShiftMap() map[uint16]rune {

	return map[uint16]rune {
		evdev.KEY_ESC: 27,
		evdev.KEY_1: '!',
		evdev.KEY_2: '"',
		evdev.KEY_3: '£',
		evdev.KEY_4: '$',
		evdev.KEY_5: '%',
		evdev.KEY_6: '^',
		evdev.KEY_7: '&',
		evdev.KEY_8: '*',
		evdev.KEY_9: '(',
		evdev.KEY_0: ')',
		evdev.KEY_MINUS: '_',
		evdev.KEY_EQUAL: '+',
		evdev.KEY_BACKSPACE: 8,
		evdev.KEY_TAB: '\t',
		evdev.KEY_Q: 'Q',
		evdev.KEY_W: 'W',
		evdev.KEY_E: 'E',
		evdev.KEY_R: 'R',
		evdev.KEY_T: 'T',
		evdev.KEY_Y: 'Y',
		evdev.KEY_U: 'U',
		evdev.KEY_I: 'I',
		evdev.KEY_O: 'O',
		evdev.KEY_P: 'P',
		evdev.KEY_LEFTBRACE: '{',
		evdev.KEY_RIGHTBRACE: '}',
		evdev.KEY_ENTER: 13,
		evdev.KEY_A: 'A',
		evdev.KEY_S: 'S',
		evdev.KEY_D: 'D',
		evdev.KEY_F: 'F',
		evdev.KEY_G: 'G',
		evdev.KEY_H: 'H',
		evdev.KEY_J: 'J',
		evdev.KEY_K: 'K',
		evdev.KEY_L: 'L',
		evdev.KEY_SEMICOLON: ':',
		evdev.KEY_APOSTROPHE: '@',
		evdev.KEY_BACKSLASH: '~',
		evdev.KEY_Z: 'Z',
		evdev.KEY_X: 'X',
		evdev.KEY_C: 'C',
		evdev.KEY_V: 'V',
		evdev.KEY_B: 'B',
		evdev.KEY_N: 'N',
		evdev.KEY_M: 'M',
		evdev.KEY_COMMA: '<',
		evdev.KEY_DOT: '>',
		evdev.KEY_SLASH: '?',
		evdev.KEY_KPASTERISK: '*',
		evdev.KEY_SPACE: ' ',
		evdev.KEY_KP7: '7',
		evdev.KEY_KP8: '8',
		evdev.KEY_KP9: '9',
		evdev.KEY_KPMINUS: '-',
		evdev.KEY_KP4: '4',
		evdev.KEY_KP5: '5',
		evdev.KEY_KP6: '6',
		evdev.KEY_KPPLUS: '+',
		evdev.KEY_KP1: '1',
		evdev.KEY_KP2: '2',
		evdev.KEY_KP3: '3',
		evdev.KEY_KP0: '0',
		evdev.KEY_KPDOT: '.',
		evdev.KEY_KPENTER: 13,
		evdev.KEY_KPSLASH: '/',
		evdev.KEY_HOME:	   102 | 0x100,
		evdev.KEY_UP:		 103 | 0x100,
		evdev.KEY_PAGEUP:	 104 | 0x100,
		evdev.KEY_LEFT:	   105 | 0x100,
		evdev.KEY_RIGHT:	  106 | 0x100,
		evdev.KEY_END:		107 | 0x100,
		evdev.KEY_DOWN:	   108 | 0x100,
		evdev.KEY_PAGEDOWN:   109 | 0x100,
		evdev.KEY_INSERT:	 110 | 0x100,
		evdev.KEY_DELETE:	 111 | 0x100,
		evdev.KEY_MUTE:	   113 | 0x100,
		evdev.KEY_VOLUMEDOWN: 114 | 0x100,
		evdev.KEY_VOLUMEUP:   115 | 0x100,
		evdev.KEY_POWER:	  116 | 0x100,
		evdev.KEY_KPEQUAL: '=',
		evdev.KEY_PAUSE:	  119 | 0x100,
		evdev.KEY_KPCOMMA: ',',
		evdev.KEY_STOP:	   128 | 0x100,
		evdev.KEY_MENU:	   139 | 0x100,
		evdev.KEY_BACK:	   158 | 0x100,
		evdev.KEY_FORWARD:	159 | 0x100,
		evdev.KEY_NEXTSONG:   163 | 0x100,
		evdev.KEY_PLAYPAUSE:  164 | 0x100,
		evdev.KEY_PREVIOUSSONG: 165 | 0x100,
		evdev.KEY_REWIND:	 168 | 0x100,
		evdev.KEY_SCROLLUP:   177 | 0x100,
		evdev.KEY_SCROLLDOWN: 178 | 0x100,
		evdev.KEY_KPLEFTPAREN: '(',
		evdev.KEY_KPRIGHTPAREN: ')',
		evdev.KEY_PLAY:	   207 | 0x100,
		evdev.KEY_FASTFORWARD: 208 | 0x100,
		evdev.KEY_BASSBOOST: 209 | 0x100,
		evdev.KEY_SEARCH:	 217 | 0x100,
		evdev.KEY_OK: 13,
		evdev.KEY_SELECT: 13,
		evdev.KEY_NUMERIC_0: '0',
		evdev.KEY_NUMERIC_1: '1',
		evdev.KEY_NUMERIC_2: '2',
		evdev.KEY_NUMERIC_3: '3',
		evdev.KEY_NUMERIC_4: '4',
		evdev.KEY_NUMERIC_5: '5',
		evdev.KEY_NUMERIC_6: '6',
		evdev.KEY_NUMERIC_7: '7',
		evdev.KEY_NUMERIC_8: '8',
		evdev.KEY_NUMERIC_9: '9',
		evdev.KEY_NUMERIC_STAR: '*',
		evdev.KEY_NUMERIC_POUND: '#',
	}
}
