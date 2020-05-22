package service

import (
	"github.com/hpifu/go-kit/hhttp"
	"testing"
	"time"
)

func TestCheckReqDevId(t *testing.T) {
	req := &ClickReq{}
	if len(CheckReqDevId(req)) == 0 {
		req.Idfa = "iphone"
		if CheckReqDevId(req) == req.Idfa {
			t.Log("Yes baby.")
		} else {
			t.Error("Are you kidding me?")
		}
	} else {
		t.Error("what are you doing!!!")
	}
}

func TestSubstr(t *testing.T)  {
	m := "17744581949"
	t.Log(m[:7])
}

func TestHttpGet(t *testing.T)  {
	client := hhttp.NewHttpClient(20, 200*time.Millisecond, 200*time.Millisecond)

	url := "http://advertiser.equblock.com/click"
	params := make(map[string]interface {})
	params["pub"] = "jim"
	params["cid"] = "cid"
	params["idfa"] = "fengming-phone"

	res := client.Do("GET", url, nil, params, nil)
	if res.Err != nil {
		t.Errorf("response [%d][%s]", res.Status, res.Err.Error())
	} else {
		t.Log(string(res.Res))
	}
}