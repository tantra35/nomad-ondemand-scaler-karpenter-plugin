## Auxiliary rpc plugin for use in karpenter nodeprovider in hashicorp nomad ondemand cluster scaler

In fact, it is just a layer for implementing the following service:

```proto
service KarpenterService {
	rpc ListInstances (ListInstancesRequest) returns (ListInstancesResponse);
	rpc AddInstances (AddInstancesRequest) returns (AddInstancesResponse);
	rpc RemoveInstances (DeleteInstancesRequest) returns (DeleteInstancesResponse);
}
```