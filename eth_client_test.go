package eth

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"testing"
)

const (
	EndPoint   = "http://192.168.3.173:18545"
	PrivateKey = "a492823c3e193d6c595f37a18e3c06650cf4c74558cc818b16130b293716106f"
)

func TestPrivateKeyToAccount(t *testing.T) {
	account := PrivateKeyToAccount("bf3beef3bd999ba9f2451e06936f0423cd62b815c9233dd3bc90f7e02a1e8673")
	if account != "0xf1424826861ffbbD25405F5145B5E50d0F1bFc90" {
		log.Fatal("计算出错")
	}
}

func TestGetBalance(t *testing.T) {
	ec := NewEthClient(EndPoint, PrivateKey)
	balance := ec.GetBalance("0xD9211042f35968820A3407ac3d80C725f8F75c14")
	fmt.Println(balance)
}

func TestGetMyBalance(t *testing.T) {
	ec := NewEthClient(EndPoint, PrivateKey)
	balance := ec.GetMyBalance()
	fmt.Println(balance)
}

func TestGetEthBalance(t *testing.T) {
	ec := NewEthClient(EndPoint, PrivateKey)
	ether := ec.GetEthBalance("0xD9211042f35968820A3407ac3d80C725f8F75c14")
	fmt.Println(ether)
}

func TestEstimateGas(t *testing.T) {
	ec := NewEthClient(EndPoint, PrivateKey)
	data := PackMethodData("deposit")
	gasLimit := ec.EstimateGas("0xab67F8Bcc56b6a96CCe5e52A7b0A81f9E2317eB9", data, 55555)
	fmt.Println(gasLimit)
}

func TestInvokeContract(t *testing.T) {
	ec := NewEthClient(EndPoint, PrivateKey)
	gasLimit := ec.InvokeContractWithoutArgs("0xbF5e210DdB8723aB184B65af785b75683f81A12C", "deposit", 55555)
	fmt.Println(gasLimit)
}

func TestGetChainId(t *testing.T) {
	ec := NewEthClient(EndPoint, PrivateKey)
	chainId := ec.GetChainId()
	fmt.Println(chainId)
}

func TestGetBlockHeaderByHeight(t *testing.T) {
	ec := NewEthClient("http://192.168.3.173:28545", PrivateKey)
	header := ec.GetBlockHeaderByHeight(0)
	fmt.Println(header)
}

func TestGeTxReceipt(t *testing.T) {
	ec := NewEthClient(EndPoint, PrivateKey)
	txHashHex := "0x32c5ec93deeb29dcd9cd5e1d19cc698f1295d93e0abf48251d29013a028150f8"
	receipt := ec.GetTxReceipt(txHashHex)
	fmt.Println(receipt.Status)
}

func TestGetByteCode(t *testing.T) {
	ec := NewEthClient(EndPoint, PrivateKey)
	contractAddr := "0x5FbDB2315678afecb367f032d93F642f64180aa3"
	//contractAddr := "0xfF4A621a6d8dC31e20100D0D2605332cb06F4e51"
	byteCode, err := ec.Client.CodeAt(context.Background(), common.HexToAddress(contractAddr), nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(byteCode)

}

func TestPackMethodData(t *testing.T) {
	//contractAddr := "0xfF4A621a6d8dC31e20100D0D2605332cb06F4e51"
	data := PackMethodData("isBatchFinalized", [2]string{"uint256", "1"})
	fmt.Println(hex.EncodeToString(data))
}

// 0x01ab27adc5d9822f92c1501595e6bd7f90f3fef6e7f5ae926461752633dca145
