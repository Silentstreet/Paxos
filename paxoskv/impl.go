/*
 * @Author: Peng.cao
 * @Date: 2021-04-30 12:01:16
 * @LastEditTime: 2021-04-30 18:27:42
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /Codes/paxoskv/impl.go
 */
package paxoskv

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/kr/pretty"
	"github.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	NotEnoughQuorum = errors.New("not enough qourum")
	AcceptorBasePort = int64(3333)
)

// GE compare two ballot number a, b and return whether a >= b in a bool
//我们怎么去定义这个bool
func (a *BallotNum) GE(b *BallotNum) bool {
	if a.N > b.N {
		return true;
	}
	if a.N < b.N {
		return false;
	}
	return a.ProposerId >= b.ProposerId
}

// RunPaxos execute the paxos phase-1 and phase-2 to establish a value.
// 'value' is the value caller wants to propose.
// It returns the established value, which may be a voted value that is not 'val'.
//之后有完整运行一次RunPaxos的方法
func (p *Proposer) RunPaxos(acceptorIds []int64, val *Value) *Value {
	quorum := len(acceptorIds)/2 + 1   //Go中, :=表示声明并赋值，并且系统自动推断类型，不需要var关键字

	for {
		p.Val = nil

		maxVotedVal, higherBal, err := p.Phase1(acceptorIds, quorum)
		if err != nil {
			pretty.Logf("Proposer: fail to run phase-1: highest ballot: %v, increment ballot and retry", higherBal)
			p.Bal.N = higherBal.N + 1
			continue
		}

		if maxVotedVal == nil {
			pretty.Logf("Proposer: no voted value seen, propose my value: %v", val)
		} else {
			val = maxVotedVal
		}

		if val == nil {
			pretty.Logf("Proposer: no value to propose in phase-2, quit")
			return nil
		}

		p.Val = val
		pretty.Logf("Proposer: proposer chose value to propose: %s", p.Val)

		higherBal, err = p.Phase2(acceptorIds, quorum)
		if err != nil {
			pretty.Logf("Proposer: fail to run phase-2: highest ballot: %v, increment ballot and retry", higherBal)
			p.Bal.N = higherBal.N + 1;
			continue
		}

		pretty.Logf("Proposer: value is voted by a quorum and has been safe: %v", maxVotedVal)
		return p.Val
	}
}


// Phase1 run paxos phase-1 on the specified acceptorIds.
// if a higher ballot number is seen and phase-1 failed to constitute a quorum         我们怎么定义这个quorum？
// one of the higher ballot number and a NotEnoughQuorum is returned.
func (p *Proposer) Phase1(acceptorIds []int64, quorum int)(*Value, *BallotNum, error) {

	replies := p.rpcToAll(acceptorIds, "Prepare")

	ok := 0
	higherBal := *p.Bal
	maxVoted := &Acceptor(VBal: &BallotNum{})

	for _, r:= range replies {
		
		pretty.Logf("Proposer: handing Prepare reply: %s", r)
		if !p.Bal.GE(r.LastBal)
	}
}

// Phase2 run paxos Phase-2 on the specified acceptorIds.
// if a higher ballot number is seen and phase-2 failed to constitute a quorum,
// one of the higher ballot number and a NotEnoughQuorum is returned.
func (p *Proposer) Phase2(acceptorIds []int64, quorum int)(*BallotNum, error) {

}

//@return replies := []*Acceptor{} 
//[]*Acceptor is string?
func (p *Proposer) rpcToAll(acceptorIds []int64, action string []*Acceptor) {
	
	replies := []*Acceptor{}

	for _, aid := range acceptorIds {
		var err error
		address := fmt.Sprintf("127.0.0.1:%d", AcceptorBasePort+int64(aid))
		// set up a connection to the server
		// 这里还要去看grpc裸露出来的接口
		conn, err := grpc.Dial(address,grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		defer conn.Close()
		// 这里NewPaxosKVClient ???
		c := NewPaxosKVClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		var reply *Acceptor
		// 这里的action是？
		if action == "Prepare" {
			reply, err = c.Prepare(ctx, p)
		} else if action == "Accept" {
			reply, err = c.Accept(ctx, p)
		}
		if err != nil {
			log.Printf("Proposer: %s failure from Acceptor-%d: %v", action, aid, err)
		}
		log.Printf("Proposer: recv %s reply from: Acceptor-%d: %v", action, aid, reply)

		replies = append(replies, reply)
	}
	return replies
}

//需要加锁
// Version defines one modification of a key-value record
// It is barely an Acceptor with a lock
type Version struct {
	mu			sync.Mutex
	acceptor 	Acceptor
}

// Versions stones all versions of a record
// The value of every version is decided by a paxos instance, e.g. an Acceptor.
type Versions map[int64]*Version

// KVServer impl the paxos Acceptor API: handing Prepare and Accept request.
type KVServer struct {
	mu		sync.Mutex
	Storage map[string]Versions
}
