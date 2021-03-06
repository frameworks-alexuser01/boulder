// Copyright 2015 ISRG.  All rights reserved
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package core

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/letsencrypt/boulder/test"
)

func TestProblemDetails(t *testing.T) {
	pd := &ProblemDetails{
		Type:   MalformedProblem,
		Detail: "Wat? o.O"}
	test.AssertEquals(t, pd.Error(), "urn:acme:error:malformed :: Wat? o.O")
}

func TestRegistrationUpdate(t *testing.T) {
	oldURL, _ := url.Parse("http://old.invalid")
	newURL, _ := url.Parse("http://new.invalid")

	reg := Registration{
		ID:        1,
		Contact:   []AcmeURL{AcmeURL(*oldURL)},
		Agreement: "",
	}
	update := Registration{
		Contact:   []AcmeURL{AcmeURL(*newURL)},
		Agreement: "totally!",
	}

	reg.MergeUpdate(update)
	test.Assert(t, len(reg.Contact) == 1 && reg.Contact[0] == update.Contact[0], "Contact was not updated %v != %v")
	test.Assert(t, reg.Agreement == update.Agreement, "Agreement was not updated")
}

func TestSanityCheck(t *testing.T) {
	tls := true
	chall := Challenge{Type: ChallengeTypeSimpleHTTP, Status: StatusValid}
	test.Assert(t, !chall.IsSane(false), "IsSane should be false")
	chall.Status = StatusPending
	test.Assert(t, !chall.IsSane(false), "IsSane should be false")
	chall.R = "bad"
	chall.S = "bad"
	chall.Nonce = "bad"
	test.Assert(t, !chall.IsSane(false), "IsSane should be false")
	chall = Challenge{Type: ChallengeTypeSimpleHTTP, Path: "bad", Status: StatusPending}
	test.Assert(t, !chall.IsSane(false), "IsSane should be false")
	chall.Token = ""
	test.Assert(t, !chall.IsSane(false), "IsSane should be false")
	chall.Token = "notlongenough"
	test.Assert(t, !chall.IsSane(false), "IsSane should be false")
	chall.Token = "evaGxfADs6pSRb2LAv9IZf17Dt3juxGJ+PCt92wr+o!"
	test.Assert(t, !chall.IsSane(false), "IsSane should be false")
	chall.Token = "KQqLsiS5j0CONR_eUXTUSUDNVaHODtc-0pD6ACif7U4"
	chall.Path = ""
	test.Assert(t, !chall.IsSane(false), "IsSane should be false")
	chall.TLS = &tls
	test.Assert(t, chall.IsSane(false), "IsSane should be true")

	test.Assert(t, !chall.IsSane(true), "IsSane should be false")
	chall.Path = "../.."
	test.Assert(t, !chall.IsSane(true), "IsSane should be false")
	chall.Path = "/asd"
	test.Assert(t, !chall.IsSane(true), "IsSane should be false")
	chall.Path = "bad//test"
	test.Assert(t, !chall.IsSane(true), "IsSane should be false")
	chall.Path = "bad/./test"
	test.Assert(t, !chall.IsSane(true), "IsSane should be false")
	chall.Path = "good"
	test.Assert(t, chall.IsSane(true), "IsSane should be true")
	chall.Path = "good/test"
	test.Assert(t, chall.IsSane(true), "IsSane should be true")

	chall = Challenge{Type: ChallengeTypeDVSNI, Status: StatusPending}
	chall.Path = "bad"
	chall.Token = "bad"
	chall.TLS = &tls
	test.Assert(t, !chall.IsSane(false), "IsSane should be false")
	chall = Challenge{Type: ChallengeTypeDVSNI, Status: StatusPending}
	test.Assert(t, !chall.IsSane(false), "IsSane should be false")
	chall.Nonce = "wutwut"
	test.Assert(t, !chall.IsSane(false), "IsSane should be false")
	chall.Nonce = "!2345678901234567890123456789012"
	test.Assert(t, !chall.IsSane(false), "IsSane should be false")
	chall.Nonce = "12345678901234567890123456789012"
	test.Assert(t, !chall.IsSane(false), "IsSane should be false")
	chall.R = "notlongenough"
	test.Assert(t, !chall.IsSane(false), "IsSane should be false")
	chall.R = "evaGxfADs6pSRb2LAv9IZf17Dt3juxGJ+PCt92wr+o!"
	test.Assert(t, !chall.IsSane(false), "IsSane should be false")
	chall.R = "KQqLsiS5j0CONR_eUXTUSUDNVaHODtc-0pD6ACif7U4"
	test.Assert(t, chall.IsSane(false), "IsSane should be true")
	chall.S = "anything"
	test.Assert(t, !chall.IsSane(false), "IsSane should be false")
	test.Assert(t, !chall.IsSane(true), "IsSane should be false")
	chall.S = "evaGxfADs6pSRb2LAv9IZf17Dt3juxGJ+PCt92wr+o!"
	test.Assert(t, !chall.IsSane(true), "IsSane should be false")
	chall.S = "KQqLsiS5j0CONR_eUXTUSUDNVaHODtc-0pD6ACif7U4"
	test.Assert(t, chall.IsSane(true), "IsSane should be true")

	chall = Challenge{Type: "bogus", Status: StatusPending}
	test.Assert(t, !chall.IsSane(false), "IsSane should be false")
	test.Assert(t, !chall.IsSane(true), "IsSane should be false")
}

func TestJSONBufferUnmarshal(t *testing.T) {
	testStruct := struct {
		Buffer JSONBuffer
	}{}

	notValidBase64 := []byte(`{"Buffer":"!!!!"}`)
	err := json.Unmarshal(notValidBase64, &testStruct)
	test.Assert(t, err != nil, "Should have choked on invalid base64")
}
