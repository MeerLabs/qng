package test

import (
	"fmt"
	"github.com/Qitmeer/qng/common/hash"
	"github.com/Qitmeer/qng/common/system"
	"github.com/Qitmeer/qng/database"
	"github.com/Qitmeer/qng/meerdag"
	"github.com/Qitmeer/qng/services/common"
	"os"
	"strconv"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
	exit()
}

func Test_GetFutureSet(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig2-blocks")
	if ibd == nil {
		t.FailNow()
	}

	//ph:=ibd.(*Phantom)
	anBlock := tbMap[testData.PH_GetFutureSet.Input]
	bset := meerdag.NewIdSet()
	bd.GetFutureSet(bset, anBlock)
	fmt.Printf("Get %s future set：\n", testData.PH_GetFutureSet.Input)
	printBlockSetTag(bset)
	//
	if !processResult(bset, changeToIDList(testData.PH_GetFutureSet.Output)) {
		t.FailNow()
	}
}

func Test_GetAnticone(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig2-blocks")
	if ibd == nil {
		t.FailNow()
	}
	//
	anBlock := tbMap[testData.PH_GetAnticone.Input]

	////////////
	bset := bd.GetAnticone(anBlock, nil)
	fmt.Printf("Get %s anticone set：\n", testData.PH_GetAnticone.Input)
	printBlockSetTag(bset)
	//
	if !processResult(bset, changeToIDList(testData.PH_GetAnticone.Output)) {
		t.FailNow()
	}

}

func Test_BlueSetFig1(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig1-blocks")
	if ibd == nil {
		t.FailNow()
	}
	ph := ibd.(*meerdag.Phantom)
	//
	blueSet := ph.GetDiffBlueSet()
	fmt.Println("Fig1 diff blue set：")
	printBlockSetTag(blueSet)
	if !processResult(blueSet, changeToIDList(testData.PH_BlueSetFig1.Output)) {
		t.FailNow()
	}
}

func Test_BlueSetFig2(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig2-blocks")
	if ibd == nil {
		t.FailNow()
	}
	ph := ibd.(*meerdag.Phantom)
	//
	blueSet := ph.GetDiffBlueSet()
	fmt.Println("Fig2 diff blue set：")
	printBlockSetTag(blueSet)
	if !processResult(blueSet, changeToIDList(testData.PH_BlueSetFig2.Output)) {
		t.FailNow()
	}
}

func Test_BlueSetFig4(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig4-blocks")
	if ibd == nil {
		t.FailNow()
	}
	ph := ibd.(*meerdag.Phantom)
	//
	blueSet := ph.GetDiffBlueSet()
	fmt.Println("Fig4 diff blue set：")
	printBlockSetTag(blueSet)
	if !processResult(blueSet, changeToIDList(testData.PH_BlueSetFig4.Output)) {
		t.FailNow()
	}
}

func Test_OrderFig1(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig1-blocks")
	if ibd == nil {
		t.FailNow()
	}
	ph := ibd.(*meerdag.Phantom)
	order := []uint{}
	var i uint
	ph.UpdateVirtualBlockOrder()
	for i = 0; i < bd.GetBlockTotal(); i++ {
		order = append(order, bd.GetBlockByOrder(uint(i)).GetID())
	}
	fmt.Printf("The Fig.1 Order: ")
	printBlockChainTag(order)

	if !processResult(order, changeToIDList(testData.PH_OrderFig1.Output)) {
		t.FailNow()
	}

	//
	da := ph.GetDiffAnticone()
	fmt.Printf("The diffanticoner: ")
	printBlockSetTag(da)
}

func Test_OrderFig2(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig2-blocks")
	if ibd == nil {
		t.FailNow()
	}
	ph := ibd.(*meerdag.Phantom)
	order := []uint{}
	var i uint
	ph.UpdateVirtualBlockOrder()
	for i = 0; i < bd.GetBlockTotal(); i++ {
		order = append(order, bd.GetBlockByOrder(uint(i)).GetID())
	}
	fmt.Printf("The Fig.2 Order: ")
	printBlockChainTag(order)

	if !processResult(order, changeToIDList(testData.PH_OrderFig2.Output)) {
		t.FailNow()
	}

	//
	da := ph.GetDiffAnticone()
	fmt.Printf("The diffanticoner: ")
	printBlockSetTag(da)
}

