package main

import (
  "net/url"
  "fmt"
  "strconv"

  pb "github.com/LCOGT/keda-celery-flower-scaler/externalscaler"
)

type Config struct {
  FlowerURL *url.URL

  DesiredMetricValue int

  ActivationThreshold int
}

func parseConfig(so *pb.ScaledObjectRef) (*Config, error) {
  c := &Config{
    DesiredMetricValue: 2,
    ActivationThreshold: 1,
  }

  if val, ok := so.ScalerMetadata["flowerURL"]; ok {
    url, err := url.Parse(val)

    if err != nil {
      return nil, fmt.Errorf("flowerURL is not a valid URL")
    }

    c.FlowerURL = url

  } else {
    return nil, fmt.Errorf("flowerURL is required")
  }

  if val, ok := so.ScalerMetadata["desiredMetricValue"]; ok {
    valInt, err := strconv.ParseInt(val, 10, 32)

    if err != nil {
      return nil, fmt.Errorf("failed to parse desiredMetricValue as an integer: %w", err)
    }

    if valInt <= 0 {
      return nil, fmt.Errorf("desiredMetricValue must be greater than 0")
    }

    c.DesiredMetricValue = int(valInt)

  }

  if val, ok := so.ScalerMetadata["activationThreshold"]; ok {
    valInt, err := strconv.ParseInt(val, 10, 64)

    if err != nil {
      return nil, fmt.Errorf("failed to parse activationThreshold as an integer: %w", err)
    }

    if valInt < 0 {
      return nil, fmt.Errorf("activationThreshold must be greater than or equal to 0")
    }

    c.ActivationThreshold = int(valInt)
  }

  return c, nil
}
