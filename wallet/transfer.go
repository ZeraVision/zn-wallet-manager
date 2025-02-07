package wallet

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zn-wallet-manager/transcode"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Inputs struct {
	B58Address         string
	KeyType            KeyType
	PublicKey          string     // Base 58 encoded
	PrivateKey         string     // Base 58 encoded
	Amount             *big.Float // full coins (not parts)
	FeePercent         float32    // 0-100 max 6 digits of precision
	ContractFeePercent *float32   // 0-100 max 6 digits of precision
}

// CreateCoinTxn creates a CoinTXN protobuf message for a given set of inputs and outputs
// Inputs: [Key] = address, [Value] = Inputs struct
// Outputs: [Key] = address, [Value] = amount of whole coins (not parts)
// symbol: contract symbol to send (ex $ZRA+0000)
// baseFeeID: fee id for the base fee (ex $ZRA+0000)
// baseFeeAmountParts: fee amount for the base fee in *parts* (ex 1000000000 = 1 ZRA)
// contractFeeID: fee id for the contract fee (ex $ZRA+0000)
// contractFeeAmountParts: fee amount for the contract fee in *parts* (ex 1000000000 = 1 ZRA)
func CreateCoinTxn(inputs []Inputs, outputs map[string]*big.Float, symbol, baseFeeID, baseFeeAmountParts string, contractFeeID, contractFeeAmountParts *string) (*pb.CoinTXN, error) {
	// TODO parts lookup on zv indexer or other source
	parts := big.NewInt(1_000_000_000) // Example: amount of parts for ZRA

	// Step 1: Process Inputs
	inputTransfers, auth, keys, totalInput, err := processInputs(inputs, parts)
	if err != nil {
		return nil, err
	}

	// Step 2: Process Outputs
	outputTransfers, totalOutput, err := processOutputs(outputs, parts)
	if err != nil {
		return nil, err
	}

	// Check to see if inputs and outputs match
	if totalInput.Cmp(totalOutput) != 0 {
		return nil, fmt.Errorf("total input does not equal total output: %s != %s", totalInput.String(), totalOutput.String())
	}

	// Step 3: Build Transfer Authentication
	transferAuth := buildTransferAuthentication(auth)

	// Step 4: Build Transaction Base
	txnBase := buildTransactionBase(baseFeeID, baseFeeAmountParts)

	// Step 5: Assemble Transaction
	txn := &pb.CoinTXN{
		Auth:              transferAuth,
		Base:              txnBase,
		ContractId:        symbol,
		InputTransfers:    inputTransfers,
		OutputTransfers:   outputTransfers,
		ContractFeeId:     contractFeeID,
		ContractFeeAmount: contractFeeAmountParts,
	}

	// Step 6: Serialize and Sign Transaction
	txn, err = signTransaction(txn, keys)
	if err != nil {
		return nil, err
	}

	// Step 7: Marshal the vote with the signature
	byteDataWithSig, err := proto.Marshal(txn)
	if err != nil {
		return nil, fmt.Errorf("error while serializing txn: %v", err)
	}

	// Step 8: Hash the serialized data with the signature
	hash := transcode.SHA3_256(byteDataWithSig)

	// Step 9: Add the hash
	txn.Base.Hash = hash

	return txn, nil
}

type keyTracking struct {
	KeyType    KeyType
	PrivateKey string
}

type authTracking struct {
	PublicKeyBytes []byte
	Signature      []byte
	Nonce          uint64
}

