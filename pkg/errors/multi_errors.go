package errors

import "strings"

type MultiError []error

func (m MultiError) Error() string {
	output := make([]string, len(m))
	for _, elem := range m {
		if elem == nil {
			continue
		}
		output = append(output, elem.Error())
	}
	return strings.Join(output, "\n")
}
