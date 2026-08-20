
	// DescribeConfiguration is keyed on the configuration ARN, which is only
	// known after creation. When the owner opted into adoption but no ARN is set
	// -- an adoption by name, or an adopt-or-create with no adoption-fields
	// annotation -- resolve the ARN from the name so the read below can proceed.
	// Configuration names are unique per account and region.
	if rm.shouldLookupARN(r) {
		r = rm.concreteResource(r.DeepCopy())
		configurationARN, err := rm.getConfigurationARN(ctx, r)
		if err != nil {
			return nil, err
		}
		if configurationARN != nil {
			if r.ko.Status.ACKResourceMetadata == nil {
				r.ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
			}
			arn := ackv1alpha1.AWSResourceName(*configurationARN)
			r.ko.Status.ACKResourceMetadata.ARN = &arn
		}
	}
