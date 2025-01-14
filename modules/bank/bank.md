# PLUGCHAIN SDK GO

## BANK MODULE

- [Query](#query)
  - [QueryAccount](#account) --Account Amount
  - [TotalSupply](#supply) --TotalSupply
- [TX](#tx)
  - [Send](#send) --Transfer
  - [MsgSend](#msgsend) --MsgSend

# realization

## Query<a name="query"></a><br/>

#### QueryAccount<a name="account"></a><br/>
>QueryAccount return account information specified address
```go
balance, err := client.Bank.QueryAccount("gx1yhf7w0sq8yn6gqre2pulnqwyy30tjfc4v08f3x")
plug:=balance.Coins.AmountOf("plug")
```

#### TotalSupply<a name="supply"></a><br/>
>TotalSupply queries the total supply of all coins.
```go
supply, err := client.Bank.TotalSupply()
```

## TX<a name="tx"></a><br/>

#### Send<a name="send"></a><br/>
>Send is responsible for transferring tokens from `From` to `to` account

**You need to import the private key before you can operate，Please see the key package for importing the private key**

```go
baseTx := types.BaseTx{
    From:     "demo", //Account name 
    Password: "123123123",
    Gas:      200000,
    Mode:     types.Commit,
    Memo:     "test",
}
baseTx.Fee, err = types.ParseDecCoins("2000plug") //Fee
coins, err := types.ParseDecCoins("100000plug")   //Transfer out quantity + currency name, for example:100000plug
to := "gx1akqhezuftdcc0eqzkq5peqpjlucgmyr7srx54j" //Receiving address
result, err := client.Bank.Send(to, coins, baseTx)
```


#### MsgSend<a name="msgsend"></a><br/>
>get TxHash before sending transactions

**You need to import the private key before you can operate，Please see the key package for importing the private key**
```go
baseTx := types.BaseTx{
    From:     "demo", //Account name 
    Password: "123123123",
    Gas:      200000,
    Mode:     types.Commit,
    Memo:     "test",
}
coins, err := types.ParseCoins("100000plug")
from := "gx1yhf7w0sq8yn6gqre2pulnqwyy30tjfc4v08f3x"
to := "gx1akqhezuftdcc0eqzkq5peqpjlucgmyr7srx54j"
msg := &bank.MsgSend{
    FromAddress: from,
    ToAddress:   to,
    Amount:      coins,
}
txhash, err := client.BuildTxHash([]types.Msg{msg}, baseTx)
```