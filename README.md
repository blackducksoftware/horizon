# Horizon

Horizon is an application framwork for horizontal scaling, monitoring, and terse descriptions of complex, distributed systems running in a Cloud Native environment.

# Why another kubernetes application framework ? 

Horizon was built based on ideas from the CoreOS operators framework, the Coki (short) framework for succinct kubernetes app definitions, and blackduck's 'protoform' framework for imperative deployment of the perceptor platform.

It doesn't represent a competitor to any of the above platforms, but rather, an opinionated cloud native deployment platform which is aligned with a proscriptive model for how one would define, and maintain, a distributed application in a kubernetes or openshift environment.

# Why would I use Horizon instead of YAML, Helm, ...

- Horizon apps replaces the need to create and provide yaml or json files for the application.   
- Horizon apps are unit testable.
- Horizon apps embed operator like semantics, from the beggining, so they don't need to be layered after the fact.
- Horizon apps give you a deployment mechanism that are human readable, without redundant schema, formatting.
- Horizon apps are embeddalbe as a library.
- Horizon apps ship with a prometheus implementation (this is coming soon).
- Horizon apps enforce an idiom wherein you to ship exactly one container, and only one container.
- Horizon apps separate the *creation* of your app from the *maintainance* of it over time, rather then forcing you to conform to watch semantics (as a raw controller or operator would do).

# How does it work ?

- Vendor horizon into your golang app.
- Programmatically define an application that will be deployed to a kubernetes based cluster using the objects we provide, using the objects in this repository.  
- Start the horizon deployer from inside your own golang based container.
- Optionally send a controller that will be run after all defined parts of the application are deployed succesfully.

## License

Apache License 2.0
