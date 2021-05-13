package mysupplyaudit

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"io/ioutil"
	"math"
	"math/big"
	"os"
	"time"

	"github.com/currencytycoon/mysupplyaudit/client"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/params"
	"github.com/sirupsen/logrus"
)

type Config struct {
	client.Config
}

type ether struct {
	*big.Int
}

// String implements the Stringer interface, returns the value in Ether
func (w *ether) String() string {
	f, _, _ := big.ParseFloat(w.Int.String(), 10, 256, big.ToNearestAway)
	ethValue := new(big.Float).Quo(f, big.NewFloat(math.Pow10(18)))
	return ethValue.String() + " ETH"
}

func asETH(x *big.Int) *ether {
	w := new(ether)
	w.Int = x
	return w
}

type Supplier struct {
	config *Config
}

func NewSupplier(path string) (*Supplier, error) {
	s := new(Supplier)
	if err := s.loadConfig(path); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Supplier) loadConfig(path string) error {
	c := &Config{}
	if path == "" {
		path = "./config.json"
	}
	if data, err := ioutil.ReadFile(path); err == nil {
		if err := json.Unmarshal(data, c); err == nil {
			s.config = c
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

// genesisBlock is the 1st block that was mined with all the allocations
// here we sum up all the allocation balances that were mined in the Genesis block
func (s *Supplier) genesisBlock() *big.Int {
	total := big.NewInt(0)
	b := core.DefaultGenesisBlock()
	for _, account := range b.Alloc {
		total = total.Add(total, account.Balance)
	}
	return total
}

func isByzantium(number *big.Int) bool {
	if r := number.Cmp(params.MainnetChainConfig.ByzantiumBlock); r > -1 {
		return true
	}
	return false
}

func isConstantinople(number *big.Int) bool {
	if r := number.Cmp(params.MainnetChainConfig.ConstantinopleBlock); r > -1 {
		return true
	}
	return false
}

func (s *Supplier) calculateBaseReward(n *big.Int) *big.Int {
	blockReward := new(big.Int)
	if isConstantinople(n) {
		blockReward.Set(ethash.ConstantinopleBlockReward)
	} else if isByzantium(n) {
		blockReward.Set(ethash.ByzantiumBlockReward)
	} else {
		blockReward.Set(ethash.FrontierBlockReward)
	}
	return blockReward

}

// Some weird constants to avoid constant memory allocs for them.
var (
	big8  = big.NewInt(8)
	big32 = big.NewInt(32)
)

// Calculates the uncle reward
func calculateUncleReward(uncleNumber *big.Int, blockNumber *big.Int, blockReward *big.Int) *big.Int {
	r := new(big.Int)
	// (U_n + 8 - B_n) * R / 8
	r.Add(uncleNumber, big8)
	r.Sub(r, blockNumber)
	r.Mul(r, blockReward)
	r.Div(r, big8)
	return r
}

func logSyncProgress(p *ethereum.SyncProgress) {
	if p != nil {
		//return errors.New("could not get node progress")
		logrus.WithFields(map[string]interface{}{
			"CurrentBlock":  p.CurrentBlock,
			"HighestBlock":  p.HighestBlock,
			"StartingBlock": p.StartingBlock,
			"KnownStates":   p.KnownStates,
			"PulledStates":  p.PulledStates,
		}).Info("Geth current sync progress")
	}
}

func (s *Supplier) isAlive(rpc *client.Client) bool {
	prompt:
	fmt.Println("Hey Pierre, it looks like your full node is down. Retry? (y/n)")
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if response == "y" {
		err := s.reconnect(rpc)
		if err != nil {
			logrus.WithError(err).Error("Still having RPC connection issues")
			goto prompt
		}
		return true
	}
	return false
}

func (s *Supplier)reconnect(rpc *client.Client) error {
	rpc.Close()
	rpc = nil
	rpc, err := client.Dial(s.config.RPCURL)
	return err
}

// DoAudit completes the audit until a specified block (highestBlock)
// if highestBlock is -1 it will stop when it reaches your node's highest block
func (s *Supplier) DoAudit(highestBlock int64) error {
	rpc, err := client.Dial(s.config.RPCURL)
	if err != nil {
		return err
	}
	defer rpc.Close()
	progress, err := rpc.SyncProgress(context.Background())
	if err != nil {
		return err
	}
	logSyncProgress(progress)

	blockNumber := big.NewInt(0)
	supply := s.genesisBlock()
	logrus.WithFields(map[string]interface{}{
		"blockNumber":  blockNumber.String(),
		"total": supply.String(),
	}).Info("total genesis allocation")
	checkpoint := time.Now().Unix()
	for {
		blockNumber = blockNumber.Add(blockNumber, big.NewInt(1)) // blockNumber++
		if highestBlock!= -1 && blockNumber.Cmp(big.NewInt(highestBlock)) > 0 {
			logrus.WithFields(map[string]interface{}{
				"issuedSupply": asETH(supply).String(),
			}).Info("audit completed - congratulations!")
			break
		}
	fetchBlock:
		block, err := rpc.BlockByNumber(context.Background(), blockNumber)
		if err != nil {
			if s.isAlive(rpc) {
				goto fetchBlock
			} else {
				return err
			}
		}
		header := block.Header()
		// determine the base block reward
		reward := s.calculateBaseReward(blockNumber)
		// count the uncle rewards
		uncles := block.Uncles()
		uncleReward := new(big.Int)
		for _, v := range uncles {
			r := calculateUncleReward(v.Number, header.Number, reward)
			uncleReward.Add(uncleReward, r)
			// 1/32 per Uncle included rewarded to miner
			r.Div(reward, big32)
			reward.Add(reward, r)
		}
		supply.Add(supply, reward)
		logrus.WithFields(map[string]interface{}{
			"blockNumber":  blockNumber.String(),
			"uncleReward":  asETH(uncleReward).String(),
			"blockReward":  asETH(reward).String(),
			"issuedSupply": asETH(supply).String(),
		}).Info("block reward")

	progress:
		if time.Now().Unix() - checkpoint >= 5 {
			// update sync progress every 5 seconds
			progress, err = rpc.SyncProgress(context.Background())
			if err != nil {
				if s.isAlive(rpc) {
					goto progress
				} else {
					return err
				}
			}
			logSyncProgress(progress)
			checkpoint = time.Now().Unix()
		}
	}
	return nil
}
