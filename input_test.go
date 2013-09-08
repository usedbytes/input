/*
 *
 * input_test.go - Tests for input.go
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
    "time"
    "testing"
    "github.com/gvalkov/golang-evdev"
)

func TestSomething(t *testing.T) {
    fmt.Println("Trying all input devices")
    devices := FindInputDevices();

    for _, d := range devices {
        id, error := evdev.Open(d);
        if (error != nil) {
            fmt.Println(error.Error());
        } else {
            fmt.Println(id.String());
            fmt.Print("\n\n");
        }
    }
    fmt.Println("Done");
}

/*
func TestManager(t *testing.T) {
    in := make(chan int);
    out := make(chan *evdev.InputEvent, 10);
    ManageInput( in, out);
}
*/

func TestProcessor(t *testing.T) {
    events := make(chan *evdev.InputEvent)
    StartListening(events)

    quit := make(chan int)
    keys := make(chan rune)
    go ProcessInputEvents(events, keys, quit)
    for {
        select {
        case k := <-keys:
            fmt.Printf("%#U ", rune(k))
        case <-quit:
            fmt.Println("Quitting")
            return
        }
    }
}

func NoTestListener(t *testing.T) {

    quit := make(chan int);
    go func() {
        time.Sleep(100000 * time.Millisecond)
        quit <- 1
    }()

    events := make(chan *evdev.InputEvent)
    StartListening(events)
    for {
        select {
        case ev := <-events:
            fmt.Println(ev.String(), evdev.KEY[int(ev.Code)])
        case <-quit:
            fmt.Println("Quitting")
            return
        }
    }
}

func NoTestKeyboard(t *testing.T) {
    id, error := evdev.Open("/dev/input/event3");

    go func(dev *evdev.InputDevice) {
        time.Sleep(5000 * time.Millisecond);
        dev.File.Close();
    }(id);

    if (error != nil) {
        fmt.Println(error.Error());
    } else {
        fmt.Println(id.String());
        for {
            ev, error := id.ReadOne();
            if (error != nil) {
                fmt.Println(error.Error());
                return;
            }
            fmt.Println(ev.String());
        }
    }
}
