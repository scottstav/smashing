![](./NIGEL.png)

# Smashing

A CLI for generating auth0 tokens.

# Usage

You can provide your credentials in a few ways

## 1. Through a "profile" file

Provide the file as a command line argument

```
$ cat ~/basic-test.env

CLIENTSECRET="09876"
CLIENTID="12345"
USERNAME=your@email.com
PASSWORD=asdfasdf
AUDIENCE=your-api

$ smashing --profile="~/basic-test.env"
```

## 2. Through command line arguments
`$ smashing --username=your@email.com --password=asdfasdf`

## 3. Input manually
If you run the program with no input, or a subset of required input, you will be prompted to enter the needed values at runtime.

## 4.
Combination of the above 3 options. Command line arguments take precedence over profiles.
