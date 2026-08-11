
	// A configuration can be adopted either by its ARN or by its name. MSK
	// configuration ARNs embed a service-generated UUID, so they are not knowable
	// ahead of creation; adopting by name lets a GitOps workflow declare the
	// resource without first discovering the ARN. sdkFind resolves the name to an
	// ARN.
	//
	// When "arn" is present we fall through to the generated handling below.
	if _, ok := fields["arn"]; !ok {
		name, ok := fields["name"]
		if !ok {
			return ackerrors.NewTerminalError(fmt.Errorf("required field missing: one of arn or name"))
		}
		r.ko.Spec.Name = &name
		return nil
	}
