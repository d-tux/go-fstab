package fstab

import (
	"strings"
	"testing"
)

var successfulParseLineExpectations map[string]Mount = map[string]Mount{
	"/dev/sda / ext4 defaults 1 2": Mount{
		"/dev/sda",
		"/",
		"ext4",
		map[string]string{
			"defaults": "",
		},
		1,
		2,
	},

	"UUID=homer / ext4 rw,uid=0": Mount{
		"UUID=homer",
		"/",
		"ext4",
		map[string]string{
			"uid": "0",
			"rw":  "",
		},
		0,
		0,
	},
}

var successfulMountStringExpectations map[string]Mount = map[string]Mount{
	"/dev/sda / ext4 defaults 1 2": Mount{
		"/dev/sda",
		"/",
		"ext4",
		map[string]string{
			"defaults": "",
		},
		1,
		2,
	},

	"UUID=homer / ext4 uid=0 0 0": Mount{
		"UUID=homer",
		"/",
		"ext4",
		map[string]string{
			"uid": "0",
		},
		0,
		0,
	},
}

func TestParseLine(t *testing.T) {
	for line, expectation := range successfulParseLineExpectations {
		mount, err := ParseLine(line)
		if nil != err {
			t.Errorf("Unexpected parse error while parsing '%s': %s", line, err)
			continue
		}

		if !mount.Equals(&expectation) {
			t.Errorf("Expected %+v, got %+v", expectation, mount)
		}

		if 0 == strings.Index(mount.Spec, "UUID") && mount.SpecType() != UUID {
			t.Errorf("Expected SpecType to be UUID")
		}
	}
}

func TestMountString(t *testing.T) {
	for expectation, mount := range successfulMountStringExpectations {
		str := mount.String()
		if str != expectation {
			t.Errorf("Expected '%s', got '%s'", expectation, str)
		}
	}
}

func TestWritetoFile(t *testing.T) {
	mounts, err := ParseSystem()
	if err != nil {
		t.Fatal(err)
	}
	if err := mounts.WritetoFile("/etc/tmp_fstab"); err != nil {
		t.Fatal(err)
	}
}
