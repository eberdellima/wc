package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"unicode/utf8"
)

func main() {
	input := os.Stdin
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[len(os.Args)-1], "-") {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal("could not get working directory", err)
		}

		file := path.Join(dir, os.Args[len(os.Args)-1])
		f, err := os.Open(file)
		if err != nil {
			log.Fatal("could not open file", err)
		}
		defer f.Close()

		input = f
	}

	opts := parseOptions()

	c := &counter{}
	err := c.count(bufio.NewReader(input), opts)
	if err != nil {
		log.Fatal(err)
	}

	if opts.countLines {
		fmt.Printf("\t%d", c.totalLines)
	}
	if opts.countWords {
		fmt.Printf("\t%d", c.totalWords)
	}
	if opts.countBytes {
		fmt.Printf("\t%d", c.totalBytes)
	}
	if opts.countCharacters {
		fmt.Printf("\t%d", c.totalCharacters)
	}
	if input != os.Stdin {
		fmt.Printf("\t%s", os.Args[len(os.Args)-1])
	}
	fmt.Println()
}

type counter struct {
	totalBytes      int
	totalLines      int
	totalWords      int
	totalCharacters int
}

func (c *counter) count(r *bufio.Reader, opts options) error {
	withMultiByteChars := supportsMultiByteChars()
	for {
		data, err := r.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		if opts.countBytes {
			c.totalBytes += len(data)
		}
		if opts.countLines {
			c.totalLines++
		}
		if opts.countWords {
			c.totalWords += len(strings.Fields(string(data)))
		}
		if opts.countCharacters {
			if withMultiByteChars {
				c.totalCharacters += utf8.RuneCount(data)
			} else {
				c.totalCharacters += len(data)
			}
		}
	}
	return nil
}

func supportsMultiByteChars() bool {
	vars := []string{"LC_ALL", "LC_CTYPE", "LANG"}
	for _, v := range vars {
		val := os.Getenv(v)
		if val != "" {
			if strings.Contains(strings.ToUpper(val), "UTF-8") ||
				strings.Contains(strings.ToUpper(val), "EUC") ||
				strings.Contains(strings.ToUpper(val), "BIG5") {
				return true
			}
		}
	}
	return false
}

type options struct {
	countBytes      bool
	countLines      bool
	countWords      bool
	countCharacters bool
}

func parseOptions() options {
	opts := options{}

	if len(os.Args) == 1 {
		return opts
	}

	switch {
	case !strings.HasPrefix(os.Args[1], "-"):
		opts.countBytes = true
		opts.countLines = true
		opts.countWords = true
	default:
		if strings.Contains(os.Args[1], "c") {
			opts.countBytes = true
		}
		if strings.Contains(os.Args[1], "l") {
			opts.countLines = true
		}
		if strings.Contains(os.Args[1], "w") {
			opts.countWords = true
		}
		if strings.Contains(os.Args[1], "m") {
			opts.countCharacters = true
		}
	}

	return opts
}
