
	// DescribeClusterV2 is keyed on the cluster ARN, which is only known after
	// creation. When the owner opted into adoption but no ARN is set -- an
	// adoption by name, or an adopt-or-create with no adoption-fields
	// annotation -- resolve the ARN from the name so the read below can
	// proceed. Cluster names are unique per account and region.
	if rm.shouldLookupARN(r) {
		r = rm.concreteResource(r.DeepCopy())
		clusterARN, err := rm.getClusterARN(ctx, r)
		if err != nil {
			return nil, err
		}
		if clusterARN != nil {
			if r.ko.Status.ACKResourceMetadata == nil {
				r.ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
			}
			arn := ackv1alpha1.AWSResourceName(*clusterARN)
			r.ko.Status.ACKResourceMetadata.ARN = &arn
		}
	}
