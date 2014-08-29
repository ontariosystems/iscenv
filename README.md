# ISCEnv
This utility allows the user to manage docker-based ISC product instances.  These instances are meant to be used as
development environments on a Linux host.

# Usage
The command consists of a single command line utility **`iscenv`**.  This utility is self-documenting.  Please refer
to its help for further documentation on each command.

# TODO
Below are some items that we wish (or may wish) to do in the future.  Those in the maybe section should be considered
only if they will not complicate the intentionally simple design of this project.

## Needed
- Add Service wrappers for the containers which will restart them on reboot
- Evaluate command execution / exits for consistency (use fatal at the same times, always loop over all instances or always fail on first error)
- Use tag-based versioning instead of committed version
- Add the ability to run a command with `iscenv csession`
- Add version command
- Figure out why the terminal emulation is messed up in `iscenv csession`
- Add command that tails the cconsole log (or any file?)
- Add an inspect command

## Maybe
- Add new commands or switches with start, stop or remove all existing instances
- Make all commands return the instance name rather than the container ID
- Have prep update the deployment service

## Rejected
- _Add a command that reads a simple configuration file to set up a specific environment_  
**Just use a simple bash script**
- _Make "prep" an external that does the ssh to the container, make another internal weird named prep that it uses_  
**The purpose of this was to allow prep to be called again on an existing instance.  Instead just recreate the instance using start --rm.**
- Add the ability to recognize when running in a repository and name the instance after the repository, use Gem to determine the version, use a .file to determine the port.  
**This could be a secondary tool but does not belong as part of this tool.  It ties us too many external systems and complicates this simple single purpose tool.**
