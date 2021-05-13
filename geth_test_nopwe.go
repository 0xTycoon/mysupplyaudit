package mysupplyaudit

/*



import (
	"fmt"

	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/bloombits"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"

	cli "gopkg.in/urfave/cli.v1"
	"testing"
)


const (
	clientIdentifier = "geth" // Client identifier to advertise over the networkr
)

var (
	gitCommit = ""
	gitDate   = ""
)


func generateTestChain() (*core.Genesis, []*types.Block) {
	db := rawdb.NewMemoryDatabase()
	config := params.AllEthashProtocolChanges
	genesis := &core.Genesis{
		Config:    config,
		Alloc:     core.GenesisAlloc{testAddr: {Balance: testBalance}},
		ExtraData: []byte("test genesis"),
		Timestamp: 9000,
	}
	generate := func(i int, g *core.BlockGen) {
		g.OffsetTime(5)
		g.SetExtra([]byte("test"))
	}
	gblock := genesis.ToBlock(db)
	engine := ethash.NewFaker()
	blocks, _ := core.GenerateChain(config, gblock, engine, db, 1, generate)
	blocks = append([]*types.Block{gblock}, blocks...)
	return genesis, blocks
}

type ethstatsConfig struct {
	URL string `toml:",omitempty"`
}


type gethConfig struct {
	Eth      eth.Config
	Node     node.Config
	Ethstats ethstatsConfig
}


func defaultNodeConfig() node.Config {
	cfg := node.DefaultConfig
	cfg.Name = clientIdentifier
	cfg.Version = params.VersionWithCommit(gitCommit, gitDate)
	cfg.HTTPModules = append(cfg.HTTPModules, "eth")
	cfg.WSModules = append(cfg.WSModules, "eth")
	cfg.IPCPath = "mysupply.ipc"
	cfg.P2P.MaxPeers = 0
	return cfg
}

func makeConfigNode(ctx *cli.Context) (*node.Node, gethConfig) {
	// Load defaults.
	cfg := gethConfig{
		Eth:  eth.DefaultConfig,
		//Shh:  whisper.DefaultConfig,
		Node: defaultNodeConfig(),
	}

	// Load config file.


	// Apply flags.
	utils.SetNodeConfig(ctx, &cfg.Node)
	stack, err := node.New(&cfg.Node)
	if err != nil {
		utils.Fatalf("Failed to create the protocol stack: %v", err)
	}
	utils.SetEthConfig(ctx, stack, &cfg.Eth)
	if ctx.GlobalIsSet(utils.EthStatsURLFlag.Name) {
		cfg.Ethstats.URL = ctx.GlobalString(utils.EthStatsURLFlag.Name)
	}
	//utils.SetShhConfig(ctx, stack, &cfg.Shh)

	return stack, cfg
}

func TestGeth(t *testing.T) {


	// Start Ethereum service.
	ctx := cli.Context{}
	stack, config := makeConfigNode(&ctx)
	chainDb, err := stack.OpenDatabaseWithFreezer("chaindata", config.Eth.DatabaseCache, config.Eth.DatabaseHandles, config.Eth.DatabaseFreezer, "eth/db/chaindata/")
	if err != nil {
		t.Error(err)
		return
	}

	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlock(chainDb, config.Eth.Genesis)
	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {

	if err := stack.Start(); err != nil {
		t.Fatalf("can't start test node: %v", err)
	}
	defer func () {
		err := stack.Close()
		if err != nil {
			t.Error(err)
		}
	}()

		eth := &Ethereum{
			config:            config,
			chainDb:           chainDb,
			eventMux:          stack.EventMux(),
			accountManager:    stack.AccountManager(),
			engine:            CreateConsensusEngine(stack, chainConfig, &config.Ethash, config.Miner.Notify, config.Miner.Noverify, chainDb),
			closeBloomHandler: make(chan struct{}),
			networkID:         config.NetworkId,
			gasPrice:          config.Miner.GasPrice,
			etherbase:         config.Miner.Etherbase,
			bloomRequests:     make(chan chan *bloombits.Retrieval),
			bloomIndexer:      NewBloomIndexer(chainDb, params.BloomBitsBlocks, params.BloomConfirms),
			p2pServer:         stack.Server(),
		}

	var (
		vmConfig = vm.Config{
			EnablePreimageRecording: config.Eth.EnablePreimageRecording,
			EWASMInterpreter:        config.Eth.EWASMInterpreter,
			EVMInterpreter:          config.Eth.EVMInterpreter,
		}
		cacheConfig = &core.CacheConfig{
			TrieCleanLimit:      config.Eth.TrieCleanCache,
			TrieCleanJournal:    stack.ResolvePath(config.Eth.TrieCleanCacheJournal),
			TrieCleanRejournal:  config.Eth.TrieCleanCacheRejournal,
			TrieCleanNoPrefetch: config.Eth.NoPrefetch,
			TrieDirtyLimit:      config.Eth.TrieDirtyCache,
			TrieDirtyDisabled:   config.Eth.NoPruning,
			TrieTimeLimit:       config.Eth.TrieTimeout,
			SnapshotLimit:       config.Eth.SnapshotCache,
		}
	)
	stack.E
	var bc *core.BlockChain
	bc, err = core.NewBlockChain(chainDb, cacheConfig, chainConfig, eth.engine, vmConfig, eth.shouldPreserve, &config.Eth.TxLookupLimit)
	if err != nil {
		return nil, err
	}


	block := bc.GetBlockByNumber(40)
	if block != nil {
		fmt.Println("here")
	}

}



 */