// Helper Function: Process Inputs
func processInputs(inputs []Inputs, parts *big.Int) ([]*pb.InputTransfers, []authTracking, map[string]keyTracking, *big.Float, error) {
	var inputTransfers []*pb.InputTransfers
	var auth []authTracking
	keys := map[string]keyTracking{}
	totalInput := big.NewFloat(0)
	index := uint64(0)

	for _, input := range inputs {
		_, _, pubKeyByte, err := transcode.Base58DecodePublicKey(input.PublicKey)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("could not decode public key: %v", err)
		}

		amountPartsBigF := new(big.Float).Mul(input.Amount, big.NewFloat(float64(parts.Int64())))

		inputTransfers = append(inputTransfers, &pb.InputTransfers{
			Index:      index,
			Amount:     amountPartsBigF.String(),
			FeePercent: uint32(input.FeePercent * 1_000_000),
		})

		nonce := GetNonce(input.B58Address)

		auth = append(auth, authTracking{
			PublicKeyBytes: pubKeyByte,
			Signature:      nil,
			Nonce:          nonce,
		})

		keys[transcode.Base58Encode(pubKeyByte)] = keyTracking{
			KeyType:    input.KeyType,
			PrivateKey: input.PrivateKey,
		}

		totalInput.Add(totalInput, amountPartsBigF)
		index++
	}

	return inputTransfers, auth, keys, totalInput, nil
}

// Helper Function: Process Outputs
func processOutputs(outputs map[string]*big.Float, parts *big.Int) ([]*pb.OutputTransfers, *big.Float, error) {
	var outputsTransfers []*pb.OutputTransfers
	totalOutput := big.NewFloat(0)

	for address, amount := range outputs {
		decodedAddr, err := transcode.Base58Decode(address)
		if err != nil {
			return nil, nil, fmt.Errorf("could not decode address: %v", err)
		}

		bigFParts := new(big.Float).Mul(amount, new(big.Float).SetInt(parts))

		outputsTransfers = append(outputsTransfers, &pb.OutputTransfers{
			WalletAddress: decodedAddr,
			Amount:        bigFParts.String(),
		})

		totalOutput.Add(totalOutput, bigFParts)
	}

	return outputsTransfers, totalOutput, nil
}

// Helper Function: Build Transfer Authentication
func buildTransferAuthentication(auth []authTracking) *pb.TransferAuthentication {
	transferAuth := &pb.TransferAuthentication{}
	for _, a := range auth {
		transferAuth.PublicKey = append(transferAuth.PublicKey, &pb.PublicKey{Single: a.PublicKeyBytes})
		transferAuth.Nonce = append(transferAuth.Nonce, a.Nonce)
	}
	return transferAuth
}

// Helper Function: Build Transaction Base
func buildTransactionBase(feeID, feeAmountParts string) *pb.BaseTXN {
	return &pb.BaseTXN{
		Timestamp: timestamppb.Now(),
		FeeAmount: feeAmountParts,
		FeeId:     feeID,
	}
}

// Helper Function: Sign Transaction
func signTransaction(txn *pb.CoinTXN, keys map[string]keyTracking) (*pb.CoinTXN, error) {
	txnBytes, err := proto.Marshal(txn)
	if err != nil {
		return nil, fmt.Errorf("could not marshal transaction: %v", err)
	}

	for _, auth := range txn.Auth.PublicKey {
		if key, ok := keys[transcode.Base58Encode(auth.Single)]; ok {
			signature, err := Sign(key.PrivateKey, txnBytes, key.KeyType)
			if err != nil {
				return nil, fmt.Errorf("could not sign transaction: %v", err)
			}
			txn.Auth.Signature = append(txn.Auth.Signature, signature)
		} else {
			return nil, fmt.Errorf("could not find private key for public key: %s", transcode.Base58Encode(auth.Single))
		}
	}
	return txn, nil
}

// struct for client implementation of grpcs
type NetworkClient struct {
	client pb.TXNServiceClient
}

// constructor for client implementation of grpcs
func NewNetworkClient(conn *grpc.ClientConn) *NetworkClient {
	client := pb.NewTXNServiceClient(conn)
	return &NetworkClient{client: client}
}

func SendCoinTXN(txn *pb.CoinTXN) (*emptypb.Empty, error) {
	// Create a gRPC connection to the server
	conn, err := grpc.Dial(os.Getenv("GRPC_ADDR")+":50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to the server: %v", err)
		return nil, err
	}
	defer conn.Close()

	// Create a new instance of ValidatorNetworkClient
	client := NewNetworkClient(conn)

	response, err := client.client.Coin(context.Background(), txn)

	if err != nil {
		return nil, err
	}

	return response, nil
}
