package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"

	"github.com/yukithm/json2csv"
	"github.com/yukithm/json2csv/jsonpointer"

	"github.com/urfave/cli"
)

const (
	// ApplicationName is the name of this application.
	ApplicationName = "json2csv"
)

// injected by build process
var version = "unknown"

var headerStyleTable = map[string]json2csv.KeyStyle{
	"jsonpointer": json2csv.JSONPointerStyle,
	"slash":       json2csv.SlashStyle,
	"dot":         json2csv.DotNotationStyle,
	"dot-bracket": json2csv.DotBracketStyle,
}

func main() {
	// Hide timestamp because this is CLI application, so just print message for users.
	log.SetFlags(0)

	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .Flags}}[OPTIONS]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}

   If no files are specified, JSON content is read from STDIN.
   {{if .Version}}{{if not .HideVersion}}
VERSION:
   {{.Version}}
   {{end}}{{end}}{{if len .Authors}}
AUTHOR(S):
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:{{range .Categories}}{{if .Name}}
  {{.Name}}{{ ":" }}{{end}}{{range .Commands}}
    {{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}{{end}}
{{end}}{{end}}{{if .Flags}}
OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}
`

	app := cli.NewApp()
	app.Name = ApplicationName
	app.Version = version
	app.Usage = "convert JSON to CSV"
	app.ArgsUsage = "[FILE]"
	app.HideHelp = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "header-style",
			Value: "jsonpointer",
			Usage: "header style (jsonpointer, slash, dot, dot-bracket)",
		},
		cli.StringFlag{
			Name:  "path",
			Usage: "target path (JSON Pointer) of the content",
		},
		cli.IntFlag{
			Name:  "slice-len",
			Value: math.MaxInt,
			Usage: "Specify the length of the slice to be processed.",
		},
		cli.BoolFlag{
			Name:  "transpose",
			Usage: "transpose rows and columns",
		},
		cli.BoolFlag{
			Name:  "stream",
			Usage: "convert data stream",
		},
		cli.HelpFlag,
	}

	app.Before = func(c *cli.Context) error {
		if _, ok := headerStyleTable[c.String("header-style")]; !ok {
			return fmt.Errorf("Invalid --header-style value %q", c.String("header-style"))
		}
		return nil
	}

	app.Action = func(c *cli.Context) {
		if c.Bool("help") {
			cli.ShowAppHelp(c)
			return
		}
		mainAction(c)
	}

	app.RunAndExitOnError()
}

func streamReaderFromFile(filename string) json2csv.JSONStreamReader {
	var reader json2csv.JSONStreamReader
	if strings.HasSuffix(filename, ".zip") {
		zipReader, err := zip.OpenReader(filename)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		// defer zipReader.Close()
		reader = json2csv.NewJSONStreamZipReader(zipReader)
	} else {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		// defer file.Close()
		reader = json2csv.NewJSONStreamLineReader(file)

	}
	return reader
}

func mainAction(c *cli.Context) {
	var data interface{}
	var err error
	headerStyle := headerStyleTable[c.String("header-style")]
	sliceLen := c.Int("slice-len")
	if c.NArg() > 0 && c.Args()[0] != "-" {
		filename := c.Args()[0]
		if c.Bool("stream") {
			reader := streamReaderFromFile(filename)
			csvHeader, err := json2csv.JSON2CSVHeader(reader, c.String("path"), sliceLen)
			if err != nil {
				log.Fatal(err)
			}
			reader.Close()
			reader = streamReaderFromFile(filename)
			err = json2csv.JSON2CSVOnline(reader, csvHeader, os.Stdout, headerStyle, false, c.String("path"), sliceLen)
			reader.Close()
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		data, err = readJSONFile(filename)
	} else {
		data, err = readJSON(os.Stdin)
	}
	if err != nil {
		log.Fatal(err)
	}

	if c.String("path") != "" {
		data, err = jsonpointer.Get(data, c.String("path"))
		if err != nil {
			log.Fatal(err)
		}
	}

	results, err := json2csv.JSON2CSV(data, nil, sliceLen)
	if err != nil {
		log.Fatal(err)
	}
	if len(results) == 0 {
		return
	}

	err = printCSV(os.Stdout, results, headerStyle, c.Bool("transpose"))
	if err != nil {
		log.Fatal(err)
	}
}

func readJSONFile(filename string) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return readJSON(f)
}

func readJSON(r io.Reader) (interface{}, error) {
	decoder := json.NewDecoder(r)
	decoder.UseNumber()

	var data interface{}
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

func printCSV(w io.Writer, results []json2csv.KeyValue, headerStyle json2csv.KeyStyle, transpose bool) error {
	csv := json2csv.NewCSVWriter(w, headerStyle, transpose)
	csv.HeaderStyle = headerStyle
	csv.Transpose = transpose
	if err := csv.WriteCSV(results); err != nil {
		return err
	}
	return nil
}
