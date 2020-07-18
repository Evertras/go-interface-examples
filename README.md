# go-interface-examples

Some examples of using Go interfaces to make cleaner, more testable code.

These assume general familiarity with Go syntax and the ability to write
simple programs.  If you haven't already, take [a tour of Go](https://tour.golang.org/welcome/1)!
You'll be introduced to interfaces there.  You should also read
[Effective Go](https://golang.org/doc/effective_go.html#interfaces).

The point of this repository is to expand on those ideas and explore them
with simple but relatable examples to explain *why* Go tells you to use
interfaces the way it does.

## The Outside World

*What are interfaces? What problems do interfaces solve? Why do I care?*

*My code is hard to test because it has to talk to a database!  What do?  What's this 'mock' thing I hear about?*

*Why should I pass in dependencies instead of just creating them or using globals?*

[Let's find out!](./outside-world)

This is pretty entry level.  We'll explore creating a simple service that
ends up being easy to reason about and easy to test thanks to Go interfaces.
We'll start as simple as possible and explore why we make each iteration of
changes towards a more robust solution.  This will be "interfaces the long way".

I think this is a good read if you're new to Go and/or relatively new to coding
in general.  You may find it useful even if you're intermediate, however, as
I hope it articulates a lot of the underlying reasons behind doing the things
you may already be doing with interfaces/IoC/DI.

## Declare interfaces locally

*Why does Go tell me to declare interfaces at the consumer level?*

[See the article here.](./local-interfaces)

There's a reason libraries don't usually give you an interface to work with,
and why you shouldn't provide one either.  Let's explore why.

This is an example of defining your interfaces locally in each package rather
than trying to do an "interface package" or providing giant interfaces from
a single package that also contains the implementation.

I think this is a good read if you're coming to Go from another language
and may have preconceived ideas of what interfaces should be because of that.
You'll need to be at least a little familiar with Go at this point and
may have written a service or two.
