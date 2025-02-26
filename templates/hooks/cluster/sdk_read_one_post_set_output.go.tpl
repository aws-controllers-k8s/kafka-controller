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
	}