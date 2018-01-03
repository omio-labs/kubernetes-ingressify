## Templating with `kubernetes-ingressify`

For templating we use the [golang template system](https://golang.org/pkg/text/template/). We enrich the set of available
functions by adding [sprig](http://masterminds.github.io/sprig/) and the following two functions:

- GroupByHost: returns a `map[string]IngressifyRule` grouping ingressify rules by host as key
- GroupByPath: returns a `map[string]IngressifyRule` grouping ingressify rules by path as key
- GroupBySvcNs: returns a `map[string]IngressifyRule` grouping ingressify rules by name key that is concatenation result  of the ServiceName and Namespace

## What data is available when rendering a template ?

We provide all information that you get when you call `kubectl get ingress --all-namespaces -o yaml` but we choose to
denormalize all the nested output into an `IngressifyRule` that contains:

- ServiceName
- ServicePort
- Host
- Path
- Namespace
- Name
- IngressRaw

the last one is the plain [ingress rule](https://godoc.org/k8s.io/api/extensions/v1beta1#Ingress) modeled by the
official kubernetes client.

## Examples

check out `nginx.tmpl` and `haproxy.tmpl` and run them with:


`kunernetes-ingressify -config config_for_nginx.yaml`

and

`kunernetes-ingressify -config config_for_haproxy.yaml`

respectively.
