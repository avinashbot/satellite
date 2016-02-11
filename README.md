<a href="https://imgur.com/4qYolKu">
    <img src="https://i.imgur.com/4qYolKu.gif">
</a>

## Satellite

### Download
[Latest Release](https://github.com/avinashbot/satellite/releases/latest)

### Building from source
1. Download and install [Go](https://golang.org/dl/).
2. `go get github.com/avinashbot/satellite`
3. `go build github.com/avinashbot/satellite`
4. `./satellite --help`

### Usage:
```bash
$ satellite -use himawari -depth 4 -every 60s
$ satellite -use dscovr -path ./images/dscovr/1.png
$ satellite -depth 20 # This will result in a massive image!
```
