package main

import (
	"bufio"
	"os"
	"path"
	"strings"
	"time"
	"encoding/json"
	"github.com/konveyor/tackle2-addon/command"
	"github.com/konveyor/tackle2-hub/api"
)

// tca application analyzer.
type Tca struct {
	appName     string
	input       string
	application *api.Application
}

// Run tca addon
func (r *Tca) Run() (err error) {

	// Run tca_cli
	addon.Activity("[TCA] Starting..")
	tcaCommand := command.Command{Path: "python3"}
	tcaCommand.Dir = "/app"
	tcaCommand.Options = command.Options{"tca_cli.py"}
	tcaCommand.Options.Add("-input_string", r.input)
	tcaCommand.Options.Add("-operation", "addon")
	start := time.Now()
	err = tcaCommand.Run()
	if err != nil {
		r.reportLog()
		return err
	}
	elapsed := time.Since(start)
	addon.Activity("[TCA] DONE - execution time %s", elapsed)

	// Mapping standardize output to the tag categories in tackle inventory
	tagKeyMap := make(map[string]string)
	tagKeyMap["OS"] = "Operating System"
	tagKeyMap["Lang"] = "Language"
	tagKeyMap["Runtime"] = "Runtime"
	tagKeyMap["App Server"] = "Application Type"
	tagKeyMap["Libs"] = "Application Type"
	tagKeyMap["Dependent Apps"] = "Application Type"

	// Creating keys for setting facts
	factKeys := []string{"Ref Dockers", "Reason", "Recommend"}

	tags := make(map[string][]string)
	myFacts := make(map[string]interface{})

	result_str, _ := GetStringInBetweenTwoString(string(tcaCommand.Output), "^^^", "$$$")

	var res map[string]interface{}
	json.Unmarshal([]byte(result_str), &res)

	for appk := range res {
		app_map := res[appk].([]interface{})
		if len(app_map) == 0 {
			continue
		}
		app_map_0 := app_map[0].(map[string]interface{})
		if appk == "standardize" {
			for k := range tagKeyMap {
				k_map := app_map_0[k].(map[string]interface{})
				for kname := range k_map {
					kname_map := k_map[kname].(map[string]interface{})
					var val_str string
					val_str = kname_map["standard_name"].(string)
					val_str_parts := strings.Split(val_str, "|")
					for i := range val_str_parts {
						if val_str_parts[i] == "*" {
							break
						}
						tags[tagKeyMap[k]] = append(tags[tagKeyMap[k]], val_str_parts[i])
					}
				}
			}
		} else if appk == "containerize" {
			for k := range factKeys {
				factKey := factKeys[k]
				myFacts[factKey] = app_map_0[factKey]
			}
		}
	}
	// tagging new entities from standardize output
	for k := range tags {
		addon.Activity("TagType:%s", k, "Numtags:%s", len(tags[k]), "TagName:%s", tags[k])
		addTags(r.application, k, tags[k])
	}
	// container recommendation is set as facts
	// setting facts using sub-resource
	facts := addon.Application.Facts(r.application.ID)
	err = facts.Set("Container_Advisory", myFacts)
	if err != nil {
		return
	}

	addon.Increment()
	return nil
}

// reportLog reports the log content.
func (r *Tca) reportLog() {
	logPath := path.Join(
		HomeDir,
		".mta",
		"log",
		"mta.log")
	f, err := os.Open(logPath)
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		addon.Activity(">> %s", scanner.Text())
	}
	_ = f.Close()
}

// addTags ensure tags created and associated with application.
// Ensure tag exists and associated with the application.
func addTags(application *api.Application, tagTypeName string, names []string) (err error) {
	addon.Activity("Adding tags: %v", names)
	var wanted []uint
	// Ensure type exists.
	tp := &api.TagType{
		Name:  tagTypeName,
		Color: "#2b9af3",
		Rank:  3,
	}
	err = addon.TagType.Ensure(tp)
	if err != nil {
		return
	}
	// Ensure tags exist.
	for _, name := range names {
		tag := &api.Tag{
			Name: name,
			TagType: api.Ref{
				ID: tp.ID,
			}}
		err = addon.Tag.Ensure(tag)
		if err == nil {
			wanted = append(wanted, tag.ID)
		} else {
			return
		}
	}
	// Associate tags.
	tags := addon.Application.Tags(application.ID)
	// tags.Replace(wanted)
	for _, id := range wanted {
		err = tags.Add(id)
		if err != nil {
			return
		}
	}
	return
}

func GetStringInBetweenTwoString(str string, startS string, endS string) (result string, found bool) {
	s := strings.Index(str, startS)
	if s == -1 {
		return result, false
	}
	newS := str[s+len(startS):]
	e := strings.Index(newS, endS)
	if e == -1 {
		return result, false
	}
	result = newS[:e]
	return result, true
}
