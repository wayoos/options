# options

## Setting up Go to work on options

If you have never worked with Go before, you will have to complete the
following steps in order to be able to compile and test Packer. These instructions target POSIX-like environments (Mac OS X, Linux, Cygwin, etc.) so you may need to adjust them for Windows or other shells.

1. [Download](https://golang.org/dl) and install Go. The instructions below
   are for go 1.7. Earlier versions of Go are no longer supported.

2. Set and export the `GOPATH` environment variable and update your `PATH`. For
   example, you can add to your `.bash_profile`.

    ```
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin
    ```

3. Download the options source (and its dependencies) by running `go get
   github.com/wayoos/options`. This will download the options source to
   `$GOPATH/src/github.com/wayoos/options`.
