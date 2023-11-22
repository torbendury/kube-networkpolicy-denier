{{/*
kube-networkpolicy-denier.selfSignedCABundleCertPEM is the self-signed CA to:
- sign the certification key pair used by webhook server
*/}}
{{- define "kube-networkpolicy-denier.selfSignedCABundleCertPEM" -}}
  {{- $caKeypair := .selfSignedCAKeypair | default (genCA "kube-networkpolicy-denier-ca" 1825) -}}
  {{- $_ := set . "selfSignedCAKeypair" $caKeypair -}}
  {{- $caKeypair.Cert -}}
{{- end -}}

{{/*
Get the caBundle for clients of the webhooks.
It would use .selfSignedCAKeypair as the place to store the generated CA keypair, it is actually kind of dirty work to prevent generating keypair with multiple times.
When using this template, it requires the top-level scope.
*/}}
{{- define "webhook.caBundleCertPEM" -}}
    {{- /* Generate ca with CN "kube-networkpolicy-denier-ca" and 5 years validity duration if not exists in the current scope.*/ -}}
    {{- $caKeypair := .selfSignedCAKeypair | default (genCA "kube-networkpolicy-denier-ca" 1825) -}}
    {{- $_ := set . "selfSignedCAKeypair" $caKeypair -}}
    {{- $caKeypair.Cert -}}
{{- end -}}

{{/*
webhook.certPEM is the cert of certification used by validating/mutating admission webhook server.
Like generating CA, it would use .webhookTLSKeypair as the place to store the generated keypair, it is actually kind of dirty work to prevent generating keypair with multiple times.
When using this template, it requires the top-level scope
*/}}
{{- define "webhook.certPEM" -}}
    {{- /* FIXME: Duplicated codes with named template "webhook.keyPEM" because of no way to nested named template.*/ -}}
    {{- /* webhookName would be the FQDN of in-cluster service kube-networkpolicy-denier.*/ -}}
    {{- $webhookName := printf "%s.%s.svc" (include "kube-networkpolicy-denier.fullname" .) .Release.Namespace }}
    {{- $webhookCA := required "self-signed CA keypair is requried" .selfSignedCAKeypair -}}
    {{- /* Generate cert keypair for webhook with 5 year validity duration. */ -}}
    {{- $webhookServerTLSKeypair := .webhookTLSKeypair | default (genSignedCert $webhookName nil (list $webhookName) 1825 $webhookCA) }}
    {{- $_ := set . "webhookTLSKeypair" $webhookServerTLSKeypair -}}
    {{- $webhookServerTLSKeypair.Cert -}}
{{- end -}}

{{/*
webhook.keyPEM is the key of certification used by validating/mutating admission webhook server.
Like generating CA, it would use .webhookTLSKeypair as the place to store the generated keypair, it is actually kind of dirty work to prevent generating keypair with multiple times.

When using this template, it requires the top-level scope
*/}}
{{- define "webhook.keyPEM" -}}
    {{- /* FIXME: Duplicated codes with named template "webhook.keyPEM" because of no way to nested named template.*/ -}}
    {{- /* webhookName would be the FQDN of in-cluster service kube-networkpolicy-denier.*/ -}}
    {{- $webhookName := printf "%s.%s.svc" (include "kube-networkpolicy-denier.fullname" .) .Release.Namespace -}}
    {{- $webhookCA := required "self-signed CA keypair is requried" .selfSignedCAKeypair -}}
    {{- /* Generate cert key pair for webhook with 5 year validity duration. */ -}}
    {{- $webhookServerTLSKeypair := .webhookTLSKeypair | default (genSignedCert $webhookName nil (list $webhookName) 1825 $webhookCA) -}}
    {{- $_ := set . "webhookTLSKeypair" $webhookServerTLSKeypair -}}
    {{- $webhookServerTLSKeypair.Key -}}
{{- end -}}
