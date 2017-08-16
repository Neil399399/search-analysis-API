package datamodel

import (
	"errors"
	"testing"
)

var (
	ErrVerifyNOPass = errors.New("Verify NO Pass")
)

func TestVerify(t *testing.T) {

	testsearch := Search{}
	testsearch.APIKEY = "qwj5661&**&"
	testsearch.KEYWORD = "apple"

	testok := testsearch.Verify(testsearch)
	if !testok {
		t.Error(ErrVerifyNOPass)
	}

}
