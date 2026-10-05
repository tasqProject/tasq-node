package offer

import "fmt"

func errField(name string) error { return fmt.Errorf("offer: missing field %q", name) }
func errMsg(msg string) error    { return fmt.Errorf("offer: %s", msg) }
