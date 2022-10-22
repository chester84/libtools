package libtools

import (
	"testing"

	"github.com/beego/beego/v2/core/logs"
)

func TestSetByFields(t *testing.T) {

	type Respond struct {
		Server int64
		Info   struct {
			Int         int
			Int64       int64
			String      string
			SliceInt    []int
			SliceString []string
			Struct      struct {
				Int    int
				String string
			}
		}
	}

	resp := Respond{}
	SetByFields(&resp, "Server", int64(1564451339282))
	if resp.Server != 1564451339282 {
		t.Errorf(`SetByFields Server no ok. [%v] `, resp.Server)
	} else {
		logs.Notice("SetByFields Server ok")
	}

	SetByFields(&resp, "Info.Int", int(1))
	if resp.Info.Int != 1 {
		t.Errorf(`SetByFields Info.Int no ok. [%v] `, resp.Info.Int)
	} else {
		logs.Notice("SetByFields Info.Int ok")
	}

	SetByFields(&resp, "Info.String", "test_string")
	if resp.Info.String != "test_string" {
		t.Errorf(`SetByFields Info.String no ok. [%v] `, resp.Info.String)
	} else {
		logs.Notice("SetByFields Info.String ok")
	}

	SetByFields(&resp, "Info.Struct.String", "test_struct_string")
	if resp.Info.Struct.String != "test_struct_string" {
		t.Errorf(`SetByFields Info.Struct.String no ok. [%v] `, resp.Info.Struct.String)
	} else {
		logs.Notice("SetByFields Info.Struct.String ok")
	}

	SetByFields(&resp, "Info.SliceString", "test_struct_string")
	if len(resp.Info.SliceString) != 1 {
		t.Errorf(`SetByFields Info.SliceString no ok. [%v] `, resp.Info.SliceString)
	} else {
		logs.Notice("SetByFields Info.SliceString  ok")
	}

	logs.Notice("resp:%v", resp)
}
