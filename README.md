# Magellan

Space flight imitation for LARP game.

Using ebiten engine (https://github.com/hajimehoshi/ebiten). Thanks you, Hajime Hoshi!

Using diskv for disk storage and packr for embedding static data

# how to install

you need go compiler https://dl.google.com/go/go1.10.3.windows-386.msi
to compile non-server parts also needs mingw32 https://sourceforge.net/projects/mingw/files/Installer/mingw-get-setup.exe/download

maybe you will need to set %PATH% variable

use command to get server part:

`go get -v github.com/Shnifer/magellan/execs/server/...`

or use to get all project:

`go get -v github.com/Shnifer/magellan/...`

then use from %UserName%/go/src/github.com/Shnifer/magellan/execs/server

`go install`
