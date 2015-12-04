# FrequonInvaders
This is a a port of [Frequon Invaders](http://www.blonzonics.us/games/frequon-invaders)
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
4. Run `go build -tags=release`

Use just `go build` to build a developer version, which has keyboard
shortcuts for testing and profiling support.

## Status (2015-Dec-4)

* Works on Windows 8 (Intel 64 processor).
* Works on MacOS 10.11.1 -- released as [Frequon Invaders 2.2](http://www.blonzonics.us/games/frequon-invaders).  

Please post issues for features that you think are missing that were in the classic version. 
