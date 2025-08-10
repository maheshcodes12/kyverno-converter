# YAML to CEL Converter
This tool is specifically designed to streamline the migration of legacy validation policies to the modern, CEL-based format. It dynamically converts any YAML-based validation rule that uses a `pattern` or `foreach` construct into a functionally equivalent, CEL-based `ValidatingPolicy`.

`Live Web App`: [yaml-to-cel.netlify.app](https://yaml-to-cel.netlify.app)

## Key Features:
`Automated Conversion`: Automatically transforms complex YAML files into their corresponding CEL expressions, preserving the intended logic and structure.

`Logic Preservation`: Ensures that the validation logic of your original policy remains fully intact after conversion.


## How It Works
The converter parses the structure and values within a given YAML input. It then maps the YAML hierarchy and data types into a valid and functionally equivalent CEL expression. This enables you to define complex logic in the familiar YAML format and translate it for use in any system that leverages CEL.

The web application provides a more direct way to see this in action: simply paste your YAML into the input field, and the corresponding CEL output is generated instantly.
