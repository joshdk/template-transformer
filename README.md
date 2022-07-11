[![License][license-badge]][license-link]
[![Actions][github-actions-badge]][github-actions-link]
[![Releases][github-release-badge]][github-release-link]

# Kustomize Template Transformer Plugin

🫥 Kustomize transformer plugin for strict templating of resources

## Motivations


Currently, there are two built-in ways to inject dynamic values into the manifests generated when running `kustomize build`.
The first is [kustomize vars](https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/vars/), which requires that the values already exist as part of a resource's properties, and is also currently planned to be deprecated.
The second is [kustomize replacements](https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/replacements/), which is the successor to vars, but has a number of usability concerns.

This repository provides a [kustomize plugin](https://kubectl.docs.kubernetes.io/guides/extending_kustomize/exec_plugins/) (specifically a transformer plugin) that can template values from the working environment into a stream of Kubernetes resource manifests.
This plugin also aims to facilitate this in a tightly controlled way in order to avoid runtime errors, user confusion, and to prevent any backsliding into open-ended template sprawl like with Helm.

## Installation

Prebuilt binaries for several architectures can be found attached to any of the available [releases][github-release-link].

```shell
$ wget https://github.com/joshdk/template-transformer/releases/download/v0.1.0/template-transformer-linux-amd64.tar.gz
$ tar -xf template-transformer-linux-amd64.tar.gz
$ mkdir -p ${XDG_CONFIG_HOME}/kustomize/plugin/jdk.sh/v1beta1/templatetransformer/
$ install TemplateTransformer ${XDG_CONFIG_HOME}/kustomize/plugin/jdk.sh/v1beta1/templatetransformer/TemplateTransformer
```

## Usage

To show how things are configured, we can demo adding a version label to an existing deployment manifest.

Inside our existing Kustomize app, we can create a `template-transformer.yaml` file.
Here we are configuring a single property named `VERSION` which gets its value from the `DRONE_COMMIT_SHA` environment variable.

The `apiVersion` and `kind` must be set to `jdk.sh/v1beta1` and `TemplateTransformer` respectively. 

```yaml
apiVersion: jdk.sh/v1beta1
kind: TemplateTransformer

metadata:
  name: example

properties:
  - name: VERSION
    description: Current application version
    source:
      - DRONE_COMMIT_SHA
```

In our example `kustomization.yaml` file, we can reference the `template-transformer.yaml` file under the `transformers` section:

```diff
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - deployment.yaml

transformers:
+ - template-transformer.yaml
```

In our existing (and very abridged) `deployment.yaml` manifest, we can add a new label that references our property using the template syntax `${{.VERSION}}`.

```diff
apiVersion: apps/v1
kind: Deployment

metadata:
  name: example
+ labels:
+   app.kubernetes.io/version: ${{.VERSION}}
```

We can now build our Kustomize app, making sure to set a value for the `DRONE_COMMIT_SHA` environment variable.

```bash
$ DRONE_COMMIT_SHA=$(git rev-parse HEAD) kustomize build --enable-alpha-plugins .
```

You should now see a rendered deployment labeled with the current git sha.

```yaml
apiVersion: apps/v1
kind: Deployment

metadata:
  name: example
  labels:
    app.kubernetes.io/version: deca...5bd9
```

### Default Values

If the template transformer can't resolve a value for every single property, it will fail immediately with an error.
To prevent errors in environments where an appropriate value might not be available, you can configure a default value instead.

```diff
apiVersion: jdk.sh/v1beta1
kind: TemplateTransformer

metadata:
  name: example

properties:
  - name: VERSION
    description: Current application version
    source:
      - DRONE_COMMIT_SHA
+   default: development
```

Now if we build our Kustomize app, without setting a value for the `DRONE_COMMIT_SHA` environment variable, we can observe that our default value is templated instead.

```bash
$ kustomize build --enable-alpha-plugins .
```

```yaml
apiVersion: apps/v1
kind: Deployment

metadata:
  name: example
  labels:
    app.kubernetes.io/version: development
```

### Mutating Values

Often the raw value of the source environment variable isn't appropriate for direct templating.
In cases like this you can configure a mutator to modify such values.

In this example we want to keep only the first 8 digits of our git SHA, instead of the entire thing.
Additionally, we want to prefix the value with `git-` since that is the source of our versioning.

> ℹ️ See the docs on [Regexp.Expand](https://pkg.go.dev/regexp#Regexp.Expand) for more information.

```diff
apiVersion: jdk.sh/v1beta1
kind: TemplateTransformer

metadata:
  name: example

properties:
  - name: VERSION
    description: Current application version
    source:
      - DRONE_COMMIT_SHA
    default: development
+   mutate:
+     pattern: ^.{8}
+     replace: git-$0
```

Now if we build our Kustomize app, we can observe the prefixed short SHA value in our deployment label.

```bash
$ DRONE_COMMIT_SHA=$(git rev-parse HEAD) kustomize build --enable-alpha-plugins .
```

```yaml
apiVersion: apps/v1
kind: Deployment

metadata:
  name: example
  labels:
    app.kubernetes.io/version: git-decafbad
```

### ArgoCD

If your are deploying Kustomize applications using ArgoCD, then please take note of the [ArgoCD build environment](https://argo-cd.readthedocs.io/en/stable/user-guide/build-environment/) as it contains a very limited set of environment variables.
Additional environment variables cannot be configured without completely rebuilding the built-in ArgoCD Kustomize integration. 

Our example `template-transformer.yaml` configuration could be updated to support ArgoCD like so.

```diff
apiVersion: jdk.sh/v1beta1
kind: TemplateTransformer

metadata:
  name: example

properties:
  - name: VERSION
    description: Current application version
    source:
      - DRONE_COMMIT_SHA
+     - ARGOCD_APP_REVISION
    default: development
    mutate:
      pattern: ^.{8}
      replace: git-$0
```

## License

This code is distributed under the [MIT License][license-link], see [LICENSE.txt][license-file] for more information.

[github-actions-badge]:  https://github.com/joshdk/template-transformer/workflows/Build/badge.svg
[github-actions-link]:   https://github.com/joshdk/template-transformer/actions
[github-release-badge]:  https://img.shields.io/github/release/joshdk/template-transformer/all.svg
[github-release-link]:   https://github.com/joshdk/template-transformer/releases
[license-badge]:         https://img.shields.io/badge/license-MIT-green.svg
[license-file]:          https://github.com/joshdk/template-transformer/blob/master/LICENSE.txt
[license-link]:          https://opensource.org/licenses/MIT
