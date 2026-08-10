
	// DescribeClusterV2 is keyed on the cluster ARN, which is only known after
	// creation. When the ARN is absent but the user gave us a name -- an
	// adoption by name, or an adopt-or-create with no adoption-fields
	// annotation -- resolve the ARN from the name so the read below can
	// proceed. Cluster names are unique per account and region.
	r.ko = r.ko.DeepCopy()
	if rm.requiredFieldsMissingFromReadOneInput(r) {
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
