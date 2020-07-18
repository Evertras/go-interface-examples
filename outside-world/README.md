# Dealing with the outside world

## Being Useful

Any program you write is going to hopefully try to do something useful in its
lifetime.

So what's 'useful'?  I/O.  Data in, data out.  Someone makes a web request,
[they see a video](https://www.youtube.com/watch?v=dQw4w9WgXcQ).  Someone clicks
a button, [a shirt shows up at their door two days later](https://www.amazon.com/Mountain-Three-Wolf-Short-Sleeve/dp/B002HJ377A/ref=cm_cr_arp_d_product_top?ie=UTF8).
Someone spends $3,500 USD on a clothes iron [and gets told how to iron clothes](https://www.laurastar.com.au/product/laurastar-smart-u/).

You know, *useful*.

The problem with making things useful is that it makes us leave our nice, cozy
land of self-contained code and makes us venture into the terrible place known
as "the outside world".  We don't like the outside world.  That's why we're programmers.

What's so bad about the outside world?  Lots of velociraptors.

In code, there's things worse than velociraptors.  There's *dependencies* and *coupling*.
Let's see what problems can appear.

## The Before Times

First, are you familiar with the [GSL](https://en.wikipedia.org/wiki/Global_StarCraft_II_League)?
You should be.  It's pretty sweet.

Let's say you're a big GSL fan, and you want to make a cool API that people can use to get
the GSL rankings and the like.

Let's also say you're starting on your Go adventure.  You write a simple web server.
This server is so simple that it has one HTTP handler.  It tells you who the current
reigning champion is: TY.

```go
// server.go

func gslCurrentChampionHandler(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("TY"))
}
```

[Check out the code here for reference.](./simple)

Look how simple that is!  It's beautiful.  Job's done, we can dust our hands and call it a day.

Oh, we need testing.  Because if you're not testing, you're going to get attacked by all those
velociraptors I warned you about earlier.  And no one likes that.  Even the velociraptors are
just phoning it in these days.

```go
// server_test.go

func TestGSLCurrentChampionIsTY(t *testing.T) {
	expectedWorldChampion := "TY"

	req := httptest.NewRequest("GET", "/champion", nil)
	res := httptest.NewRecorder()

	gslCurrentChampionHandler(res, req)

	gotWorldChampion := string(res.Body.Bytes())

	if expectedWorldChampion != gotWorldChampion {
		t.Errorf("Expected world champion to be %q but got %q", expectedWorldChampion, gotWorldChampion)
	}
}
```

Great.  Our GSL champion is TY, and we know this because we tested it.
Now we're not getting eaten by velociraptors and we have a useful API that
someone can hit to find the GSL champion!  Everyone's a winner.

The astute readers among you may have realized that, worryingly, this API
will not last.  Defending a title is actually pretty rare, so while we could
cheer for TY in his next tournament appearance and hope he wins so that we
don't have to do a new deploy, that's probably a losing strategy. In fact,
by the time you read this TY is probably no longer the current champion!

Instead, we must go to the *outside world*.

## Here be velociraptors

[Code here](./velociraptors)

You don your safari hat and take a step outside.

The simplest thing we could do here is just store the current champion in a file.
Whenever a new champion is crowned, we can update the file and the server will
return the new champion.

And so we add `champion.txt`:

```
TY
```

And now we can have our handler read the file every time someone asks (*aside: never do this for real!*).

```golang
// server.go

func gslCurrentChampionHandler(res http.ResponseWriter, req *http.Request) {
	contents, err := ioutil.ReadFile("./champion.txt")

	if err != nil {
		log.Println("Failed to read file:", err)
		res.WriteHeader(500)
		return
	}

	res.Write(contents)
}
```

Ok great, now let's update our tes-... oh.  Uh.  Hmm.

```golang
func TestGSLCurrentChampionIsTY(t *testing.T) {
	// I hate everything about this.  Writing this has caused my keyboard
	// to rebel in anger.  Do not use this.  Do not even think about it
	// for too long or adverse health effects may arise.
	contents, err := ioutil.ReadFile("./champion.txt")

	if err != nil {
		t.Fatal("Failed to read file:", err)
	}

	expectedWorldChampion := string(contents)

	req := httptest.NewRequest("GET", "/champion", nil)
	res := httptest.NewRecorder()

	gslCurrentChampionHandler(res, req)

	gotWorldChampion := string(res.Body.Bytes())
	gotCode := res.Code

	if 200 != gotCode {
		t.Errorf("Expected code 200 but got %d", gotCode)
	}

	if expectedWorldChampion != gotWorldChampion {
		t.Errorf("Expected world champion to be %q but got %q", expectedWorldChampion, gotWorldChampion)
	}
}
```

Oof.  This is bad.

...but why?

Read this test as if you had no idea what `gslCurrentChampionHandler` did.
If your head combusts the moment you see "ioutil.ReadFile", then you're
doing well.

Seriously, take a moment.  Look at the test.  Imagine you're a fresh new developer
and you had never seen `gslCurrentChampionHandler`.  What questions would you have?

Here are the questions that would be going through my head.

- Why do I need to read champion.txt?  Where did I know that?
- What happens if the file moves?
- How was I supposed to know the format of champion.txt?  What if it was JSON?
- How can I test for any code other than 200?  Do I... do I delete the file first?  That'll screw up future tests...
- ~~Seriously what the f~~

This is just for a single, simple file.  Actual real world code is going to be
dealing with various databases, caching, etc.  While end to end testing is nice
and all, you still want unit tests that are sane for both the happy path and
checking how you handle errors.  This will also multiply very quickly once you
start adding more handlers that have to deal with more data in more and more ways.
This is a fragile test because of all this extra knowledge, and fragile tests
means tests break, and tests breaking means people don't test, and people not
testing means the velociraptors come at night for you and your loved ones.

And yet... AND YET... this is about as good as it gets to test the handler as it
currently exists.  We must know the current champion in order to test the current
champion is getting returned.  If we wanted to test that we got a 500 when the
file didn't exist, we can't do that without physically deleting the file.  I want
you to briefly imagine trying to write a test that reads the file, deletes the file,
runs the handler, then rewrites the file in order to preserve state.  Now stop
imagining that because velociraptors will find your house if you ever write that.

Let's fend off the velociraptors and try to do better.

## The first abstraction

Let's take another look at the code as-is.

```go
func gslCurrentChampionHandler(res http.ResponseWriter, req *http.Request) {
	contents, err := ioutil.ReadFile("./champion.txt")

	if err != nil {
		log.Println("Failed to read file:", err)
		res.WriteHeader(500)
		return
	}

	res.Write(contents)
}
```

While this is very few lines of actual code, the current handler has to deal
with the following:

- What even is HTTP? (res/req parameters)
- How do I get the data I need? (read a file with ioutil.ReadFile)
- What even is a file? (ioutil.ReadFile)
- What file do I need to read? (champion.txt)
- What is the format of that file? (literally just the name in this case, but we still need to know that!)
- What do I do if I can't get the data? (return 500)
- What do I return to the user? (plain text name)

And because the handler has to deal with that, *you* have to deal with all that whenever
you look at this code.

This is the result of coupling.  We are *coupling* our handler to the file system.
We are *coupling* our handler to the fact that we're storing our data in a file called
"champion.txt".

While we can't get rid of all of these, we can get some more organization going.  It looks
like some of these are dealing with HTTP and others are dealing with files.

HTTP:

- What even is HTTP? (res/req parameters)
- What do I do if I can't get the data? (return 500)
- What do I return to the user? (plain text name)

Data:

- What even is a file? (ioutil.ReadFile)
- How do I get the data I need? (read a file with ioutil.ReadFile)
- What file do I need to read? (champion.txt)
- What is the format of that file? (literally just the name in this case, but we still need to know that!)

**The less a block of code has to know, the less WE have to know when working on that code. The less we have to know when working on code, the easier it is to maintain.**

We'll come back to that point later. For now let's break out the stuff that deals
with data into its own thing.

```golang
// data.go

// GSLDataStore knows how to get GSL data
type GSLDataStore struct {
	championFile string
}

// NewGSLDataStore returns a GSLDataStore ready to tell us about the GSL
func NewGSLDataStore(championFile string) *GSLDataStore {
	return &GSLDataStore{
		championFile,
	}
}

// GetCurrentChampion returns the name of the current GSL champion
func (s *GSLDataStore) GetCurrentChampion() (string, error) {
	contents, err := ioutil.ReadFile(s.championFile)

	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(contents), nil
}
```

Unlike the old HTTP handler, this code is much more focused.  It says precisely
what it needs (a filename for the champion) and the function GetCurrentChampion()
will handle any formatting/encoding.  It just returns a string and an error we can
check if something went wrong.

Ok great, but what about our handler?  When we check the signature, we immediately
see a problem.

```golang
func gslCurrentChampionHandler(res http.ResponseWriter, req *http.Request) {
```

In order to be a proper HTTP handler, it *must* have those two parameters and only
those two parameters, and return nothing.  So how do we use our fancy new data store?

There's a few ways we can do this, but I'm going to stick with one that I've found
to be very simple and easy to maintain.

```golang
func gslCurrentChampionHandler(dataStore *GSLDataStore) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		champion, err := dataStore.GetCurrentChampion()

		if err != nil {
			log.Println("Failed to get current champion:", err)
			res.WriteHeader(500)
			return
		}

		res.Write([]byte(champion))
	}
}
```

Ok woah, what happened here?  What is this madness?  If you didn't know that
functions are just like any other type in Go... welcome to Go!

Instead of a handler just existing on its own, we've created a function that
returns a function.  Notice the difference between the old version and the new:

```golang
// Old
mux.HandleFunc("/champion", gslCurrentChampionHandler)

// New                                               vvvvvvvvvvv
mux.HandleFunc("/champion", gslCurrentChampionHandler(dataStore))
```

The created function returns a function with the proper signature, but that
function acts as a *closure* that has access to the data store.  This is a
useful pattern that lets us add *dependencies* in without having to create them.

```golang
func runServer(address string, dataStore *GSLDataStore) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/champion", gslCurrentChampionHandler(dataStore))

	return http.ListenAndServe(address, mux)
}
```

We also added a parameter to `runServer`.  We don't want to create the data store
ourselves, because we are very lazy and it'd be nice if that was someone else's problem.

### A small note on dependencies

Actually there's a few real reasons I want to touch on here for why we're not creating
the data store in `runServer` or the handler itself.

If multiple things are going to use the data store, we don't want to create multiple
data stores!  We want to reuse them.  If we had multiple handlers that all wanted
to use this data store, we wouldn't create separate data stores in each handler, right?
So let's just pass it in.

But why defer it even outside `runServer`?  What if `runServer` was really the only thing
using this?  What's the harm in creating the data store here?

Configuration, for one thing.  Even this ridiculously simple example requires the data store
to be told where the file is.  In a real world example, you'd be dealing with database
connection information, including authentication and other madness.  You do *not* want to
deal with that here.  That sort of thing belongs in `main()`, generally.

This is touching on a much larger subject that I encourage you to read up on called
Inversion of Control (IoC) and Dependency Injection (DI).  Going on about this is a
topic for other articles.

**The short version is that any instance or function should always explicitly declare what it needs to run, and whatever creates it must provide those things.**

### Back to the code

Ok, we've made some changes.  Check out the code [for the server](./less-velociraptors/server.go)
and [for the data store](./less-velociraptors/data.go) for a recap.  Separating out
these two files gave us the following benefits.

- The server doesn't need to know where our data is stored or what format any of our data is stored in, so it is easier to reason about and maintain
- The data store can focus entirely on the data itself
- We can configure our data store in main more easily

I want you to pretend that you're a new dev again coming in on this project.  Look at these
two lines of code.

```golang
contents, err := ioutil.ReadFile("./champion.txt")
```

```golang
champion, err := dataStore.GetCurrentChampion()
```

Which of these feels cleaner?  Which of these would you rather work with?  Which of these
do you need to know more about in order to understand?

Which of these would you rather have to worry about changing at 2 AM because your server
is on fire and you need to ship a fix *right now* or the company's going under?

I've asked you to put yourself in the spot of a new dev multiple times now.  The reason is
that in a few months (or even weeks/days) of not looking at your code, you actually are
a new dev again because you will have forgotten about the details of your code.  Make things
easy for yourself.  Keep it simple.  Keep it focused.

That second line of code isn't going to change if you change how your data is stored.
That first line of code may end up turning into more like 20 if you start dealing with a
database or add other real world concerns like metrics.

Speaking of things being on fire, how's our test doing?

[Still awful.](./less-velociraptors/server_test.go)

Oof.  If anything we've made things worse.  Our server may not need to worry about what
the data store is doing, but our test still does.  While our code is now easier to reason
about, we are still *coupled* to our data store as a *dependency*.  This is particularly
highlighted by the fact that our test still needs to know about all the stuff we didn't
want our server to worry about!

On the bright side, we can at least add a reasonable test to check that we're returning
an error code properly. That's just a sad consolation prize for the moment.  We can do better.
We *will* do better.

You've made it this far.  Now is the time for interfaces to shine.

Behold! ~~Corn!~~

## Interfaces in Go

I haven't told you what an interface is yet or how to write one.  I still won't.  Not yet.
I want you to understand *why* they exist before we look at how to add one.

Right now I want you to think about the [current state of our little service](./less-velociraptors).

We have a data store that gets the current GSL champion's name.  It uses the file system
and can be configured to read a certain file to find out who the champion is.

We have a server that has an HTTP handler that writes the current GSL champion's name back
to whoever called it.  It uses the data store we defined above.

Why does our server need the data store?  It needs something that can get the current GSL champion's name.

This is, in plain English, what an interface is.  Read this again.

*It needs something that can get the current GSL champion's name.*

There's nothing about a database there.  No file systems.  No encodings.  No code.  Just
a simple statement of *intent*.

This is an immensely powerful statement.  "What do I need to run this?  Ah, it needs to be
something that can get the current GSL champion's name."

The code is close to doing this already.

```golang
func gslCurrentChampionHandler(dataStore *GSLDataStore) http.HandlerFunc {
```

This is saying something slightly different.  This is saying, "It needs a GSLDataStore to run."
We can then go to GSLDataStore and see that it gets a GSL champion name and maybe infer from there.
It's an improvement over what we started with, but we can do better.

So now, finally, after this long, arduous journey you've come on with me, we will define
an interface in Go.

```golang
type CurrentChampionGetter interface {
	GetCurrentChampion() (string, error)
}
```

An interface in Go is a list of capabilities.  We can declare an interface wherever we want, and
anything that fulfills all the entries of the interface can be used as an instance of that interface.

That's a lot of words, but we'll see what they mean here.

Forget the data store.  Forget about it.  It's not important.  Not anymore.  What's important
is this interface.  This interface says something very specific.  It says "I can get the current champion."
It says it right there, look!  `GetCurrentChampion()`

We can use an interface as a variable type.  Let's see what our handler signature looks like
if we use this interface instead of the data store type.

```golang
func gslCurrentChampionHandler(championNameGetter ChampionNameGetter) http.HandlerFunc {
```

We're still sending in a dependency, but this is worlds apart in how it reads.  Before the handler
said "I need a GSLDataStore, because... reasons."

Now it says "I need something that can get the current champion." That's *huge* in terms of what
it tells the reader without having to read a single line of code inside the handler itself.
And because it's just a simple list of methods, it also means we don't need to know anything else
about what's being passed in here.

Remember: **The less a block of code has to know, the less WE have to know when working on that code. The less we have to know when working on code, the easier it is to maintain.**

In fact, nothing's changed on the inside besides the name.

```golang
champion, err := currentChampionGetter.GetCurrentChampion()
```

Ok, you can remember the data store again.  It hasn't changed.

```golang
func (s *GSLDataStore) GetCurrentChampion() (string, error) {
```

Notice this method matches our interface.  How convenient!  That was definitely not an accident.
While the data store doesn't need to know about the interface, it does need to have one method
for every method listed in the interface and they must match exactly.  It can have other methods
that aren't in the interface, but it must have at *least* all the methods in the interface
in order to be considered a match.

If it is a match, then any time we use the type `CurrentChampionGetter` we could pass it a `*GSLDataStore`.

```golang
var getter CurrentChampionGetter = NewGSLDataStore("champion.txt")
```

As a small detail, note that I said `*GSLDataStore`, not `GSLDataStore`.  It must be a pointer,
because the receiver for the `GetCurrentChampion` method takes `(s *GSLDataStore)` and not `(s GSLDataStore)`.

Ok, so we see that we can pass in a data store to match our interface.  Interestingly, our `main.go` doesn't change.

```golang
dataStore := NewGSLDataStore("./champion.txt")
err := runServer(":8080", dataStore)
```

Again, we can pass in a `*GSLDataStore` as `CurrentChampionGetter` because it matches that interface.

But the real prize is our test.  [Our test is beautiful](./no-velociraptors/server_test.go).

Because we're no longer tied to our data store type, we can provide whatever we want
as long as it fulfills our interface.  We have *decoupled* our server from our data store.
This lets us define our own simple type that we can set up test scenarios with.

```golang
type mockCurrentChampionGetter struct {
	current      string
	pendingError error
}

func (g *mockCurrentChampionGetter) GetCurrentChampion() (string, error) {
	if g.pendingError != nil {
		return "", g.pendingError
	}

	return g.current, nil
}
```

Using this mock, we can test both happy paths (where things go right) as well
as sad paths (where things return errors) without touching our file system
or any external things at all.  We have come full circle.  We are no longer
dealing with the outside world.  The velociraptors scratch at the door in vain.
We are safe.

## Summary 

We started with a simple service that just returned a string.

We then started to make it touch the outside world just enough to show the problems that quickly arise.

When our code started to do too much, we saw that we could split it up.

When we split it up, we found that we were still *coupled* because of our concrete struct type.

Finally, we saw how a Go interface allowed us to *decouple* our code and describe its dependency precisely while letting us cleanly test our service.

The power of interfaces is to *decouple* our code.  You've now seen one of the
uses of interfaces in decoupling the outside world, but there's a lot more that
interfaces can do.  Go forth and explore!

But maybe bring some velociraptor repellant with you.
