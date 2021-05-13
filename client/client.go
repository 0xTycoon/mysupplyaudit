package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"strings"
)

type Client struct {
	*ethclient.Client
	c *rpc.Client
}

type Config struct {
	RPCURL string `json:"rpc_url"`
}

type ExecutionResult struct {
	Gas         uint64         `json:"gas"`
	Failed      bool           `json:"failed"`
	ReturnValue string         `json:"returnValue"`
	StructLogs  []StructLogRes `json:"structLogs"`
}

// StructLogRes stores a structured log emitted by the EVM while replaying a
// transaction in debug mode
type StructLogRes struct {
	Pc      uint64             `json:"pc"`
	Op      string             `json:"op"`
	Gas     uint64             `json:"gas"`
	GasCost uint64             `json:"gasCost"`
	Depth   int                `json:"depth"`
	Error   error              `json:"error,omitempty"`
	Stack   *[]string          `json:"stack,omitempty"`
	Memory  *[]string          `json:"memory,omitempty"`
	Storage *map[string]string `json:"storage,omitempty"`
}

func (c *Client) TraceTransaction(ctx context.Context, txHash common.Hash) error {
	//var v interface{}
	var t *ExecutionResult
	err := c.c.CallContext(ctx, &t, "debug_traceTransaction", txHash)
	if err != nil {
		fmt.Println(err)
		return err
	}
	//fmt.Println(v, err)
	for _, op := range t.StructLogs {
		if strings.Index(op.Op, "CALL") != -1 {
			// to address of call is 2nd stack value, value 3rd
			fmt.Println(op.Op)
		}
	}

	_ = t
	return nil
	//json.Unmarshal()
	/*

			https://ethereum.stackexchange.com/questions/3417/how-to-get-contract-internal-transactions

			debug.traceTransaction(txHash, {tracer: "callTracer"})

		"result":{
		         "block":"0x37",
		         "hash":"0xdb16f0d4465f2fd79f10ba539b169404a3e026db1be082e7fd6071b4c5f37db7",
		         "traces":[
		            {
		               "from":"0x31b98d14007bdee637298086988a0bbd31184523",
		               "gas":"0x0",
		               "gasUsed":"0x0",
		               "input":"0x",
		               "output":"0x",
		               "time":"1.077µs",
		               "to":"0x2ed530faddb7349c1efdbf4410db2de835a004e4",
		               "type":"CALL",
		               "value":"0xde0b6b3a7640000"
		            }
		         ]
		      }

			var num hexutil.Uint
			err := ec.c.CallContext(ctx, &num, "eth_getBlockTransactionCountByHash", blockHash)
			return uint(num), err
	*/
}

type T struct {
	From     string
	To       string
	Gas      int
	GasPrice *big.Int
	Value    *big.Int
	Data     string
	Nonce    int
}

func (t T) MarshalJSON() ([]byte, error) {
	params := map[string]interface{}{
		"from": t.From,
	}
	if t.To != "" {
		params["to"] = t.To
	}
	if t.Gas > 0 {
		params["gas"] = IntToHex(t.Gas)
	}
	if t.GasPrice != nil {
		params["gasPrice"] = BigToHex(*t.GasPrice)
	}
	if t.Value != nil {
		params["value"] = BigToHex(*t.Value)
	}
	if t.Data != "" {
		params["data"] = t.Data
	}
	if t.Nonce > 0 {
		params["nonce"] = IntToHex(t.Nonce)
	}

	return json.Marshal(params)
}

// BigToHex covert big.Int to hexadecimal representation
func BigToHex(bigInt big.Int) string {
	if bigInt.BitLen() == 0 {
		return "0x0"
	}

	return "0x" + strings.TrimPrefix(fmt.Sprintf("%x", bigInt.Bytes()), "0")
}

// IntToHex convert int to hexadecimal representation
func IntToHex(i int) string {
	return fmt.Sprintf("0x%x", i)
}

func (c *Client) SendTransaction(ctx context.Context, tx *T) (txID string, err error) {
	err = c.c.CallContext(ctx, &txID, "eth_sendTransaction", tx)
	return
}

func NewClient(ec *ethclient.Client, c *rpc.Client) *Client {
	return &Client{ec, c}
}

func Dial(rawurl string) (*Client, error) {
	c, err := rpc.DialContext(context.Background(), rawurl)
	if err != nil {
		return nil, err
	}
	ec := ethclient.NewClient(c)
	return NewClient(ec, c), nil
}
