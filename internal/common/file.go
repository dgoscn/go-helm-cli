package common

import (
    "io/ioutil"
)

func ReadFile(filename string) ([]byte, error) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return data, nil
}
