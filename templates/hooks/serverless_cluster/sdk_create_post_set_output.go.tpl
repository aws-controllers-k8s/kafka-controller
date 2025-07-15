    if !serverlessClusterActive(&resource{ko}) {
        // This causes a requeue and scram secrets will be synced on the next
        // reconciliation loop
        ackcondition.SetSynced(&resource{ko}, corev1.ConditionFalse, nil, nil)
        return &resource{ko}, nil
    }