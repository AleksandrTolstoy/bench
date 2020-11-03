package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type User struct {
	Browsers []string
	Company  string
	Country  string
	Email    string
	Job      string
	Name     string
	Phone    string
}

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
		id         int
		foundUsers string
	)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		user := User{}
		err = json.Unmarshal([]byte(line), &user)
		if err != nil {
			panic(err)
		}

		var isAndroid, isMSIE bool

		for _, browser := range user.Browsers {
			AndroidOK := strings.Contains(browser, "Android")
			if AndroidOK {
				isAndroid = true
			}

			MSIEOK := strings.Contains(browser, "MSIE")
			if MSIEOK {
				isMSIE = true
			}

			if (AndroidOK || MSIEOK) && !isSeenBefore(browser) {
				seenBrowsers = append(seenBrowsers, browser)
			}
		}

		if isAndroid && isMSIE {
			email := strings.ReplaceAll(user.Email, "@", " [at] ")
			foundUsers += fmt.Sprintf("[%d] %s <%s>\n", id, user.Name, email)
		}

		id++
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
