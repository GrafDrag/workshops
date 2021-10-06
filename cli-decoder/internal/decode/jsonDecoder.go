package decode

import (
	"cli-decoder/internal/store"
	"crypto/md5"
	"encoding/json"
	"fmt"
)

type JSONDecoder struct {
	Store *store.DataStore
	str   []byte
}

func (d *JSONDecoder) SetHash(str string) error {
	d.str = []byte(str)

	if !d.validator() {
		return fmt.Errorf("problem json is not valid\n")
	}

	if err := d.writeToFile(); err != nil {
		return fmt.Errorf("problem wreate json to file, %v\n", err)
	}
	return nil
}

func (d *JSONDecoder) validator() bool {
	return json.Valid(d.str)
}

func (d *JSONDecoder) getMD5() string {
	return fmt.Sprintf("%x", md5.Sum(d.str))
}

func (d *JSONDecoder) writeToFile() error {
	return d.Store.SavaFile(d.getMD5(), ".json", d.str)
}
