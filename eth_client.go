package service

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"strings"
	"time"
)

func BalanceToEther(balance *big.Int) *big.Float {
	balanceInEther := new(big.Float).SetInt(balance)
	etherValue := new(big.Float).Quo(balanceInEther, big.NewFloat(1e18))
	return etherValue
}

func PackMethodData(method string, args ...[2]string) []byte {

	//var args [][]string
	//arg1 := []string{"address", "0x17435ccE3d1B4fA2e5f8A08eD921D57C6762A180"}
	//arg2 := []string{"uint256", "3333"}
	//args = append(args, arg1)
	//args = append(args, arg2)

	//fmt.Println(args) // TODO

	methodSigArgs := ""
	if len(args) != 0 {
		var argTypes []string
		var argValues []string
		for _, arg := range args {
			argTypes = append(argTypes, arg[0])
			argValues = append(argValues, arg[1])
		}
		methodSigArgs = strings.Join(argTypes[:], ",")
	}

	methodSig := []byte(fmt.Sprintf("%s(%s)", method, methodSigArgs))

	methodId := crypto.Keccak256Hash(methodSig).Hex()[2:10] // 去掉前面的0x
	data, err := hex.DecodeString(methodId)
	if err != nil {
		log.Fatalf("解析MethodId %s 失败: %s", methodId, err)
	}
	return data
}

func PrivateKeyToAccount(privateKey string) string {
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatalf("解析用户私钥 %s 失败: %s\n", privateKey, err)
	}
	address := crypto.PubkeyToAddress(key.PublicKey)
	return address.Hex()
}

type EthClient struct {
	PrivateKey *ecdsa.PrivateKey // 用户私钥对象
	Address    common.Address    // 用户账户地址
	Client     *ethclient.Client // 以太坊客户端
}

func NewEthClient(endpoint, privateKey string) *EthClient {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		log.Fatalf("连接以太坊 %s 失败: %s\n", endpoint, err)
	}
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatalf("解析用户私钥 %s 失败: %s\n", privateKey, err)
	}
	address := crypto.PubkeyToAddress(key.PublicKey)

	return &EthClient{
		key,
		address,
		client,
	}
}

func (ec *EthClient) Account() string {
	return ec.Address.Hex()
}

func (ec *EthClient) GetCurrentBlockHeight() uint64 {
	height, err := ec.Client.BlockNumber(context.Background())
	if err != nil {
		log.Fatalf("获取当前区块高度失败: %s", err)
	}
	return height
}

func (ec *EthClient) GetBlockByHeight(blockHeight uint64) *types.Block {
	block, err := ec.Client.BlockByNumber(context.Background(), big.NewInt(int64(blockHeight)))
	if err != nil {
		log.Fatalf("根据区块高度 %d 获取区块失败: %s", blockHeight, err)
	}
	return block
}

func (ec *EthClient) GetBlockByHash(blockHash string) *types.Block {
	block, err := ec.Client.BlockByHash(context.Background(), common.HexToHash(blockHash))
	if err != nil {
		log.Fatalf("根据区块哈希 %s 获取区块失败: %s", blockHash, err)
	}
	return block
}

func (ec *EthClient) GetBlockHeaderByHeight(blockHeight uint64) *types.Header {
	header, err := ec.Client.HeaderByNumber(context.Background(), big.NewInt(int64(blockHeight)))
	if err != nil {
		log.Fatalf("根据区块高度 %d 获取区块头失败: %s", blockHeight, err)
	}
	return header
}

func (ec *EthClient) GetBlockHeaderByHash(blockHash string) *types.Header {
	header, err := ec.Client.HeaderByHash(context.Background(), common.HexToHash(blockHash))
	if err != nil {
		log.Fatalf("根据区块哈希 %s 获取区块头失败: %s", blockHash, err)
	}
	return header
}

func (ec *EthClient) GetNonce() uint64 {
	nonce, err := ec.Client.PendingNonceAt(context.Background(), ec.Address)
	if err != nil {
		log.Fatalf("获取账户 %s nonce失败: %s\n", ec.Address, err)
	}
	return nonce
}

func (ec *EthClient) GetGasPrice() *big.Int {
	gasPrice, err := ec.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("获取GasPrice失败: %s\n", err)
	}
	return gasPrice
}

