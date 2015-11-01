# FrequonInvaders
This is a work in progress to port [Frequon Invaders](http://www.blonzonics.us/games/frequon-invaders)
from C++ to the Go Language

## Prerequisites
* [Go](https://golang.org/doc/install)
* SDL-2
* [Go bindings to SDL-2](https://github.com/veandco/go-sdl2)

The instructions for the Go bindings describe how to install SDL-2.

## To Build 
1. Set your `GOPATH`
2. Run 
```
go get github.com/ArchRobison/FrequonInvaders
```
This should fetch the three repositories and put them relative to your `GOPATH` as follows:
```
src/github.com/ArchRobison/FrequonInvaders
src/github.com/ArchRobison/Gophetica
src/github.com/veandco/go-sdl2
```
3. cd to `src/github.com/ArchRobison/FrequonInvaders`
4. Run `go build`

## Status (2015-Oct-31)

Works on Windows 8.  Not tested on other platforms.  

Please post issues for features that you think are missing from the classic version. 
