package main

import (
	"fmt"

	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		for i := 0; i < 10; i++ {
			apiProduct, err := apiextensions.NewCustomResource(
				ctx,
				fmt.Sprintf("auto-%04d", i),
				&apiextensions.CustomResourceArgs{
					ApiVersion: pulumi.String("apigee.dev/v1alpha1"),
					Kind:       pulumi.String("ApiProduct"),
					Metadata: metav1.ObjectMetaArgs{
						Namespace: pulumi.String("registry"),
					},
					OtherFields: kubernetes.UntypedArgs{
						"spec": kubernetes.UntypedArgs{
							"title":       fmt.Sprintf("My API %d", i),
							"description": "This is my API.",
						},
					},
				},
			)
			if err != nil {
				return err
			}
			ctx.Export("name", apiProduct.Metadata.Name())
		}
		return nil
	})
}
