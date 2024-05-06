package falgosdk_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	falgosdk "github.com/stein-f/algosdk-with-fallback"
	"github.com/stretchr/testify/assert"
)

const address = "ABCEDCMH2IQXMU37WR7SJH4WXXUGC2TB35WVMQRH3S5TOZE3VQRZEFJE5E"

type testLogger struct {
	wasInvoked bool
}

func (t *testLogger) Log(msg string) {
	t.wasInvoked = true
	fmt.Println(msg)
}

func TestUsePrimaryClient(t *testing.T) {
	accInfo, err := os.ReadFile("testdata/account_information.json")
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(accInfo)
	}))
	defer srv.Close()

	testLogger := &testLogger{}
	primaryAlgod, err := algod.MakeClient(srv.URL, "")
	if err != nil {
		t.Fatal(err)
	}
	secondaryAlgod, err := algod.MakeClient("http://localhost:80", "")
	if err != nil {
		t.Fatal(err)
	}
	algoD := falgosdk.AlgodClient{
		Client:          primaryAlgod,
		FallbackClient:  secondaryAlgod,
		FallbackEnabled: true,
		Logger:          testLogger,
	}

	accountInformation, err := algoD.AccountInformation(context.Background(), address)

	assert.NoError(t, err)
	assert.Equal(t, address, accountInformation.Address)
	assert.False(t, testLogger.wasInvoked)
}

func TestFailsOver(t *testing.T) {
	accInfo, err := os.ReadFile("testdata/account_information.json")
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(accInfo)
	}))
	defer srv.Close()

	testLogger := &testLogger{}
	primaryAlgod, err := algod.MakeClient("http://localhost:80", "")
	if err != nil {
		t.Fatal(err)
	}
	secondaryAlgod, err := algod.MakeClient(srv.URL, "")
	if err != nil {
		t.Fatal(err)
	}
	algoD := falgosdk.AlgodClient{
		Client:          primaryAlgod,
		FallbackClient:  secondaryAlgod,
		FallbackEnabled: true,
		Logger:          testLogger,
	}

	accountInformation, err := algoD.AccountInformation(context.Background(), address)

	assert.NoError(t, err)
	assert.Equal(t, address, accountInformation.Address)
	assert.True(t, testLogger.wasInvoked)
}
