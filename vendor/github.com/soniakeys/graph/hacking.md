#Hacking

Basic use of the package is just go get, or git clone; go install.  There are
no dependencies outside the standard library.

The primary to-do list is the issue tracker on Github.  I maintained a
journal on google drive for a while but at some point filed issues for all
remaining ideas in that document that still seemed relevant.  So currently
there is no other roadmap or planning document.

CI is currently on travis-ci.org.  The .travis.yml builds for go 1.2.1
following https://github.com/soniakeys/graph/issues/49, and it currently builds
for go 1.6 as well.  The travis script calls a shell script right away because
I didn’t see a way to get it to do different steps for the different go
versions.  For 1.2.1, I just wanted the basic tests.  For a current go version
such as 1.6, there’s a growing list of checks.

The GOARCH=386 test is for https://github.com/soniakeys/graph/issues/41.
The problem is the architecture specific code in bits32.go and bits64.go.
Yes, there are architecture independent algorithms.  There is also assembly
to access machine instructions.  Anyway, it’s the way it is for now.

Im not big on making go vet happy just for a badge but I really like the
example check that I believe appeared with go 1.6.  (I think it will be a
standard check with 1.7, so the test script will have to change then.)

https://github.com/client9/misspell has been valuable.

Also I wrote https://github.com/soniakeys/vetc to validate that each source
file has copyright/license statement.

Then, it’s not in the ci script, but I wrote https://github.com/soniakeys/rcv
to put coverage stats in the readme.  Maybe it could be commit hook or
something but for now I’ll try just running it manually now and then.

Go fmt is not in the ci script, but I have at least one editor set up to run
it on save, so code should stay formatted pretty well.
