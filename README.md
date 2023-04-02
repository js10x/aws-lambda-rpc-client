# AWS Lambda RPC Client 
A simple RPC client that can be used to debug AWS lambda functions locally that are written in Go. This is made possible by the fact that the AWS SDK for Golang implements lambda functions as RPC servers, making them very easy to invoke in a multitude of different environments.

### Compile From Source (default build with flag for removing the debug symbols from the binary)

```powershell
go build -ldflags=-w -o "lrpc.exe"
```

### Prerequisites

* Must have a lambda to test against
* Must have the `_LAMBDA_SERVER_PORT` environment variable set to a value of `9988` in either your local configuration or in your lambda source code before `lambda.Start(handler)` is invoked (i.e. `os.Setenv("_LAMBDA_SERVER_PORT", "9988")`).

### Usage

First, you'll need to run your lambda function by invoking the `lambda.Start(handler)` method (i.e. `go run <module-name>`), this will start the RPC server for that particular lambda function and it will begin to run as a server.

### Example Output (Terminal)

```
./lrpc.exe -e "./events/apigateway-authorizer.json" -o "resp.json"
LAMBDA RESPONSE:
 {
    "policyDocument": {
        "Statement": [
            {
                "Action": [
                    "execute-api:Invoke"     
                ],
                "Effect": "Deny",
                "Resource": [
                    "*"
                ]
            }
        ],
        "Version": "2012-10-17"
    },
    "principalId": "apigateway.amazonaws.com"
}
```

Event used in the above example (stored locally in `"./events/apigateway-authorizer.json"`)
```
{
  "type": "TOKEN",
  "authorizationToken": "incoming-client-token",
  "methodArn": "arn:aws:execute-api:us-east-1:123456789012:example/prod/POST/{proxy+}"
}
```

### AWS SDK References

- [_LAMBDA_SERVER_PORT Environment Variable](https://github.com/aws/aws-lambda-go/blob/bc1ec47cb1670c0d5eaca47c10d89789d8507c3d/lambda/rpc.go#L16)
- [InvokeRequest/InvokeResponse](https://github.com/aws/aws-lambda-go/blob/bc1ec47cb1670c0d5eaca47c10d89789d8507c3d/lambda/messages/messages.go#L20)

### Help

The following arguments can be passed to the CLI to invoke the help documentation:

* -h

* --h

* -help

* --help
