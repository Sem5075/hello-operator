package controllers
func (r *HelloReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var hello demov1.Hello

	if err := r.Get(ctx, req.NamespacedName, &hello); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "hello-" + hello.Name,
			Namespace: hello.Namespace,
		},
		Data: map[string]string{
			"message": hello.Spec.Message,
		},
	}

	err := ctrl.SetControllerReference(&hello, configMap, r.Scheme)
	if err != nil {
		return ctrl.Result{}, err
	}

	var existing corev1.ConfigMap
	err = r.Get(ctx, types.NamespacedName{
		Name:      configMap.Name,
		Namespace: configMap.Namespace,
	}, &existing)

	if err != nil {
		return ctrl.Result{}, r.Create(ctx, configMap)
	}

	existing.Data = configMap.Data
	return ctrl.Result{}, r.Update(ctx, &existing)
}
