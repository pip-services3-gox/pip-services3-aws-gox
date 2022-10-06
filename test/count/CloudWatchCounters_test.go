package test

import (
	"context"
	"os"
	"testing"

	awscount "github.com/pip-services3-gox/pip-services3-aws-gox/count"
	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	cinfo "github.com/pip-services3-gox/pip-services3-components-gox/info"
)

func TestCloudWatchCounters(t *testing.T) {

	var counters *awscount.CloudWatchCounters
	var fixture *CountersFixture

	AWS_REGION := os.Getenv("AWS_REGION")
	AWS_ACCESS_ID := os.Getenv("AWS_ACCESS_ID")
	AWS_ACCESS_KEY := os.Getenv("AWS_ACCESS_KEY")

	if AWS_REGION == "" || AWS_ACCESS_ID == "" || AWS_ACCESS_KEY == "" {
		return
	}

	counters = awscount.NewCloudWatchCounters()
	fixture = NewCountersFixture(counters.CachedCounters)

	config := cconf.NewConfigParamsFromTuples(
		"interval", "5000",
		"connection.region", AWS_REGION,
		"credential.access_id", AWS_ACCESS_ID,
		"credential.access_key", AWS_ACCESS_KEY,
	)
	counters.Configure(context.Background(), config)

	contextInfo := cinfo.NewContextInfo()
	contextInfo.Name = "Test"
	contextInfo.Description = "This is a test container"

	var references = cref.NewReferencesFromTuples(context.Background(),
		cref.NewDescriptor("pip-services", "context-info", "default", "default", "1.0"), contextInfo,
		cref.NewDescriptor("pip-services", "counters", "cloudwatch", "default", "1.0"), counters,
	)
	counters.SetReferences(context.Background(), references)
	counters.Open(context.Background(), "")
	defer counters.Close(context.Background(), "")

	t.Run("Simple Counters", fixture.TestSimpleCounters)
	t.Run("Measure Elapsed Time", fixture.TestMeasureElapsedTime)
}
