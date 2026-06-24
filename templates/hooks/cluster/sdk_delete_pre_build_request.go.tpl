	if clusterDeleting(r) {
		return nil, requeueWaitWhileDeleting
	}
	if clusterCreating(r) {
		return nil, requeueWaitWhileCreating
	}
	// When SCRAM secret associations are managed externally, the controller
	// must not disassociate any secret on delete either - association
	// management is entirely the customer's responsibility. Skipping this is
	// safe because MSK automatically disassociates SCRAM secrets when the
	// cluster is deleted.
	if !isSCRAMSecretsManagedExternally(r.ko) {
		groupCpy := r.ko.DeepCopy()
		groupCpy.Spec.AssociatedSCRAMSecrets = nil
		if err := rm.syncAssociatedScramSecrets(ctx, &resource{ko: groupCpy}, r); err != nil {
			return nil, err
		}
	}
	if !clusterActive(r) {
		// doing this to avoid BadRequestException
		return r, ackrequeue.NeededAfter(
			fmt.Errorf("waiting for cluster to be active before deletion"),
			ackrequeue.DefaultRequeueAfterDuration,
		)
	}
