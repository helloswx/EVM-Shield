
Through Txtrace as the full node, the local private chain or the full synchronization with Ethereum, the opcode and operand sequence of each transaction can be obtained in mongodb, which can be used to evaluate SODA and Aegis.

![image](https://github.com/helloswx/EVM-Shield/assets/52559342/1f6cb417-db4a-4da6-a6e5-ae66caa79979)


**For example, a trace fragment of a transaction looks like this:**
```
0;PUSH1;96&2;PUSH1;64&4;MSTORE;&5;CALLDATASIZE;4&6;ISZERO;&7;PUSH2;185&10;JUMPI;&11;PUSH1;224&13;PUSH1;2&15;EXP;&16;PUSH1;0&18;CALLDATALOAD;44493657185760383935816513384858880510381257242149215912015128952099573858304&19;DIV;&20;PUSH4;330252341&25;DUP2;&26;EQ;&27;PUSH2;414&30;JUMPI;&31;DUP1;&32;PUSH4;653633737&37;EQ;&38;PUSH2;449&41;JUMPI;&42;DUP1;&43;PUSH4;924821588&48;EQ;&49;PUSH2;458&52;JUMPI;&53;DUP1;&54;PUSH4;1096947359&59;EQ;&60;PUSH2;467
```

