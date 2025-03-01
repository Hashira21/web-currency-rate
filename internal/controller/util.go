package controller

import (
	"regexp"
	"strings"
)

const isoCodePattern = "^[A-Z]{3}$"

func (ctr *controller) validateIsoCode(isoCodes ...*string) (string, bool) {
	for _, isoCode := range isoCodes {
		*isoCode = strings.ToUpper(*isoCode)

		re := regexp.MustCompile(isoCodePattern)
		if !re.MatchString(*isoCode) {
			return *isoCode, false
		}

		if _, ok := ctr.validIsoCodes[*isoCode]; !ok {
			return *isoCode, false
		}
	}

	return "", true
}

func (ctr *controller) getValidIsoCodesString() string {
	var res strings.Builder
	for code := range ctr.validIsoCodes {
		res.WriteString(code)
		res.WriteString(" ")
	}
	return res.String()[:len(res.String())-1]
}
