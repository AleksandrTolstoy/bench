package main

import (
	"bench/data"
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// 1) easyjson -all data/user.go
// 2) go test -bench . -benchmem
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	seenBrowsers := make([]string, 0)
	isSeenBefore := func(browser string) bool {
		for _, item := range seenBrowsers {
			if item == browser {
				return true
			}
		}

		return false
	}

	var (
		id   int
		user = new(data.User)
	)

	fmt.Fprintln(out, "found users:")
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		err = user.UnmarshalJSON([]byte(line))
		if err != nil {
			panic(err)
		}

		var isAndroid, isMSIE bool
		for _, browser := range user.Browsers {
			okAndroid := strings.Contains(browser, "Android")
			if okAndroid {
				isAndroid = true
			}

			okMSIE := strings.Contains(browser, "MSIE")
			if okMSIE {
				isMSIE = true
			}

			if (okAndroid || okMSIE) && !isSeenBefore(browser) {
				seenBrowsers = append(seenBrowsers, browser)
			}
		}

		if isAndroid && isMSIE {
			email := strings.ReplaceAll(user.Email, "@", " [at] ")
			foundUser := fmt.Sprintf("[%d] %s <%s>\n", id, user.Name, email)
			fmt.Fprint(out, foundUser)
		}

		id++
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}
