# chenv
chenv is a completely overegineered program to manage env files.

For example, times when developing you will want the change the
environment of your application from **dev**, **staging**,
**stress-test**, etc. It is as simple as:
```
$ chenv dev
loaded .chenv/dev.env
using:
  * dev

$
```

## Usage

### Save
Store your current **.env** to a profile
```
$ chenv save dev
saved current .env to .chenv/dev.env
  dev saved!

$
```

### Change
Change to **dev**

```
$ chenv dev
loaded .chenv/dev.env
using:
  * dev

$
```

### List
List your envs

```
$ chenv
view usage with: chenv -h
available envs:
  * dev
    prod
    test

$
```

## Example Envs
[Examples](.chenv)
