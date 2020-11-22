# domaintaken
> CLI which checks if a domain is taken via Cloudflare DoH

## Install

```sh
go get -u github.com/ethanent/domaintaken
go install github.com/ethanent/domaintaken
```

## Usage

```sh
domaintaken ethanent.me example.com google.com ...
# Checks individual domains

domaintaken "mydomain(alpha).com"
# Checks domains, substituting in alphabetical characters

domaintaken "mydomain.(tld,3)"
# Checks domains, substituting in TLDs of length 3

domaintaken mydomain(alpha).(tld,3)
# Combining the above variation methods
```