func Test_OrderFig4(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig4-blocks")
	if ibd == nil {
		t.FailNow()
	}
	ph := ibd.(*meerdag.Phantom)
	order := []uint{}
	var i uint
	ph.UpdateVirtualBlockOrder()
	for i = 0; i < bd.GetBlockTotal(); i++ {
		order = append(order, bd.GetBlockByOrder(uint(i)).GetID())
	}
	fmt.Printf("The Fig.4 Order: ")
	printBlockChainTag(order)

	if !processResult(order, changeToIDList(testData.PH_OrderFig4.Output)) {
		t.FailNow()
	}

	//
	da := ph.GetDiffAnticone()
	fmt.Printf("The diffanticoner: ")
	printBlockSetTag(da)
}

func Test_GetLayer(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig2-blocks")
	if ibd == nil {
		t.FailNow()
	}
	var result string = ""
	var i uint
	ph := ibd.(*meerdag.Phantom)
	ph.UpdateVirtualBlockOrder()
	for i = 0; i < bd.GetBlockTotal(); i++ {
		l := bd.GetLayer(bd.GetBlockByOrder(uint(i)).GetID())
		result = fmt.Sprintf("%s%d", result, l)
	}
	if result != testData.PH_GetLayer.Output[0] {
		t.FailNow()
	}
}

func Test_IsOnMainChain(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig2-blocks")
	if ibd == nil {
		t.FailNow()
	}
	if strconv.FormatBool(bd.IsOnMainChain(tbMap[testData.PH_IsOnMainChain.Input].GetID())) != testData.PH_IsOnMainChain.Output[0] {
		t.FailNow()
	}
}

func Test_LocateBlocks(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig2-blocks")
	if ibd == nil {
		t.FailNow()
	}
	gs := meerdag.NewGraphState()
	gs.SetTips([]*hash.Hash{bd.GetGenesisHash()})
	gs.SetTotal(1)
	gs.SetLayer(0)
	lb := bd.LocateBlocks(gs, 100)

	lbids := meerdag.NewIdSet()
	for _, v := range lb {
		lbids.Add(bd.GetBlockId(v))
	}
	if !processResult(lbids, changeToIDList(testData.PH_LocateBlocks.Output)) {
		t.FailNow()
	}
}

func Test_LocateMaxBlocks(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig2-blocks")
	if ibd == nil {
		t.FailNow()
	}
	gs := meerdag.NewGraphState()
	gs.SetTips([]*hash.Hash{bd.GetGenesisHash(), tbMap["G"].GetHash()})
	gs.SetTotal(4)
	gs.SetLayer(2)
	lb := bd.LocateBlocks(gs, 4)
	lbids := meerdag.NewIdSet()
	for _, v := range lb {
		lbids.Add(bd.GetBlockId(v))
	}
	if !processResult(lbids, changeToIDList(testData.PH_LocateMaxBlocks.Output)) {
		t.FailNow()
	}
}

func Test_Confirmations(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig2-blocks")
	if ibd == nil {
		t.FailNow()
	}
	mainTip := bd.GetMainChainTip()
	mainChain := []uint{}
	for cur := mainTip; cur != nil; cur = bd.GetBlockById(cur.GetMainParent()) {
		mainChain = append(mainChain, cur.GetID())
	}
	printBlockChainTag(reverseBlockList(mainChain))

	ph := ibd.(*meerdag.Phantom)
	ph.UpdateVirtualBlockOrder()
	for i := uint(0); i < bd.GetBlockTotal(); i++ {
		blockHash := bd.GetBlockByOrder(uint(i)).GetID()
		fmt.Printf("%s : %d\n", getBlockTag(blockHash), bd.GetConfirmations(blockHash))
	}
}

func Test_IsDAG(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig2-blocks")
	if ibd == nil {
		t.FailNow()
	}
	//ph:=ibd.(*Phantom)
	//
	parentsTag := []string{"I", "G"}
	parents := []*hash.Hash{}
	for _, parent := range parentsTag {
		parents = append(parents, tbMap[parent].GetHash())
	}
	_, err := buildBlock("L", parents)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_IsHourglass(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "CP_Blocks")
	if ibd == nil {
		t.FailNow()
	}
	if !bd.IsHourglass(tbMap["J"].GetID()) {
		t.Fatal()
	}
}

