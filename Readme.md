# Lightning Strikes

Package Management tool for Azure Functions

# Motivation

When I want to share the Azure Function code, it usualy require to set up other resources, like CosmosDB, Storage Account, EventHubs and so on.

# Usage

NOTE: This repo is under construction. Now I'm start working it. 

I'd like to install everything in just one command. 

```
strikes init 
strikes install stable/drabeState
```

These command install / Setup Lightning Strikes. Then install `stable/durableState`  package which includes Azure Functions, Storage Account, Function App, and CosmosDB with all configured. 

Also, you can host your own package

```
strikes package --name powerplant
strikes push --name powerplant
```

It will be package your code into new package called `powerplant` then push into the repository. Now everyone can get the environment like this.

```
strikes install stable/powerplant
```

This is experimental. The command/subcommand might be changed in the future.

