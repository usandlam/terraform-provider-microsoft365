#!/bin/bash

# Example 1: Import an iOS general device restriction template
echo "terraform import microsoft365_graph_beta_device_management_ios_device_configuration_templates.general_restriction_example \"12345678-1234-1234-1234-123456789012\""

# Example 2: Import an iOS trusted root certificate template
echo "terraform import microsoft365_graph_beta_device_management_ios_device_configuration_templates.trusted_certificate_example \"87654321-4321-4321-4321-210987654321\""
