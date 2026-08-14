	if resp.ClusterInfo.CurrentBrokerSoftwareInfo != nil {
		if resp.ClusterInfo.CurrentBrokerSoftwareInfo.KafkaVersion != nil {
			ko.Spec.KafkaVersion = resp.ClusterInfo.CurrentBrokerSoftwareInfo.KafkaVersion
		}
		if resp.ClusterInfo.CurrentBrokerSoftwareInfo.ConfigurationArn != nil &&
			resp.ClusterInfo.CurrentBrokerSoftwareInfo.ConfigurationRevision != nil {
			ko.Spec.ConfigurationInfo = &svcapitypes.ConfigurationInfo{
				ARN:      resp.ClusterInfo.CurrentBrokerSoftwareInfo.ConfigurationArn,
				Revision: resp.ClusterInfo.CurrentBrokerSoftwareInfo.ConfigurationRevision,
			}
		} else {
			ko.Spec.ConfigurationInfo = nil
		}
	}
	if resp.ClusterInfo.CurrentVersion != nil {
		ko.Status.CurrentVersion = resp.ClusterInfo.CurrentVersion
	} else {
		ko.Status.CurrentVersion = nil
	}
	if !clusterActive(&resource{ko}) {
		// Setting resource synced condition to false will trigger a requeue of
		// the resource. No need to return a requeue error here.
		ackcondition.SetSynced(&resource{ko}, corev1.ConditionFalse, nil, nil)
	} else {
		ackcondition.SetSynced(&resource{ko}, corev1.ConditionTrue, nil, nil)
		ko.Spec.AssociatedSCRAMSecrets, err = rm.getAssociatedScramSecrets(ctx, &resource{ko})
		if err != nil {
			return nil, err
		}
		err = rm.setResourceAdditionalFields(ctx, ko)
		if err != nil {
			return nil, err
		}
		// Cluster resource policies are managed via a separate side-API
		// (GetClusterPolicy). A NotFoundException simply means no policy is
		// attached and must not surface as a ReadOne 404 for the Cluster.
		if err = rm.setClusterPolicy(ctx, ko); err != nil {
			return nil, err
		}
	}