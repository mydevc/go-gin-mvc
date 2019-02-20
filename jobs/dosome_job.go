package jobs

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Subscribe struct {
	Name string
}

func (s *Subscribe) Dosome(args interface{}) error {

	var scb Subscribe
	b:=args.([]byte)
	decoder := gob.NewDecoder(bytes.NewReader(b))
	decoder.Decode(&scb)

	fmt.Println(scb.Name)

	return nil
}
