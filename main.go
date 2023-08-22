package main

import (
	"context"
	"fmt"

	"cosmossdk.io/simapp/params"
	amino "github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/evmos/evmos/v14/app"
	enccodec "github.com/evmos/evmos/v14/encoding/codec"
	"google.golang.org/grpc"
)

func main() {
	grpcConn, err := grpc.Dial(
		"127.0.0.1:9090",    // Or your gRPC server address.
		grpc.WithInsecure(), // The Cosmos SDK doesn't support any transport security mechanism.
	)
	if err != nil {
		return
	}
	defer grpcConn.Close()

}

func GetSequence(ctx context.Context, gcConn grpc.ClientConn, addr string) (uint64, error) {
	authClient := authtypes.NewQueryClient(&gcConn)
	res, err := authClient.Account(ctx, &authtypes.QueryAccountRequest{
		Address: addr,
	})

	if err != nil {
		fmt.Println("error getting account", err)
		return 0, err
	}

	encodingCfg := makeEncodingConfig()
	var account authtypes.AccountI
	err = encodingCfg.InterfaceRegistry.UnpackAny(res.Account, &account)
	if err != nil {
		fmt.Println("error unpacking account", err)
		return 0, err
	}

	return account.GetSequence(), nil
}

func makeEncodingConfig() params.EncodingConfig {
	mb := app.ModuleBasics
	cdc := amino.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	codec := amino.NewProtoCodec(interfaceRegistry)

	encodingConfig := params.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             codec,
		TxConfig:          tx.NewTxConfig(codec, tx.DefaultSignModes),
		Amino:             cdc,
	}

	enccodec.RegisterLegacyAminoCodec(encodingConfig.Amino)
	mb.RegisterLegacyAminoCodec(encodingConfig.Amino)
	enccodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	mb.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
