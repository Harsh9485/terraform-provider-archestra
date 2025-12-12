---
# archestra_prompt (Data Source)

Fetches an existing Archestra prompt by ID or name.

## Example Usage

```hcl
data "archestra_prompt" "example" {
  name = "Coding Assistant Prompt"
}
```

## Argument Reference

- `id` (Optional) - Prompt identifier.
- `name` (Optional) - The name of the prompt. Either `id` or `name` must be provided.

## Attribute Reference

- `id` (Computed) - Prompt identifier.
- `name` (Computed) - The name of the prompt.
- `description` (Computed) - Description of the prompt.
- `content` (Computed) - The content of the prompt.
- `tags` (Computed) - List of tags for the prompt.
- `visibility` (Computed) - Visibility of the prompt.
