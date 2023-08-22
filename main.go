package main

import (
	"context"
	"fmt"

	"cosmossdk.io/simapp/params"
	amino "github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/ethereum/go-ethereum/common"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/evmos/evmos/v14/app"
	"github.com/evmos/evmos/v14/cmd/config"
	"github.com/evmos/evmos/v14/crypto/ethsecp256k1"
	enccodec "github.com/evmos/evmos/v14/encoding/codec"
	"google.golang.org/grpc"
)

const PRIV_KEY = "710145C4E48A4F31F00E5FEE3849ED816E5E06E6239947199E658AAA816285AA"
const addr2 = "evmos15r9hnaeflmzse5nzpnz2nffs4vfe6xgdr5q330"

func main() {
	grpcConn, err := grpc.Dial(
		"127.0.0.1:9090",    // Or your gRPC server address.
		grpc.WithInsecure(), // The Cosmos SDK doesn't support any transport security mechanism.
	)
	if err != nil {
		return
	}
	defer grpcConn.Close()

	sender, _ := generateSenderAccount()

	accSeq, err := GetSequence(context.Background(), *grpcConn, sender.String())
	if err != nil {
		return
	}

    fmt.Println("Sequence: ", accSeq)
}

func generateSenderAccount() (accAddr sdktypes.AccAddress, priv cryptotypes.PrivKey) {
	privKey := &ethsecp256k1.PrivKey{
		Key: common.FromHex(PRIV_KEY),
	}

	sdktypes.GetConfig().SetBech32PrefixForAccount(config.Bech32Prefix, "")
	addr1 := privKey.PubKey().Address()
	sender := sdktypes.AccAddress(addr1)
	return sender, privKey
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
