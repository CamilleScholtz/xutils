package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/go2c/optparse"
)

func main() {
	// Define valid arguments.
	argn := optparse.Bool("next", 'n', false)
	argp := optparse.Bool("previous", 'p', false)
	args := optparse.IntList("set", 's')
	argt := optparse.Bool("total", 't', false)
	argw := optparse.Bool("watch", 'w', false)
	argh := optparse.Bool("help", 'h', false)

	// Parse arguments.
	_, err := optparse.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invaild argument, use -h for a list of arguments!")
		os.Exit(1)
	}

	// Print help.
	if *argh {
		fmt.Println("Usage: xdesktop [arguments]")
		fmt.Println("")
		fmt.Println("arguments:")
		fmt.Println("  -n,   --next            switch to next desktop")
		fmt.Println("  -p,   --previous        switch to previous desktop")
		fmt.Println("  -s,   --set             switch to specified desktop")
		fmt.Println("  -t,   --total           print total number of desktops")
		fmt.Println("  -w,   --watch           watch current desktop number")
		fmt.Println("  -h,   --help            print help and exit")
		os.Exit(0)
	}

	// Connect to the X server using the DISPLAY environment variable.
	X, err := xgbutil.NewConn()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get the current desktop number.
	cd, err := ewmh.CurrentDesktopGet(X)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print the current desktop number.
	if *argw || len(os.Args) == 1 {
		fmt.Println(cd + 1)

		if *argw {
			r := xwindow.New(X, X.RootWin())
			r.Listen(xproto.EventMaskPropertyChange)

			xevent.PropertyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
				// Only listen to desktop change events.
				// TODO: Can I somehow do this in r.Listen?
				if ev.Atom != 372 {
					return
				}

				// Get the current desktop number.
				cd, err := ewmh.CurrentDesktopGet(X)
				if err != nil {
					return
				}

				fmt.Println(cd + 1)
			}).Connect(X, r.Id)
			xevent.Main(X)
		}

		os.Exit(0)
	}

	// Get the total number of desktops.
	td, err := ewmh.NumberOfDesktopsGet(X)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print the total number of desktops.
	if *argt {
		fmt.Println(td)
		os.Exit(0)
	}

	var d int

	// Switch to next/previous desktop.
	if *argn || *argp {
		if *argn {
			if cd == td-1 {
				d = 0
			} else {
				d = int(cd) + 1
			}
		} else {
			if cd == 0 {
				d = int(td - 1)
			} else {
				d = int(cd) - 1
			}
		}
	}

	// Switch to specified desktop.
	if len(*args) != 0 {
		d = (*args)[0] - 1
	}

	if err := ewmh.CurrentDesktopReq(X, d); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
