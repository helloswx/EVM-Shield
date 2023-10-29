// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"bytes"
	"io/ioutil"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"

	//【*】
	//"fmt"
	"encoding/json"
	"strings"
	//"strconv"
	//"reflect"
	//mapset "github.com/deckarep/golang-set"
	//"github.com/ethereum/go-ethereum/crypto"
)

// ContractRef is a reference to the contract's backing object
type ContractRef interface {
	Address() common.Address
}

// AccountRef implements ContractRef.
//
// Account references are used during EVM initialisation and
// it's primary use is to fetch addresses. Removing this object
// proves difficult because of the cached jump destinations which
// are fetched from the parent contract (i.e. the caller), which
// is a ContractRef.
type AccountRef common.Address

// Address casts AccountRef to a Address
func (ar AccountRef) Address() common.Address { return (common.Address)(ar) }

//【*】变量名对应的绑定信息
const (
	Common = iota
	Package
	Mapping
	Dynamic
)

type Variable struct {
	Slot         []uint256.Int //所有slot
	StartSlot    uint256.Int
	PackageStart int //用于紧凑存储的情况,其他类型即为0
	DataType     int //Common,Package,Dynamic,Mapping
	DataSize     int //对于静态类型是数据长度，对于动态数组是元素长度，单位字节

	DynamicHash uint256.Int //
	DynamicEnd  uint256.Int //

	SubValueType string //value的类型
	SubVars      []Variable
}

func (v *Variable) InitSlot() {
	if v.DataSize <= 32 {
		if v.DataSize == 0 {
			v.DataSize = 32
		}
		v.Slot = make([]uint256.Int, 1)
		v.Slot[0] = v.StartSlot

	} else {
		v.Slot = make([]uint256.Int, (v.DataSize-1)/32+1)
		s := v.StartSlot
		for i := 0; i < (v.DataSize-1)/32+1; i++ {
			v.Slot[i] = s
			s.Add(&s, uint256.NewInt(1))
		}
	}
	if v.DataType == Mapping || v.DataType == Dynamic {
		v.SubVars = make([]Variable, 0)
	}
}

func (v *Variable) Contains(who uint256.Int) bool {
	for i := 0; i < len(v.Slot); i++ {
		if v.Slot[i] == who {
			return true
		}
	}
	if v.DataType == Dynamic {
		res1 := v.DynamicHash.Cmp(&who)
		res2 := v.DynamicEnd.Cmp(&who)
		if res1 <= 0 && res2 == 1 {
			return true
		}
	}
	return false
}

type Firewall struct {
	Hash int
	Vars []int
}
type Rule struct {
	Firewalls []Firewall
	FocusVars []Variable
	// CanSendEther bool
}

// Contract represents an ethereum contract in the state database. It contains
// the contract code, calling arguments. Contract implements ContractRef
type Contract struct {
	// CallerAddress is the result of the caller which initialised this
	// contract. However when the "call method" is delegated this value
	// needs to be initialised to that of the caller's caller.
	CallerAddress common.Address
	caller        ContractRef
	self          ContractRef

	jumpdests map[common.Hash]bitvec // Aggregated result of JUMPDEST analysis.
	analysis  bitvec                 // Locally cached result of JUMPDEST analysis

	Code     []byte
	CodeHash common.Hash
	CodeAddr *common.Address
	Input    []byte

	Gas   uint64
	value *big.Int

	//【*】屏蔽规则
	rule        Rule
	CheckedHash int
	TmpVars     []Variable
}

func getRuleByAddr(adr ContractRef) Rule {
	var rule Rule
	str, err := ioutil.ReadFile("rules/" + adr.Address().String())
	if err != nil {
		rule.FocusVars = make([]Variable, 0)
		rule.Firewalls = make([]Firewall, 0)

	} else {
		_ = json.Unmarshal([]byte(str), &rule) //第二个参数用来存放第一个参数的内容,c2需要被修改，想要被需改成功必须得传入指针
	}
	return rule
}
func Write(rule *Rule, adr ContractRef) {
	data, err := json.MarshalIndent(rule, "", "	")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("rules/"+adr.Address().String(), data, 0777)
	if err != nil {
		panic(err)
	}
}

