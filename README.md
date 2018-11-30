# Magellan

Space flight imitation for LARP game.

Game is over and this project is finished.

Using [ebiten engine](https://github.com/hajimehoshi/ebiten). Thanks you, Hajime Hoshi!

Using [diskv](https://github.com/peterbourgon/diskv) for disk storage and [packr](https://github.com/gobuffalo/packr) for embedding static data

# how to install

## install compilers

you need [go compiler](https://dl.google.com/go/go1.10.3.windows-386.msi)

to compile non-server parts also need [mingw32](https://sourceforge.net/projects/mingw/files/Installer/mingw-get-setup.exe/download)

maybe you will need to add into %PATH% variable `/go/bin/` , `%UserName%/go/bin/` , `mingw32/bin/`

## install packr

```
go get -u github.com/gobuffalo/packr
go install github.com/gobuffalo/packr
```
## generate static

_%UserName%/go/src/github.com/shnifer/magellan/_
```
packr clean
packr
```

## compile

use command to get server part:

`go get -u -v github.com/shnifer/magellan/execs/server/...`

or use to get all project:

`go get -u -v github.com/shnifer/magellan/...`

then run

```
go install github.com/shnifer/magellan/execs/server/
go install github.com/shnifer/magellan/execs/pilot/
go install github.com/shnifer/magellan/execs/navi/
go install github.com/shnifer/magellan/execs/engi/
```

Congratulations! You got your exes in _%UserName%/go/bin/_
