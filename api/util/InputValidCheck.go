package util

import (
    "errors"
)

func CheckValidPageObj(perPage int, pageNum int) error {
    if perPage == 0 || perPage < -1 {
        return errors.New("perPage is more than 1")
    }
    if pageNum < 1 {
        return errors.New("pageNum is more than 0")
    }
    return nil
}

