package falgosdk

import (
	"context"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
)

type Logger interface {
	Log(msg string)
}

type AlgodClient struct {
	*algod.Client
	FallbackClient  *algod.Client
	FallbackEnabled bool
	Logger          Logger
}

func (a *AlgodClient) AccountInformation(ctx context.Context, address string) (models.Account, error) {
	information, err := a.Client.AccountInformation(address).Do(ctx)
	if err != nil {
		if a.FallbackEnabled {
			if a.Logger != nil {
				a.Logger.Log("failed to get account information from primary client, using failover")
			}
			return a.FallbackClient.AccountInformation(address).Do(ctx)
		}
	}
	return information, nil
}

func (a *AlgodClient) EnableFallback() {
	a.FallbackEnabled = true
}

func (a *AlgodClient) DisableFallback() {
	a.FallbackEnabled = false
}
