package dksync

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Args struct {
	Field string
	Value string

	// If it is to remove desired key[ Value ]
	// instead of set it
	Del   bool
}

var dKronHost string

// Set custom args
func (args *Args) Set(field, value string, toDelete bool) {
	args.Field = field
	args.Value = value
	args.Del = toDelete
}

// setFieldValue gets each job instance from struct
// then sets new desired value for the field, using Types
// in case it is of type Map the desired key can be removed
func (args *Args) setFieldValue(job interface{}) {
	log.Debugf("Assigning attrs: %v to job %v", args, job)
	s := reflect.ValueOf(job)
	if res := s.Elem().FieldByName(args.Field); res.IsValid() {
		if res.CanSet() {
			switch res.Kind() {
			case reflect.String:
				res.SetString(args.Value)
			case reflect.Map:
				if args.Del {
					delete(res.Interface().(map[string]map[string]string), args.Value)
					break
				}
				if res.Type().String() == "map[string]map[string]string" {
					var vMap map[string]map[string]string
					if err := json.Unmarshal([]byte(args.Value), &vMap); err != nil {
						log.Error(err)
						break
					}
					v := reflect.ValueOf(vMap)
					k := v.MapKeys()
					for i := range k {
						res.SetMapIndex(k[i], v.MapIndex(k[i]))
					}
				} else {
					res.SetMapIndex(reflect.ValueOf(args.Field), reflect.ValueOf(args.Value))
				}
			case reflect.Bool:
				val, err := strconv.ParseBool(args.Value)
				if err != nil {
					log.Error(val)
					break
				}
				res.SetBool(val)
			case reflect.Int:
				val, err := strconv.ParseInt(args.Value, 0, 0)
				if err != nil {
					log.Error(val)
					break
				}
				res.SetInt(val)
			default:
				log.Fatal("Unknown value type! Needs to be a String, Int, Bool or Map")
				break
			}
		}
	}
}

// parser parses file content with defined regex
func parser(content []byte) string {
	res := regexp.MustCompile("(#[a-zA-Z0-9- ].*)").ReplaceAll(content, []byte(""))
	return strings.TrimSpace(string(res))
}