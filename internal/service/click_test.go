package service

import "testing"

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
