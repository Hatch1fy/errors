package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestErrorListJSON(t *testing.T) {
	var (
		ms  testMarshalStruct
		err error
	)

	ms.Errs.Push(Error("test error"))
	buf := bytes.NewBuffer(nil)
	if err = json.NewEncoder(buf).Encode(ms); err != nil {
		t.Fatal(err)
	}

	var nms testMarshalStruct
	if err = json.NewDecoder(buf).Decode(&nms); err != nil {
		t.Fatal(err)
	}

	if mse, nmse := ms.Errs.Err().Error(), nms.Errs.Err().Error(); mse != nmse {
		t.Fatalf("invalid error, expected \"%s\" and recieved \"%s\"", mse, nmse)
	}
	fmt.Println(buf.String())
}

type testMarshalStruct struct {
	Errs ErrorList `json:"errs"`
}
