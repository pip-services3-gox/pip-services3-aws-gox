package clients

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	awscon "github.com/pip-services3-gox/pip-services3-aws-gox/connect"
	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	cerr "github.com/pip-services3-gox/pip-services3-commons-gox/errors"
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	ccount "github.com/pip-services3-gox/pip-services3-components-gox/count"
	clog "github.com/pip-services3-gox/pip-services3-components-gox/log"
	ctrace "github.com/pip-services3-gox/pip-services3-components-gox/trace"
	rpcsrv "github.com/pip-services3-gox/pip-services3-rpc-gox/services"
)

// Abstract client that calls AWS Lambda Functions.
//
// When making calls "cmd" parameter determines which what action shall be called, while
// other parameters are passed to the action itself.
//
// Configuration parameters:
//
//  - connections:
//      - discovery_key:               (optional) a key to retrieve the connection from IDiscovery
//      - region:                      (optional) AWS region
//  - credentials:
//      - store_key:                   (optional) a key to retrieve the credentials from ICredentialStore
//      - access_id:                   AWS access/client id
//      - access_key:                  AWS access/client id
//  - options:
//      - connect_timeout:             (optional) connection timeout in milliseconds (default: 10 sec)
//
// References:
//
//  - \*:logger:\*:\*:1.0            (optional) ILogger components to pass log messages
//  - \*:counters:\*:\*:1.0          (optional) ICounters components to pass collected measurements
//  - \*:discovery:\*:\*:1.0         (optional) IDiscovery services to resolve connection
//  - \*:credential-store:\*:\*:1.0  (optional) Credential stores to resolve credentials
//
//  See LambdaFunction
//  See CommandableLambdaClient
//
// Example:
//
//	type MyLambdaClient struct  {
//		*LambdaClient
//		...
//	}
//
//	func (c* MyLambdaClient) getData(ctx context.Context, correlationId string, id string)(result MyData, err error){
//		timing := c.Instrument(ctx, correlationId, "myclient.get_data");
//		callRes, callErr := c.Call(ctx ,"get_data" correlationId, map[string]interface{ "id": id })
//		if callErr != nil {
//			return callErr
//		}
//		defer timing.EndTiming(ctx, nil)
//		return awsclient.HandleLambdaResponse[*cdata.DataPage[MyData]](calValue)
//	}
//	...
//
//
//	client = NewMyLambdaClient();
//	client.Configure(context.Background(), NewConfigParamsFromTuples(
//	    "connection.region", "us-east-1",
//	    "connection.access_id", "XXXXXXXXXXX",
//	    "connection.access_key", "XXXXXXXXXXX",
//	    "connection.arn", "YYYYYYYYYYYYY"
//	))
//
//	data, err := client.GetData(context.Background(), "123", "1")
//	...
type LambdaClient struct {
	// The reference to AWS Lambda Function.
	Lambda *lambda.Lambda
	// The opened flag.
	Opened bool
	// The AWS connection parameters
	Connection     *awscon.AwsConnectionParams
	connectTimeout int
	// The dependencies resolver.
	DependencyResolver *cref.DependencyResolver
	// The connection resolver.
	ConnectionResolver *awscon.AwsConnectionResolver
	// The logger.
	Logger *clog.CompositeLogger
	//The performance counters.
	Counters *ccount.CompositeCounters
	// The tracer.
	Tracer *ctrace.CompositeTracer
}

func NewLambdaClient() *LambdaClient {
	c := &LambdaClient{
		Opened:             false,
		connectTimeout:     10000,
		DependencyResolver: cref.NewDependencyResolver(),
		ConnectionResolver: awscon.NewAwsConnectionResolver(),
		Logger:             clog.NewCompositeLogger(),
		Counters:           ccount.NewCompositeCounters(),
		Tracer:             ctrace.NewCompositeTracer(),
	}
	return c
}

// Configures component by passing configuration parameters.
//	Parameters:
//		- ctx context.Context	operation context.
//		- config    configuration parameters to be set.
func (c *LambdaClient) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.ConnectionResolver.Configure(ctx, config)
	c.DependencyResolver.Configure(ctx, config)
	c.connectTimeout = config.GetAsIntegerWithDefault("options.connect_timeout", c.connectTimeout)
}

// Sets references to dependent components.
//	Parameters:
//		- ctx context.Context	operation context.
//		- references	references to locate the component dependencies.
func (c *LambdaClient) SetReferences(ctx context.Context, references cref.IReferences) {
	c.Logger.SetReferences(ctx, references)
	c.Counters.SetReferences(ctx, references)
	c.ConnectionResolver.SetReferences(ctx, references)
	c.DependencyResolver.SetReferences(ctx, references)
}

