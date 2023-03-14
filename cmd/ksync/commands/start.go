package commands

import (
	"KYVENetwork/ksync/collector"
	cfg "KYVENetwork/ksync/config"
	"KYVENetwork/ksync/executor"
	"KYVENetwork/ksync/executor/db"
	log "KYVENetwork/ksync/logger"
	"KYVENetwork/ksync/pool"
	"KYVENetwork/ksync/types"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	logger = log.Logger()
)

var (
	home         string
	poolId       int64
	targetHeight int64
)

func init() {
	startCmd.Flags().StringVar(&home, "home", "", "home directory")
	if err := startCmd.MarkFlagRequired("home"); err != nil {
		panic(err)
	}

	startCmd.Flags().Int64Var(&poolId, "pool-id", 0, "pool id")
	if err := startCmd.MarkFlagRequired("pool-id"); err != nil {
		panic(err)
	}

	startCmd.Flags().Int64Var(&targetHeight, "target-height", 0, "target sync height (including)")

	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start fast syncing blocks",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := cfg.LoadConfig(home)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		blockDB, blockStore, err := db.GetBlockstoreDBs(config)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		blockHeight := blockStore.Height()

		if err := blockDB.Close(); err != nil {
			panic(err)
		}

		if targetHeight <= blockHeight {
			logger.Error(fmt.Sprintf("target height %d is not greater than current block height %d", targetHeight, blockHeight))
			os.Exit(1)
		}

		logger.Info(fmt.Sprintf("continuing from block height = %d", blockHeight+1))

		pool.VerifyPool(poolId, blockHeight)

		// process
		// - find out current height from data/ folder
		// - verify pool supports this runtime
		// - verify pool has the min height already archived
		// - find kyve bundle with corresponding height
		// - start downloading bundles from storage provider from that height
		// - apply blocks against blockchain application

		blockCh := make(chan *types.BlockPair)
		quitCh := make(chan int)

		// collector
		go collector.StartBlockCollector(blockCh, quitCh, poolId, blockHeight+1, targetHeight)
		// executor
		go executor.StartBlockExecutor(blockCh, quitCh, home)

		<-quitCh

		logger.Info("done")
	},
}
