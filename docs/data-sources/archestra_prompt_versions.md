---
# archestra_prompt_versions (Data Source)

Lists versions of a specific Archestra prompt.

## Example Usage

```hcl
data "archestra_prompt_versions" "example" {
  prompt_id = "prompt-123"
}
```

## Argument Reference

- `prompt_id` (Required) - The ID of the prompt.

## Attribute Reference

- `versions` (Computed) - List of prompt versions.
  - `id` - Version identifier.
  - `version_number` - Version number.
  - `created_at` - Creation timestamp.
