set -l subcommands completions.create messages.create messages.count_tokens messages.batches.create messages.batches.retrieve messages.batches.list messages.batches.delete messages.batches.cancel models.retrieve models.list beta.models.retrieve beta.models.list beta.messages.create beta.messages.count_tokens beta.messages.batches.create beta.messages.batches.retrieve beta.messages.batches.list beta.messages.batches.delete beta.messages.batches.cancel
complete -c anthropic-cli --no-files \
  -n "not __fish_seen_subcommand_from $subcommands" \
  -a "$subcommands"

complete -c anthropic-cli --no-files \
  -n "__fish_seen_subcommand_from completions.create" \
  -a "--max-tokens-to-sample --model --prompt --metadata.user_id --stop-sequences --+stop_sequence --temperature --top-k --top-p"
complete -c anthropic-cli --no-files \
  -n "__fish_seen_subcommand_from messages.create" \
  -a "--max-tokens --messages.content.text --messages.content.type --messages.content.cache_control.type --messages.content.source.data --messages.content.source.media_type --messages.content.source.type --messages.content.id --messages.content.name --messages.content.tool_use_id --messages.content.content.text --messages.content.content.type --messages.content.content.cache_control.type --messages.content.content.source.data --messages.content.content.source.media_type --messages.content.content.source.type --messages.content.+content --messages.content.is_error --messages.+content --messages.role --+message --model --metadata.user_id --stop-sequences --+stop_sequence --system.text --system.type --system.cache_control.type --+system --temperature --tool-choice.type --tool-choice.disable_parallel_tool_use --tool-choice.name --tools.name --tools.cache_control.type --tools.description --+tool --top-k --top-p"
complete -c anthropic-cli --no-files \
  -n "__fish_seen_subcommand_from messages.count_tokens" \
  -a "--messages.content.text --messages.content.type --messages.content.cache_control.type --messages.content.source.data --messages.content.source.media_type --messages.content.source.type --messages.content.id --messages.content.name --messages.content.tool_use_id --messages.content.content.text --messages.content.content.type --messages.content.content.cache_control.type --messages.content.content.source.data --messages.content.content.source.media_type --messages.content.content.source.type --messages.content.+content --messages.content.is_error --messages.+content --messages.role --+message --model --system --system.text --system.type --system.cache_control.type --+system --tool-choice.type --tool-choice.disable_parallel_tool_use --tool-choice.name --tools.name --tools.cache_control.type --tools.description --+tool"
complete -c anthropic-cli --no-files \
  -n "__fish_seen_subcommand_from messages.batches.create" \
  -a "--requests.custom_id --requests.params.max_tokens --requests.params.messages.content.text --requests.params.messages.content.type --requests.params.messages.content.cache_control.type --requests.params.messages.content.source.data --requests.params.messages.content.source.media_type --requests.params.messages.content.source.type --requests.params.messages.content.id --requests.params.messages.content.name --requests.params.messages.content.tool_use_id --requests.params.messages.content.content.text --requests.params.messages.content.content.type --requests.params.messages.content.content.cache_control.type --requests.params.messages.content.content.source.data --requests.params.messages.content.content.source.media_type --requests.params.messages.content.content.source.type --requests.params.messages.content.+content --requests.params.messages.content.is_error --requests.params.messages.+content --requests.params.messages.role --requests.params.+message --requests.params.model --requests.params.metadata.user_id --requests.params.stop_sequences --requests.params.+stop_sequence --requests.params.stream --requests.params.system.text --requests.params.system.type --requests.params.system.cache_control.type --requests.params.+system --requests.params.temperature --requests.params.tool_choice.type --requests.params.tool_choice.disable_parallel_tool_use --requests.params.tool_choice.name --requests.params.tools.name --requests.params.tools.cache_control.type --requests.params.tools.description --requests.params.+tool --requests.params.top_k --requests.params.top_p --+request"
complete -c anthropic-cli --no-files \
  -n "__fish_seen_subcommand_from messages.batches.retrieve" \
  -a "--message-batch-id"
complete -c anthropic-cli --no-files \
  -n "__fish_seen_subcommand_from messages.batches.list" \
  -a "--after-id --before-id --limit"
complete -c anthropic-cli --no-files \
  -n "__fish_seen_subcommand_from messages.batches.delete" \
  -a "--message-batch-id"
