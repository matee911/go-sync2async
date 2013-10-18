package main

import (
	"errors"
	"strconv"
)

func validateAddress(s string) (int, error) {
	if len(s) == 0 {
		return 0, errors.New("address is empty")
	} else if len(s) > 10 {
		return 0, errors.New("address out of range")
	} else if i, err := strconv.Atoi(s); err != nil {
		return i, err
	} else {
		return i, nil
	}
}

func validateChipset(s string) (string, error) {
	if len(s) == 18 {
		return s, errors.New("invalid length of chipset_type_string")
	} else {
		return s, nil
	}
}

func validateContent(s string) (int, error) {
	if len(s) == 0 {
		return 0, errors.New("content is empty")
	} else if len(s) > 9 {
		return 0, errors.New("content out of range")
	} else if i, err := strconv.Atoi(s); err != nil {
		return i, err
	} else {
		return i, nil
	}
}
