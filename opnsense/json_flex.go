package opnsense

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// flexInt is an integer that tolerates the inconsistent JSON shapes the
// OPNsense MVC API uses for numeric Kea lease fields. Depending on the
// OPNsense release and Kea plugin build, a field such as valid_lifetime,
// expire, state, subnet_id or pool_id may arrive as a bare JSON number
// (86400), a quoted number ("86400"), an empty string ("") or null. All of
// these decode into an int.
//
// Bounds: the token is parsed at the width of the platform int, so the
// conversion below can never truncate. All release targets are 64-bit,
// where every legitimate DHCP lease value fits. valid_lifetime tops out
// at 0xffffffff (4294967295, the "infinite" lease) and expire is a unix
// timestamp. A value that overflows int, or any non-numeric token, is
// returned as an error rather than silently coerced to 0, so a malformed
// payload surfaces loudly instead of emitting a wrong metric. Empty
// string and null are treated as a legitimate "absent" value and decode
// to 0.
type flexInt int

func (f *flexInt) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || bytes.Equal(b, []byte("null")) {
		*f = 0
		return nil
	}

	// Unwrap a quoted value -- "86400" -> 86400 -- decoding the JSON string
	// so any escaping is handled correctly before we parse it.
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return fmt.Errorf("flexInt: %w", err)
		}
		s = strings.TrimSpace(s)
		if s == "" {
			*f = 0
			return nil
		}
		b = []byte(s)
	}

	// Bit size 0 parses at the width of int, so the conversion below is
	// always lossless. Values that overflow int fail here instead.
	n, err := strconv.ParseInt(string(b), 10, 0)
	if err != nil {
		return fmt.Errorf("flexInt: cannot parse %q as integer: %w", string(b), err)
	}
	*f = flexInt(n)
	return nil
}

// flexStringSlice is a []string that tolerates the OPNsense MVC API
// returning the Kea is_reserved field either as a JSON array (["hwaddr"])
// or as a bare string ("hwaddr" for a reserved lease, "" otherwise). A
// non-empty string becomes a single-element slice; an empty string or null
// becomes an empty slice, so callers can keep using len() to detect whether
// a lease is reserved.
type flexStringSlice []string

func (s *flexStringSlice) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || bytes.Equal(b, []byte("null")) {
		*s = nil
		return nil
	}

	if b[0] == '[' {
		var arr []string
		if err := json.Unmarshal(b, &arr); err != nil {
			return fmt.Errorf("flexStringSlice: %w", err)
		}
		*s = arr
		return nil
	}

	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return fmt.Errorf("flexStringSlice: %w", err)
	}
	if strings.TrimSpace(str) == "" {
		*s = nil
		return nil
	}
	*s = flexStringSlice{str}
	return nil
}