complete -c anthropic-cli --no-files \
  -n "__fish_seen_subcommand_from messages.batches.cancel" \
  -a "--message-batch-id"
complete -c anthropic-cli --no-files \
  -n "__fish_seen_subcommand_from models.retrieve" \
  -a "--model-id"
complete -c anthropic-cli --no-files \
  -n "__fish_seen_subcommand_from models.list" \
  -a "--after-id --before-id --limit"
complete -c anthropic-cli --no-files \
  -n "__fish_seen_subcommand_from beta.models.retrieve" \
  -a "--model-id"
complete -c anthropic-cli --no-files \
  -n "__fish_seen_subcommand_from beta.models.list" \
  -a "--after-id --before-id --limit"
complete -c anthropic-cli --no-files \
  -n "__fish_seen_subcommand_from beta.messages.create" \
  -a "--max-tokens --messages.content.text --messages.content.type --messages.content.cache_control.type --messages.content.source.data --messages.content.source.media_type --messages.content.source.type --messages.content.id --messages.content.name --messages.content.tool_use_id --messages.content.content.text --messages.content.content.type --messages.content.content.cache_control.type --messages.content.content.source.data --messages.content.content.source.media_type --messages.content.content.source.type --messages.content.+content --messages.content.is_error --messages.+content --messages.role --+message --model --metadata.user_id --stop-sequences --+stop_sequence --system.text --system.type --system.cache_control.type --+system --temperature --tool-choice.type --tool-choice.disable_parallel_tool_use --tool-choice.name --tools.input_schema.type --tools.name --tools.cache_control.type --tools.description --tools.type --tools.display_height_px --tools.display_width_px --tools.display_number --+tool --top-k --top-p --betas --+beta"
complete -c anthropic-cli --no-files \
  -n "__fish_seen_subcommand_from beta.messages.count_tokens" \
  -a "--messages.content.text --messages.content.type --messages.content.cache_control.type --messages.content.source.data --messages.content.source.media_type --messages.content.source.type --messages.content.id --messages.content.name --messages.content.tool_use_id --messages.content.content.text --messages.content.content.type --messages.content.content.cache_control.type --messages.content.content.source.data --messages.content.content.source.media_type --messages.content.content.source.type --messages.content.+content --messages.content.is_error --messages.+content --messages.role --+message --model --system --system.text --system.type --system.cache_control.type --+system --tool-choice.type --tool-choice.disable_parallel_tool_use --tool-choice.name --tools.input_schema.type --tools.name --tools.cache_control.type --tools.description --tools.type --tools.display_height_px --tools.display_width_px --tools.display_number --+tool --betas --+beta"
complete -c anthropic-cli --no-files \
  -n "__fish_seen_subcommand_from beta.messages.batches.create" \
  -a "--requests.custom_id --requests.params.max_tokens --requests.params.messages.content.text --requests.params.messages.content.type --requests.params.messages.content.cache_control.type --requests.params.messages.content.source.data --requests.params.messages.content.source.media_type --requests.params.messages.content.source.type --requests.params.messages.content.id --requests.params.messages.content.name --requests.params.messages.content.tool_use_id --requests.params.messages.content.content.text --requests.params.messages.content.content.type --requests.params.messages.content.content.cache_control.type --requests.params.messages.content.content.source.data --requests.params.messages.content.content.source.media_type --requests.params.messages.content.content.source.type --requests.params.messages.content.+content --requests.params.messages.content.is_error --requests.params.messages.+content --requests.params.messages.role --requests.params.+message --requests.params.model --requests.params.metadata.user_id --requests.params.stop_sequences --requests.params.+stop_sequence --requests.params.stream --requests.params.system.text --requests.params.system.type --requests.params.system.cache_control.type --requests.params.+system --requests.params.temperature --requests.params.tool_choice.type --requests.params.tool_choice.disable_parallel_tool_use --requests.params.tool_choice.name --requests.params.tools.input_schema.type --requests.params.tools.name --requests.params.tools.cache_control.type --requests.params.tools.description --requests.params.tools.type --requests.params.tools.display_height_px --requests.params.tools.display_width_px --requests.params.tools.display_number --requests.params.+tool --requests.params.top_k --requests.params.top_p --+request --betas --+beta"
