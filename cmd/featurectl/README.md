# featurectl CLI

See `featurectl --help` for all commands.

## Usage Examples

### Login

```sh
featurectl login --server https://api.featureflags.example.com --token <JWT>
```

### Create a Feature Flag

```sh
featurectl flag create --name "my-feature" --description "My new feature" --enabled
```

### List Flags

```sh
featurectl flag list
```

### Get Flag Details

```sh
featurectl flag get <flag_id>
```

### Update a Flag

```sh
featurectl flag update <flag_id> --enabled true
```

### Delete a Flag

```sh
featurectl flag delete <flag_id>
```

### Add Targeting Rule

```sh
featurectl flag rule add --flag-id <flag_id> --target "user.country=US" --percentage 10
```

### Watch for Flag Updates (real-time)

```sh
featurectl flag watch --id <flag_id>
```

### View Audit Logs

```sh
featurectl audit list --flag-id <flag_id>
```
