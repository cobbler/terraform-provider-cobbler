package cobbler

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/go-homedir"
	"hash/crc32"
	"os"
)

// If the argument is a path, Read loads it and returns the contents,
// otherwise the argument is assumed to be the desired contents and is simply
// returned.
//
// The boolean second return value can be called `wasPath` - it indicates if a
// path was detected and a file loaded.
func Read(poc string) (string, bool, error) {
	if len(poc) == 0 {
		return poc, false, nil
	}

	path := poc
	if path[0] == '~' {
		var err error
		path, err = homedir.Expand(path)
		if err != nil {
			return path, true, err
		}
	}

	if _, err := os.Stat(path); err == nil {
		contents, err := os.ReadFile(path)
		if err != nil {
			return string(contents), true, err
		}
		return string(contents), true, nil
	}

	return poc, false, nil
}

// String hashes a string to a unique hashcode.
//
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
func String(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

// Strings hashes a list of strings to a unique hashcode.
func Strings(strings []string) string {
	var buf bytes.Buffer

	for _, s := range strings {
		buf.WriteString(fmt.Sprintf("%s-", s))
	}

	return fmt.Sprintf("%d", String(buf.String()))
}

// GetStringSlice is a helper which safely retrieves the data of a given key and casts it to a string slice.
func GetStringSlice(d *schema.ResourceData, key string) ([]string, error) {
	result := make([]string, 0)
	keyData, ok := d.Get(key).([]interface{})
	if !ok {
		return nil, fmt.Errorf("key `%s` is not an array", key)
	}
	for _, element := range keyData {
		var castedElement string
		castedElement, ok = element.(string)
		if !ok {
			return nil, fmt.Errorf("key `%s` is not a string", key)
		}
		result = append(result, castedElement)
	}
	return result, nil
}

func GetInterfaceMap(d *schema.ResourceData, key string) (map[string]interface{}, error) {
	interfaceData, ok := d.Get(key).(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("key `%s` is not a map", key)
	}
	return interfaceData, nil
}
