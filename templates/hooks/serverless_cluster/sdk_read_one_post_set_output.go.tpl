	if resp.ClusterInfo != nil {
		if resp.ClusterInfo.Provisioned != nil &&
			resp.ClusterInfo.Provisioned.CurrentBrokerSoftwareInfo != nil &&
			resp.ClusterInfo.Provisioned.CurrentBrokerSoftwareInfo.KafkaVersion != nil {
			ko.Spec.Provisioned.KafkaVersion = resp.ClusterInfo.Provisioned.CurrentBrokerSoftwareInfo.KafkaVersion
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
