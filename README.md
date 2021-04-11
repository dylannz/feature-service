# feature-service

[![Go](https://github.com/dylannz/feature-service/actions/workflows/go.yml/badge.svg)](https://github.com/dylannz/feature-service/actions/workflows/go.yml)

**Note:** This is a work in progress. I'll attempt to keep backward compatibility with all changes, but like many things in life, there are no guarantees.

This is a stateless feature flag service built in Go. It's designed to be operated in Kubernetes with configuration mounted to the volume via a [ConfigMap](https://kubernetes.io/docs/concepts/configuration/configmap/), however there is no Kubernetes-specific code so it should work fine on other platforms.

The repository includes:

- A Dockerfile for easy building
- A docker-compose file 
- An openapi 3 spec (spec/spec.yaml) which is used to generate request/response objects. This can be used to generate API clients in many languages, although I have not done so yet.

## API

See spec/spec.yaml for the API specification, and the examples below to see how to use the API.

## Configuration

The service allows some configuration via environment variables:

- **LOG_LEVEL** [logrus log level](https://github.com/sirupsen/logrus#level-logging). 'debug' level will tell you exactly why a feature was enabled/disabled in the log output.
- **CONFIG_DIR** specifies the directory containing YAML files to load. You can split your configuration across multiple YAML files and the service will read/combine all of them. This can help prevent merge conflicts if you are managing these files across multiple teams.
- **HTTP_ADDR** sets the IP address and port to listen for connections on. This defaults to 127.0.0.1:3000 to prevent the macOS warning that you get when you listen to :3000, but you probably want this set to :3000 when running within your chosen orchestration system.

## Examples

See config/example.yml for example YAML configurations, they have some annotations in there explaining what's going on}. Some example requests are below, the responses were generated using the example configuration in config/example.yml.

### Get the status of a feature

One of the features of the service is that it allows arbitrary key/value pairs in the "vars" parameter within the JSON request body. You can use this to include information about the user such as their user ID or name, which can then be used to evaluate whether a feature is enabled or not.

Here is an example using custom variable 'customer_id':

```bash
curl -XPOST localhost:3000/features/status/stripe_billing -d '{"vars":{"customer_id":"1"}}' | jq
{
  "features": {
    "stripe_billing": {
      "enabled": true
    }
  }
}
```

Using custom variable 'customer_name':

```bash
curl -XPOST localhost:3000/features/status/stripe_billing -d '{"vars":{"customer_name":"Alex"}}' | jq
{
  "features": {
    "stripe_billing": {
      "enabled": true
    }
  }
}
```

Using custom variable 'customer_id'. Note it also returns the custom variables which are only enabled for customer ID '123':

```bash
curl -XPOST localhost:3000/features/status/stripe_billing -d '{"vars":{"customer_id":"123"}}' | jq
{
  "features": {
    "stripe_billing": {
      "enabled": true,
      "vars": {
        "foo": "bar"
      }
    }
  }
}
```

### Get a list of all enabled features
```bash
curl -XPOST localhost:3000/features/status -d '{"vars":{"customer_id":"1"}}' | jq
{
  "features": {
    "profile_page_v2": {
      "enabled": true
    },
    "stripe_billing": {
      "enabled": true
    }
  }
}
```

## Run

You can run using docker/docker-compose with:
```bash
docker-compose up
```

Or you can of course run it on your local system:
```bash
go run main.go
```

## Test

```bash
go test ./...
```

## FAQ

### Why should I use this?
- You want to avoid writing your own feature flag service. Yes these services are generally pretty 'easy' to write, but it's not a business differentiator - i.e. it's probably not part of your core product, so you probably have better things to spend time on.
- You don't want to pay thousands of dollars a year for access to a SaaS product, or the hassle of introducing another vendor to your organization. Or you want to run it yourself for some other reason.
- You want a solution which prioritises availability.

### Why not use a database?
A database creates an operational dependency which is (in my opinion) unnecessary. Performance and uptime are the biggest concerns for a service like this, and generally the configuration data should be small enough to keep a copy of it on every instance. Eventually I'd like to build a UI/DB layer that knows how to generate + distribute the YAML files periodically, but that's a separate project.

### Why 'feature-service'? Can't you be a little more imaginative?
Naming things is hard. Open to suggestions for catchier names :)