func Test_GetMaturity(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig2-blocks")
	if ibd == nil {
		t.FailNow()
	}
	if bd.GetMaturity(tbMap["D"].GetID(), []uint{tbMap["I"].GetID()}) != 2 {
		t.Fatal()
	}
}

func Test_GetMainParentConcurrency(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig2-blocks")
	if ibd == nil {
		t.FailNow()
	}

	//ph:=ibd.(*Phantom)
	anBlock := bd.GetBlock(tbMap[testData.PH_MPConcurrency.Input].GetHash())
	//fmt.Println(bd.GetMainParentConcurrency(anBlock))
	if bd.GetMainParentConcurrency(anBlock) != testData.PH_MPConcurrency.Output {
		t.Fatal()
	}
}

func Test_GetBlockConcurrency(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig2-blocks")
	if ibd == nil {
		t.FailNow()
	}

	//ph:=ibd.(*Phantom)
	blueNum, err := bd.GetBlockConcurrency(tbMap[testData.PH_MPConcurrency.Input].GetHash())
	if err != nil {
		t.Fatal(err)
	}
	if blueNum != uint(testData.PH_BConcurrency.Output) {
		t.Fatal()
	}
}

func Test_MainChainTip(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig2-blocks")
	if ibd == nil {
		t.FailNow()
	}
	ph := ibd.(*meerdag.Phantom)
	ph.UpdateVirtualBlockOrder()
	for _, v := range testData.PH_MainChainTip {
		err := bd.CheckSubMainChainTip(getBlocksByTag(v.Input))
		if err != nil {
			t.Log(err)
		}
	}
}

func Test_Rollback(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig2-blocks")
	if ibd == nil {
		t.FailNow()
	}
	orders := meerdag.NewIdSet()
	total := bd.GetBlockTotal()
	tips := bd.GetTipsSet().Clone()

	for i := uint(0); i < bd.GetBlockTotal(); i++ {
		ib := bd.GetBlockById(i)
		orders.AddPair(ib.GetID(), ib.GetOrder())
	}

	parents := []*hash.Hash{}
	parents = append(parents, tbMap["I"].GetHash())
	parents = append(parents, tbMap["G"].GetHash())

	_, _, err := addBlock("L", parents)
	if err != nil {
		t.Fatal(err)
	}

	bd.Rollback()

	if bd.GetBlockTotal() != total {
		t.Fatalf("Roll back error")
	}
	for i := uint(0); i < bd.GetBlockTotal(); i++ {
		ib := bd.GetBlockById(i)
		v := orders.Get(i)
		o, ok := v.(uint)
		if !ok {
			t.Fatalf("Roll back error")
		}
		if o != ib.GetOrder() {
			t.Fatalf("Roll back error")
		}
	}

	if !bd.GetTipsSet().IsEqual(tips) {
		t.Fatalf("Roll back error")
	}
}

func Test_tips(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig2-blocks")
	if ibd == nil {
		t.FailNow()
	}

	//ph := ibd.(*Phantom)
	bd.SetTipsDisLimit(1)

	parents := []*hash.Hash{}
	parents = append(parents, tbMap["J"].GetHash())

	_, err := buildBlock("L", parents)
	if err != nil {
		t.Fatal(err)
	}

	parents = []*hash.Hash{}
	parents = append(parents, tbMap["L"].GetHash())

	_, err = buildBlock("M", parents)
	if err != nil {
		t.Fatal(err)
	}

	parents = []*hash.Hash{}
	parents = append(parents, tbMap["M"].GetHash())

	_, err = buildBlock("N", parents)
	if err != nil {
		t.Fatal(err)
	}

	bd.DB().Close()

	checkLoad(t)
}

