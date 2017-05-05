package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgbutil"
	"github.com/go2c/optparse"
)

func main() {
	// Define valid arguments.
	argh := optparse.Bool("help", 'h', false)

	// Parse arguments.
	vals, err := optparse.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr,
			"Invaild argument, use -h for a list of arguments!")
		os.Exit(1)
	}

	// Print help.
	if *argh {
		fmt.Println("Usage: xcomf [arguments] [temperature]")
		fmt.Println("")
		fmt.Println("arguments:")
		fmt.Println("  -h,   --help            print help and exit")
		os.Exit(0)
	}

	// Connect to the X server using the DISPLAY environment variable.
	X, err := xgbutil.NewConn()
	if err != nil {
		panic(err)
	}

	// Initialize randr.
	if err := randr.Init(X.Conn()); err != nil {
		panic(err)
	}

	sr, err := randr.GetScreenResources(X.Conn(), X.RootWin()).Reply()
	if err != nil {
		panic(err)
	}

	r, err := randr.GetCrtcGamma(X.Conn(), sr.Crtcs[0]).Reply()
	if err != nil {
		panic(err)
	}

	// Reset gamma to default value if no vals have been given.
	if len(vals) == 0 {
		randr.SetCrtcGamma(X.Conn(), sr.Crtcs[0], r.Size, r.Red,
			r.Green, r.Blue)
		os.Exit(0)
	}

	// Only alow intergers in vals.
	temp, err := strconv.Atoi(vals[0])
	if err != nil || temp < 1 || temp > 10 {
		fmt.Fprintln(os.Stderr,
			"Please choose a temperature between 1 and 10!")
		os.Exit(1)
	}

	var nr []uint16
	for i := range r.Red {
		v := int(r.Red[i])
		v += int(r.Red[len(r.Red)-i-1]) / 8
		v -= (len(r.Red) / 100 * (i * 2)) / 4
		v = ((v * temp) + (int(r.Red[i]) * (10 - temp))) / 10

		nr = append(nr, uint16(v))
	}
	var ng []uint16
	for i := range r.Green {
		v := int(r.Green[i])
		v -= (len(r.Green) / 100 * (i * 2)) / 3
		v = ((v * temp) + (int(r.Green[i]) * (10 - temp))) / 10

		ng = append(ng, uint16(v))
	}
	var nb []uint16
	for i := range r.Blue {
		v := int(r.Blue[i])
		v -= (len(r.Blue) / 100 * (i * 2)) / 3
		v = ((v * temp) + (int(r.Blue[i]) * (10 - temp))) / 10

		nb = append(nb, uint16(v))
	}

	randr.SetCrtcGamma(X.Conn(), sr.Crtcs[0], r.Size, nr, ng, nb)
}
