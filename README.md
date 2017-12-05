# convox/dmz

Proxy requests to internal apps

## Usage

Deploy as a Convox app with the following environment variables:

    ALLOW=^/$|^/test/.*$
    REMOTE_URL=https://other.example.org/path
