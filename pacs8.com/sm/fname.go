// fname
package sm

import (
	//	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func getFunctionName(i interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	l := strings.LastIndex(name, ".")
	name = name[l+1:]
	//	fmt.Printf("FName: %v\n", name)
	return name
}