complete -c anthropic-cli --no-files \
  -n "__fish_seen_subcommand_from beta.messages.batches.retrieve" \
  -a "--message-batch-id --betas --+beta"
complete -c anthropic-cli --no-files \
  -n "__fish_seen_subcommand_from beta.messages.batches.list" \
  -a "--after-id --before-id --limit --betas --+beta"
complete -c anthropic-cli --no-files \
  -n "__fish_seen_subcommand_from beta.messages.batches.delete" \
  -a "--message-batch-id --betas --+beta"
complete -c anthropic-cli --no-files \
  -n "__fish_seen_subcommand_from beta.messages.batches.cancel" \
  -a "--message-batch-id --betas --+beta"

 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from completions.create" \
   -l model \
   -ra "claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.create" \
   -l messages.content.type \
   -ra "text image tool_use tool_result document"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.create" \
   -l messages.content.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.create" \
   -l messages.content.source.media_type \
   -ra "image/jpeg image/png image/gif image/webp application/pdf"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.create" \
   -l messages.content.source.type \
   -ra "base64"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.create" \
   -l messages.content.content.type \
   -ra "text image"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.create" \
   -l messages.content.content.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.create" \
   -l messages.content.content.source.media_type \
   -ra "image/jpeg image/png image/gif image/webp"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.create" \
   -l messages.content.content.source.type \
   -ra "base64"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.create" \
   -l messages.role \
   -ra "user assistant"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.create" \
   -l model \
   -ra "claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.create" \
   -l system.type \
   -ra "text"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.create" \
   -l system.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.create" \
   -l tool-choice.type \
   -ra "auto any tool"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.create" \
   -l tools.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.count_tokens" \
   -l messages.content.type \
   -ra "text image tool_use tool_result document"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.count_tokens" \
   -l messages.content.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.count_tokens" \
   -l messages.content.source.media_type \
   -ra "image/jpeg image/png image/gif image/webp application/pdf"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.count_tokens" \
   -l messages.content.source.type \
   -ra "base64"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.count_tokens" \
   -l messages.content.content.type \
   -ra "text image"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.count_tokens" \
   -l messages.content.content.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.count_tokens" \
   -l messages.content.content.source.media_type \
   -ra "image/jpeg image/png image/gif image/webp"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.count_tokens" \
   -l messages.content.content.source.type \
   -ra "base64"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.count_tokens" \
   -l messages.role \
   -ra "user assistant"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.count_tokens" \
   -l model \
   -ra "claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.count_tokens" \
   -l system.type \
   -ra "text"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.count_tokens" \
   -l system.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.count_tokens" \
   -l tool-choice.type \
   -ra "auto any tool"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.count_tokens" \
   -l tools.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.batches.create" \
   -l requests.params.messages.content.type \
   -ra "text image tool_use tool_result document"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.batches.create" \
   -l requests.params.messages.content.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.batches.create" \
   -l requests.params.messages.content.source.media_type \
   -ra "image/jpeg image/png image/gif image/webp application/pdf"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.batches.create" \
   -l requests.params.messages.content.source.type \
   -ra "base64"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.batches.create" \
   -l requests.params.messages.content.content.type \
   -ra "text image"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.batches.create" \
   -l requests.params.messages.content.content.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.batches.create" \
   -l requests.params.messages.content.content.source.media_type \
   -ra "image/jpeg image/png image/gif image/webp"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.batches.create" \
   -l requests.params.messages.content.content.source.type \
   -ra "base64"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.batches.create" \
   -l requests.params.messages.role \
   -ra "user assistant"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.batches.create" \
   -l requests.params.model \
   -ra "claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.batches.create" \
   -l requests.params.system.type \
   -ra "text"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.batches.create" \
   -l requests.params.system.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.batches.create" \
   -l requests.params.tool_choice.type \
   -ra "auto any tool"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from messages.batches.create" \
   -l requests.params.tools.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.create" \
   -l messages.content.type \
   -ra "text image tool_use tool_result document"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.create" \
   -l messages.content.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.create" \
   -l messages.content.source.media_type \
   -ra "image/jpeg image/png image/gif image/webp application/pdf"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.create" \
   -l messages.content.source.type \
   -ra "base64"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.create" \
   -l messages.content.content.type \
   -ra "text image"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.create" \
   -l messages.content.content.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.create" \
   -l messages.content.content.source.media_type \
   -ra "image/jpeg image/png image/gif image/webp"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.create" \
   -l messages.content.content.source.type \
   -ra "base64"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.create" \
   -l messages.role \
   -ra "user assistant"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.create" \
   -l model \
   -ra "claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.create" \
   -l system.type \
   -ra "text"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.create" \
   -l system.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.create" \
   -l tool-choice.type \
   -ra "auto any tool"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.create" \
   -l tools.input_schema.type \
   -ra "object"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.create" \
   -l tools.name \
   -ra "computer bash str_replace_editor"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.create" \
   -l tools.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.create" \
   -l tools.type \
   -ra "custom computer_20241022 bash_20241022 text_editor_20241022"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.create" \
   -l betas \
   -ra "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.create" \
   -l +beta \
   -ra "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.count_tokens" \
   -l messages.content.type \
   -ra "text image tool_use tool_result document"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.count_tokens" \
   -l messages.content.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.count_tokens" \
   -l messages.content.source.media_type \
   -ra "image/jpeg image/png image/gif image/webp application/pdf"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.count_tokens" \
   -l messages.content.source.type \
   -ra "base64"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.count_tokens" \
   -l messages.content.content.type \
   -ra "text image"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.count_tokens" \
   -l messages.content.content.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.count_tokens" \
   -l messages.content.content.source.media_type \
   -ra "image/jpeg image/png image/gif image/webp"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.count_tokens" \
   -l messages.content.content.source.type \
   -ra "base64"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.count_tokens" \
   -l messages.role \
   -ra "user assistant"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.count_tokens" \
   -l model \
   -ra "claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.count_tokens" \
   -l system.type \
   -ra "text"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.count_tokens" \
   -l system.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.count_tokens" \
   -l tool-choice.type \
   -ra "auto any tool"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.count_tokens" \
   -l tools.input_schema.type \
   -ra "object"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.count_tokens" \
   -l tools.name \
   -ra "computer bash str_replace_editor"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.count_tokens" \
   -l tools.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.count_tokens" \
   -l tools.type \
   -ra "custom computer_20241022 bash_20241022 text_editor_20241022"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.count_tokens" \
   -l betas \
   -ra "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.count_tokens" \
   -l +beta \
   -ra "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.create" \
   -l requests.params.messages.content.type \
   -ra "text image tool_use tool_result document"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.create" \
   -l requests.params.messages.content.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.create" \
   -l requests.params.messages.content.source.media_type \
   -ra "image/jpeg image/png image/gif image/webp application/pdf"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.create" \
   -l requests.params.messages.content.source.type \
   -ra "base64"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.create" \
   -l requests.params.messages.content.content.type \
   -ra "text image"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.create" \
   -l requests.params.messages.content.content.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.create" \
   -l requests.params.messages.content.content.source.media_type \
   -ra "image/jpeg image/png image/gif image/webp"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.create" \
   -l requests.params.messages.content.content.source.type \
   -ra "base64"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.create" \
   -l requests.params.messages.role \
   -ra "user assistant"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.create" \
   -l requests.params.model \
   -ra "claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.create" \
   -l requests.params.system.type \
   -ra "text"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.create" \
   -l requests.params.system.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.create" \
   -l requests.params.tool_choice.type \
   -ra "auto any tool"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.create" \
   -l requests.params.tools.input_schema.type \
   -ra "object"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.create" \
   -l requests.params.tools.name \
   -ra "computer bash str_replace_editor"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.create" \
   -l requests.params.tools.cache_control.type \
   -ra "ephemeral"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.create" \
   -l requests.params.tools.type \
   -ra "custom computer_20241022 bash_20241022 text_editor_20241022"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.create" \
   -l betas \
   -ra "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.create" \
   -l +beta \
   -ra "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.retrieve" \
   -l betas \
   -ra "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.retrieve" \
   -l +beta \
   -ra "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.list" \
   -l betas \
   -ra "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.list" \
   -l +beta \
   -ra "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.delete" \
   -l betas \
   -ra "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.delete" \
   -l +beta \
   -ra "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.cancel" \
   -l betas \
   -ra "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01"
 complete -c anthropic-cli --no-files \
   -n "__fish_seen_subcommand_from beta.messages.batches.cancel" \
   -l +beta \
   -ra "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01"