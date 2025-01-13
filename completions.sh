#!/usr/bin/env bash

# This also works for zsh: https://zsh.sourceforge.io/Doc/Release/Completion-System.html#Completion-System
_main()
{
    COMPREPLY=()

    local subcommands="completions.create messages.create messages.count_tokens messages.batches.create messages.batches.retrieve messages.batches.list messages.batches.delete messages.batches.cancel models.retrieve models.list beta.models.retrieve beta.models.list beta.messages.create beta.messages.count_tokens beta.messages.batches.create beta.messages.batches.retrieve beta.messages.batches.list beta.messages.batches.delete beta.messages.batches.cancel"

    if [[ "$COMP_CWORD" -eq 1 ]]
    then
      local cur="${COMP_WORDS[COMP_CWORD]}"
      COMPREPLY=( $(compgen -W "$subcommands" -- "$cur") )
      return
    fi

    local subcommand="${COMP_WORDS[1]}"
    local flags
    case "$subcommand" in
      completions.create)
        flags="--max-tokens-to-sample --model --prompt --metadata.user_id --stop-sequences --+stop_sequence --temperature --top-k --top-p"
        ;;
      messages.create)
        flags="--max-tokens --messages.content.text --messages.content.type --messages.content.cache_control.type --messages.content.source.data --messages.content.source.media_type --messages.content.source.type --messages.content.id --messages.content.name --messages.content.tool_use_id --messages.content.content.text --messages.content.content.type --messages.content.content.cache_control.type --messages.content.content.source.data --messages.content.content.source.media_type --messages.content.content.source.type --messages.content.+content --messages.content.is_error --messages.+content --messages.role --+message --model --metadata.user_id --stop-sequences --+stop_sequence --system.text --system.type --system.cache_control.type --+system --temperature --tool-choice.type --tool-choice.disable_parallel_tool_use --tool-choice.name --tools.name --tools.cache_control.type --tools.description --+tool --top-k --top-p"
        ;;
      messages.count_tokens)
        flags="--messages.content.text --messages.content.type --messages.content.cache_control.type --messages.content.source.data --messages.content.source.media_type --messages.content.source.type --messages.content.id --messages.content.name --messages.content.tool_use_id --messages.content.content.text --messages.content.content.type --messages.content.content.cache_control.type --messages.content.content.source.data --messages.content.content.source.media_type --messages.content.content.source.type --messages.content.+content --messages.content.is_error --messages.+content --messages.role --+message --model --system --system.text --system.type --system.cache_control.type --+system --tool-choice.type --tool-choice.disable_parallel_tool_use --tool-choice.name --tools.name --tools.cache_control.type --tools.description --+tool"
        ;;
      messages.batches.create)
        flags="--requests.custom_id --requests.params.max_tokens --requests.params.messages.content.text --requests.params.messages.content.type --requests.params.messages.content.cache_control.type --requests.params.messages.content.source.data --requests.params.messages.content.source.media_type --requests.params.messages.content.source.type --requests.params.messages.content.id --requests.params.messages.content.name --requests.params.messages.content.tool_use_id --requests.params.messages.content.content.text --requests.params.messages.content.content.type --requests.params.messages.content.content.cache_control.type --requests.params.messages.content.content.source.data --requests.params.messages.content.content.source.media_type --requests.params.messages.content.content.source.type --requests.params.messages.content.+content --requests.params.messages.content.is_error --requests.params.messages.+content --requests.params.messages.role --requests.params.+message --requests.params.model --requests.params.metadata.user_id --requests.params.stop_sequences --requests.params.+stop_sequence --requests.params.stream --requests.params.system.text --requests.params.system.type --requests.params.system.cache_control.type --requests.params.+system --requests.params.temperature --requests.params.tool_choice.type --requests.params.tool_choice.disable_parallel_tool_use --requests.params.tool_choice.name --requests.params.tools.name --requests.params.tools.cache_control.type --requests.params.tools.description --requests.params.+tool --requests.params.top_k --requests.params.top_p --+request"
        ;;
      messages.batches.retrieve)
        flags="--message-batch-id"
        ;;
      messages.batches.list)
        flags="--after-id --before-id --limit"
        ;;
      messages.batches.delete)
        flags="--message-batch-id"
        ;;
      messages.batches.cancel)
        flags="--message-batch-id"
        ;;
      models.retrieve)
        flags="--model-id"
        ;;
      models.list)
        flags="--after-id --before-id --limit"
        ;;
      beta.models.retrieve)
        flags="--model-id"
        ;;
      beta.models.list)
        flags="--after-id --before-id --limit"
        ;;
      beta.messages.create)
        flags="--max-tokens --messages.content.text --messages.content.type --messages.content.cache_control.type --messages.content.source.data --messages.content.source.media_type --messages.content.source.type --messages.content.id --messages.content.name --messages.content.tool_use_id --messages.content.content.text --messages.content.content.type --messages.content.content.cache_control.type --messages.content.content.source.data --messages.content.content.source.media_type --messages.content.content.source.type --messages.content.+content --messages.content.is_error --messages.+content --messages.role --+message --model --metadata.user_id --stop-sequences --+stop_sequence --system.text --system.type --system.cache_control.type --+system --temperature --tool-choice.type --tool-choice.disable_parallel_tool_use --tool-choice.name --tools.input_schema.type --tools.name --tools.cache_control.type --tools.description --tools.type --tools.display_height_px --tools.display_width_px --tools.display_number --+tool --top-k --top-p --betas --+beta"
        ;;
      beta.messages.count_tokens)
        flags="--messages.content.text --messages.content.type --messages.content.cache_control.type --messages.content.source.data --messages.content.source.media_type --messages.content.source.type --messages.content.id --messages.content.name --messages.content.tool_use_id --messages.content.content.text --messages.content.content.type --messages.content.content.cache_control.type --messages.content.content.source.data --messages.content.content.source.media_type --messages.content.content.source.type --messages.content.+content --messages.content.is_error --messages.+content --messages.role --+message --model --system --system.text --system.type --system.cache_control.type --+system --tool-choice.type --tool-choice.disable_parallel_tool_use --tool-choice.name --tools.input_schema.type --tools.name --tools.cache_control.type --tools.description --tools.type --tools.display_height_px --tools.display_width_px --tools.display_number --+tool --betas --+beta"
        ;;
      beta.messages.batches.create)
        flags="--requests.custom_id --requests.params.max_tokens --requests.params.messages.content.text --requests.params.messages.content.type --requests.params.messages.content.cache_control.type --requests.params.messages.content.source.data --requests.params.messages.content.source.media_type --requests.params.messages.content.source.type --requests.params.messages.content.id --requests.params.messages.content.name --requests.params.messages.content.tool_use_id --requests.params.messages.content.content.text --requests.params.messages.content.content.type --requests.params.messages.content.content.cache_control.type --requests.params.messages.content.content.source.data --requests.params.messages.content.content.source.media_type --requests.params.messages.content.content.source.type --requests.params.messages.content.+content --requests.params.messages.content.is_error --requests.params.messages.+content --requests.params.messages.role --requests.params.+message --requests.params.model --requests.params.metadata.user_id --requests.params.stop_sequences --requests.params.+stop_sequence --requests.params.stream --requests.params.system.text --requests.params.system.type --requests.params.system.cache_control.type --requests.params.+system --requests.params.temperature --requests.params.tool_choice.type --requests.params.tool_choice.disable_parallel_tool_use --requests.params.tool_choice.name --requests.params.tools.input_schema.type --requests.params.tools.name --requests.params.tools.cache_control.type --requests.params.tools.description --requests.params.tools.type --requests.params.tools.display_height_px --requests.params.tools.display_width_px --requests.params.tools.display_number --requests.params.+tool --requests.params.top_k --requests.params.top_p --+request --betas --+beta"
        ;;
      beta.messages.batches.retrieve)
        flags="--message-batch-id --betas --+beta"
        ;;
      beta.messages.batches.list)
        flags="--after-id --before-id --limit --betas --+beta"
        ;;
      beta.messages.batches.delete)
        flags="--message-batch-id --betas --+beta"
        ;;
      beta.messages.batches.cancel)
        flags="--message-batch-id --betas --+beta"
        ;;
      *)
        # Unknown subcommand
        return
        ;;
    esac

    local cur="${COMP_WORDS[COMP_CWORD]}"
    if [[ "$COMP_CWORD" -eq 2 || $cur == -* ]] ; then
        COMPREPLY=( $(compgen -W "$flags" -- $cur) )
        return 0
    fi

    local prev="${COMP_WORDS[COMP_CWORD-1]}"
    case "$subcommand" in
      completions.create)
        case "$prev" in
          --model)
            COMPREPLY=( $(compgen -W "claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0" -- $cur) )
            ;;
        esac
        ;;
      messages.create)
        case "$prev" in
          --messages.content.type)
            COMPREPLY=( $(compgen -W "text image tool_use tool_result document" -- $cur) )
            ;;
          --messages.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --messages.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp application/pdf" -- $cur) )
            ;;
          --messages.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --messages.content.content.type)
            COMPREPLY=( $(compgen -W "text image" -- $cur) )
            ;;
          --messages.content.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --messages.content.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp" -- $cur) )
            ;;
          --messages.content.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --messages.role)
            COMPREPLY=( $(compgen -W "user assistant" -- $cur) )
            ;;
          --model)
            COMPREPLY=( $(compgen -W "claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0" -- $cur) )
            ;;
          --system.type)
            COMPREPLY=( $(compgen -W "text" -- $cur) )
            ;;
          --system.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --tool-choice.type)
            COMPREPLY=( $(compgen -W "auto any tool" -- $cur) )
            ;;
          --tools.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
        esac
        ;;
      messages.count_tokens)
        case "$prev" in
          --messages.content.type)
            COMPREPLY=( $(compgen -W "text image tool_use tool_result document" -- $cur) )
            ;;
          --messages.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --messages.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp application/pdf" -- $cur) )
            ;;
          --messages.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --messages.content.content.type)
            COMPREPLY=( $(compgen -W "text image" -- $cur) )
            ;;
          --messages.content.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --messages.content.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp" -- $cur) )
            ;;
          --messages.content.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --messages.role)
            COMPREPLY=( $(compgen -W "user assistant" -- $cur) )
            ;;
          --model)
            COMPREPLY=( $(compgen -W "claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0" -- $cur) )
            ;;
          --system.type)
            COMPREPLY=( $(compgen -W "text" -- $cur) )
            ;;
          --system.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --tool-choice.type)
            COMPREPLY=( $(compgen -W "auto any tool" -- $cur) )
            ;;
          --tools.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
        esac
        ;;
      messages.batches.create)
        case "$prev" in
          --requests.params.messages.content.type)
            COMPREPLY=( $(compgen -W "text image tool_use tool_result document" -- $cur) )
            ;;
          --requests.params.messages.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --requests.params.messages.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp application/pdf" -- $cur) )
            ;;
          --requests.params.messages.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --requests.params.messages.content.content.type)
            COMPREPLY=( $(compgen -W "text image" -- $cur) )
            ;;
          --requests.params.messages.content.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --requests.params.messages.content.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp" -- $cur) )
            ;;
          --requests.params.messages.content.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --requests.params.messages.role)
            COMPREPLY=( $(compgen -W "user assistant" -- $cur) )
            ;;
          --requests.params.model)
            COMPREPLY=( $(compgen -W "claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0" -- $cur) )
            ;;
          --requests.params.system.type)
            COMPREPLY=( $(compgen -W "text" -- $cur) )
            ;;
          --requests.params.system.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --requests.params.tool_choice.type)
            COMPREPLY=( $(compgen -W "auto any tool" -- $cur) )
            ;;
          --requests.params.tools.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
        esac
        ;;
      beta.messages.create)
        case "$prev" in
          --messages.content.type)
            COMPREPLY=( $(compgen -W "text image tool_use tool_result document" -- $cur) )
            ;;
          --messages.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --messages.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp application/pdf" -- $cur) )
            ;;
          --messages.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --messages.content.content.type)
            COMPREPLY=( $(compgen -W "text image" -- $cur) )
            ;;
          --messages.content.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --messages.content.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp" -- $cur) )
            ;;
          --messages.content.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --messages.role)
            COMPREPLY=( $(compgen -W "user assistant" -- $cur) )
            ;;
          --model)
            COMPREPLY=( $(compgen -W "claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0" -- $cur) )
            ;;
          --system.type)
            COMPREPLY=( $(compgen -W "text" -- $cur) )
            ;;
          --system.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --tool-choice.type)
            COMPREPLY=( $(compgen -W "auto any tool" -- $cur) )
            ;;
          --tools.input_schema.type)
            COMPREPLY=( $(compgen -W "object" -- $cur) )
            ;;
          --tools.name)
            COMPREPLY=( $(compgen -W "computer bash str_replace_editor" -- $cur) )
            ;;
          --tools.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --tools.type)
            COMPREPLY=( $(compgen -W "custom computer_20241022 bash_20241022 text_editor_20241022" -- $cur) )
            ;;
          --betas)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01" -- $cur) )
            ;;
          --+beta)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01" -- $cur) )
            ;;
        esac
        ;;
      beta.messages.count_tokens)
        case "$prev" in
          --messages.content.type)
            COMPREPLY=( $(compgen -W "text image tool_use tool_result document" -- $cur) )
            ;;
          --messages.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --messages.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp application/pdf" -- $cur) )
            ;;
          --messages.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --messages.content.content.type)
            COMPREPLY=( $(compgen -W "text image" -- $cur) )
            ;;
          --messages.content.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --messages.content.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp" -- $cur) )
            ;;
          --messages.content.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --messages.role)
            COMPREPLY=( $(compgen -W "user assistant" -- $cur) )
            ;;
          --model)
            COMPREPLY=( $(compgen -W "claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0" -- $cur) )
            ;;
          --system.type)
            COMPREPLY=( $(compgen -W "text" -- $cur) )
            ;;
          --system.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --tool-choice.type)
            COMPREPLY=( $(compgen -W "auto any tool" -- $cur) )
            ;;
          --tools.input_schema.type)
            COMPREPLY=( $(compgen -W "object" -- $cur) )
            ;;
          --tools.name)
            COMPREPLY=( $(compgen -W "computer bash str_replace_editor" -- $cur) )
            ;;
          --tools.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --tools.type)
            COMPREPLY=( $(compgen -W "custom computer_20241022 bash_20241022 text_editor_20241022" -- $cur) )
            ;;
          --betas)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01" -- $cur) )
            ;;
          --+beta)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01" -- $cur) )
            ;;
        esac
        ;;
      beta.messages.batches.create)
        case "$prev" in
          --requests.params.messages.content.type)
            COMPREPLY=( $(compgen -W "text image tool_use tool_result document" -- $cur) )
            ;;
          --requests.params.messages.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --requests.params.messages.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp application/pdf" -- $cur) )
            ;;
          --requests.params.messages.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --requests.params.messages.content.content.type)
            COMPREPLY=( $(compgen -W "text image" -- $cur) )
            ;;
          --requests.params.messages.content.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --requests.params.messages.content.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp" -- $cur) )
            ;;
          --requests.params.messages.content.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --requests.params.messages.role)
            COMPREPLY=( $(compgen -W "user assistant" -- $cur) )
            ;;
          --requests.params.model)
            COMPREPLY=( $(compgen -W "claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0" -- $cur) )
            ;;
          --requests.params.system.type)
            COMPREPLY=( $(compgen -W "text" -- $cur) )
            ;;
          --requests.params.system.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --requests.params.tool_choice.type)
            COMPREPLY=( $(compgen -W "auto any tool" -- $cur) )
            ;;
          --requests.params.tools.input_schema.type)
            COMPREPLY=( $(compgen -W "object" -- $cur) )
            ;;
          --requests.params.tools.name)
            COMPREPLY=( $(compgen -W "computer bash str_replace_editor" -- $cur) )
            ;;
          --requests.params.tools.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --requests.params.tools.type)
            COMPREPLY=( $(compgen -W "custom computer_20241022 bash_20241022 text_editor_20241022" -- $cur) )
            ;;
          --betas)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01" -- $cur) )
            ;;
          --+beta)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01" -- $cur) )
            ;;
        esac
        ;;
      beta.messages.batches.retrieve)
        case "$prev" in
          --betas)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01" -- $cur) )
            ;;
          --+beta)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01" -- $cur) )
            ;;
        esac
        ;;
      beta.messages.batches.list)
        case "$prev" in
          --betas)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01" -- $cur) )
            ;;
          --+beta)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01" -- $cur) )
            ;;
        esac
        ;;
      beta.messages.batches.delete)
        case "$prev" in
          --betas)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01" -- $cur) )
            ;;
          --+beta)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01" -- $cur) )
            ;;
        esac
        ;;
      beta.messages.batches.cancel)
        case "$prev" in
          --betas)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01" -- $cur) )
            ;;
          --+beta)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 pdfs-2024-09-25 token-counting-2024-11-01" -- $cur) )
            ;;
        esac
        ;;
    esac
}
complete -F _main anthropic-cli