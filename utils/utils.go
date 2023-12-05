package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

// session id is a github session id (plucked from cookies in browser)
type Session struct {
	SessionID string
	GID       string
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func AsInt(m string) int {
	i, err := strconv.Atoi(m)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

// FetchInputs gets the input data for a given year and day
func (session Session) FetchInputs(year, day int) []byte {
	fmt.Println("reading from web")
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day), nil)
	fmt.Println("session =", session.SessionID)
	fmt.Println("gid =", session.GID)
	panicIf(err)
	// ...
	req.AddCookie(
		&http.Cookie{
			Name:  "session",
			Value: session.SessionID,
		},
	)
	req.AddCookie(
		&http.Cookie{
			Name:  "_gid",
			Value: session.GID,
		},
	)
	resp, err := client.Do(req)
	panicIf(err)

	if resp.StatusCode != 200 {
		panic(fmt.Errorf("Unhandled response %s", fmt.Sprint(resp)))
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	panicIf(err)
	return data
}

func GetInputs(year int, day int) string {
	localRoot := "/tmp/aoc/"
	localCopy := fmt.Sprintf("%s/%d/%d", localRoot, year, day)

	fh, err := os.Open(localCopy)
	if err != nil {
		session := Session{SessionID: os.Getenv("AOC_SESSION_ID"), GID: os.Getenv("AOC_GID")}
		response := session.FetchInputs(year, day)
		err := os.MkdirAll(fmt.Sprintf("%s/%d", localRoot, year), 0750)
		panicIf(err)
		os.WriteFile(localCopy, response, 0750)
		return string(response)
	} else {
		bytes, err := io.ReadAll(fh)
		panicIf(err)
		return string(bytes)
	}
}

func Subslice(arr []interface{}, i int) []interface{} {
	ret := make([]interface{}, 0)
	ret = append(ret, arr[:i]...)
	ret = append(ret, arr[i+1:]...)
	return ret
}

func AsInterfaceSlice(s []string) []interface{} {
	o := make([]interface{}, len(s), len(s))
	for i := range s {
		o[i] = s[i]
	}
	return o
}
func AsStringSlice(s []interface{}) []string {
	o := make([]string, len(s), len(s))
	for i := range s {
		o[i] = s[i].(string)
	}
	return o
}
