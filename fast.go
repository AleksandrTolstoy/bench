package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

func main() {
	fastOut := new(bytes.Buffer)
	FastSearch(fastOut)
	fastResult := fastOut.String()
	fmt.Println(fastResult)
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath) // change it to filePath
	if err != nil {
		panic(err)
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	seenBrowsers := make([]string, 0)
	foundUsers := ""

	lines := strings.Split(string(fileContents), "\n")

	users := make([]*User, 0, len(lines))
	for _, line := range lines {
		user := new(User)
		err := json.Unmarshal([]byte(line), &user)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}

	isSeenBefore := func(browser string) bool {
		for _, item := range seenBrowsers {
			if item == browser {
				return true
			}
		}

		return false
	}

	for i, user := range users {
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

		if !(isAndroid && isMSIE) {
			continue
		}

		email := strings.ReplaceAll(user.Email, "@", " [at] ")
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email)
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
