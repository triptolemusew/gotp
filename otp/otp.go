package otp

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/triptolemusew/gotp/db"
)

type Type string

const (
	TOTP Type = "totp"
	HOTP Type = "hotp"
)

// type URL struct {
// 	Type      string
// 	Issuer    string
// 	Account   string
// 	Secret    string
// 	Algorithm string
// 	Counter   uint64
// 	Digits    int
// 	Period    int
// }

func ParseURL(s string) (*db.Key, error) {
	var typeLabel string

	out := new(db.Key)

	if ps, err := url.Parse(s); err == nil {
		if ps.Scheme != "otpauth" {
			return nil, fmt.Errorf("invalid scheme %q", ps.Scheme)
		}

		s = strings.ReplaceAll(s, "otpauth://", "")
		if x := strings.SplitN(s, "?", 2); len(x) == 2 {
			typeLabel = x[0]
		}
		if val, ok := ps.Query()["secret"]; ok {
			out.Secret = val[0]
		}
		if val, ok := ps.Query()["algorithm"]; ok {
			out.Algorithm = strings.ToUpper(val[0])
		}
		if val, ok := ps.Query()["counter"]; ok {
			n, _ := strconv.ParseUint(val[0], 10, 64)
			out.Counter = n
		}
		if val, ok := ps.Query()["digits"]; ok {
			n, _ := strconv.ParseInt(val[0], 10, 64)
			out.Digits = int(n)
		}
		if val, ok := ps.Query()["period"]; ok {
			n, _ := strconv.ParseInt(val[0], 10, 64)
			out.Period = int(n)
		}
	}

	ps := strings.SplitN(strings.TrimPrefix(typeLabel, "//"), "/", 2)
	if len(ps) != 2 || ps[0] == "" || ps[1] == "" {
		return nil, errors.New("invalid type/label")
	}

	out.Type = ps[0]
	account, err := url.PathUnescape(ps[1])
	if err != nil {
		return nil, err
	}

	if i := strings.Index(account, ":"); i >= 0 {
		out.Issuer = strings.TrimSpace(account[:i])
		if out.Issuer == "" {
			return nil, errors.New("empty issuer")
		}
		account = account[i+1:]
	}
	out.Account = strings.TrimSpace(account)

	return out, nil
}

func GetType(s string) Type {
	switch s {
	case "hotp":
		return HOTP
	case "totp":
		return TOTP
	}
	return HOTP
}