// NewContract returns a new contract environment for the execution of EVM.
func NewContract(caller ContractRef, object ContractRef, value *big.Int, gas uint64) *Contract {
	c := &Contract{CallerAddress: caller.Address(), caller: caller, self: object}

	if parent, ok := caller.(*Contract); ok {
		// Reuse JUMPDEST analysis from parent context if available.
		c.jumpdests = parent.jumpdests
	} else {
		c.jumpdests = make(map[common.Hash]bitvec)
	}

	// Gas should be a pointer so it can safely be reduced through the run
	// This pointer will be off the state transition
	c.Gas = gas
	// ensures a value is set
	c.value = value

	//【*】
	c.rule = getRuleByAddr(object)
	c.TmpVars = make([]Variable, 0)

	return c
}
func add(slots *[]uint256.Int, slot uint256.Int) {
	for _, s := range *slots {
		if bytes.Equal(slot.Bytes(), s.Bytes()) {
			return
		}
	}
	*slots = append(*slots, slot)
}

//【*】屏蔽逻辑,SSTORE时调用
func (v *Variable) Shield(loc uint256.Int, val uint256.Int, interpreter *EVMInterpreter, scope *ScopeContext) bool {
	if v.Contains(loc) {
		if v.DataType == Package {
			valueothers := val.Bytes32()
			valueoriginal := interpreter.evm.StateDB.GetState(scope.Contract.Address(), common.BytesToHash(loc.Bytes()))
			return bytes.Equal(valueothers[v.PackageStart:v.PackageStart+v.DataSize], valueoriginal[v.PackageStart:v.PackageStart+v.DataSize])
		}
		return false
	}
	return true
}

//【*】SHA3识别
// 本函数的功能在于
//给定slot，寻找是否为标记的变量
//记录hash(slot)
func (v *Variable) SSI(slot uint256.Int, hash uint256.Int, interpreter *EVMInterpreter, scope *ScopeContext, isShield bool) {
	if v.DataType != Mapping && v.DataType != Dynamic {
		return
	}
	if bytes.Equal(slot.Bytes(), v.StartSlot.Bytes()) {
		clan := strings.Split(v.SubValueType, "->")
		if v.DataType == Mapping {
			if len(clan) == 1 && isShield {
				for _, Evar := range scope.Contract.TmpVars {
					if Evar.StartSlot == hash {
						return
					}
				}
				scope.Contract.TmpVars = append(scope.Contract.TmpVars, Variable{
					StartSlot: hash,
					DataType:  Common,
				})
				scope.Contract.TmpVars[len(scope.Contract.TmpVars)-1].InitSlot()
				//fmt.Println("TmpVars",scope.Contract.TmpVars)

			} else if len(clan) != 1 {
				for _, Evar := range v.SubVars {
					if Evar.StartSlot == hash {
						return
					}
				}
				Value := clan[0]
				var Value_Number int
				if Value == "Mapping" {
					Value_Number = Mapping
				} else if Value == "Dynamic" {
					Value_Number = Dynamic
				}
				newVar := Variable{
					StartSlot:    hash,
					PackageStart: 0,
					DataType:     Value_Number,
					DataSize:     v.DataSize,
					SubValueType: clan[1],
				}
				newVar.InitSlot()
				//if(Value_Number == Dynamic){
				//	newVar.SSI(slot, hash, interpreter, scope, isShield)
				//}
				v.SubVars = append(v.SubVars, newVar)
			}
		} else {
			var Len uint256.Int
			Len.SetBytes(interpreter.evm.StateDB.GetState(scope.Contract.Address(), common.BytesToHash(slot.Bytes())).Bytes())
			Zero := uint256.NewInt(0).Bytes32()
			LenB := Len.Bytes32()
			if bytes.Equal(LenB[0:32], Zero[0:32]) {
				return
			}
			if v.DataSize >= 32 {
				per := uint64((v.DataSize-1)/32 + 1)
				Len.Mul(&Len, uint256.NewInt(per))

			} else {
				per := uint64(32 / v.DataSize)
				Len.Sub(&Len, uint256.NewInt(1))
				Len.Div(&Len, uint256.NewInt(per))
				Len.Add(&Len, uint256.NewInt(1))
			}
			if v.DynamicHash.Cmp(&hash) != 0 {
				v.SubVars = make([]Variable, 0)
				v.DynamicHash.SetBytes(hash.Bytes())
			}
			v.DynamicEnd.SetBytes(hash.Bytes())
			v.DynamicEnd.Add(&v.DynamicEnd, &Len)
			if len(clan) > 1 {
				Value := clan[0]
				var Value_Number int
				if Value == "Mapping" {
					Value_Number = Mapping
				} else if Value == "Dynamic" {
					Value_Number = Dynamic
				}
				Count := v.DynamicEnd
				for Count.Cmp(&v.DynamicEnd) != 0 {
					v.SubVars = append(v.SubVars, Variable{
						StartSlot:    Count,
						PackageStart: 0,
						DataType:     Value_Number,
						DataSize:     v.DataSize,
						SubValueType: clan[1],
					})
					Count.Sub(&Count, uint256.NewInt(1))
				}
			}
		}

	} else if len(v.SubVars) != 0 {
		for _, Evar := range v.SubVars {
			Evar.SSI(slot, hash, interpreter, scope, isShield)
		}
	}
}

