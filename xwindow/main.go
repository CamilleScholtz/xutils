package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/go2c/optparse"
)

func main() {
	// Define valid arguments.
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
		fmt.Println("Usage: xwindow [arguments]")
		fmt.Println("")
		fmt.Println("arguments:")
		fmt.Println("  -w,   --watch           watch active window")
		fmt.Println("  -h,   --help            print help and exit")
		os.Exit(0)
	}

	// Connect to the X server using the DISPLAY environment variable.
	X, err := xgbutil.NewConn()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get the active window ID.
	w, err := ewmh.ActiveWindowGet(X)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get active window name.
	n, err := ewmh.WmNameGet(X, w)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(os.Args) == 1 {
		fmt.Println(n)
		os.Exit(0)
	}

	// Print the current desktop number.
	if *argw {
		fmt.Println(n)
		r := xwindow.New(X, X.RootWin())
		r.Listen(xproto.EventMaskPropertyChange)
		var on string
		xevent.PropertyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
			a, err := xprop.AtomName(XU, ev.Atom)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if a == "_NET_ACTIVE_WINDOW" {
				// Get the active window ID.
				w, err := ewmh.ActiveWindowGet(X)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				// Get active window name.
				n, err := ewmh.WmNameGet(X, w)
				if err != nil {
					return
				}

				if n != on {
					fmt.Println(n)
				}
				on = n
			}
		}).Connect(X, r.Id)
		xevent.Main(X)

		os.Exit(0)
	}
}
