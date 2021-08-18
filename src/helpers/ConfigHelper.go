package helpers

import (
	"github.com/shenyisyn/goft-gin/goft"
	"io/ioutil"
	"os"
)
func CertData(path string ) []byte{
	f,err:=os.Open(path )
	goft.Error(err)
	defer f.Close()
	b,err:=ioutil.ReadAll(f)
	goft.Error(err)
	return b
}