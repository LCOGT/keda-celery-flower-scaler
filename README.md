keda-celery-flower-scaler
=========================

[Keda](https://keda.sh/docs/latest/scalers/external) scaler for [Celery](https://docs.celeryq.dev/) workers.

This project provides an [External Scaler GRPC service](https://keda.sh/docs/2.8/concepts/external-scalers/)
that exposes scaling metrics using [Flower](https://flower.readthedocs.io).


Usage
-----

Deploy the scaler (see `k8/`) and expose a service to it:

```yaml

apiVersion: v1
kind: Service
metadata:
  name: keda-celery-flower-scaler-api
  namespace: scalers
spec:
  ...
  ports:
    - name: api
      port: 8000
      targetPort: api
```

Assuming you have a Deployment of Celery workers:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: celery-worker
  namespace: app
spec:
  replicas: 1
  ...
```

And Flower is correctly deployed (outside the scope of this doc, but it MUST be
able to connect the RabbitMQ Management API) and accessible via a Service:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: flower
  namespace: app
spec:
  selector:
    ...
  ports:
    - name: api
      port: 8000
      targetPort: api
```

You can then create a `ScaledObject` in the same namespace as the Celery worker
Deployment:


```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: celery-worker
  namespace: app
spec:
  scaleTargetRef:
    kind: Deployment
    name: celery-worker
  pollingInterval: 30
  cooldownPeriod: 30
  minReplicaCount: 0
  maxReplicaCount: 100
  fallback:
    failureThreshold: 3
    replicas: 5
  triggers:
  - type: external
    metadata:
      scalerAddress: keda-celery-flower-scaler-api.scalers.svc.cluster.local:8000
      flowerAddress: http://flower.app.svc.cluster.local:8000 # required

      # This is the scaling factor as described at
      # https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#algorithm-details
      # Set this to the number tasks you expect each Worker Pod to consume on average.
      # For example, With CELERY_WORKER_PREFETCH_MULTIPLIER=1 each worker can
      # claim up to 2 tasks (1 runnning & 1 prefetched). With `desiredMetricValue: 2`
      # `desiredReplias` will be half the number of total queued tasks.
      # Must be greater than or equal to 1. Default: 2
      desiredMetricValue: 2

      # Minimum number of tasks pending or running to scale the Celery Deployment
      # from 0. For example, a value of 1 will scale it the moment a task is queued.
      # A value of 2 would wait for at-least 2 tasks to be queued before the
      # deployment is scaled. Must be greather than or equal to 0. Default: 1
      activationThreshold: 1
```

Development
-----------

Setup dev tools with [devenv](https://devenv.sh/getting-started/):

```shell
devenv shell
```

Start a development K8s cluster:

```shell
kind create cluster --name keda
```

Insall [Keda](https://keda.sh/docs/2.8/deploy/):

```shell
helm upgrade keda keda \
  --install \
  --namespace keda-system \
  --create-namespace \
  --repo https://kedacore.github.io/charts
```

Use skaffold to develop:

```shell
skaffold dev --port-forward
```

You can simulate work by using `service/celery-app-api` to submit fake Celery
tasks.
