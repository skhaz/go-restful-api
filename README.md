This project uses the repository pattern with gorm, it uses gin for handling the HTTP requests and has unit tests


Running

```shell
make compose
```

Creating a workspace

```shell
http localhost:8000/workspaces name=foobar
```

List workspaces

```shell
http localhost:8000/workspaces
```

Update a workspace

```shell
http PATCH localhost:8000/workspaces/82877b28-93d5-4b37-83cb-2bf8e090caf5 name=newname
```

Get a workspace

```
http localhost:8000/workspaces/82877b28-93d5-4b37-83cb-2bf8e090caf5
```

Delete a workspace

```
http DELETE localhost:8000/workspaces/82877b28-93d5-4b37-83cb-2bf8e090caf5
```