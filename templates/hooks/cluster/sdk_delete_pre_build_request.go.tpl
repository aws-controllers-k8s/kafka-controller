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
