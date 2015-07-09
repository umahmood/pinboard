/*
Package pinboard uses the Pinboard API to interact programatically with your
bookmarks, notes and other Pinboard data.

To start, you must authenticate using the users API token:

    import (
       "github.com/umahmood/pinboard"
    )

    pin := pinboard.New()

    token, err := pin.Auth("username:TOKEN")
    ...

Last time users account had activity:

    t, err := pin.LastUpdate()

    if err != nil {
        ...
    }

    fmt.Println("Last Update:", t)

Add a bookmark:

    b := pinboard.Bookmark{URL: "https://www.eff.org/",
        Title:   "Electronic Frontier Foundation",
        Desc:    "Defending your rights in the digital world",
        Tags:    []string{"privacy", "rights", "encryption"},
        Created: time.Now(),
        Replace: true,
        Shared:  true,
        ToRead:  false,
    }

    ok, err := pin.Add(b)

    if err != nil {
        ...
    }

    if ok {
        fmt.Println("bookmark added!")
    }

Deleting a bookmark:

    ok, err := pin.Del("https://www.eff.org/")

    if err != nil {
        ...
    }

    if ok {
        fmt.Println("bookmark deleted!")
    }

Get the bookmarks added today:

    bmarks, err := pin.Get(time.Time{}, "", nil, false)

    if err != nil {
        ...
    }

    for _, b := range bmarks {
        fmt.Println(b.Title)
    }

Get bookmarks added on April 1st 2015:

    date := time.Date(2015, 4, 1, 0, 0, 0, 0, time.UTC)
    bmarks, err := pin.Get(date, "", nil, false)
    ...

Get bookmarks added on April 1st 2015 which have the tags:

    date := time.Date(2015, 4, 1, 0, 0, 0, 0, time.UTC)
    tags := []string{"python", "maths"}
    bmarks, err := pin.Get(date, "", tags, false)
    ...

Get bookmark for this url:

    bmarks, err := pin.Get(time.Time{}, "https://www.eff.org/", nil, false)
    ...

Number of posts on each date:

    posts, err := pin.Dates(nil)

    if err != nil {
        ...
    }

    for _, p := range posts {
        fmt.Println("Date:", p.Date, "# of posts:", p.Count)
    }

Number of posts on each date for given tags:

    tags := []string{"recipe", "cooking"}

    posts, err := pin.Dates(tags)
    ...

Users 5 most recent posts:

    bmarks, err := pin.Recent(nil, 5)

    if err != nil {
        ...
    }

    for _, b := range bmarks {
        fmt.Println(b.Title)
    }

Users most recent posts filtered by tag:

    tags := []string{"economics"}
    bmarks, err := pin.Recent(tags, 0)
    ...

Get all bookmarks in the users account:

    bmarks, err := pin.Bookmarks(nil, 0, 0, time.Time{}, time.Time{}, false)

    if err != nil {
        ...
    }

    for _, b := range bmarks {
        fmt.Println(b.Title)
    }

Get 10 days worth of bookmarks:

    start := time.Now().AddDate(0, 0, -10)
    end := time.Now()
    bmarks, err := pin.Bookmarks(nil, 0, 0, start, end, false)
    ...

Get all bookmarks which have the tags:

    tags := []string{"ruby", "tutorials"}
    bmarks, err := pin.Bookmarks(tags, 0, 0, time.Time{}, time.Time{}, false)
    ...

Get suggestions for tags based on a URL:

    pop, rec, err := pin.Suggest("https://www.eff.org/")

    if err != nil {
        ...
    }

    fmt.Println("Popular tags:")
    for _, i := range pop {
        fmt.Println(i)
    }

    fmt.Println("Recommended tags:")
    for _, i := range rec {
        fmt.Println(i)
    }

Get a list of all tags in the users account:

    tags, err := pin.Tags()

    if err != nil {
        ...
    }

    for _, t := range tags {
        fmt.Println("Name:", t.Name, "# of tagged:", t.Count)
    }

Delete a tag:

    ok, err := pin.DelTag("fonts")

    if err != nil {
        ...
    }

    if ok {
        fmt.Println("tag deleted.")
    }

Rename a tag:

    ok, err := pin.RenTag("fonts", "typography")

    if err != nil {
        ...
    }

    if ok {
        fmt.Println("tag renamed.")
    }

Get a list of all notes in the users account:

    notes, err := pin.Notes()

    if err != nil {
        ...
    }

    for _, n := range notes {
        fmt.Println("ID", n.ID, "Title", n.Title)
    }

Get an individual note:

    notes, err := pin.Notes()

    if err != nil {
        ...
    }

    for _, m := range notes {
        n, err := pin.NoteID(m.ID)
        if err != nil {
            ...
        }
        fmt.Println(n.Title)
        fmt.Println(n.Text)
    }
*/
package pinboard