func checkLoad(t *testing.T) {
	cfg := common.DefaultConfig(os.TempDir())
	cfg.DevNextGDB = false
	db, err := database.New(cfg, system.InterruptListener())
	if err != nil {
		t.Fatal(err)
	}

	getBlockData := func(h *hash.Hash) meerdag.IBlockData {
		tb, err := fetchBlock(h)
		if err != nil {
			t.Fatal(err)
		}
		return tb
	}
	bd = meerdag.New(meerdag.PHANTOM, -1, db, getBlockData)
	total, err := dbGetTotal()
	if err != nil {
		t.Fatal(err)
	}
	geneis, err := dbGetGenesis()
	if err != nil {
		t.Fatal(err)
	}

	err = bd.Load(uint(total), geneis)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ForeachFig1(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig1-blocks")
	if ibd == nil {
		t.FailNow()
	}
	ph := ibd.(*meerdag.Phantom)
	order := []uint{bd.GetMainChainTip().GetID()}

	ph.UpdateVirtualBlockOrder()

	err := bd.Foreach(bd.GetMainChainTip(), meerdag.MaxId, meerdag.All, func(block meerdag.IBlock) (bool, error) {
		t.Logf("block id:%d hash:%s order:%d", block.GetID(), block.GetHash().String(), block.GetOrder())
		order = append(order, block.GetID())
		return true, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	order = reverseBlockList(order)
	target := changeToIDList(testData.PH_OrderFig1.Output)
	for i := 0; i < len(order); i++ {
		if order[i] != target[i] {
			t.FailNow()
		}
	}
	fmt.Printf("The Fig.1 Order from mainTip: ")
	printBlockChainTag(order)
}

func Test_ForeachFig2(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig2-blocks")
	if ibd == nil {
		t.FailNow()
	}
	ph := ibd.(*meerdag.Phantom)
	order := []uint{bd.GetMainChainTip().GetID()}

	ph.UpdateVirtualBlockOrder()

	err := bd.Foreach(bd.GetMainChainTip(), meerdag.MaxId, meerdag.All, func(block meerdag.IBlock) (bool, error) {
		t.Logf("block id:%d hash:%s order:%d", block.GetID(), block.GetHash().String(), block.GetOrder())
		order = append(order, block.GetID())
		return true, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	order = reverseBlockList(order)
	target := changeToIDList(testData.PH_OrderFig2.Output)
	for i := 0; i < len(order); i++ {
		if order[i] != target[i] {
			t.FailNow()
		}
	}
	fmt.Printf("The Fig.2 Order from mainTip: ")
	printBlockChainTag(order)
}

func Test_ForeachFig4(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig4-blocks")
	if ibd == nil {
		t.FailNow()
	}
	ph := ibd.(*meerdag.Phantom)
	order := []uint{bd.GetMainChainTip().GetID()}

	ph.UpdateVirtualBlockOrder()

	err := bd.Foreach(bd.GetMainChainTip(), meerdag.MaxId, meerdag.All, func(block meerdag.IBlock) (bool, error) {
		t.Logf("block id:%d hash:%s order:%d", block.GetID(), block.GetHash().String(), block.GetOrder())
		order = append(order, block.GetID())
		return true, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	order = reverseBlockList(order)
	target := changeToIDList(testData.PH_OrderFig4.Output)
	for i := 0; i < len(order); i++ {
		if order[i] != target[i] {
			t.FailNow()
		}
	}
	fmt.Printf("The Fig.1 Order from mainTip: ")
	printBlockChainTag(order)
}

func Test_ForeachDepth(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig2-blocks")
	if ibd == nil {
		t.FailNow()
	}

	mt := bd.GetMainChainTip()
	for i := uint(0); i <= mt.GetOrder(); i++ {
		count := uint(0)
		err := bd.Foreach(mt, i, meerdag.All, func(block meerdag.IBlock) (bool, error) {
			//t.Logf("depth:%d,block id:%d hash:%s order:%d", i, block.GetID(), block.GetHash().String(), block.GetOrder())
			count++
			return true, nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if count != i {
			t.Fatalf("expect:%d != %d", i, count)
		}
		//t.Log("-------------")
	}

}

func Test_LastBlock(t *testing.T) {
	ibd := InitBlockDAG(meerdag.PHANTOM, "PH_fig2-blocks")
	if ibd == nil {
		t.FailNow()
	}
	lastBlock := bd.GetLastBlock()
	lastBlockID := bd.GetLastBlockID()
	if lastBlock.GetID() != lastBlockID {
		t.FailNow()
	}
	if !tbMap["K"].GetHash().IsEqual(lastBlock.GetHash()) {
		t.FailNow()
	}
}
