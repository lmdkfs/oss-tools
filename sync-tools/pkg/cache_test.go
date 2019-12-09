package pkg

import (
	"github.com/gogo/protobuf/proto"
	"io/ioutil"
	"log"
	"testing"
)

func TestProto(t *testing.T) {
	ucache := &UfileCache{
		Ufilecache:          map[string]bool{"zrq": true},
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
	out, err := proto.Marshal(ucache)
	if err != nil {
		log.Fatalln("Failed to encode ucache:", err)
	}
	if err := ioutil.WriteFile("ucache", out, 0644); err != nil {
		log.Fatalln("Failed to write ucache:", err)
	}
}

func TestUnmarshalProto(t *testing.T) {
	in, err := ioutil.ReadFile("ucache")
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	ucache := &UfileCache{}
	if err := proto.Unmarshal(in, ucache); err !=nil {
		log.Fatalln("Failed to parse ucache:", err)
	}

	if _, ok := ucache.Ufilecache["zrq1"]; ok {
		log.Println("zrq 存在")
	}else {
		log.Println("zrq 不存在")
	}

}


