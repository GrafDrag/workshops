package decode

import (
	"cli-decoder/internal/store"
	"crypto/md5"
	"encoding/xml"
	"fmt"
)

type XMLDecoder struct {
	Store *store.DataStore
	str   []byte
}

func (d *XMLDecoder) SetHash(str string) error {
	d.str = []byte(str)

	if !d.validator() {
		return fmt.Errorf("problem xml is not valid\n")
	}

	if err := d.writeToFile(); err != nil {
		return fmt.Errorf("problem wreate json to file, %v\n", err)
	}
	return nil
}

func (d *XMLDecoder) validator() bool {
	return xml.Unmarshal(d.str, new(interface{})) == nil
}

func (d *XMLDecoder) getMD5() string {
	return fmt.Sprintf("%x", md5.Sum(d.str))
}

func (d *XMLDecoder) writeToFile() error {
	return d.Store.SavaFile(d.getMD5(), ".xml", d.str)
}
