package cfgargs

import (
	"framework/file"
	"testing"
)

func TestGetSrvConfig(t *testing.T) {
	var in = file.GetAbsPath("../../etc/config-example.yaml")
	_, err := GetSrvConfig(in)
	if err != nil {
		t.Error("Test Get Config File Error. Err: ", err)
	}
}
