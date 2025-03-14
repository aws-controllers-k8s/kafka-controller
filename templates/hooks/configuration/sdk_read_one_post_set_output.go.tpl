	response, err := rm.sdkapi.DescribeConfigurationRevision(ctx, &svcsdk.DescribeConfigurationRevisionInput{
		Arn: (*string)(ko.Status.ACKResourceMetadata.ARN),
		Revision: ko.Status.LatestRevision.Revision,
	})
	rm.metrics.RecordAPICall("READ_ONE", "DescribeConfigurationRevision", err)
	if err != nil {
		return &resource{ko}, err
	}
	if response.ServerProperties != nil {
		ko.Spec.ServerProperties = response.ServerProperties
	}