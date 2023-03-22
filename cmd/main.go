package main

import (
	hub "github.com/konveyor/tackle2-hub/addon"
	"os"
)

var (
	// hub integration.
	addon = hub.Addon
	// HomeDir directory.
	HomeDir = ""
)

type SoftError = hub.SoftError

// addon data passed in secret
type Data struct {
	Name string `json:"application_name"`
	Desc string `json:"application_description"`
	Tech string `json:"technology_summary"`
}

// main
func main() {

	addon.Run(func() error {

		HomeDir, _ = os.UserHomeDir()

		// Get the addon data associated with the task.
		d := &Data{}
		if err := addon.DataWith(d); err != nil {
			return &SoftError{Reason: err.Error()}
		}
		addon.Activity("Fetching data.")
		addon.Activity("%s", d.Tech)
		input_string := "\"application_name\":\"" + d.Name + "\",\"application_description\":\"" + d.Desc + "\",\"technology_summary\":\"" + d.Tech + "\""
		addon.Activity("%s", input_string)
		// Setup tca
		tca := Tca{}

		// Fetch application.
		addon.Activity("Fetching application.")
		application, err := addon.Task.Application()
		if err != nil {
			return err
		}

		tca.appName = "TCA"
		tca.input = input_string
		tca.application = application

		addon.Total(1)
		// Run tca.
		if err = tca.Run(); err != nil {
			return &SoftError{Reason: err.Error()}
		}
		addon.Increment()

		return nil
	})
}
