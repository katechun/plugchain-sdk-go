# PLUGCHAIN SDK GO

## TOKEN MODULE

- [Query](#query)
    - [QueryToken](#token) --QueryToken
    - [QueryTokens](#tokens) --QueryTokens
    - [QueryFees](#fees) --QueryFees
    - [QueryParams](#params) --QueryParams
- [TX](#tx)
    - [IssueToken](#issue) --IssueToken
    - [EditToken](#edit) --EditToken
    - [TransferToken](#transfer) --TransferToken
    - [MintToken](#mint) --MintToken


# realization

## Query<a name="query"></a><br/>

#### QueryToken<a name="token"></a><br/>
>Query a single token
```go
token, err := client.Token.QueryToken("test1")
```

#### QueryTokens<a name="tokens"></a><br/>
>Query all tokens
```go
token, err := client.Token.QueryTokens("")
```

#### QueryFees<a name="fees"></a><br/>
>Inquiry fee
```go
token, err := client.Token.QueryFees("test1")
```

#### QueryParams<a name="params"></a><br/>
>Query parameters
```go
res, err := client.Token.QueryParams()
```


## TX<a name="tx"></a><br/>

#### IssueToken<a name="issue"></a><br/>
>Issue tokens
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
issueTokenReq := token.IssueTokenRequest{
	Symbol:        "test1",
	Name:          "testToken",
	Scale:         8,
	MinUnit:       "tt1",
	InitialSupply: 10000000,
	MaxSupply:     21000000,
	Mintable:      true,
}
rs, err := client.Token.IssueToken(issueTokenReq, baseTx)
```

#### EditToken<a name="edit"></a><br/>
>Modify token
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
editTokenReq := token.EditTokenRequest{
    Symbol:    "test1",
    Name:      "testToken66",
    MaxSupply: 22000000,
}
rs, err := client.Token.EditToken(editTokenReq, baseTx)
```

#### TransferToken<a name="transfer"></a><br/>
>Transfer token ownership
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
rs, err := client.Token.TransferToken("gx1akqhezuftdcc0eqzkq5peqpjlucgmyr7srx54j", "test1", baseTx)
```

#### MintToken<a name="mint"></a><br/>
>Coinage token
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
rs, err = client.Token.MintToken("test1", 11000000, "gx1yhf7w0sq8yn6gqre2pulnqwyy30tjfc4v08f3x", baseTx)
```

