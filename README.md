<a href="https://imgur.com/4qYolKu">
    <img src="https://i.imgur.com/4qYolKu.gif">
</a>

**satellite** by [*@avinashbot*](https://github.com/avinashbot)  
![Travis](https://img.shields.io/travis/avinashbot/satellite.svg?style=flat-square)
![Release](https://img.shields.io/github/release/avinashbot/satellite.svg?style=flat-square)
![License](https://img.shields.io/github/license/mashape/apistatus.svg?style=flat-square)

**Building from source**

1. Download and install [Go](https://golang.org/dl/).
2. `go get github.com/avinashbot/satellite`
3. `go build github.com/avinashbot/satellite`
4. `./satellite --help`

**Usage:**

```bash
$ satellite -use himawari -depth 4 -every 60s
$ satellite -use dscovr -path ./images/dscovr/1.png
$ satellite -depth 20 # This will result in a massive image!
```
