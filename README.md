<a href="https://imgur.com/4qYolKu">
    <img src="https://i.imgur.com/4qYolKu.gif">
</a>

![Travis](https://img.shields.io/travis/avinashbot/satellite.svg?style=flat-square)
![Release](https://img.shields.io/github/release/avinashbot/satellite.svg?style=flat-square)
![License](https://img.shields.io/github/license/mashape/apistatus.svg?style=flat-square)

**satellite** by [*@avinashbot*](https://github.com/avinashbot)  

**Download:**

1. Get the [latest release](https://github.com/avinashbot/satellite/releases/latest)
2. Move it somewhere like ~/bin or /usr/local/bin
3. Use your preferred method to run it on login or on a timer (cron, /etc/profile.d, ...)
4. Enjoy!

**Building from source**

1. Download and install [Go](https://golang.org/dl/).
2. `go get github.com/avinashbot/satellite`
3. `go build github.com/avinashbot/satellite`
4. `./satellite`

**Examples:**

```bash
$ satellite -use himawari -depth 4 -every 60s
$ satellite -use dscovr -path ./images/dscovr/1.png -dontset
$ satellite -depth 20 # This will result in a massive image!
```

**Note:** i3 requires `feh` to be installed.
