package etcd

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/kubernetes/pkg/printers"
	printerstorage "k8s.io/kubernetes/pkg/printers/storage"

	"github.com/openshift/api/oauth"

	oauthapi "github.com/openshift/openshift-apiserver/pkg/oauth/apis/oauth"
	"github.com/openshift/openshift-apiserver/pkg/oauth/apiserver/registry/oauthclient"
	oauthprinters "github.com/openshift/openshift-apiserver/pkg/oauth/printers/internalversion"
)

// rest implements a RESTStorage for oauth clients against etcd
type REST struct {
	*registry.Store
}

var _ rest.StandardStorage = &REST{}

// NewREST returns a RESTStorage object that will work against oauth clients
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &registry.Store{
		NewFunc:                  func() runtime.Object { return &oauthapi.OAuthClient{} },
		NewListFunc:              func() runtime.Object { return &oauthapi.OAuthClientList{} },
		DefaultQualifiedResource: oauth.Resource("oauthclients"),

		TableConvertor: printerstorage.TableConvertor{TableGenerator: printers.NewTableGenerator().With(oauthprinters.AddOAuthOpenShiftHandler)},

		CreateStrategy: oauthclient.Strategy,
		UpdateStrategy: oauthclient.Strategy,
		DeleteStrategy: oauthclient.Strategy,
	}

	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	return &REST{store}, nil
}
