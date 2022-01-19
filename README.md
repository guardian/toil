# Toil

Tiny program to track toil on key services.

Whenever you start to work on a health task for a service, run:

    $ toil [service] [-m msg]

Note, `service` should autocomplete.

Motivations:

- provide actual data on time spent on toil for key services
- visibility into who works on what and gaps in team knowledge (e.g. if service
  X is only ever fixed by person Y)

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
