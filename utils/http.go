package utils

import (
	"github.com/aronlt/toolkit/thttp"
	"github.com/sirupsen/logrus"
)

func DoQuery(args []byte, cookie ...string) ([]byte, error) {
	u := "https://leetcode.cn/graphql/"
	header := map[string]string{
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 Edg/123.0.0.0",
	}
	if len(cookie) > 0 {
		header["Cookie"] = cookie[0]
	}
	resp, err := thttp.PostJSON(u, args, thttp.Option{
		Decompress: false,
		Header:     header,
	})
	if err != nil {
		logrus.WithError(err).Errorf("call thttp.PostJSON fail")
		return []byte{}, err
	}
	return resp, nil
}
