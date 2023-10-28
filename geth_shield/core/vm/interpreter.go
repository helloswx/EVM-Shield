// Copyright 2014 The go-ethereum Authors
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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	//【*】
	"github.com/holiman/uint256"
	"strings"
	"encoding/json"
	//"strconv"
	"fmt"
)

// Config are the configuration options for the Interpreter
type Config struct {
	Debug                   bool      // Enables debugging
	Tracer                  EVMLogger // Opcode logger
	NoBaseFee               bool      // Forces the EIP-1559 baseFee to 0 (needed for 0 price calls)
	EnablePreimageRecording bool      // Enables recording of SHA3/keccak preimages

	JumpTable *JumpTable // EVM instruction table, automatically populated if unset

	ExtraEips []int // Additional EIPS that are to be enabled
}

// ScopeContext contains the things that are per-call, such as stack and memory,
// but not transients like pc and gas
type ScopeContext struct {
	Memory   *Memory
	Stack    *Stack
	Contract *Contract
}

// EVMInterpreter represents an EVM interpreter
type EVMInterpreter struct {
	evm *EVM
	cfg Config

	hasher    crypto.KeccakState // Keccak256 hasher instance shared across opcodes
	hasherBuf common.Hash        // Keccak256 hasher result array shared aross opcodes

	readOnly   bool   // Whether to throw on stateful modifications
	returnData []byte // Last CALL's return data for subsequent reuse
}

// NewEVMInterpreter returns a new instance of the Interpreter.
func NewEVMInterpreter(evm *EVM, cfg Config) *EVMInterpreter {
	// If jump table was not initialised we set the default one.
	if cfg.JumpTable == nil {
		switch {
		case evm.chainRules.IsMerge:
			cfg.JumpTable = &mergeInstructionSet
		case evm.chainRules.IsLondon:
			cfg.JumpTable = &londonInstructionSet
		case evm.chainRules.IsBerlin:
			cfg.JumpTable = &berlinInstructionSet
		case evm.chainRules.IsIstanbul:
			cfg.JumpTable = &istanbulInstructionSet
		case evm.chainRules.IsConstantinople:
			cfg.JumpTable = &constantinopleInstructionSet
		case evm.chainRules.IsByzantium:
			cfg.JumpTable = &byzantiumInstructionSet
		case evm.chainRules.IsEIP158:
			cfg.JumpTable = &spuriousDragonInstructionSet
		case evm.chainRules.IsEIP150:
			cfg.JumpTable = &tangerineWhistleInstructionSet
		case evm.chainRules.IsHomestead:
			cfg.JumpTable = &homesteadInstructionSet
		default:
			cfg.JumpTable = &frontierInstructionSet
		}
		var extraEips []int
		for _, eip := range cfg.ExtraEips {
			copy := *cfg.JumpTable
			if err := EnableEIP(eip, &copy); err != nil {
				// Disable it, so caller can check if it's activated or not
				log.Error("EIP activation failed", "eip", eip, "error", err)
			} else {
				extraEips = append(extraEips, eip)
			}
			cfg.JumpTable = &copy
		}
		cfg.ExtraEips = extraEips
	}

	return &EVMInterpreter{
		evm: evm,
		cfg: cfg,
	}
}

