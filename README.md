# ISCEnv
[![Latest Release](https://img.shields.io/github/release/ontariosystems/iscenv.svg)](https://github.com/ontariosystems/iscenv/releases)
[![CLA assistant](https://cla-assistant.io/readme/badge/ontariosystems/iscenv)](https://cla-assistant.io/ontariosystems/iscenv)
[![Build Status](https://travis-ci.org/ontariosystems/iscenv.svg?branch=master)](https://travis-ci.org/ontariosystems/iscenv)
[![Go Report Card](https://goreportcard.com/badge/github.com/ontariosystems/iscenv)](https://goreportcard.com/report/github.com/ontariosystems/iscenv)
[![GoDoc](https://godoc.org/github.com/ontariosystems/iscenv?status.svg)](https://godoc.org/github.com/ontariosystems/iscenv)

This utility allows the user to manage docker-based ISC product instances.  These instances are meant to be used as
development environments on a Linux host.

## Usage
The command consists of a single command line utility **`iscenv`**.  This utility is self-documenting.  Please refer
to its help for further documentation on each command.

## Caveats
- Do not add a default namespace to the root user on the instance, this will break features that wrap csession as it can no longer use the -U switch

## Known issues
- Orphaned plugins are occasionally left running after the primary process exits.  This prevents upgrades.  They can be killed by (after stopping all running iscenv containers) executing `killall iscenv`

## Future changes
- Rework plugin system to use [go plugins](https://golang.org/pkg/plugin/) ([simple example](https://jeremywho.com/go-1.8---plugins/))

## Rejected Features
- _Make all commands return the instance name rather than the container ID_ **We're using full logging now rather than specific items being printed to stdout**
- _Add Service wrappers for the containers which will restart them on reboot_ **We want iscenv to remain thin, users should do this themselves**
- _Have prep update the deployment service_ **Removing the deployment service from ISCEnv altogether**
- _Add a command that reads a simple configuration file to set up a specific environment_  
**Just use a simple bash script**
- _Make "prep" an external that does the ssh to the container, make another internal weird named prep that it uses_  
**The purpose of this was to allow prep to be called again on an existing instance.  Instead just recreate the instance using start --rm.**
- Add the ability to recognize when running in a repository and name the instance after the repository, use Gem to determine the version, use a .file to determine the port.  
**This could be a secondary tool but does not belong as part of this tool.  It ties us too many external systems and complicates this simple single purpose tool.**
