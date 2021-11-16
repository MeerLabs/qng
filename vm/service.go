package vm

import (
	"encoding/hex"
	"fmt"
	"github.com/Qitmeer/qng/config"
	"github.com/Qitmeer/qng/consensus"
	"github.com/Qitmeer/qng/core/address"
	"github.com/Qitmeer/qng/core/blockchain"
	"github.com/Qitmeer/qng/core/blockchain/opreturn"
	"github.com/Qitmeer/qng/core/event"
	"github.com/Qitmeer/qng/core/types"
	"github.com/Qitmeer/qng/engine/txscript"
	"github.com/Qitmeer/qng/node/service"
	"github.com/Qitmeer/qng/params"
	"github.com/Qitmeer/qng/vm/chainvm"

	"io/ioutil"
	"path/filepath"
)

// meerevm ID of the platform
const (
	MeerEVMID = "meerevm"
)

type Factory interface {
	New() (consensus.ChainVM, error)
	GetVM() consensus.ChainVM
	Context() *consensus.Context
}

type Service struct {
	service.Service

	events *event.Feed

	factories map[string]Factory

	versions map[string]string

	cfg *config.Config
}

func (s *Service) Start() error {
	log.Info("Starting Virtual Machines Service")
	if err := s.Service.Start(); err != nil {
		return err
	}
	vm, err := s.GetFactory(MeerEVMID)
	if err != nil {
		log.Debug(fmt.Sprintf("no %s", MeerEVMID))
	} else {
		err := vm.GetVM().Initialize(vm.Context())
		if err != nil {
			log.Warn(err.Error())
		} else {
			err := vm.GetVM().Bootstrapping()
			if err != nil {
				log.Warn(err.Error())
			} else {
				err := vm.GetVM().Bootstrapped()
				if err != nil {
					log.Warn(err.Error())
				}
			}
		}
	}
	s.subscribe()
	return nil
}

func (s *Service) Stop() error {
	log.Info("Stopping Virtual Machines Service")
	if err := s.Service.Stop(); err != nil {
		return err
	}
	vm, err := s.GetFactory(MeerEVMID)
	if err == nil {
		vm.GetVM().Shutdown()
	}
	return nil
}

func (s *Service) GetFactory(id string) (Factory, error) {
	f, ok := s.factories[id]
	if !ok {
		return nil, fmt.Errorf("No factory:%s", id)
	}
	return f, nil
}

func (s *Service) HasFactory(id string) bool {
	f, err := s.GetFactory(id)
	return err == nil && f != nil
}

func (s *Service) RegisterFactory(vmID string, factory Factory) error {
	if s.HasFactory(vmID) {
		return fmt.Errorf(fmt.Sprintf("Already exists:%s", vmID))
	}

	s.factories[vmID] = factory

	log.Debug(fmt.Sprintf("Adding factory for vm %s", vmID))

	vm, err := factory.New()
	if err != nil {
		return err
	}

	commonVM, ok := vm.(consensus.VM)
	if !ok {
		return nil
	}

	version, err := commonVM.Version()
	if err != nil {
		log.Error(fmt.Sprintf("fetching version for %q errored with: %s", vmID, err))

		if err := commonVM.Shutdown(); err != nil {
			return fmt.Errorf("shutting down VM errored with: %s", err)
		}
		return nil
	}
	s.versions[vmID] = version
	return nil
}

func (s *Service) Versions() (map[string]string, error) {
	return s.versions, nil
}

func (s *Service) registerVMs() error {
	if len(s.cfg.PluginDir) <= 0 {
		return nil
	}
	files, err := ioutil.ReadDir(s.cfg.PluginDir)
	if err != nil {
		return err
	}
	log.Debug(fmt.Sprintf("Register Virtual Machines from:%s num:%d", s.cfg.PluginDir, len(files)))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		// Strip any extension from the file. This is to support windows .exe
		// files.
		name = name[:len(name)-len(filepath.Ext(name))]
		// Skip hidden files.
		if len(name) == 0 {
			continue
		}

		if s.HasFactory(name) {
			continue
		}

		if err = s.RegisterFactory(name, &chainvm.Factory{
			Path: filepath.Join(s.cfg.PluginDir, file.Name()),
			Ctx: &consensus.Context{Context: s.Context(),
				Datadir: s.cfg.DataDir, LogLevel: s.cfg.DebugLevel, NetworkID: params.ActiveNetParams.Net, LogLocate: s.cfg.DebugPrintOrigins},
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) subscribe() {
	ch := make(chan *event.Event)
	sub := s.events.Subscribe(ch)
	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case ev := <-ch:
				if ev.Data != nil {
					switch value := ev.Data.(type) {
					case *blockchain.Notification:
						s.handleNotifyMsg(value)
					}
				}
				if ev.Ack != nil {
					ev.Ack <- struct{}{}
				}
			}

			if s.IsShutdown() {
				log.Info("Close Miner Event Subscribe")
				return
			}
		}
	}()
}

func (s *Service) handleNotifyMsg(notification *blockchain.Notification) {
	switch notification.Type {
	case blockchain.BlockAccepted:
		ban, ok := notification.Data.(*blockchain.BlockAcceptedNotifyData)
		if !ok {
			return
		}
		vm, err := s.GetFactory(MeerEVMID)
		if err == nil {
			txs := []*consensus.Tx{}
			for _, tx := range ban.Block.Transactions() {
				if types.IsCrossChainExportTx(tx.Tx) {
					ctx := &consensus.Tx{Type: consensus.TxTypeExport}
					_, pksAddrs, _, err := txscript.ExtractPkScriptAddrs(tx.Tx.TxOut[0].PkScript, params.ActiveNetParams.Params)
					if err != nil {
						log.Error(err.Error())
						return
					}

					if len(pksAddrs) > 0 {
						secpPksAddr, ok := pksAddrs[0].(*address.SecpPubKeyAddress)
						if !ok {
							log.Error(fmt.Sprintf("Not SecpPubKeyAddress:%s", pksAddrs[0].String()))
							return
						}
						ctx.To = hex.EncodeToString(secpPksAddr.PubKey().SerializeUncompressed())
						ctx.Value = uint64(tx.Tx.TxOut[0].Amount.Value)
						txs = append(txs, ctx)
					}

				} else {
					for _, out := range tx.Tx.TxOut {
						if !opreturn.IsMeerEVM(out.PkScript) {
							continue
						}
						me, err := opreturn.NewOPReturnFrom(out.PkScript)
						if err != nil {
							log.Error(err.Error())
							continue
						}
						ctx := &consensus.Tx{Type: consensus.TxTypeNormal, Data: []byte(me.(*opreturn.MeerEVM).GetHex())}
						txs = append(txs, ctx)
					}
				}
			}
			if len(txs) <= 0 {
				return
			}
			_, err := vm.GetVM().BuildBlock(txs)
			if err != nil {
				log.Warn(err.Error())
			}
		}
	}
}

func NewService(cfg *config.Config, events *event.Feed) (*Service, error) {
	ser := Service{
		events:    events,
		factories: make(map[string]Factory),
		versions:  make(map[string]string),
		cfg:       cfg,
	}
	ser.InitContext()

	if err := ser.registerVMs(); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &ser, nil
}
