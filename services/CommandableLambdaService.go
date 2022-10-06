package services

import (
	"context"

	ccomands "github.com/pip-services3-gox/pip-services3-commons-gox/commands"
	crun "github.com/pip-services3-gox/pip-services3-commons-gox/run"
)

// Abstract service that receives commands via AWS Lambda protocol
// to operations automatically generated for commands defined in ICommandable components.
// Each command is exposed as invoke method that receives command name and parameters.
//
// Commandable services require only 3 lines of code to implement a robust external
// Lambda-based remote interface.
//
// This service is intended to work inside LambdaFunction container that
// exploses registered actions externally.
//
// Configuration parameters
//
// - dependencies:
//   - controller:            override for Controller dependency
//
// References
//
//	- *:logger:*:*:1.0               (optional) ILogger components to pass log messages
//	- *:counters:*:*:1.0             (optional) ICounters components to pass collected measurements
//
// See CommandableLambdaClient
// See LambdaService
//
// Example:
//
//	type MyCommandableLambdaService struct  {
//		*CommandableLambdaService
//	}
//
//	func NewMyCommandableLambdaService() *MyCommandableLambdaService {
//		c:= &MyCommandableLambdaService{
//			CommandableLambdaService: NewCommandableLambdaService("v1.service")
//		}
//		c.DependencyResolver.Put(context.Background(),
//			"controller",
//			cref.NewDescriptor("mygroup","controller","*","*","1.0")
//		)
//		return c
//	}
//
//	service := NewMyCommandableLambdaService();
//	service.SetReferences(context.Background(), NewReferencesFromTuples(
//	   NewDescriptor("mygroup","controller","default","default","1.0"), controller
//	))
//
//	service.Open(context.Background(),"123")
//	fmt.Println("The AWS Lambda 'v1.service' service is running")
//
type CommandableLambdaService struct {
	*LambdaService
	commandSet *ccomands.CommandSet
}

// Creates a new instance of the service.
// - name a service name.
func InheritCommandableLambdaService(overrides ILambdaServiceOverrides, name string) *CommandableLambdaService {
	c := &CommandableLambdaService{
		LambdaService: InheritLambdaService(overrides, name),
	}

	c.DependencyResolver.Put(context.Background(), "controller", "none")
	return c
}

//   Registers all actions in AWS Lambda function.
func (c *CommandableLambdaService) Register() {
	resCtrl, depErr := c.DependencyResolver.GetOneRequired("controller")
	if depErr != nil {
		return
	}
	controller, ok := resCtrl.(ccomands.ICommandable)
	if !ok {
		c.Logger.Error(context.Background(), "CommandableHttpService", nil, "Can't cast Controller to ICommandable")
		return
	}

	c.commandSet = controller.GetCommandSet()

	commands := c.commandSet.Commands()
	for index := 0; index < len(commands); index++ {
		command := commands[index]
		name := command.Name()
		c.RegisterAction(name, nil, func(ctx context.Context, params map[string]any) (interface{}, error) {
			correlationId, _ := params["correlation_id"].(string)

			args := crun.NewParametersFromValue(params)
			args.Remove("correlation_id")

			timing := c.Instrument(ctx, correlationId, name)
			result, err := command.Execute(ctx, correlationId, args)
			timing.EndTiming(ctx, err)
			return result, err

		})
	}
}
