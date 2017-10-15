// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// InvokeFunction a function
func InvokeFunction(gateway string, name string, bytesIn *[]byte, contentType string, query []string) (*[]byte, error) {
	var resBytes []byte

	gateway = strings.TrimRight(gateway, "/")

	reader := bytes.NewReader(*bytesIn)
	client := http.Client{}
	qs := ""

	if len(query) > 0 {
		qs = "?"
		for _, queryValue := range query {
			qs = qs + queryValue + "&"
			if strings.Contains(queryValue, "=") == false {
				return nil, fmt.Errorf("The --query flags must take the form of key=value (= not found)")
			}
			if strings.HasSuffix(queryValue, "=") {
				return nil, fmt.Errorf("The --query flag must take the form of: key=value (empty value given, or value ends in =)")
			}
		}
		qs = strings.TrimRight(qs, "&")
	}

	gatewayURL := gateway + "/function/" + name + qs
	// fmt.Println(gatewayURL)
	req, _ := http.NewRequest(http.MethodPost, gatewayURL, reader)
	req.Header.Add("Content-Type", contentType)

	res, err := client.Do(req)

	if err != nil {
		fmt.Println()
		fmt.Println(err)
		return nil, fmt.Errorf("cannot connect to OpenFaaS on URL: %s", gateway)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	switch res.StatusCode {
	case 200:
		var readErr error
		resBytes, readErr = ioutil.ReadAll(res.Body)
		if readErr != nil {
			return nil, fmt.Errorf("cannot read result from OpenFaaS on URL: %s %s", gateway, readErr)
		}

	default:
		bytesOut, err := ioutil.ReadAll(res.Body)
		if err == nil {
			return nil, fmt.Errorf("Server returned unexpected status code: %d - %s", res.StatusCode, string(bytesOut))
		}
	}

	return &resBytes, nil
}
