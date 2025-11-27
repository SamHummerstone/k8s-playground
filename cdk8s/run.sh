#!/bin/bash

cdk8s import k8s -l go
cdk8s synth

kubectl apply -f dist/