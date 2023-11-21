// SPDX-FileCopyrightText: 2023 VinyMeuh
// SPDX-License-Identifier: MIT
package engine

import (
	"fmt"
)

type usiOption interface {
	fmt.Stringer
	set(value string) error
}

// type checkOption struct {
// 	callback func(value bool)
// 	value    bool
// }

// func (co checkOption) String() string {
// 	return fmt.Sprintf("type check default %s", strconv.FormatBool(co.value))
// }

// func (co checkOption) set(value string) error {
// 	switch value {
// 	case "true":
// 		co.callback(true)
// 	case "false":
// 		co.callback(false)
// 	default:
// 		return fmt.Errorf("valid values are [true, false]")
// 	}
// 	return nil
// }

type comboOption struct {
	callback func(value string)
	value    string
	values   []string
}

func (co comboOption) set(value string) error {
	for _, v := range co.values {
		if v == value {
			co.callback(value)
			return nil
		}
	}
	return fmt.Errorf("valid values are %v", co.values)
}

func (co comboOption) String() string {
	s := fmt.Sprintf("type combo default %s", co.value)
	for _, v := range co.values {
		s += fmt.Sprintf(" var %s", v)
	}
	return s
}

// type spinOption struct {
// 	callback func(value int)
// 	value    int
// 	min      int
// 	max      int
// }

// func (so spinOption) String() string {
// 	return fmt.Sprintf("type spin default %d min %d max %d", so.value, so.min, so.max)
// }

// func (so spinOption) set(value string) error {
// 	ivalue, err := strconv.Atoi(value)

// 	if err != nil {
// 		return fmt.Errorf("not a number")
// 	}

// 	if ivalue < so.min || ivalue > so.max {
// 		return fmt.Errorf("out of range [%d, %d]", so.min, so.max)
// 	}

// 	so.callback(ivalue)
// 	return nil
// }

// Noop Callbacks
// func noopBoolCallback(_ bool) {}

// func noopIntCallback(_ int) {}

func noopStringCallback(_ string) {}
