# Toil

Not to be confused with TOIL (Time off in lieu)

> Toil is the kind of work tied to running a production service that tends to be manual, repetitive, 
> automatable, tactical, devoid of enduring value, and that scales linearly as a service grows.
> 
> -- [Goole Site Reliability Engineering book](https://sre.google/sre-book/eliminating-toil/)

Whenever you start to work on a health task for a service, run:

    $ toil [-m msg] [service]

Note, `toil` is a daily tracker, so if you work on a health task for more than
one day in a row, run it again for each day of work.

Motivations:

- provide actual data on time spent on toil for key services
- visibility into who works on what and gaps in team knowledge (e.g. if service
  X is only ever fixed by person Y)

## Installation

Ensure you have Go 1.17+ installed (`brew install go`) and then run:

    go install github.com/guardian/toil

The same command will apply any updates.

Note, at the moment this is written for DevX use only, but please shout if you
want to use it and we can sort something!

# What actually happens?

Toil is basically like the ticketing system described
[here](https://joearms.github.io/published/2014-06-25-minimal-viable-program.html).
When Toil is invoked, a record is created under `~/toil/YYYY-MM-DD-HH-MM-DD`
with the following structure:

```
responsible: [git config user.email]
service: [service]
----
[msg | 'Describe your problem here.']
```

These tickets are pushed to a `toil-records` repo to track things centrally.
