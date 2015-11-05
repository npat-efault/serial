
#serial [![GoDoc](https://godoc.org/github.com/npat-efault/serial?status.png)](https://godoc.org/github.com/npat-efault/serial)
 
Download:
```shell
go get github.com/npat-efault/serial
```

***

Package serial provides a simple, high-level, system-independent
interface for accessing asynchronous serial ports.

It provides functions and methods for opening serial ports,
configuring their basic parameters (baudrate, character format,
flow-control, etc.), for reading and writing data from and to them,
and for a few other miscellaneous operations (e.g. sending a break
signal, flushing the I/O buffers).

Data transfer operations support deadlines (timeouts)
and safe cancelation; a blocked read or write operation can be
safely and reliably canceled from another goroutine by closing the
port.

###Supported systems

Most unix-like systems are supported.

Package serial uses package
[termios](https://godoc.org/github.com/npat-efault/serial/termios) to
supports systems that provide the POSIX Terminal Interface for
configuring terminal devices.

For data-transfer operations (Read, Write) it uses package
[poller](https://github.com/npat-efault/poller), which provides I/O
operations with timeouts and safe cancelation.

Depending on the specific system, both of these packages (termios and
poller) can be compiled either to use CGo or as pure-Go packages. See
their documentation for details.

Additional of support for other systems is certainly possible, and
mostly welcome. Patches and pull requests for this will be greatly
appreciated.

***

#serial/termios
 [![GoDoc](https://godoc.org/github.com/npat-efault/serial/termios?status.png)](https://godoc.org/github.com/npat-efault/serial/termios)

Package termios is a simple Go wrapper to the POSIX Terminal
Interface (POSIX Termios). It's more low-level and system-specific
than package serial and can be used to configure terminal
devices on most Unix-like systems that support POSIX Termios.

###Supported systems
Package termios should work on all systems that support the POSIX
terminal interface, that is, on most Unix-like systems.  Depending
on the system, package termios can either be built to use the
system's LIBC functions and macros through CGo, or as a pure-go
package that issues system calls directly. In either case, the API
presented by the package is identical.

By default package termios is built to use CGo to access the system's
LIBC termios functions and macros. This is the most portable option.

If you wish, you can build package termios as a pure-Go package that
issues system-calls directly. To do this define the "nocgo" build-tag
when building/installing the package, like this:

```shell
cd $GOPATH/github.com/npat-efault/serial/termios
go install -tags 'nocgo'
```

Building termios as a pure-Go package is *not* supported for all
systems.





* * *
Automatically generated by [autoreadme](https://github.com/jimmyfrasche/autoreadme) on 2015.11.05
