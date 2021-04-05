### Knative install 

- using already present isito service mesh 

- using [xip.io](xip.io) as magic DNS


<!-- In civo , when istio is installed, the external ip assigned is not 1 but 2 ! yes, this is indeed strange.

upon close examination of those 2 IPs, one is the IP of the pod and the other is the external IP respectively.

```bash
kubectl get svc istio-ingressgateway -n istio-system
NAME                   TYPE           CLUSTER-IP    EXTERNAL-IP                 PORT(S)                                                                      AGE
istio-ingressgateway   LoadBalancer   10.43.25.88   192.168.1.3,212.2.241.254   15021:31931/TCP,80:32152/TCP,443:31128/TCP,15012:30451/TCP,15443:30038/TCP   25h
```
-->
<!-- ## why it gets 2 ip question was raised and awaiting answer from the sre team -->



now, onto setting DNS for knative serving , there are few options which are mentioned in the knative docs page 


 [https://knative.dev/docs/install/any-kubernetes-cluster/#serving_dns-0](https://knative.dev/docs/install/any-kubernetes-cluster/#serving_dns-0)

this will do the trick. 

Upon digging deeper in the docs , one can arrive to this doc page too 

[https://knative.dev/docs/serving/using-a-custom-domain/](https://knative.dev/docs/serving/using-a-custom-domain/)

so, here is the fix :

```bash
 kubectl patch configmap/config-domain -n knative-serving --type merge --patch '{"data":{"212.542.368.021.xip.io":""}}'
```


> ### Note:
>
>`212.542.368.021` this IP should be the external ip of the istio-ingressgateway on our cluster [this `212.542.368.021` is a arbitrary example]

<!-- 
references:

https://medium.com/@poy/custom-domains-on-knative-for-dev-deployments-8bb1fe4a8e43
 -->