func (c *Contract) validJumpdest(dest *uint256.Int) bool {
	udest, overflow := dest.Uint64WithOverflow()
	// PC cannot go beyond len(code) and certainly can't be bigger than 63bits.
	// Don't bother checking for JUMPDEST in that case.
	if overflow || udest >= uint64(len(c.Code)) {
		return false
	}
	// Only JUMPDESTs allowed for destinations
	if OpCode(c.Code[udest]) != JUMPDEST {
		return false
	}
	return c.isCode(udest)
}

// isCode returns true if the provided PC location is an actual opcode, as
// opposed to a data-segment following a PUSHN operation.
func (c *Contract) isCode(udest uint64) bool {
	// Do we already have an analysis laying around?
	if c.analysis != nil {
		return c.analysis.codeSegment(udest)
	}
	// Do we have a contract hash already?
	// If we do have a hash, that means it's a 'regular' contract. For regular
	// contracts ( not temporary initcode), we store the analysis in a map
	if c.CodeHash != (common.Hash{}) {
		// Does parent context have the analysis?
		analysis, exist := c.jumpdests[c.CodeHash]
		if !exist {
			// Do the analysis and save in parent context
			// We do not need to store it in c.analysis
			analysis = codeBitmap(c.Code)
			c.jumpdests[c.CodeHash] = analysis
		}
		// Also stash it in current contract for faster access
		c.analysis = analysis
		return analysis.codeSegment(udest)
	}
	// We don't have the code hash, most likely a piece of initcode not already
	// in state trie. In that case, we do an analysis, and save it locally, so
	// we don't have to recalculate it for every JUMP instruction in the execution
	// However, we don't save it within the parent context
	if c.analysis == nil {
		c.analysis = codeBitmap(c.Code)
	}
	return c.analysis.codeSegment(udest)
}

// AsDelegate sets the contract to be a delegate call and returns the current
// contract (for chaining calls)
func (c *Contract) AsDelegate() *Contract {
	// NOTE: caller must, at all times be a contract. It should never happen
	// that caller is something other than a Contract.
	parent := c.caller.(*Contract)
	c.CallerAddress = parent.CallerAddress
	c.value = parent.value

	return c
}

// GetOp returns the n'th element in the contract's byte array
func (c *Contract) GetOp(n uint64) OpCode {
	if n < uint64(len(c.Code)) {
		return OpCode(c.Code[n])
	}

	return STOP
}

// Caller returns the caller of the contract.
//
// Caller will recursively call caller when the contract is a delegate
// call, including that of caller's caller.
func (c *Contract) Caller() common.Address {
	return c.CallerAddress
}

// UseGas attempts the use gas and subtracts it and returns true on success
func (c *Contract) UseGas(gas uint64) (ok bool) {
	if c.Gas < gas {
		return false
	}
	c.Gas -= gas
	return true
}

// Address returns the contracts address
func (c *Contract) Address() common.Address {
	return c.self.Address()
}

// Value returns the contract's value (sent to it from it's caller)
func (c *Contract) Value() *big.Int {
	return c.value
}

// SetCallCode sets the code of the contract and address of the backing data
// object
func (c *Contract) SetCallCode(addr *common.Address, hash common.Hash, code []byte) {
	c.Code = code
	c.CodeHash = hash
	c.CodeAddr = addr
}

// SetCodeOptionalHash can be used to provide code, but it's optional to provide hash.
// In case hash is not provided, the jumpdest analysis will not be saved to the parent context
func (c *Contract) SetCodeOptionalHash(addr *common.Address, codeAndHash *codeAndHash) {
	c.Code = codeAndHash.code
	c.CodeHash = codeAndHash.hash
	c.CodeAddr = addr
}
