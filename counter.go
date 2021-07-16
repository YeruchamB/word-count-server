/*
Based mostly on the word count example in https://github.com/vladimirvivien/automi/blob/master/examples/wordcount/wordcount.go
 */
package main

import (
	"fmt"
	"github.com/vladimirvivien/automi/collectors"
	"github.com/vladimirvivien/automi/stream"
	"regexp"
	"strings"
)

// Initialize the regex used to split string into slice of words
var reg *regexp.Regexp
func init() {
	// All non alphabetical chars (allow for words with apostrophes)
	reg = regexp.MustCompile(`[^a-zA-Z']+`)
}

// Read into stream, calculate the word count and input it to memory
func InputWordCount(src interface{}) error {
	stream := stream.New(src)

	// Convert to String and make lower case
	stream.Process(func(line interface{}) string {
		return strings.ToLower(fmt.Sprintf("%s", line))
	})

	// Split into slice of words which are then streamed individually
	stream.FlatMap(func(line string) []string {
		return reg.Split(line, -1)
	})

	// Filter empty strings or word's of length 1 that aren't valid one letter words
	stream.Filter(func(word string) bool {
		return len(word) > 1 || word == "a" || word == "i"
	})

	// Map word to an array where array[0]=word and array[1]=1 to mark the occurrence of the word
	stream.Map(func(word string) [2]interface{} {
		return [2]interface{}{word, 1}
	})

	// Next:
	// 1) batch the stream of arrays from last step
	// 2) group by position 0 which has the word which returns a map[word]count
	// 3) Sum by the keys of the group
	stream.Batch().GroupByPos(0).SumAllKeys()

	// Collect the results and store them in cache
	stream.Into(collectors.Func(func(items interface{}) error {
		words := items.([]map[interface{}]float64)
		// For each word, increment the value in the cache by the count
		for _, wmap := range words {
			for word, count := range wmap {
				Increment(word.(string), int64(count))
			}
		}
		return nil
	}))

	// open the stream
	return <-stream.Open()
}
