	if resp.ClusterInfo != nil &&
		resp.ClusterInfo.Provisioned != nil &&
		resp.ClusterInfo.Provisioned.CurrentBrokerSoftwareInfo != nil {
		if resp.ClusterInfo.Provisioned.CurrentBrokerSoftwareInfo.KafkaVersion != nil {
			ko.Spec.Provisioned.KafkaVersion = resp.ClusterInfo.Provisioned.CurrentBrokerSoftwareInfo.KafkaVersion
		}
		if resp.ClusterInfo.Provisioned.CurrentBrokerSoftwareInfo.ConfigurationArn != nil &&
			resp.ClusterInfo.Provisioned.CurrentBrokerSoftwareInfo.ConfigurationRevision != nil {
			ko.Spec.Provisioned.ConfigurationInfo = &svcapitypes.ConfigurationInfo{
				ARN:      resp.ClusterInfo.Provisioned.CurrentBrokerSoftwareInfo.ConfigurationArn,
				Revision: resp.ClusterInfo.Provisioned.CurrentBrokerSoftwareInfo.ConfigurationRevision,
			}
		} else {
			ko.Spec.Provisioned.ConfigurationInfo = nil
		}
	}
	if !serverlessClusterActive(&resource{ko}) {
		// Setting resource synced condition to false will trigger a requeue of
		// the resource. No need to return a requeue error here.
		ackcondition.SetSynced(&resource{ko}, corev1.ConditionFalse, nil, nil)
		return &resource{ko}, nil
	}

	// Unprovisioned Clusters do not use Scram Secrets
	if ko.Spec.Provisioned != nil {
		ko.Spec.AssociatedSCRAMSecrets, err = rm.getAssociatedScramSecrets(ctx, &resource{ko})
		if err != nil {
			return nil, err
		}
	}
