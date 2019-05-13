package bootstrap

import (
	"fmt"
	"mobingi/ocean/pkg/constants"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	bootstrapapi "k8s.io/cluster-bootstrap/token/api"
	bootstraputil "k8s.io/cluster-bootstrap/token/util"
)

// BootstrapToken describes one bootstrap token, stored as a Secret in the cluster
// TODO: The BootstrapToken object should move out to either k8s.io/client-go or k8s.io/api in the future
// (probably as part of Bootstrap Tokens going GA). It should not be staged under the kubeadm API as it is now.
type BootstrapToken struct {
	// Token is used for establishing bidirectional trust between nodes and control-planes.
	// Used for joining nodes in the cluster.
	Token *BootstrapTokenString
	// Description sets a human-friendly message why this token exists and what it's used
	// for, so other administrators can know its purpose.
	Description string
	// TTL defines the time to live for this token. Defaults to 24h.
	// Expires and TTL are mutually exclusive.
	TTL *metav1.Duration
	// Expires specifies the timestamp when this token expires. Defaults to being set
	// dynamically at runtime based on the TTL. Expires and TTL are mutually exclusive.
	Expires *metav1.Time
	// Usages describes the ways in which this token can be used. Can by default be used
	// for establishing bidirectional trust, but that can be changed here.
	Usages []string
	// Groups specifies the extra groups that this token will authenticate as when/if
	// used for authentication
	Groups []string
}

func NewBootstrapToken() (*BootstrapToken, error) {
	tokenString, err := bootstraputil.GenerateBootstrapToken()
	if err != nil {
		return nil, err
	}

	bootstrapTokenString, err := NewBootstrapTokenString(tokenString)
	if err != nil {
		return nil, err
	}

	return &BootstrapToken{
		Token:  bootstrapTokenString,
		Usages: []string{"authentication", "signing"},
		Groups: []string{constants.NodeBootstrapTokenAuthGroup},
	}, nil
}

func (bt *BootstrapToken) ToSecret() *v1.Secret {
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s%s", bootstrapapi.BootstrapTokenSecretPrefix, bt.Token.ID),
			Namespace: metav1.NamespaceSystem,
		},
		Type: v1.SecretType(bootstrapapi.SecretTypeBootstrapToken),
		Data: bt.encodeSecretData(time.Now()),
	}
}

func (bt *BootstrapToken) encodeSecretData(now time.Time) map[string][]byte {
	data := map[string][]byte{
		bootstrapapi.BootstrapTokenIDKey:     []byte(bt.Token.ID),
		bootstrapapi.BootstrapTokenSecretKey: []byte(bt.Token.Secret),
	}

	if len(bt.Description) > 0 {
		data[bootstrapapi.BootstrapTokenDescriptionKey] = []byte(bt.Description)
	}

	// If for some strange reason both token.TTL and token.Expires would be set
	// (they are mutually exlusive in validation so this shouldn't be the case),
	// token.Expires has higher priority, as can be seen in the logic here.
	if bt.Expires != nil {
		// Format the expiration date accordingly
		// TODO: This maybe should be a helper function in bootstraputil?
		expirationString := bt.Expires.Time.Format(time.RFC3339)
		data[bootstrapapi.BootstrapTokenExpirationKey] = []byte(expirationString)

	} else if bt.TTL != nil && bt.TTL.Duration > 0 {
		// Only if .Expires is unset, TTL might have an effect
		// Get the current time, add the specified duration, and format it accordingly
		expirationString := now.Add(bt.TTL.Duration).Format(time.RFC3339)
		data[bootstrapapi.BootstrapTokenExpirationKey] = []byte(expirationString)
	}

	for _, usage := range bt.Usages {
		data[bootstrapapi.BootstrapTokenUsagePrefix+usage] = []byte("true")
	}

	if len(bt.Groups) > 0 {
		data[bootstrapapi.BootstrapTokenExtraGroupsKey] = []byte(strings.Join(bt.Groups, ","))
	}
	return data
}

func CreateSecret(client clientset.Interface, secret *v1.Secret) error {
	if _, err := client.CoreV1().Secrets(secret.ObjectMeta.Namespace).Create(secret); err != nil {
		return err
	}

	return nil
}

type BootstrapTokenString struct {
	ID     string
	Secret string
}

func NewBootstrapTokenString(token string) (*BootstrapTokenString, error) {
	substrs := bootstraputil.BootstrapTokenRegexp.FindStringSubmatch(token)
	// TODO: Add a constant for the 3 value here, and explain better why it's needed (other than because how the regexp parsin works)
	if len(substrs) != 3 {
		return nil, fmt.Errorf("token string is not valid:%s", token)
	}

	return &BootstrapTokenString{ID: substrs[1], Secret: substrs[2]}, nil
}

func (bts BootstrapTokenString) String() string {
	if len(bts.ID) > 0 && len(bts.Secret) > 0 {
		return fmt.Sprintf("%s.%s", bts.ID, bts.Secret)
	}

	return ""
}
