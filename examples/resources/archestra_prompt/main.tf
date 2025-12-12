resource "archestra_prompt" "example" {
  name        = "Coding Assistant Prompt"
  description = "System prompt for coding assistance"
  content     = "You are a helpful coding assistant..."
  tags        = ["coding", "assistant"]
  visibility  = "public"
}
