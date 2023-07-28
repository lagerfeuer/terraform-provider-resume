#!/bin/zsh

readonly TFRC="/tmp/dev.tfrc"

cat <<EOF >"${TFRC}"
provider_installation {

  dev_overrides {
      "deutz.io/provider/resume" = "$(go env GOBIN)"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
EOF

echo "export TF_CLI_CONFIG_FILE=${TFRC}"
