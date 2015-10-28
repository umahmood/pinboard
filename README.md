# Pinboard
Package pinboard is a GO library, which uses the Pinboard API to interact 
programatically with your bookmarks, notes and other Pinboard data.

**Note**: This project is not affiliated with [Pinboard.in](https://pinboard.in/about/) in any way.

# Installation

You will need to set up and configure the [Go](https://golang.org/doc/install) 
tool chain.

> go get github.com/umahmood/pinboard

> cd $GOPATH/src/github.com/umahmood/pinboard
    
> go test ./...

# Usage

    import (
        "fmt"
        "github.com/umahmood/pinboard"
    )
    pin := pinboard.New()
    token, err := pin.Auth("username:TOKEN")
    if err := nil {
        ...
    }
    // get all tags in users account.
    tags, err := pin.Tags()
    if err := nil {
        ...
    }
    for _, t := range tags {
        fmt.Println("Name:", t.Name, "# of tagged:", t.Count)
    }

# Documentation

> http://godoc.org/github.com/umahmood/pinboard

# License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).
