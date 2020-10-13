// Code generated by "stringer -type syzcall"; DO NOT EDIT.

package namespace

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[BIND-2]
	_ = x[CHDIR-3]
	_ = x[UNMOUNT-35]
	_ = x[MOUNT-46]
	_ = x[RFORK-19]
	_ = x[IMPORT-7]
	_ = x[INCLUDE-14]
}

const (
	_syzcall_name_0 = "BINDCHDIR"
	_syzcall_name_1 = "IMPORT"
	_syzcall_name_2 = "INCLUDE"
	_syzcall_name_3 = "RFORK"
	_syzcall_name_4 = "UNMOUNT"
	_syzcall_name_5 = "MOUNT"
)

var (
	_syzcall_index_0 = [...]uint8{0, 4, 9}
)

func (i syzcall) String() string {
	switch {
	case 2 <= i && i <= 3:
		i -= 2
		return _syzcall_name_0[_syzcall_index_0[i]:_syzcall_index_0[i+1]]
	case i == 7:
		return _syzcall_name_1
	case i == 14:
		return _syzcall_name_2
	case i == 19:
		return _syzcall_name_3
	case i == 35:
		return _syzcall_name_4
	case i == 46:
		return _syzcall_name_5
	default:
		return "syzcall(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}