// Instrument method are adds instrumentation to log calls and measure call time.
// It returns a services.InstrumentTiming object that is used to end the time measurement.
//	Parameters:
//		- ctx context.Context	operation context.
//		- correlationId string (optional) transaction id to trace execution through call chain.
//		- name string a method name.
//	Returns: services.InstrumentTiming object to end the time measurement.
func (c *LambdaClient) Instrument(ctx context.Context, correlationId string, name string) *rpcsrv.InstrumentTiming {
	c.Logger.Trace(ctx, correlationId, "Calling %s method", name)
	c.Counters.IncrementOne(ctx, name+".call_count")
	counterTiming := c.Counters.BeginTiming(ctx, name+".call_time")
	traceTiming := c.Tracer.BeginTrace(ctx, correlationId, name, "")
	return rpcsrv.NewInstrumentTiming(correlationId, name, "call",
		c.Logger, c.Counters, counterTiming, traceTiming)
}

//  Checks if the component is opened.
//  Returns true if the component has been opened and false otherwise.
func (c *LambdaClient) IsOpen() bool {
	return c.Opened
}

// Opens the component.
//	Parameters:
//		- ctx context.Context	operation context.
//		- correlationId 	(optional) transaction id to trace execution through call chain.
//		- Return 			 error or nil no errors occured.
func (c *LambdaClient) Open(ctx context.Context, correlationId string) error {
	if c.IsOpen() {
		return nil
	}

	wg := sync.WaitGroup{}
	var errGlobal error

	wg.Add(1)
	go func() {
		defer wg.Done()
		connection, err := c.ConnectionResolver.Resolve(correlationId)
		c.Connection = connection
		errGlobal = err

		awsCred := credentials.NewStaticCredentials(c.Connection.GetAccessId(), c.Connection.GetAccessKey(), "")
		sess := session.Must(session.NewSession(&aws.Config{
			MaxRetries:  aws.Int(3),
			Region:      aws.String(c.Connection.GetRegion()),
			Credentials: awsCred,
		}))
		// Create new cloudwatch client.
		c.Lambda = lambda.New(sess)
		c.Lambda.Config.HTTPClient.Timeout = time.Duration((int64)(c.connectTimeout)) * time.Millisecond
		c.Logger.Debug(ctx, correlationId, "Lambda client connected to %s", c.Connection.GetArn())

	}()
	wg.Wait()
	if errGlobal != nil {
		c.Opened = false
		return errGlobal
	}
	return nil
}

// Closes component and frees used resources.
//	Parameters:
//		- ctx context.Context	operation context.
//		- correlationId 	(optional) transaction id to trace execution through call chain.
//		- Returns 			 error or null no errors occured.
func (c *LambdaClient) Close(ctx context.Context, correlationId string) error {
	// Todo: close listening?
	c.Opened = false
	c.Lambda = nil
	return nil
}

// Performs AWS Lambda Function invocation.
//	Parameters:
//		- ctx context.Context	operation context.
//		- invocationType    an invocation type: "RequestResponse" or "Event"
//		- cmd               an action name to be called.
//		- correlationId 	(optional) transaction id to trace execution through call chain.
//		- args              action arguments
// Returns           result or error.
func (c *LambdaClient) Invoke(ctx context.Context, invocationType string, cmd string, correlationId string, args map[string]any) (result *lambda.InvokeOutput, err error) {

	if cmd == "" {
		err = cerr.NewUnknownError("", "NO_COMMAND", "Missing cmd")
		c.Logger.Error(ctx, correlationId, err, "Failed to call %s", cmd)
		return nil, err
	}

	//args = _.clone(args)

	args["cmd"] = cmd
	if correlationId != "" {
		args["correlation_id"] = correlationId
	} else {
		args["correlation_id"] = cdata.IdGenerator.NextLong()
	}

	payloads, jsonErr := json.Marshal(args)

	if jsonErr != nil {
		c.Logger.Error(ctx, correlationId, jsonErr, "Failed to call %s", cmd)
		return nil, jsonErr
	}

	params := &lambda.InvokeInput{
		FunctionName:   aws.String(c.Connection.GetArn()),
		InvocationType: aws.String(invocationType),
		LogType:        aws.String("None"),
		Payload:        payloads,
	}

	data, lambdaErr := c.Lambda.InvokeWithContext(ctx, params)

	if lambdaErr != nil {
		err = cerr.NewInvocationError(
			correlationId,
			"CALL_FAILED",
			"Failed to invoke lambda function").WithCause(err)
		return nil, err
	}

	return data, nil
}

// Calls a AWS Lambda Function action.
//	Parameters:
//		- ctx context.Context	operation context.
//		- cmd               an action name to be called.
//		- correlationId     (optional) transaction id to trace execution through call chain.
//		- params            (optional) action parameters.
//		- Returns           result and error.
func (c *LambdaClient) Call(ctx context.Context, cmd string, correlationId string, params map[string]any) (result *lambda.InvokeOutput, err error) {
	return c.Invoke(ctx, "RequestResponse", cmd, correlationId, params)
}

// Calls a AWS Lambda Function action asynchronously without waiting for response.
//	Parameters:
//		- ctx context.Context	operation context.
//		- cmd               an action name to be called.
//		- correlationId     (optional) transaction id to trace execution through call chain.
//		- params            (optional) action parameters.
//		- Returns           error or null for success.
func (c *LambdaClient) CallOneWay(ctx context.Context, cmd string, correlationId string, params map[string]any) error {
	_, err := c.Invoke(ctx, "Event", cmd, correlationId, params)
	return err
}
