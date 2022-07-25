package blockpkg

import (
	"fmt"
	"reflect"
	"time"
)

type Block struct {
	Hash      [32]byte //블록 해시
	PrevHash  [32]byte //이전 블록 해시
	PoW       [32]byte //PoW
	Txid      [32]byte //트랜잭션 해시
	Nonce     int      // nonce
	Height    int      // 현재 블록의 인덱스
	Data      []byte   //Copyright 등등..
	Timestamp []byte   //블록 생성 시간
	Sig       []byte   //서명
}

func NewBlock(prevHash [32]byte, height int, txID [32]byte) *Block {
	newBlock := &Block{}
	loc, _ := time.LoadLocation("Asia/Seoul")
	now := time.Now()
	t := now.In(loc)
	newBlock.PrevHash = prevHash
	newBlock.Timestamp = []byte(t.String())
	newBlock.Data = []byte("data")
	newBlock.Txid = txID
	newPoW := NewProofOfWork(newBlock)
	newBlock.Nonce, newBlock.Hash = newPoW.Run()

	//fmt.Printf("%d번째 블록 생성완료: %s\n\n", height, time.Now().String())
	return newBlock
}

func GenesisBlock() *Block {
	newBlock := &Block{}
	newBlock.Height = 1
	loc, _ := time.LoadLocation("Asia/Seoul")
	now := time.Now()
	t := now.In(loc)
	newBlock.Timestamp = []byte(t.String())
	newBlock.Data = []byte("GenesisBlock")
	newPoW := NewProofOfWork(newBlock)
	newBlock.Nonce, newBlock.Hash = newPoW.Run()

	//fmt.Printf("제네시스 블록체인 생성완료: %s\n\n", time.Now().String())
	return newBlock
}
func (b *Block) PrintBlock() {
	fmt.Println("==========블록체인 정보============")
	fmt.Printf("Hash: %x\nHeight: %d\nPrev Hash: %x\nNonce: %d\nPoW: %d\nTimeStamp: %d\nData: %s\nSign: %b\n", b.Hash, b.Height, b.PrevHash, b.Nonce, b.PoW, b.Timestamp, b.Data, b.Sig)
	fmt.Printf("트랜잭션ID:  %d\n", b.Txid)
}
func (b *Block) GetBlockID() [32]byte {
	return b.Hash
}
func (b *Block) GetHeight() int {
	return b.Height
}
func (b *Block) FindTx(txid [32]byte) [32]byte {
	if b.IsExisted(txid) {
		return b.Hash
	} else {
		return [32]byte{}
	}
}

// 1 Tx -> 1 Block
// n Txs -> 1 Block
func (b *Block) IsExisted(txid [32]byte) bool {
	// ToDo
	// n Txs -> 1 Block
	if reflect.DeepEqual(b.Txid, txid) {
		return true
	}
	return false
}
