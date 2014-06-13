package fstab

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Mount represetnts the filesystem info
type Mount struct {
	// The block special device or remote filesystem to be mounted
	Spec string

	// The mount point for the filesystem
	File string

	// The type of the filesystem
	VfsType string

	// Mount options associated with the filesystem
	MntOps map[string]string

	// Used by dump to determine which filesystems need to be dumped.
	Freq int

	// Used by the fsck(8) program to determine the order in which filesystem checks are done at reboot time
	PassNo int
}

const (
	Path = 0
	Label
	UUID
	PartUUID
	PartLabel
)

type DeviceIdentifierType int

// parseOptions parses the options field into an array of strings
func parseOptions(optionsString string) (options map[string]string) {
	options = make(map[string]string)
	for _, option := range strings.Split(optionsString, ",") {
		bits := strings.Split(strings.TrimSpace(option), "=")
		if len(bits) > 1 {
			options[bits[0]] = bits[1]
		} else {
			options[bits[0]] = ""
		}
	}
	return
}

func (mount *Mount) MntOpsString() (opsstring string) {
	first := true
	for key, value := range mount.MntOps {
		if first {
			first = false
		} else {
			opsstring += ","
		}

		opsstring += key

		if "" != value {
			opsstring += "=" + value
		}
	}
	return
}

// String serializes the object into fstab format
func (mount *Mount) String() string {
	return fmt.Sprintf("%-21s %-21s %-21s %-21s %2d %2d", mount.Spec, mount.File, mount.VfsType, mount.MntOpsString(), mount.Freq, mount.PassNo)
}

func (mount *Mount) IsSwap() bool {
	return "swap" == mount.VfsType
}

func (mount *Mount) IsNFS() bool {
	return "nfs" == mount.VfsType
}

func (mount *Mount) Equals(other *Mount) bool {
	return reflect.DeepEqual(*mount, *other)
}

// SpecType returns the device identifier type
func (mount *Mount) SpecType() (spectype DeviceIdentifierType) {
	bits := strings.Split(mount.Spec, "=")
	switch strings.ToUpper(bits[0]) {
	case "UUID":
		spectype = UUID

	case "LABEL":
		spectype = Label

	case "PARTUUID":
		spectype = PartUUID

	case "PARTLABEL":
		spectype = PartLabel

	default:
		spectype = Path
	}
	return
}

// SpecType returns the device identifier value; that is if Spec is
// "UUID=vogons-ate-my-sandwich", it will return "vogons-ate-my-sandwich"
func (mount *Mount) SpecValue() string {
	bits := strings.Split(mount.Spec, "=")
	if 1 == len(bits) {
		return mount.Spec
	} else {
		return bits[1]
	}
}

// ParseLine parses a single line (of an fstab).
// It will most frequently return a Mount; however,
// If a parsing error occurs, `err` will be non-nil and provide an error message.
// If the line is either empy or a comment line, `mount` will also be nil.
func ParseLine(line string) (mount *Mount, err error) {
	line = strings.TrimSpace(line)

	// Lines starting with a pound sign (#) are comments, and are ignored. So are empty lines.
	if ("" == line) || (line[0] == '#') {
		return nil, nil
	}

	fields := strings.Fields(line)
	if len(fields) < 4 {
		return nil, fmt.Errorf("too few fields (%i), at least 4 are expected", len(fields))
	} else {
		mount = new(Mount)
		mount.Spec = fields[0]
		mount.File = fields[1]
		mount.VfsType = fields[2]
		mount.MntOps = parseOptions(fields[3])

		var convErr error

		if len(fields) > 4 {
			mount.Freq, convErr = strconv.Atoi(fields[4])
			if nil != convErr {
				return nil, fmt.Errorf("%s is not a number", fields[4])
			}
		}

		if len(fields) > 5 {
			mount.PassNo, convErr = strconv.Atoi(fields[5])
			if nil != convErr {
				return nil, fmt.Errorf("%s it not a number", fields[5])
			}
		}
	}

	return
}
