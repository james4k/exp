# ui

Package ui provides a work-in-progress model for taking user input.
Everything graphical is up to the user to deal with based on the view
hierarchy.

To be written:
Why not concurrent? Well...this will be hard to explain. Will need to
introduce each approach attempted in detail. Boils down to not
gaining anything from it except complexity. Introduce your concurrency
at a higher level where appropriate. Go has a few things for that.

## Examples

[GLFW](http://www.glfw.org/) is required to run the examples.

Installing and running the examples is simple, assuming you have your
$GOPATH/bin setup in $PATH:

```
$ go get j4k.co/exp/ui/examples/...
$ ui-wip
```
