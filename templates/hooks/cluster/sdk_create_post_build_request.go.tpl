	// MSK rejects CreateCluster when any vpcConnectivity auth scheme is enabled
	// ('enabled' : true) with a terminal BadRequestException. VPC connectivity
	// auth must be created disabled and then enabled asynchronously via
	// UpdateConnectivity once the cluster is ACTIVE. Force every vpcConnectivity
	// auth flag on the built request to false while leaving desired.ko.Spec
	// untouched, so the post-ACTIVE reconcile still detects the enable delta.
	if input.BrokerNodeGroupInfo != nil &&
		input.BrokerNodeGroupInfo.ConnectivityInfo != nil &&
		input.BrokerNodeGroupInfo.ConnectivityInfo.VpcConnectivity != nil &&
		input.BrokerNodeGroupInfo.ConnectivityInfo.VpcConnectivity.ClientAuthentication != nil {
		ca := input.BrokerNodeGroupInfo.ConnectivityInfo.VpcConnectivity.ClientAuthentication
		if ca.Sasl != nil {
			if ca.Sasl.Iam != nil {
				ca.Sasl.Iam.Enabled = aws.Bool(false)
			}
			if ca.Sasl.Scram != nil {
				ca.Sasl.Scram.Enabled = aws.Bool(false)
			}
		}
		if ca.Tls != nil {
			ca.Tls.Enabled = aws.Bool(false)
		}
	}
