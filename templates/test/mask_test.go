package test

import (
	"fmt"
	"log"
	"testing"
)

func MaskNumber(msisdn int64) string {

	str := fmt.Sprintf("%d",msisdn)
	return fmt.Sprintf("%s *** %s",str[:6],str[6+3:])

}

func TestMaskNumber(t *testing.T) {

	msisdn := int64(254726120256)
	log.Printf("%d -> %s ",msisdn,MaskNumber(msisdn))
}
