package rgequalscenario

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var RestPrepareResult map[string]string = make(map[string]string)
var GrpcPrepareResult map[string]string = make(map[string]string)

func TestRestGoEqualFull(t *testing.T) {
	t.Run("rest go equal full test", func(t *testing.T) {

		for key, _ := range RestPrepareResult {
			t.Run(key, func(t *testing.T) {
				expect := strings.TrimSpace(RestPrepareResult[key])
				expect = strings.Replace(expect, "[]", "null", -1)
				expect = strings.Replace(expect, "\\u003c", "<", -1)
				expect = strings.Replace(expect, "\\u003e", ">", -1)
				expect = strings.Replace(expect, "\\u0026", "&", -1)

				actual := strings.TrimSpace(GrpcPrepareResult[key])
				actual = strings.Replace(actual, "[]", "null", -1)

				if !assert.True(t, expect == actual) {
					fmt.Fprintf(os.Stderr, "\n                Not Equal: \n"+
						"      REST : %s\n"+
						"      GO API  : %s\n", expect, actual)

					for idx, _ := range expect {
						if expect[idx] != actual[idx] {
							fmt.Fprintf(os.Stderr, "\n                Diff Found Index(%d): \n"+
								"      REST : %s\n"+
								"      GO API  : %s\n", idx, expect[idx:], actual[idx:])
							break
						}
					}
				}
			})
		}
	})

}
