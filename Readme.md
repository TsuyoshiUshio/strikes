# Lightning Strikes

Package manager for Azure Functions.

[![Build Status](https://6packdevops.visualstudio.com/StrikesRepository/_apis/build/status/TsuyoshiUshio.strikes)](https://6packdevops.visualstudio.com/StrikesRepository/_build/latest?definitionId=16)

# Motivation

When I want to share the Azure Function code, it usualy require to set up other resources, like CosmosDB, Storage Account, EventHubs and so on. This project is inspired by [helm](https://docs.helm.sh/) which is the popular package manager for kubernetes. I'd like to do the same thing for Azure Functions.

# Usage

NOTE: This repo is under construction. Now I'm start working it. 

I'd like to install everything in just one command. 

## initialize

These command install / Setup Lightning Strikes.

```
strikes init 
```

## install package

```
strikes install stable/drabeState
```

These command install / Setup Lightning Strikes. Then install `stable/durableState`  package which includes Azure Functions, Storage Account, Function App, and CosmosDB with all configured. 

Also, you can host your own package

## packaging 

```
strikes package --name powerplant
```

## push  

```
strikes push --name powerplant
```

It will be package your code into new package called `powerplant` then push into the repository. Now everyone can get the environment like this.

```
strikes install stable/powerplant
```

## package list 

You can also search packages for Azure Functions.

```
strikes package --list

stable/sendgrid 1.2.0
stable/migrator 2.1.0
   :
```

This is experimental. The command/subcommand might be changed in the future.

