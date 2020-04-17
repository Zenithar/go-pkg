package types

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/structs"
)

// AsEnvVariables sets struct values from environment variables
func AsEnvVariables(o interface{}, prefix string, skipCommented bool) (map[string]string, error) {
	// Check input
	if IsNil(o) {
		return nil, fmt.Errorf("nil input can't be exported")
	}

	// Prepare result
	r := map[string]string{}
	prefix = strings.ToUpper(prefix)
	delim := "_"
	if prefix == "" {
		delim = ""
	}
	fields := structs.Fields(o)
	for _, f := range fields {
		// If attribute is commented, ignore it
		if skipCommented {
			tag := f.Tag("commented")
			if tag != "" {
				commented, err := strconv.ParseBool(tag)
				if err != nil {
					return nil, fmt.Errorf("unable to parse tag value '%s' for field '%s': %w", tag, f.Name(), err)
				}
				if commented {
					continue
				}
			}
		}

		// If value is a struct
		if structs.IsStruct(f.Value()) {
			rf, err := AsEnvVariables(f.Value(), prefix+delim+f.Name(), skipCommented)
			if err != nil {
				return nil, err
			}
			for k, v := range rf {
				r[k] = v
			}
		} else {
			r[prefix+"_"+strings.ToUpper(f.Name())] = fmt.Sprintf("%v", f.Value())
		}
	}

	// No error
	return r, nil
}
