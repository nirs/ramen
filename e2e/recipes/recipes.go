// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

package recipes

import (
	"fmt"

	recipe "github.com/ramendr/recipe/api/v1alpha1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	"github.com/ramendr/ramen/e2e/config"
	"github.com/ramendr/ramen/e2e/types"
)

func Create(ctx types.TestContext, recipeConfig *config.Recipe) error {
	log := ctx.Logger()

	recipe := prepareRecipe(ctx, recipeConfig)
	if err := createRecipeOnClusters(ctx, recipe); err != nil {
		return err
	}

	log.Debugf("Recipe %s created for discovered app \"%s/%s\"",
		ctx.Name(), ctx.AppNamespace(), ctx.Workload().GetAppName())

	return nil
}

func Delete(ctx types.TestContext) error {
	key := k8stypes.NamespacedName{
		Name:      ctx.Name(),
		Namespace: ctx.AppNamespace(),
	}

	return deleteRecipeOnClusters(ctx, key)
}

func prepareRecipe(ctx types.TestContext, recipeConfig *config.Recipe) *recipe.Recipe {
	recipe := prepareBaseRecipe(ctx)
	ns := ctx.AppNamespace()

	if recipeConfig.CheckHook || recipeConfig.ExecHook {
		recipe.Spec.Hooks = prepareHooks(ns)
	}

	recipe.Spec.Groups = prepareGroups(ns)
	recipe.Spec.Workflows = prepareWorkflows(recipeConfig)

	return recipe
}

func createRecipeOnClusters(ctx types.TestContext, recipe *recipe.Recipe) error {
	rCopy := recipe.DeepCopy()
	for _, cluster := range []*types.Cluster{ctx.Env().C1, ctx.Env().C2} {
		if err := cluster.Client.Create(ctx.Context(), recipe); err != nil {
			return fmt.Errorf("failed to create recipe \"%s/%s\" on cluster %q",
				ctx.AppNamespace(), recipe.Name, cluster.Name)
		}

		recipe = rCopy
	}

	return nil
}

func deleteRecipeOnClusters(ctx types.TestContext, key k8stypes.NamespacedName) error {
	for _, cluster := range []*types.Cluster{ctx.Env().C1, ctx.Env().C2} {
		r := &recipe.Recipe{}

		err := cluster.Client.Get(ctx.Context(), key, r)
		if err != nil {
			if !k8serrors.IsNotFound(err) {
				return err
			}

			return nil
		}

		if err := cluster.Client.Delete(ctx.Context(), r); err != nil {
			return err
		}
	}

	return nil
}

func prepareBaseRecipe(ctx types.TestContext) *recipe.Recipe {
	return &recipe.Recipe{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Recipe",
			APIVersion: "recipe.ramendr.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ctx.Name(),
			Namespace: ctx.AppNamespace(),
		},
	}
}

func prepareGroups(namespace string) []*recipe.Group {
	return []*recipe.Group{
		{
			Name:      "rg1",
			Type:      "resource",
			BackupRef: "rg1",
			IncludedNamespaces: []string{
				namespace,
			},
			LabelSelector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "appname",
						Operator: metav1.LabelSelectorOpIn,
						Values:   []string{"busybox"},
					},
				},
			},
		},
	}
}

// nolint:mnd
func prepareHooks(namespace string) []*recipe.Hook {
	return []*recipe.Hook{
		{
			Name:           "check-hook",
			Type:           "check",
			Namespace:      namespace,
			NameSelector:   "busybox",
			SelectResource: "deployment",
			Timeout:        300,
			Chks: []*recipe.Check{
				{
					Name:      "check-replicas",
					Condition: "{$.spec.replicas} == {$.status.readyReplicas}",
				},
			},
		},
		{
			Name:           "exec-hook",
			Type:           "exec",
			Namespace:      namespace,
			NameSelector:   "busybox",
			SelectResource: "deployment",
			Timeout:        300,
			Ops: []*recipe.Operation{
				{
					Name:    "ls",
					Command: "/bin/sh -c ls",
				},
			},
		},
	}
}

func prepareWorkflows(recipeSpec *config.Recipe) []*recipe.Workflow {
	backup := &recipe.Workflow{Name: "backup"}
	restore := &recipe.Workflow{Name: "restore"}

	checkHook := map[string]string{"hook": "check-hook/check-replicas"}
	execHook := map[string]string{"hook": "exec-hook/ls"}
	group := map[string]string{"group": "rg1"}

	if recipeSpec.ExecHook && recipeSpec.CheckHook {
		backup.Sequence = []map[string]string{checkHook, execHook, group}
		restore.Sequence = []map[string]string{group, checkHook, execHook}
	} else if recipeSpec.CheckHook {
		backup.Sequence = []map[string]string{checkHook, group}
		restore.Sequence = []map[string]string{group, checkHook}
	} else if recipeSpec.ExecHook {
		backup.Sequence = []map[string]string{execHook, group}
		restore.Sequence = []map[string]string{group, execHook}
	} else {
		backup.Sequence = []map[string]string{group}
		restore.Sequence = []map[string]string{group}
	}

	return []*recipe.Workflow{backup, restore}
}
