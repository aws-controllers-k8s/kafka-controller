	if clusterDeleting(r) {
		return nil, requeueWaitWhileDeleting
	}
	if clusterCreating(r) {
		return nil, requeueWaitWhileCreating
	}