func (ec *EthClient) EstimateGas(toAccount string, data []byte, value int64) uint64 {
	toAddress := common.HexToAddress(toAccount)
	msg := ethereum.CallMsg{
		From:  ec.Address,
		To:    &toAddress,
		Data:  data,
		Value: big.NewInt(value),
	}
	gasLimit, err := ec.Client.EstimateGas(context.Background(), msg)
	if err != nil {
		log.Fatalf("预估Gas失败: %s\n", err)
	}
	return gasLimit
}

func (ec *EthClient) GetChainId() *big.Int {
	chainId, err := ec.Client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("获取链NetworkID失败: %s\n", err)
	}
	return chainId
}

func (ec *EthClient) GetBalance(account string) *big.Int {
	address := common.HexToAddress(account)
	balance, err := ec.Client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatalf("查询账户 %s 余额失败: %s\n", account, err)
	}
	return balance
}

func (ec *EthClient) GetMyBalance() *big.Int {
	return ec.GetBalance(ec.Account())
}

func (ec *EthClient) GetEthBalance(account string) *big.Float {
	balance := ec.GetBalance(account)
	etherValue := BalanceToEther(balance)
	return etherValue
}

func (ec *EthClient) GetTxByHash(txHash string) *types.Transaction {
	tx, isPending, err := ec.Client.TransactionByHash(context.Background(), common.HexToHash(txHash))
	if err != nil {
		log.Fatalf("根据交易哈希 %s 获取交易失败: %s", txHash, err)
	}
	if isPending == true {
		fmt.Println("交易执行中")
	} else {
		fmt.Println("交易执行完毕")
	}
	return tx
}

func (ec *EthClient) GetTxInBlock(blockHash string, index uint) *types.Transaction {
	tx, err := ec.Client.TransactionInBlock(context.Background(), common.HexToHash(blockHash), index)
	if err != nil {
		log.Fatalf("根据区块哈希 %s 获取区块失败: %s", blockHash, err)
	}
	return tx
}

func (ec *EthClient) GetByteCode(contractAddr string) string {
	byteCode, err := ec.Client.CodeAt(context.Background(), common.HexToAddress(contractAddr), nil)
	if err != nil {
		log.Fatalf("根据合约地址 %s 获取合约ByteCode失败: %s", contractAddr, err)
	}
	return hex.EncodeToString(byteCode)
}

func (ec *EthClient) GetTxReceipt(txHash string) *types.Receipt {
	_txHash := common.HexToHash(txHash)
	ctx := context.Background()
	for i := 0; i <= 30; i++ {
		_, isPending, err := ec.Client.TransactionByHash(ctx, _txHash)
		if err != nil {

			log.Fatalf("通过哈希获取交易失败: %s\n", err)
		}
		if isPending == false {
			break
		}
		time.Sleep(1 * time.Second)
	}
	receipt, err := ec.Client.TransactionReceipt(ctx, _txHash)
	if err != nil {

		log.Fatalf("通过哈希获取交易回执失败: %s\n", err)
	}
	return receipt
}

func (ec *EthClient) CreateTx(toAccount string, data []byte, value int64) *types.Transaction {
	nonce := ec.GetNonce()
	gasPrice := ec.GetGasPrice()
	gasLimit := ec.EstimateGas(toAccount, data, value)
	toAddress := common.HexToAddress(toAccount)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    big.NewInt(value),
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})
	return tx
}

func (ec *EthClient) SignTx(tx *types.Transaction) *types.Transaction {
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(ec.GetChainId()), ec.PrivateKey)
	if err != nil {

		log.Fatalf("签名交易失败: %s\n", err)
	}
	return signedTx
}

func (ec *EthClient) SendTx(tx *types.Transaction) *types.Receipt {
	signedTx := ec.SignTx(tx)
	txHash := signedTx.Hash().Hex()
	fmt.Printf("已签名交易Hash: %s\n", txHash)
	err := ec.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("发送调用合约交易失败: %s\n", err)
	}
	receipt := ec.GetTxReceipt(txHash)
	if receipt.Status == 0 {
		fmt.Printf("交易执行失败: %s\n", txHash)
	} else if receipt.Status == 1 {
		fmt.Printf("交易执行成功: %s\n", txHash)
	}
	return receipt
}

func (ec *EthClient) InvokeContractWithoutArgs(contractAddr, method string, value int64) *types.Receipt {
	data := PackMethodData(method)
	tx := ec.CreateTx(contractAddr, data, value)
	receipt := ec.SendTx(tx)
	return receipt
}

func (ec *EthClient) Stop() {
	ec.Client.Close()
}
