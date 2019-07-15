Go Decouple (based [python-decouple](https://github.com/henriquebastos/python-decouple)): Strict separation of settings from code
=================================================================================================================================

*Decouple* helps you to organize your settings so that you can change parameters without having to redeploy your app.

It also makes it easy for you to:

1. store parameters in *ini* or *.env* files;
1. define comprehensive default values;
1. properly convert values to the correct data type;
1. have **only one** configuration module to rule all your instances.


[![Build Status](https://travis-ci.org/alvesgabriel/go-decouple.svg?branch=master)](https://travis-ci.org/alvesgabriel/go-decouple)


### Summary
* [Why?](#why?)
    * [Why not just use environment variables?](#why-not-just-use-environment-variables?)
* [Usage](#usage)
    * [Where the settings data are stored?](#where-the-settings-data-are-stored?)
        * [Ini file](#ini-file)
        * [Env file](#env-file)
* [How it works?](#how-it-works?)
* [Contribute](#contribute)
* [License](#license)

# Why?

Web framework's settings stores many different kinds of parameters:

* Locale and i18n;
* Middlewares and Installed Apps;
* Resource handles to the database, Memcached, and other backing services;
* Credentials to external services such as Amazon S3 or Twitter;
* Per-deploy values such as the canonical hostname for the instance.

The first 2 are *project settings* the last 3 are *instance settings*.

You should be able to change *instance settings* without redeploying your app.

## Why not just use environment variables?

*Envvars* works, but since `os.LookupEnv` only returns strings and bool respectively, it's tricky.

Let's say you have an *envvar* `DEBUG=False`. If you run:

```go
if _, ok := os.LookupEnv("DEBUG"); ok {
    fmt.Println(true)
} else {
    fmt.Println(false)
}
```
It will print **true**, because `os.LookupEnv("DEBUG")` returns the **string** `"False"` and **bool** `true`. Since `ok` it's is `true`, it will be evaluated as true.

*Decouple* provides a solution that doesn't look like a workaround: ``decouple.Config("DEBUG", nil, "bool")``.

# Usage

Install:

```shell
go get github.com/alvesgabriel/go-decouple
```


Then use it on your `settings.py`.

1. Import the `decouple` object:

```go
import "decouple"
```

2. Retrieve the configuration parameters:

```go
SecretKey := decouple.Config("SECRET_KEY", nil, nil)
Debug := decouple.Config("DEBUG", false, "bool")
EmailHost := decouple.Config("EMAIL_HOST", "localhost", nil)
EmailPort := decouple.Config("EMAIL_PORT", 25, "int")
```

## Where the settings data are stored?

*Decouple* supports both *.ini* and *.env* files.

### Ini file

Simply create a `settings.ini` next to your configuration module in the form:

```ini
[settings]
DEBUG=True
TEMPLATE_DEBUG=%(DEBUG)s
SECRET_KEY=ARANDOMSECRETKEY
DATABASE_URL=mysql://myuser:mypassword@myhost/mydatabase
PERCENTILE=90%
#COMMENTED=42
```

### Env file

Simply create a `.env` text file on your repository's root directory in the form:

```console
DEBUG=True
TEMPLATE_DEBUG=True
SECRET_KEY=ARANDOMSECRETKEY
DATABASE_URL=mysql://myuser:mypassword@myhost/mydatabase
PERCENTILE=90%
#COMMENTED=42
```

# How it works?

*Decouple* always searches for *Options* in this order:

1. Environment variables;
1. Repository: ini or .env file;
1. default argument passed to config.

There are 3 structs doing the magic:


- `RepositoryEmpty`

    Can read values from `os.LookupEnv`.

- `RepositoryIni`

    Can read values from `os.LookupEnv` and ini files, in that order.

- `RepositoryEnv`

    Can read values from `os.LookupEnv` and `.env` files.

The **Config** function is an instance of `AutoConfig` that instantiates a `Config` with the proper `Repository`
on the first time it is used.


## Understanding the CAST argument

By default, all values returned by `decouple` are `strings`, after all they are read from `text files` or the `envvars`.

However, your Python code may expect some other value type, for example:

* Django's `DEBUG` expects a boolean `True` or `False`.
* Django's `EMAIL_PORT` expects an `integer`.
* Django's `ALLOWED_HOSTS` expects a `list` of hostnames.
* Django's `SECURE_PROXY_SSL_HEADER` expects a `tuple` with two elements, the name of the header to look for and the required value.

To meet this need, the `Config` function accepts a `cast` argument which receives any *string* name of type, that will be used to *transform* the string value into something else.

Let's see some examples for the above mentioned cases:

```go
os.SetEnv("DEBUG", "False")
decouple.Config("DEBUG", nil, "bool")
// False

os.SetEnv("EMAIL_PORT", "42")
decouple.Config("EMAIL_PORT", nil, "int")
// 42
```


# Contribute

Your contribution is welcome.

Setup your development environment:

```console
git clone git@github.com:alvesgabriel/go-decouple.git
cd go-decouple
go get -d ./...
```

You can submit pull requests and issues for discussion. However I only consider merging tested code.


# License

MIT License

Copyright (c) 2019 Gabriel Alves <gabriel dot alves dot monteiro1 at gmail dot com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.