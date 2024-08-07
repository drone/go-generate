Compile and install the binary:

```
go install github.com/drone/go-generate
```

# Generate

Generate pipeline configuration for a local repository:

```
go-generate generate /path/to/local/repo
```

Generate pipeline configuration for a local repository:

```
go-generate generate /path/to/local/repo
```

Generate pipeline configuration for a remote repository:

Go:

```
go-generate generate https://github.com/drone/go-scm.git
```

Node:

```
go-generate generate https://github.com/facebook/create-react-app.git 
```

Rails:

```
go-generate generate https://github.com/railstutorial/sample_app_2nd_ed.git
```

Ruby:

```
go-generate generate https://github.com/slim-template/slim.git 
```
