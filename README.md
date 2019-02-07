# Horizon

Horizon is an application framework for horizontal scaling, monitoring, and terse descriptions of complex, distributed systems running in a Cloud Native environment.

# Why another Kubernetes application framework? 

Horizon was built based on ideas from the CoreOS operators framework, the Koki (short) framework for succinct Kubernetes app definitions, and Black Duck's 'protoform' framework for imperative deployment of the Perceptor platform.

It doesn't represent a competitor to any of the above platforms, but rather, an opinionated cloud-native deployment platform that is aligned with a proscriptive model for how one would define, and maintain, a distributed application in a Kubernetes or OpenShift environment.

# Why would I use Horizon instead of YAML or Helm?

- Horizon apps replace the need to create and provide YAML or JSON files for the application.   
- Horizon apps are unit testable.
- Horizon apps embed operator-like semantics, from the beggining, so they don't need to be layered after the fact.
- Horizon apps give you a deployment mechanism that is human readable, without redundant schema formatting.
- Horizon apps are embeddable as a library.
- Horizon apps ship with a Prometheus implementation (coming soon).
- Horizon apps enforce an idiom wherein you ship exactly one container, and only one container.
- Horizon apps separate the *creation* of your app from the *maintainance* of it over time, rather then forcing you to conform to watch semantics (as a raw controller or operator would do).
- Horizon apps can be used to quickly put together an *operator* which can be deployed via Helm, so, this entire bulleted list is a false trichotimy :).

# How do I use it?

- Vendor Horizon into your Golang app.
- Programmatically define an application that will be deployed to a Kubernetes-based cluster using the objects we provide, using the objects in this repository.  
- Start the Horizon deployer from inside your own Golang-based container.
- Optionally send a controller that will be run after all defined parts of the application are deployed successfully.

## License

Apache License 2.0
