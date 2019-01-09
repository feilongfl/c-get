package core

import (
	"github.com/robertkrimen/otto"
	log "github.com/sirupsen/logrus"
	"strings"
)

func EvalDecode(code string) (string, error) {
	code = strings.Replace(code, "eval", "tempval=", 1)
	vm := otto.New()
	vm.Run(code)
	value, err := vm.Get("tempval")
	if err != nil {
		return "", err
	}
	value_str, err := value.ToString()
	if err != nil {
		return "", err
	}
	log.WithFields(log.Fields{
		"code":   code,
		"decode": value_str,
	}).Debug("eval js decode")
	return value_str, nil
}

func EvalDecodeNew(code string) (string, error) {
	code = strings.Replace(code, "eval", "", 1)
	vm := otto.New()
	value, err := vm.Run(code)
	//value, err := vm.Get("tempval")
	if err != nil {
		return "", err
	}
	value_str, err := value.ToString()
	if err != nil {
		return "", err
	}
	log.WithFields(log.Fields{
		"code":   code,
		"decode": value_str,
	}).Debug("eval js decode")
	return value_str, nil
}
