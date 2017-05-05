[![Go Report Card](https://goreportcard.com/badge/github.com/onodera-punpun/xutils)](https://goreportcard.com/report/github.com/onodera-punpun/xutils)

xutils - X utilities written in Go, mostly for personal rice and because I have issues with literally
all other x utils.

## SYNOPSIS

xcomf [arguments] [temperature]

xdesktop [arguments]

xscrot [arguments] [location]

xtitle [arguments]


## DESCRIPTION

#### xcomf

Changes the gamma level. Inspired by Redshift and f.lux, but especially
by Red Moon for Android. Unlike Redshift and f.lux do Red Moon and
xcomf not only change the display to be slightly orange, but actually
try to dim blue LEDs. xcomf also slightly lowers the contrast for more
comfy.

#### xdesktop

Changes and print the current desktop (workspace).

### xscrot

Takes screenshot and selection screenshot using
[naelstrof/slop](https://github.com/naelstrof/slop).

#### xtitle

Prints the focuesed window. Pretty much a clone of
[baskerville/xtile](https://github.com/baskerville/xtitle), but in Go.


## COMMANDS

TODO


## INSTALLATION

`go get github.com/onodera-punpun/xutils`


## AUTHORS

Camille Scholtz


## NOTES

Since this is my uhh... third Go project I'm probably making some
mistakes, feedback is highly appreciated!
