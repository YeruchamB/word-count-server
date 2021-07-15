package main

import (
	"fmt"
	"github.com/vladimirvivien/automi/collectors"
	"github.com/vladimirvivien/automi/stream"
	"regexp"
	"strings"
)

var reg *regexp.Regexp

func init() {
	reg = regexp.MustCompile(`\W+`)
}

func ReadToStream(src interface{}) error {
	stream := stream.New(src)

	// coerce incoming data to string
	// then remove extranuous spaces.
	stream.Process(func(line interface{}) string {
		repl := reg.ReplaceAllLiteralString(fmt.Sprintf("%s", line), " ")
		//log.Println(fmt.Sprintf("%s", line))
		//log.Println(repl)
		return strings.Trim(repl, " ")
	})

	// split lines into slice of words which
	// are then streamed individually
	stream.FlatMap(func(line string) []string {
		return strings.Split(line, " ")
	})

	// map word to an array where array[0]=word
	// array[1]=1 to mark the occurence of word
	stream.Map(func(word string) [2]interface{} {
		return [2]interface{}{word, 1}
	})

	// Next:
	// 1) batch the stream of arrays from last step
	// 2) group by position 0 which has the word which returns a map[word]count
	// 3) Sum by the keys of the group
	stream.Batch().GroupByPos(0).SumAllKeys()

	// sink resunt to collector function
	stream.Into(collectors.Func(func(items interface{}) error {
		words := items.([]map[interface{}]float64)
		for _, wmap := range words {
			for word, count := range wmap {
				Increment(word.(string), int64(count))
				//fmt.Printf("%s:%0.0f\n", word, count)
			}
		}
		return nil
	}))

	// open the stream
	return <-stream.Open()
}
