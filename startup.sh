#!/bin/sh

# First update the DB if required
/app/migrate --action=upgrade

# Run the service
/app/iam-manager