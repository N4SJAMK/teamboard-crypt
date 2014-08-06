package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

const (
	URL = "http://localhost:9004"
)

func request(resource string, p *Password) (*http.Response, error) {
	jb, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return http.Post((URL + "/" + resource),
		"application/json", bytes.NewReader(jb))
}

func TestHash(test *testing.T) {
	password := Password{Plain: "kakka"}
	res, err := request("hash", &password)
	if err != nil {
		test.Fatal(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		test.Fatal(err.Error())
	}

	// if status is ok, response should be json
	if res.StatusCode == http.StatusOK {
		var jsonres Password
		if err := json.Unmarshal(body, &jsonres); err != nil {
			test.Fatal(err.Error())
		}
		if string(jsonres.Hash[0]) != "$" {
			test.Fatal("Expected first character of hash to be $")
		}
	} else {
		test.Fatal(string(body))
	}
}

func TestCompare(test *testing.T) {
	expectations := []*Password{
		&Password{
			Hash:  "$2a$10$OwumTlyVhkSdlzD52IC4bODp7JoXH5.ZTVLh41gdY4/uyVJ1x5Roy",
			Plain: "ääööääöö",
			Match: true,
		},
		&Password{
			Hash:  "$2a$10$F8dMVuvWiZJT62yOgSNgYevd2SBSyDPPEC.7E8fPJQKJTlNhtljdC",
			Plain: "ääööääöö",
			Match: true,
		},
		&Password{
			Hash:  "iiririiriririr123",
			Plain: "cascada",
			Match: false,
		},
	}

	for _, expectation := range expectations {
		res, err := request("compare", expectation)
		if err != nil {
			test.Fatal(err.Error())
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			test.Fatal(err.Error())
		}

		if res.Header.Get("Content-Type") == "application/json" {
			var jsonres Password
			if err := json.Unmarshal(body, &jsonres); err != nil {
				test.Fatal(err.Error())
			}
			if jsonres.Match != expectation.Match {
				test.Fail()
				test.Logf("expected %v to produce match %t, got %t\n",
					expectation, expectation.Match, jsonres.Match)
			}
		} else {
			test.Fail()
			test.Logf("[%d] %s\n", res.StatusCode, string(body))
		}
	}
}
