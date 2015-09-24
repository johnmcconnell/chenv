# chenv
Change your envs with a single command.

## Usage

### Save
Store your current .env file to a profile name of your choice.
For this example dev.

```
chenv -save dev
```

### Change
Change your current .env file to a stored env.
For this example dev.

```
chenv dev
```

### List
List your stored envs

```
chenv
```

## Gotchas
One gotcha is that you must use -save before the environment name.

For example `chenv -save hello` works.

`chenv hello -save` does not.
