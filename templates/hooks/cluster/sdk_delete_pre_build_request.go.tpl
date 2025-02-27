	if clusterDeleting(r) {
		return nil, requeueWaitWhileDeleting
	}
	if clusterCreating(r) {
		return nil, requeueWaitWhileCreating
	}
	groupCpy := r.ko.DeepCopy()
	groupCpy.Spec.AssociatedSCRAMSecrets = nil
	if err := rm.syncAssociatedScramSecrets(ctx, &resource{ko: groupCpy}, r); err != nil {
		return nil, err
	}
	if !clusterActive(r) {
		// doing this to avoid BadRequestException
		return r, ackrequeue.NeededAfter(
			fmt.Errorf("waiting for cluster to be active before deletion"),
			ackrequeue.DefaultRequeueAfterDuration,
		)
	}
