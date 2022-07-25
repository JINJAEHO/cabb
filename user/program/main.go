package main

import (
	"cabb/user/blockpkg"
	"cabb/user/txpkg"
	"cabb/user/walletpkg"
	"fmt"
)

var prevH [32]byte // 최근 블록 해시 관리

func main() {
	//-----------wallet---------------------------
	ws := walletpkg.CreateWallets()
	var address1 string
	var address2 string
	for i := 0; i < 2; i++ {
		alias := "wallet_" + fmt.Sprint(i)
		w := walletpkg.NewWallet(alias)
		ws.SaveWallet(w)
		if i == 0 {
			address1 = w.Address
		} else {
			address2 = w.Address
		}
		w.WalletPrint()
	}
	//제네시스 블록 생성
	gBlock := blockpkg.GenesisBlock()
	//제네시스 해시 값 저장
	prevH = gBlock.Hash
	//제네시스 블록으로 블록체인 생성
	chain := blockpkg.NewBlockchain(gBlock)
	//트랙잭션 맵(DB대용) 생성
	txs := txpkg.CreateTxDB()
	n := &blockpkg.Block{}
	tx := &txpkg.Tx{}
	for i := 0; i < 3; i++ {
		//트랜잭션 생성
		if i < 2 {
			tx = txpkg.NewTx("user1", "company", "1년", "카드", "백엔드", "sign.png", address1)
		} else {
			tx = txpkg.NewTx("user2", "company", "1년", "카드", "백엔드", "sign.png", address2)
		}
		//트랜잭션 저장
		txs.AddTx(tx)
		//블록 생성
		n = blockpkg.NewBlock(prevH, len(chain.BlockChain)+1, tx.TxID) // 이전 블록 해시, 블록체인 크기, 트랜잭션 해시
		//해쉬 관리
		prevH = n.Hash
		//블록체인에 연결
		chain.AddBlock(n)
		txpkg.FindTxByTxid(tx.TxID, txs).PrintTx()
	}
	fmt.Print("=====================================================================\n")
	list := ws.FindWallet(address1).FindAllTx(chain, txs)
	//fmt.Printf("%d", len(list))
	for i := 0; i < len(list); i++ {
		list[i].PrintTx()
	}
}
