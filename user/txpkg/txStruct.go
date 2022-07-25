package txpkg

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"reflect"
	"time"

	"cabb/user/blockpkg"
)

type Tx struct {
	TxID      [32]byte
	TimeStamp []byte // 블럭 생성 시간
	Applier   []byte // 신청자
	Company   []byte // 경력회사
	Career    []byte // 경력기간
	Payment   []byte // 결제수단
	Job       []byte // 직종, 업무
	Proof     []byte // 경력증명서 pdf
	WAddr     string // 지갑 주소
}

//TX Hash 데이터 생성
func (tx *Tx) prepareData() []byte {
	data := bytes.Join([][]byte{
		tx.TimeStamp,
		tx.Payment,
		tx.Applier,
		tx.Company,
		tx.Career,
		tx.Job,
		tx.Proof,
		[]byte(tx.WAddr),
	}, []byte{})
	return data
}

//새로운 트랜잭션 생성
func NewTx(applier, company, career, payment, job, proof, wAddr string) *Tx {
	newTx := &Tx{}

	newTx.Applier = []byte(applier)
	newTx.Company = []byte(company)
	newTx.Career = []byte(career)
	newTx.Payment = []byte(payment)
	newTx.Job = []byte(job)
	newTx.Proof = []byte(proof)
	newTx.WAddr = wAddr

	loc, _ := time.LoadLocation("Asia/Seoul")
	now := time.Now()
	t := now.In(loc)
	newTx.TimeStamp = []byte(t.String())

	data := newTx.prepareData()
	newTx.TxID = sha256.Sum256(data)

	return newTx
}

// 트랜잭션 ID를 이용해 Block 조회
func FindBlockByTx(txID [32]byte, bs *blockpkg.Blocks) *blockpkg.Block {
	// 최신부터 돌려보자
	//최신 블록체인의 높이를 구한다
	current_height := len(bs.BlockChain)

	// 최신 블록ID를 찾는다
	curBlockID := [32]byte{}
	for _, v := range bs.BlockChain {
		if v.Height == current_height {
			curBlockID = v.Hash
			break
		}
	}

	for {
		blk := bs.BlockChain[curBlockID]
		if blk.IsExisted(txID) {
			return blk
		} else {
			if reflect.DeepEqual(blk.PrevHash, [32]byte{}) {
				return nil
			}
			curBlockID = blk.PrevHash
		}
	}
}

// 트랜잭션 ID를 이용해 트랜잭션 조회
func FindTxByTxid(txID [32]byte, txs *Txs) *Tx {
	return txs.TxMap[txID]
}

func FindTxByAddr(wAddr string, bs *blockpkg.Blocks, txs *Txs) []*Tx {
	// 최신부터 돌려보자
	//최신 블록체인의 높이를 구한다
	current_height := len(bs.BlockChain)

	// 최신 블록ID를 찾는다
	curBlockID := [32]byte{}
	for _, v := range bs.BlockChain {
		if v.Height == current_height {
			curBlockID = v.Hash
			break
		}
	}

	res := []*Tx{}
	for {
		blk := bs.BlockChain[curBlockID]
		if blk.Height != 1 {
			if txs.TxMap[blk.Txid].WAddr == wAddr {
				res = append(res, txs.TxMap[blk.Txid])
			}
			curBlockID = blk.PrevHash
		} else {
			break
		}
	}
	return res
}

// 트랜잭션 정보 출력
func (t *Tx) PrintTx() {
	fmt.Println("==========Transaction Info=============")
	fmt.Printf("TxId: %x\n Applier: %s\n Company: %s\n Career: %s\n Payment: %s\n TimeStamp: %d\n Job: %s\n Proof: %s\nWallet Address: %s\n\n", t.TxID, t.Applier, t.Company, t.Career, t.Payment, t.TimeStamp, t.Job, t.Proof, t.WAddr)
}
