package util

import (
    "reflect"
    "strconv"
)

func SliceAtoi64(strArr []string) ([]int64, error) {
    // NOTE:  Read Arr as Slice as you like
    var str string // O
    var i int64    // O
    var err error  // O

    iArr := make([]int64, 0, len(strArr))
    for _, str = range strArr {
        i, err = strconv.ParseInt(str, 10, 64)
        if err != nil {
            return nil, err // O
        }
        iArr = append(iArr, i)
    }
    return iArr, nil
}

func InterfaceSlice(slice interface{}) *[]interface{} {
    s := reflect.ValueOf(slice)
    if s.Kind() != reflect.Slice {
        panic("InterfaceSlice() given a non-slice type")
    }

    ret := make([]interface{}, s.Len())

    for i:=0; i<s.Len(); i++ {
        ret[i] = s.Index(i).Interface()
    }

    return &ret
}