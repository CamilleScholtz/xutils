package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/go2c/optparse"
)

func main() {
	// Define valid arguments.
	args := optparse.Bool("select", 's', false)
	argh := optparse.Bool("help", 'h', false)

	// Parse arguments.
	_, err := optparse.Parse()
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

	// Connect to the X server using the DISPLAY environment variable.
	X, err := xgbutil.NewConn()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get the geometry of the RootWin.
	g, err := xproto.GetGeometry(X.Conn(), xproto.Drawable(X.RootWin())).Reply()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *args {
		cmd := exec.Command("slop", "-b", "10", "-c", "0.952,0.952,0.952")
		var b bytes.Buffer
		cmd.Stdout = &b

		if err := cmd.Run(); err != nil {
			os.Exit(1)
		}

		s := bufio.NewScanner(strings.NewReader(b.String()))
		for s.Scan() {
			if s.Text()[0] != 'G' {
				continue
			}

			// Get selection geometry.
			f := func(c rune) bool {
				return !unicode.IsNumber(c)
			}
			sg := strings.FieldsFunc(s.Text(), f)

			x, _ := strconv.Atoi(sg[2])
			y, _ := strconv.Atoi(sg[3])
			width, _ := strconv.Atoi(sg[0])
			height, _ := strconv.Atoi(sg[1])
			g.X, g.Y, g.Width, g.Height = int16(x), int16(y), uint16(width), uint16(height)
		}
	}

	fmt.Println(g.X, g.Y, g.Width, g.Height)
	// Get the image data of the pixmap.
	pix, err := xproto.GetImage(X.Conn(), xproto.ImageFormatZPixmap, xproto.Drawable(X.RootWin()),
		g.X, g.Y, g.Width, g.Height, (1<<32)-1).Reply()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	src := pix.Data
	dat := make([]uint8, int(g.Width)*int(g.Height)*4)
	var o int
	for row := 0; row < int(g.Height); row++ {
		for col := 0; col < int(g.Width); col++ {
			o = (row*int(g.Width) + col) * 4
			dat[o+0] = src[2]
			dat[o+1] = src[1]
			dat[o+2] = src[0]
			dat[o+3] = 0xFF

			src = src[4:]
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, int(g.Width), int(g.Height)))
	img.Pix = dat

	f, err := os.Create("test.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	png.Encode(f, img)
}
