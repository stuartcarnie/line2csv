package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/stuartcarnie/line2csv/escape"
	"github.com/stuartcarnie/line2csv/models"
)

var (
	output = flag.String("o", "", "output filename. Default is stdout")
)

func main() {
	flag.Parse()

	var (
		in, out *os.File
		err     error
	)

	switch name := flag.Arg(0); {
	case name == "":
		in = os.Stdin
	default:
		if in, err = os.Open(name); err != nil {
			log.Fatal(err)
		}
	}

	switch name := *output; {
	case name == "":
		out = os.Stdout
	default:
		if in, err = os.Create(name); err != nil {
			log.Fatalf("unable to create output file: %s", err.Error())
		}
		defer out.Close()
	}

	wr := csv.NewWriter(out)
	wr.Write([]string{"timestamp", "key", "field_name", "field_value"})

	b := bufio.NewScanner(in)
	var row [4]string

	for b.Scan() {
		ps, err := models.ParsePoints(b.Bytes())
		if err != nil {
			// log it?
			continue
		}

		for _, p := range ps {
			row[0] = strconv.FormatInt(p.Time().UnixNano(), 10)
			row[1] = string(p.Key())

			for iter := p.FieldIterator(); iter.Next(); {
				row[2] = string(escape.Bytes(iter.FieldKey()))
				switch iter.Type() {
				case models.Float:
					if v, err := iter.FloatValue(); err != nil {
						continue
					} else {
						row[3] = strconv.FormatFloat(v, 'f', -1, 64)
					}

				case models.Integer:
					if v, err := iter.IntegerValue(); err != nil {
						continue
					} else {
						row[3] = strconv.FormatInt(v, 10)
					}

				case models.Unsigned:
					if v, err := iter.UnsignedValue(); err != nil {
						continue
					} else {
						row[3] = strconv.FormatUint(v, 10)
					}

				case models.String:
					row[3] = iter.StringValue()

				case models.Boolean:
					if v, err := iter.BooleanValue(); err != nil {
						continue
					} else if v {
						row[3] = "1"
					} else {
						row[3] = "0"
					}
				}
				if err := wr.Write(row[:]); err != nil {
					log.Fatalf("unable to write to '%s': %v", out.Name(), err)
				}
			}
		}
	}
	wr.Flush()
}
