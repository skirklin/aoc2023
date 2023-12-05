package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// User session id is a github session id (plucked from cookies in browser)
type User struct {
	SessionID string
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
func (user User) FetchInputs(year, day int) string {
	fmt.Println("reading from web")
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day), nil)
	panicIf(err)
	// ...
	req.AddCookie(
		&http.Cookie{
			Name:  "session",
			Value: user.SessionID,
		},
	)
	resp, err := client.Do(req)
	panicIf(err)

	if resp.StatusCode != 200 {
		panic(fmt.Errorf("Unhandled response %s", fmt.Sprint(resp)))
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	panicIf(err)
	return string(data)
}

func Subslice(arr []interface{}, i int) []interface{} {
	ret := make([]interface{}, 0)
	ret = append(ret, arr[:i]...)
	ret = append(ret, arr[i+1:]...)
	return ret
}

func permute(base []interface{}, values []interface{}, c chan<- []interface{}) {
	if len(values) == 0 {
		c <- append(make([]interface{}, 0), base...)
	} else {
		for i := range values {
			base = append(base, values[i])
			permute(base, Subslice(values, i), c)
			base = base[:len(base)-1]
		}

	}
	if len(base) == 0 {
		close(c)
	}
}

func Permutations(arr []interface{}) <-chan []interface{} {
	c := make(chan []interface{})

	go permute(make([]interface{}, 0), arr, c)

	return c
}

func combine(base []int, values []int, c chan<- []int, first bool) {
	if len(values) == 0 {
		c <- append(make([]int, 0), base...)
	} else {
		combine(base, values[1:], c, false)
		base = append(base, values[0])
		combine(base, values[1:], c, false)
	}
	if first {
		close(c)
	}
}

func Combinations(arr []int) <-chan []int {
	c := make(chan []int)

	go combine(make([]int, 0), arr, c, true)

	return c
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
