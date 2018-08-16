package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgbutil"
	"github.com/go2c/optparse"
	"github.com/maruel/temperature"
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
		fmt.Println("Usage: xcomf [arguments] [intensity]")
		fmt.Println("")
		fmt.Println("arguments:")
		fmt.Println("  -h,   --help            print help and exit")
		os.Exit(0)
	}

	// Reset gamma to default value if no vals have been given.
	if len(vals) == 0 {
		vals = []string{"0"}
	}

	// Only alow intergers in vals.
	i, err := strconv.Atoi(vals[0])
	if err != nil || i < 0 || i > 3 {
		fmt.Fprintln(os.Stderr,
			"Please choose an intensity between 0 and 3!")
		os.Exit(1)
	}

	var t uint16
	switch i {
	case 0:
		t = 6500
	case 1:
		t = 4000
	case 2:
		t = 3500
	case 3:
		t = 3000
	}

	if err := set(temperature.ToRGB(t)); err != nil {
		panic(err)
	}
}

func gamma(r, g, b uint8, size uint16, comf bool) ([]uint16, []uint16,
	[]uint16) {
	gr := make([]uint16, size)
	gg := make([]uint16, size)
	gb := make([]uint16, size)

	for i := uint16(0); i < size; i++ {
		fr := float64(r) / 255
		fg := float64(g) / 255
		fb := float64(b) / 255
		fs := float64(size)
		if comf {
			fb *= 1.618
			fs *= 0.85
			fs += float64(i) * 0.5
		}

		gamma := 65535 * float64(i) / fs

		gr[i] = uint16(gamma * fr)
		gg[i] = uint16(gamma * fg)
		gb[i] = uint16(gamma * fb)
	}

	return gr, gg, gb
}

func set(r, g, b uint8) error {
	// Connect to the X server using the DISPLAY environment variable.
	X, err := xgbutil.NewConn()
	if err != nil {
		return err
	}

	// Initialize randr.
	if err := randr.Init(X.Conn()); err != nil {
		return err
	}

	sr, err := randr.GetScreenResourcesCurrent(X.Conn(), X.RootWin()).Reply()
	if err != nil {
		return err
	}

	for _, c := range sr.Crtcs {
		gs, err := randr.GetCrtcGammaSize(X.Conn(), c).Reply()
		if err != nil {
			return err
		}

		var gg, gr, gb []uint16
		if r == 255 && g == 255 && b == 255 {
			gr, gg, gb = gamma(r, g, b, gs.Size, false)
		} else {
			gr, gg, gb = gamma(r, g, b, gs.Size, true)
		}

		if err := randr.SetCrtcGammaChecked(X.Conn(), c, gs.Size, gr, gg, gb).
			Check(); err != nil {
			return err
		}
	}

	return nil
}
