package util

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CreateServiceAccountWithToken creates a service account with a given name with a service account token.
// Need to use this function to simulate the actual behavior of Kubernetes API server with the fake client.
func CreateServiceAccountWithToken(clientset kubernetes.Interface, namespace, name, tokenName string) (*corev1.ServiceAccount, error) {
	sa, err := clientset.CoreV1().ServiceAccounts(namespace).Create(&corev1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: name}})
	if err != nil {
		return nil, err
	}
	token, err := clientset.CoreV1().Secrets(namespace).Create(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: tokenName,
			Annotations: map[string]string{
				corev1.ServiceAccountNameKey: sa.Name,
				corev1.ServiceAccountUIDKey:  string(sa.UID),
			}}, Type: corev1.SecretTypeServiceAccountToken})
	if err != nil {
		return nil, err
	}
	sa.Secrets = []corev1.ObjectReference{{Name: token.Name}}
	return clientset.CoreV1().ServiceAccounts(namespace).Update(sa)
}
