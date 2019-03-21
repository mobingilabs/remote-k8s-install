# ocean

**Note**: it is on developing and testing, so it is can't use for production.
**Note**: we develop at dev branch,when test finish,we merge dev to master

## TODO
- HA supported(now it is a just one master and one node)
- no need config file,it exposed as a rest server
- not ssh cat, but like scp,copy many files once time to speed up install

## fix up
- log add(now no any log),may we will use zap(by uber)
- ssh broken line reconnection
- More reasonable directory structure
- more check and more config
- dep manage add,because the china net, we may be used vendor mechanism,not go mod
