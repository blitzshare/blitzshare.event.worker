
[![CircleCI](https://circleci.com/gh/blitzshare/blitzshare.event.worker/tree/main.svg?style=svg&circle-token=7800b5e3b65b70a5498c5965c502470ee0af23a1)](https://circleci.com/gh/blitzshare/blitzshare.event.worker/tree/main)

![logo](./assets/logo.png)

# blitzshare.event.worker
Responsible for processing events from the subscribed topics in kubemq and updating redis records accordingly.

## Getting started

```shell
# install dependencies
$ make install
# start local server
$ make start
```

## Tests
```shell
# unit tests
$ make test
# re/build mocks
$ make build-mocks
# generate test coverage report
$ make coverage-report-html
```

### Acceptance test
Make sure your network is connected to cluster: 
```shell
# connect
$ tellepresence connect
# disconnect
$ tellepresence quit
```
Run acceptance tests:
```shell
$ make acceptance-tests
```

## Debugging events locally (kumebqctl):

```shell
# peer registry cmd
$ kubemqctl queues send p2p-peer-register-cmd '{"multiAddr": "multiAddr", "otp":"otp", "mode": "mode", "token":"token"}'
# node registry cmd
$ kubemqctl queues send p2p-bootstrap-node-registry-cmd '{"nodeId":"nodeId", "port": 123}'
# peer deregister cmd
$ kubemqctl queues send p2p-peer-deregister-cmd  '{"token":"token", "otp": "otp"}'
```
## Tools
[kubemqctl](https://docs.kubemq.io/getting-started/quick-start)

[kubectl](https://kubernetes.io/docs/reference/kubectl/overview/)

[telepresence](https://www.telepresence.io/)
