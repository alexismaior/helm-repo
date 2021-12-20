{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "template.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "template.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "template.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- define "template.build" -}}
{{- printf "%s" .Chart.Version | regexFind "[^.]+$" -}}
{{- end -}}


{{/*
Common labels
*/}}
{{- define "template.labels" -}}
helm.chart: {{ include "template.chart" . }}
{{ include "template.selectorLabels" . }}
app.managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "template.selectorLabels" -}}
app: {{ include "template.name" . }}
app.instance: {{ .Release.Name }}
{{- end -}}
{{- define "template.selectorLabelsDaemonset" -}}
app: {{ include "template.name" . }}-daemonset
app.instance: {{ .Release.Name }}-daemonset
{{- end -}}

{{/*
Pod labels
*/}}
{{- define "template.podLabels" -}}
{{ include "template.selectorLabels" . }}
version: {{ include "template.build" . | quote  }}
{{- end -}}
{{- define "template.podLabelsDaemonset" -}}
{{ include "template.selectorLabelsDaemonset" . }}
version: {{ include "template.build" . | quote }}
{{- end -}}


{{/*
Create the name of the service account to use
*/}}
{{- define "template.serviceAccountName" -}}
{{$serviceAccount := default dict .Values.serviceAccount}}
{{- if $serviceAccount.create -}}
    {{ default (include "template.fullname" .) $serviceAccount.name }}
{{- else -}}
    {{ default "default" $serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the service to use
*/}}
{{- define "template.serviceName" -}}
    {{ default (include "template.fullname" .) .Values.service.name }}
{{- end -}}

{{/*
Define envFrom if image.name ends with "-ui" (SIOP frontend services)
or set the default otherwise (SIOP backend services)
*/}}
{{- define "template.envFrom" -}}
{{- if eq (.Values.image.name | regexFind "-ui") "-ui"}}
- secretRef:
    name: secrets-frontend
{{- else -}}
- configMapRef:
    name: configmap
- secretRef:
    name: secrets-backend
{{- end -}}
{{- end -}}