//【*】策略部署
// function 读取文件,json反序列化，添加到contract对象内
func Example_Rule(adr ContractRef){
	var StartSlot uint256.Int
	var PackageStart int 
	var flag bool
	
	var firewall Firewall
	firewall.Vars=make([]int,5)
	rule := getRuleByAddr(adr)
	
	
	StartSlot.SetBytes([]byte{0x0})
	PackageStart = 0
	flag = true
	for i:=0;i<len(rule.FocusVars);i++{	
		if(rule.FocusVars[i].StartSlot==StartSlot && rule.FocusVars[i].PackageStart == PackageStart){
			flag = false
			firewall.Vars[0] = i
			break
		}
	}
	if(flag){
		firewall.Vars[0]=len(rule.FocusVars)
		rule.FocusVars=append(rule.FocusVars,Variable{
			StartSlot: StartSlot,
			DataType: Common,
			DataSize: 64,
	        })
	        rule.FocusVars[len(rule.FocusVars)-1].InitSlot()
	}
	
	StartSlot.SetBytes([]byte{0x4})
	PackageStart = 0
	flag = true
	for i:=0;i<len(rule.FocusVars);i++{	
		if(rule.FocusVars[i].StartSlot==StartSlot && rule.FocusVars[i].PackageStart == PackageStart){
			flag = false
			firewall.Vars[1] = i
			break
		}
	}
	if(flag){
		firewall.Vars[1]=len(rule.FocusVars)
		rule.FocusVars=append(rule.FocusVars,Variable{
			StartSlot: StartSlot,
			DataType: Common,
	        })
	        rule.FocusVars[len(rule.FocusVars)-1].InitSlot()
	}
	
	
	StartSlot.SetBytes([]byte{0x2})
	PackageStart = 31
	flag = true
	for i:=0;i<len(rule.FocusVars);i++{	
		if(rule.FocusVars[i].StartSlot==StartSlot && rule.FocusVars[i].PackageStart == PackageStart){
			flag = false
			firewall.Vars[2] = i
			break
		}
	}
	if(flag){
		firewall.Vars[2]=len(rule.FocusVars)
		rule.FocusVars=append(rule.FocusVars,Variable{
			StartSlot: StartSlot,
			DataSize:  1,
			PackageStart: 31,
			DataType: Package,
	        })
	        rule.FocusVars[len(rule.FocusVars)-1].InitSlot()
	}
	
	
	StartSlot.SetBytes([]byte{0x7})
	PackageStart = 0
	flag = true
	for i:=0;i<len(rule.FocusVars);i++{	
		if(rule.FocusVars[i].StartSlot==StartSlot && rule.FocusVars[i].PackageStart == PackageStart){
			flag = false
			firewall.Vars[3] = i
			break
		}
	}
	if(flag){
		firewall.Vars[3]=len(rule.FocusVars)
		rule.FocusVars=append(rule.FocusVars,Variable{
			StartSlot: StartSlot,
			DataType: Mapping,
			SubValueType: "uint256",
	        })
	        rule.FocusVars[len(rule.FocusVars)-1].InitSlot()
	}
	
	StartSlot.SetBytes([]byte{0x10})
	PackageStart = 0
	flag = true
	for i:=0;i<len(rule.FocusVars);i++{	
		if(rule.FocusVars[i].StartSlot==StartSlot && rule.FocusVars[i].PackageStart == PackageStart){
			flag = false
			firewall.Vars[4] = i
			break
		}
	}
	if(flag){
		firewall.Vars[4]=len(rule.FocusVars)
		rule.FocusVars=append(rule.FocusVars,Variable{
			StartSlot: StartSlot,
			DataType: Dynamic,
			SubValueType: "uint256",
	        })
	        rule.FocusVars[len(rule.FocusVars)-1].InitSlot()
	}
	
	for i:=0;i<5;i++{
		rule.FocusVars[i].InitSlot()
	}

	rule.Firewalls = make([]Firewall,1)
	firewall.Hash=0xc0406226
	rule.Firewalls[0] = firewall	
	
	Write(&rule,adr)
}
func Example2_Rule(adr ContractRef){
	var StartSlot uint256.Int
	var PackageStart int 
	var flag bool
	
	var firewall Firewall
	firewall.Vars=make([]int,1)
	rule := getRuleByAddr(adr)
	
	StartSlot.SetBytes([]byte{0x0})
	PackageStart = 0
	flag = true
	for i:=0;i<len(rule.FocusVars);i++{	
		if(rule.FocusVars[i].StartSlot==StartSlot && rule.FocusVars[i].PackageStart == PackageStart){
			flag = false
			firewall.Vars[0] = i
			break
		}
	}
	if(flag){
		firewall.Vars[0]=len(rule.FocusVars)
		rule.FocusVars=append(rule.FocusVars,Variable{
			StartSlot: StartSlot,
			DataType: Mapping,
			SubValueType: "Mapping->uint256",
	        })
	        rule.FocusVars[len(rule.FocusVars)-1].InitSlot()
	}
	
	for i:=0;i<1;i++{
		rule.FocusVars[i].InitSlot()
	}

	rule.Firewalls = make([]Firewall,1)
	firewall.Hash=0xc0406226
	rule.Firewalls[0] = firewall	
	
	Write(&rule,adr)
}
//【*】NewRule returns a new contract environment with the rule information for the execution of EVM.
func UpdateRule(adr ContractRef,data []byte){
	rule := getRuleByAddr(adr)
	
	fmt.Println("Original")
	str :=string(data)
	fmt.Println(str)
	fmt.Println("Splite")
	paras := strings.Split(string(data),"|")
	FidData := paras[0]
	VarsData := paras[1]
	//fmt.Println(str)
	//fmt.Println(string(Fid))
	//fmt.Println(string(VarsData))
	
	var Fid int
	var vars []Variable
	_ = json.Unmarshal([]byte(FidData), &Fid)
	_ = json.Unmarshal([]byte(VarsData), &vars)
	fmt.Println(Fid)
	fmt.Println(vars)
	data, err := json.MarshalIndent(vars, "", "	")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
	orders := make([]int,len(vars))
	for id,Evar := range vars{
		flag := true
		for oid,old := range rule.FocusVars{
			if(old.StartSlot == Evar.StartSlot && old.PackageStart == Evar.PackageStart){
				orders[id] = oid
				flag = false
				break
			}
		}
		if(flag){
			orders[id] = len(rule.FocusVars)
			Evar.InitSlot()
			rule.FocusVars = append(rule.FocusVars,Evar)
		}
	}
	flag := true
	for id,firewall := range rule.Firewalls{
		if(firewall.Hash == Fid){
			rule.Firewalls[id] = Firewall{
				Hash:Fid,
				Vars:orders,
			}
			flag = false
			break
		}
	}
	if(flag){
		rule.Firewalls = append(rule.Firewalls,Firewall{
			Hash:Fid,
			Vars:orders,
		})
	}
	
	Write(&rule,adr)
}
// Run loops and evaluates the contract's code with the given input data and returns
// the return byte-slice and an error if one occurred.
//
// It's important to note that any errors returned by the interpreter should be
// considered a revert-and-consume-all-gas operation except for
// ErrExecutionReverted which means revert-and-keep-gas-left.
func (in *EVMInterpreter) Run(contract *Contract, input []byte, readOnly bool) (ret []byte, err error) {
	// Increment the call depth which is restricted to 1024
	in.evm.depth++
	defer func() { in.evm.depth-- }()

	// Make sure the readOnly is only set if we aren't in readOnly yet.
	// This also makes sure that the readOnly flag isn't removed for child calls.
	if readOnly && !in.readOnly {
		in.readOnly = true
		defer func() { in.readOnly = false }()
	}

	// Reset the previous call's return data. It's unimportant to preserve the old buffer
	// as every returning call will return new data anyway.
	in.returnData = nil

	// Don't bother with the execution if there's no code.
	if len(contract.Code) == 0 {
		return nil, nil
	}
	//【*】
	if(len(input)>=4){
		var sum = 0;
		for i:=0;i<4;i++{
			sum =sum * 256 + int(input[i])
		}
		contract.CheckedHash = sum
		if(len(input)>=5&&sum==0&&input[4]==0x11){
			UpdateRule(contract.self,input[5:])
			//Example_Rule(contract.self)
		}
	}

	var (
		op          OpCode        // current opcode
		mem         = NewMemory() // bound memory
		stack       = newstack()  // local stack
		callContext = &ScopeContext{
			Memory:   mem,
			Stack:    stack,
			Contract: contract,
		}
		// For optimisation reason we're using uint64 as the program counter.
		// It's theoretically possible to go above 2^64. The YP defines the PC
		// to be uint256. Practically much less so feasible.
		pc   = uint64(0) // program counter
		cost uint64
		// copies used by tracer
		pcCopy  uint64 // needed for the deferred EVMLogger
		gasCopy uint64 // for EVMLogger to log gas remaining before execution
		logged  bool   // deferred EVMLogger should ignore already logged steps
		res     []byte // result of the opcode execution function
	)
	// Don't move this deferred function, it's placed before the capturestate-deferred method,
	// so that it get's executed _after_: the capturestate needs the stacks before
	// they are returned to the pools
	defer func() {
		returnStack(stack)
	}()
	contract.Input = input

	if in.cfg.Debug {
		defer func() {
			if err != nil {
				if !logged {
					in.cfg.Tracer.CaptureState(pcCopy, op, gasCopy, cost, callContext, in.returnData, in.evm.depth, err)
				} else {
					in.cfg.Tracer.CaptureFault(pcCopy, op, gasCopy, cost, callContext, in.evm.depth, err)
				}
			}
		}()
	}
	// The Interpreter main run loop (contextual). This loop runs until either an
	// explicit STOP, RETURN or SELFDESTRUCT is executed, an error occurred during
	// the execution of one of the operations or until the done flag is set by the
	// parent context.
	for {
		if in.cfg.Debug {
			// Capture pre-execution values for tracing.
			logged, pcCopy, gasCopy = false, pc, contract.Gas
		}
		// Get the operation from the jump table and validate the stack to ensure there are
		// enough stack items available to perform the operation.
		op = contract.GetOp(pc)
		operation := in.cfg.JumpTable[op]
		cost = operation.constantGas // For tracing
		// Validate stack
		if sLen := stack.len(); sLen < operation.minStack {
			return nil, &ErrStackUnderflow{stackLen: sLen, required: operation.minStack}
		} else if sLen > operation.maxStack {
			return nil, &ErrStackOverflow{stackLen: sLen, limit: operation.maxStack}
		}
		if !contract.UseGas(cost) {
			return nil, ErrOutOfGas
		}
		if operation.dynamicGas != nil {
			// All ops with a dynamic memory usage also has a dynamic gas cost.
			var memorySize uint64
			// calculate the new memory size and expand the memory to fit
			// the operation
			// Memory check needs to be done prior to evaluating the dynamic gas portion,
			// to detect calculation overflows
			if operation.memorySize != nil {
				memSize, overflow := operation.memorySize(stack)
				if overflow {
					return nil, ErrGasUintOverflow
				}
				// memory is expanded in words of 32 bytes. Gas
				// is also calculated in words.
				if memorySize, overflow = math.SafeMul(toWordSize(memSize), 32); overflow {
					return nil, ErrGasUintOverflow
				}
			}
			// Consume the gas and return an error if not enough gas is available.
			// cost is explicitly set so that the capture state defer method can get the proper cost
			var dynamicCost uint64
			dynamicCost, err = operation.dynamicGas(in.evm, contract, stack, mem, memorySize)
			cost += dynamicCost // for tracing
			if err != nil || !contract.UseGas(dynamicCost) {
				return nil, ErrOutOfGas
			}
			// Do tracing before memory expansion
			if in.cfg.Debug {
				in.cfg.Tracer.CaptureState(pc, op, gasCopy, cost, callContext, in.returnData, in.evm.depth, err)
				logged = true
			}
			if memorySize > 0 {
				mem.Resize(memorySize)
			}
		} else if in.cfg.Debug {
			in.cfg.Tracer.CaptureState(pc, op, gasCopy, cost, callContext, in.returnData, in.evm.depth, err)
			logged = true
		}
		// execute the operation
		res, err = operation.execute(&pc, in, callContext)
		if err != nil {
			break
		}
		pc++
	}

	if err == errStopToken {
		err = nil // clear stop token error
	}

	return res, err
}
