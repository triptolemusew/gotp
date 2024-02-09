package otp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

type ValidateOpts struct {
	Period    uint
	Digits    Digits
}

type Digits int

const (
	DigitsSix    Digits = 6
	DigitsEieght Digits = 8
)

func (d Digits) Length() int {
	return int(d)
}

func (d Digits) Format(in int32) string {
	f := fmt.Sprintf("%%0%dd", d)
	return fmt.Sprintf(f, in)
}

func GeneratePasscode(secret string, opts ValidateOpts) (string, error) {
	if opts.Period == 0 {
		opts.Period = 30
	}

	secret = strings.TrimSpace(secret)
	if n := len(secret) % 8; n != 0 {
		secret = secret + strings.Repeat("=", 8-n)
	}
	counter := uint64(math.Floor(float64(time.Now().UTC().Unix()) / float64(opts.Period)))

	secret = strings.ToUpper(secret)

	secretBytes, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", errors.New("Error validating secret invalid base32")
	}

	buf := make([]byte, 8)
	mac := hmac.New(sha1.New, secretBytes)
	binary.BigEndian.PutUint64(buf, counter)

	mac.Write(buf)
	sum := mac.Sum(nil)

	offset := sum[len(sum)-1] & 0xf

	value := int64(((int(sum[offset]) & 0x7f) << 24) |
		((int(sum[offset+1] & 0xff)) << 16) |
		((int(sum[offset+2] & 0xff)) << 8) |
		(int(sum[offset+3]) & 0xff))

	l := opts.Digits.Length()

	mod := int32(value % int64(math.Pow10(l)))

	passcode := opts.Digits.Format(mod)

	return passcode, nil
}
