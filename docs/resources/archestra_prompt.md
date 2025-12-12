---
# archestra_prompt (Resource)

Manages an Archestra prompt, including creation, versioning, and rollback.

## Example Usage

```hcl
resource "archestra_prompt" "example" {
  name        = "Coding Assistant Prompt"
  description = "System prompt for coding assistance"
  content     = "You are a helpful coding assistant..."
  tags        = ["coding", "assistant"]
  visibility  = "public"
}
```

## Argument Reference

- `name` (Required) - The name of the prompt.
- `description` (Optional) - Description of the prompt.
- `content` (Required) - The content of the prompt.
- `tags` (Optional) - List of tags for the prompt.
- `visibility` (Optional) - Visibility of the prompt (e.g., public, private).
- `version_id` (Optional) - Version ID for rollback.

## Attribute Reference

- `id` (Computed) - Prompt identifier.

## Import

Import is supported using the prompt ID.

```shell
terraform import archestra_prompt.example <prompt_id>
```
