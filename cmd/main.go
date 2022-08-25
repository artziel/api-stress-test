package main

import (
	"fmt"
	"strings"

	Stress "github.com/artziel/api-stress-test/pkg/api"
	Console "github.com/artziel/go-console"
)

type CommandFlags struct {
	EndpointsJson string `GoConsole:"name:file,usage:List of endpoints in Json format"`
}

func main() {

	root := Console.Root{
		Commands: map[string]Console.Command{
			"": {
				Help:    "Main Command",
				Example: "$ sample",
				Flags:   &CommandFlags{},
				Run: func(args interface{}) error {
					flags := args.(*CommandFlags)

					if flags.EndpointsJson == "" {
						return fmt.Errorf("no Endpoints file specified, please use the flag -file [PATH]/[TO]/[FILE]/[FILE-NAME].json")
					}

					data, err := Stress.ReadJSON(flags.EndpointsJson)
					if err != nil {
						return fmt.Errorf("error loading Endpoints JSON file: %s", err)
					}

					for _, e := range data.Endpoints {
						spn := Console.Spinner09()
						spn.SetColor("blue").SetPrefix("Request " + e.URL + " ... ")

						spn.Start()
						r, err := e.Exec()
						spn.Stop()
						if err != nil {
							fmt.Printf("Error: %s\n", err)
						}

						fmt.Printf("[%s] %s\n%s\n%s\n\n", e.Method, e.URL, strings.Repeat("-", 60), r.String())
					}

					return nil
				},
			},
		},
	}

	if err := root.Run(); err != nil {
		fmt.Println(err.Error())
	}
}
