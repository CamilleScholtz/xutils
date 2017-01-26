package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/go2c/optparse"
)

func main() {
	// Define valid arguments.
	args := optparse.Bool("select", 's', false)
	argh := optparse.Bool("help", 'h', false)

	// Parse arguments.
	vals, err := optparse.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invaild argument, use -h for a list of arguments!")
		os.Exit(1)
	}

	// Print help.
	if *argh {
		fmt.Println("Usage: xscrot [arguments] [location]")
		fmt.Println("")
		fmt.Println("arguments:")
		fmt.Println("  -s,   --select          make selection screenshot")
		fmt.Println("  -h,   --help            print help and exit")
		os.Exit(0)
	}
	fmt.Println(vals)

	// Connect to the X server using the DISPLAY environment variable.
	X, err := xgbutil.NewConn()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create selection.
	var sel xproto.Drawable
	if *args {
		// TODO: make this a selection
		sel = xproto.Drawable(X.RootWin())
	} else {
		sel = xproto.Drawable(X.RootWin())
	}

	// Get the geometry of the pixmap
	g, err := xwindow.RawGeometry(X, sel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get the image data of the pixmap.
	pix, err := xproto.GetImage(X.Conn(), xproto.ImageFormatZPixmap, sel, 0, 0,
		uint16(g.Width()), uint16(g.Height()), (1<<32)-1).Reply()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	src := pix.Data
	dat := make([]uint8, g.Width()*g.Height()*4)
	var o int
	for row := 0; row < g.Height(); row++ {
		for col := 0; col < g.Width(); col++ {
			o = (row*g.Width() + col) * 4
			dat[o+0] = src[2]
			dat[o+1] = src[1]
			dat[o+2] = src[0]
			dat[o+3] = 0xFF

			src = src[4:]
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, g.Width(), g.Height()))
	img.Pix = dat

	f, err := os.Create("test.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	png.Encode(f, img)
}
