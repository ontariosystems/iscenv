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
- Make all commands return the instance name rather than the container ID
- Use tag-based versioning instead of committed version
- Add the ability to run a command with csession

## Maybe
- Add the ability to recognize when running in a repository and...
    - Name the instance after the repository
    - Use Gem to determine the version
    - Use a .file to determine the port
    - Downside of this is it ties us more closely to ruby and to our Gem versioning
- Add a command that reads a simple configuration file to set up a specific environment
- Have prep update the deployment service
- Make "prep" an external that does the ssh to the container, make another internal weird named prep that it uses