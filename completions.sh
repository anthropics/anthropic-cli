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
        flags="--max-tokens --messages.content.text --messages.content.type --messages.content.cache_control.type --messages.content.citations.cited_text --messages.content.citations.document_index --messages.content.citations.document_title --messages.content.citations.end_char_index --messages.content.citations.start_char_index --messages.content.citations.type --messages.content.citations.end_page_number --messages.content.citations.start_page_number --messages.content.citations.end_block_index --messages.content.citations.start_block_index --messages.content.+citation --messages.content.source.data --messages.content.source.media_type --messages.content.source.type --messages.content.id --messages.content.name --messages.content.tool_use_id --messages.content.content.text --messages.content.content.type --messages.content.content.cache_control.type --messages.content.content.citations.cited_text --messages.content.content.citations.document_index --messages.content.content.citations.document_title --messages.content.content.citations.end_char_index --messages.content.content.citations.start_char_index --messages.content.content.citations.type --messages.content.content.citations.end_page_number --messages.content.content.citations.start_page_number --messages.content.content.citations.end_block_index --messages.content.content.citations.start_block_index --messages.content.content.+citation --messages.content.content.source.data --messages.content.content.source.media_type --messages.content.content.source.type --messages.content.+content --messages.content.is_error --messages.content.source.content --messages.content.source.content.text --messages.content.source.content.type --messages.content.source.content.cache_control.type --messages.content.source.content.citations.cited_text --messages.content.source.content.citations.document_index --messages.content.source.content.citations.document_title --messages.content.source.content.citations.end_char_index --messages.content.source.content.citations.start_char_index --messages.content.source.content.citations.type --messages.content.source.content.citations.end_page_number --messages.content.source.content.citations.start_page_number --messages.content.source.content.citations.end_block_index --messages.content.source.content.citations.start_block_index --messages.content.source.content.+citation --messages.content.source.content.source.data --messages.content.source.content.source.media_type --messages.content.source.content.source.type --messages.content.source.+content --messages.content.citations.enabled --messages.content.context --messages.content.title --messages.content.signature --messages.content.thinking --messages.content.data --messages.+content --messages.role --+message --model --metadata.user_id --stop-sequences --+stop_sequence --system.text --system.type --system.cache_control.type --system.citations.cited_text --system.citations.document_index --system.citations.document_title --system.citations.end_char_index --system.citations.start_char_index --system.citations.type --system.citations.end_page_number --system.citations.start_page_number --system.citations.end_block_index --system.citations.start_block_index --system.+citation --+system --temperature --thinking.budget_tokens --thinking.type --tool-choice.type --tool-choice.disable_parallel_tool_use --tool-choice.name --tools.name --tools.type --tools.cache_control.type --tools.description --+tool --top-k --top-p"
        ;;
      messages.count_tokens)
        flags="--messages.content.text --messages.content.type --messages.content.cache_control.type --messages.content.citations.cited_text --messages.content.citations.document_index --messages.content.citations.document_title --messages.content.citations.end_char_index --messages.content.citations.start_char_index --messages.content.citations.type --messages.content.citations.end_page_number --messages.content.citations.start_page_number --messages.content.citations.end_block_index --messages.content.citations.start_block_index --messages.content.+citation --messages.content.source.data --messages.content.source.media_type --messages.content.source.type --messages.content.id --messages.content.name --messages.content.tool_use_id --messages.content.content.text --messages.content.content.type --messages.content.content.cache_control.type --messages.content.content.citations.cited_text --messages.content.content.citations.document_index --messages.content.content.citations.document_title --messages.content.content.citations.end_char_index --messages.content.content.citations.start_char_index --messages.content.content.citations.type --messages.content.content.citations.end_page_number --messages.content.content.citations.start_page_number --messages.content.content.citations.end_block_index --messages.content.content.citations.start_block_index --messages.content.content.+citation --messages.content.content.source.data --messages.content.content.source.media_type --messages.content.content.source.type --messages.content.+content --messages.content.is_error --messages.content.source.content --messages.content.source.content.text --messages.content.source.content.type --messages.content.source.content.cache_control.type --messages.content.source.content.citations.cited_text --messages.content.source.content.citations.document_index --messages.content.source.content.citations.document_title --messages.content.source.content.citations.end_char_index --messages.content.source.content.citations.start_char_index --messages.content.source.content.citations.type --messages.content.source.content.citations.end_page_number --messages.content.source.content.citations.start_page_number --messages.content.source.content.citations.end_block_index --messages.content.source.content.citations.start_block_index --messages.content.source.content.+citation --messages.content.source.content.source.data --messages.content.source.content.source.media_type --messages.content.source.content.source.type --messages.content.source.+content --messages.content.citations.enabled --messages.content.context --messages.content.title --messages.content.signature --messages.content.thinking --messages.content.data --messages.+content --messages.role --+message --model --system --system.text --system.type --system.cache_control.type --system.citations.cited_text --system.citations.document_index --system.citations.document_title --system.citations.end_char_index --system.citations.start_char_index --system.citations.type --system.citations.end_page_number --system.citations.start_page_number --system.citations.end_block_index --system.citations.start_block_index --system.+citation --+system --thinking.budget_tokens --thinking.type --tool-choice.type --tool-choice.disable_parallel_tool_use --tool-choice.name --tools.name --tools.type --tools.cache_control.type --tools.description --+tool"
        ;;
      messages.batches.create)
        flags="--requests.custom_id --requests.params.max_tokens --requests.params.messages.content.text --requests.params.messages.content.type --requests.params.messages.content.cache_control.type --requests.params.messages.content.citations.cited_text --requests.params.messages.content.citations.document_index --requests.params.messages.content.citations.document_title --requests.params.messages.content.citations.end_char_index --requests.params.messages.content.citations.start_char_index --requests.params.messages.content.citations.type --requests.params.messages.content.citations.end_page_number --requests.params.messages.content.citations.start_page_number --requests.params.messages.content.citations.end_block_index --requests.params.messages.content.citations.start_block_index --requests.params.messages.content.+citation --requests.params.messages.content.source.data --requests.params.messages.content.source.media_type --requests.params.messages.content.source.type --requests.params.messages.content.id --requests.params.messages.content.name --requests.params.messages.content.tool_use_id --requests.params.messages.content.content.text --requests.params.messages.content.content.type --requests.params.messages.content.content.cache_control.type --requests.params.messages.content.content.citations.cited_text --requests.params.messages.content.content.citations.document_index --requests.params.messages.content.content.citations.document_title --requests.params.messages.content.content.citations.end_char_index --requests.params.messages.content.content.citations.start_char_index --requests.params.messages.content.content.citations.type --requests.params.messages.content.content.citations.end_page_number --requests.params.messages.content.content.citations.start_page_number --requests.params.messages.content.content.citations.end_block_index --requests.params.messages.content.content.citations.start_block_index --requests.params.messages.content.content.+citation --requests.params.messages.content.content.source.data --requests.params.messages.content.content.source.media_type --requests.params.messages.content.content.source.type --requests.params.messages.content.+content --requests.params.messages.content.is_error --requests.params.messages.content.source.content --requests.params.messages.content.source.content.text --requests.params.messages.content.source.content.type --requests.params.messages.content.source.content.cache_control.type --requests.params.messages.content.source.content.citations.cited_text --requests.params.messages.content.source.content.citations.document_index --requests.params.messages.content.source.content.citations.document_title --requests.params.messages.content.source.content.citations.end_char_index --requests.params.messages.content.source.content.citations.start_char_index --requests.params.messages.content.source.content.citations.type --requests.params.messages.content.source.content.citations.end_page_number --requests.params.messages.content.source.content.citations.start_page_number --requests.params.messages.content.source.content.citations.end_block_index --requests.params.messages.content.source.content.citations.start_block_index --requests.params.messages.content.source.content.+citation --requests.params.messages.content.source.content.source.data --requests.params.messages.content.source.content.source.media_type --requests.params.messages.content.source.content.source.type --requests.params.messages.content.source.+content --requests.params.messages.content.citations.enabled --requests.params.messages.content.context --requests.params.messages.content.title --requests.params.messages.content.signature --requests.params.messages.content.thinking --requests.params.messages.content.data --requests.params.messages.+content --requests.params.messages.role --requests.params.+message --requests.params.model --requests.params.metadata.user_id --requests.params.stop_sequences --requests.params.+stop_sequence --requests.params.stream --requests.params.system.text --requests.params.system.type --requests.params.system.cache_control.type --requests.params.system.citations.cited_text --requests.params.system.citations.document_index --requests.params.system.citations.document_title --requests.params.system.citations.end_char_index --requests.params.system.citations.start_char_index --requests.params.system.citations.type --requests.params.system.citations.end_page_number --requests.params.system.citations.start_page_number --requests.params.system.citations.end_block_index --requests.params.system.citations.start_block_index --requests.params.system.+citation --requests.params.+system --requests.params.temperature --requests.params.thinking.budget_tokens --requests.params.thinking.type --requests.params.tool_choice.type --requests.params.tool_choice.disable_parallel_tool_use --requests.params.tool_choice.name --requests.params.tools.name --requests.params.tools.type --requests.params.tools.cache_control.type --requests.params.tools.description --requests.params.+tool --requests.params.top_k --requests.params.top_p --+request"
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
        flags="--max-tokens --messages.content.text --messages.content.type --messages.content.cache_control.type --messages.content.citations.cited_text --messages.content.citations.document_index --messages.content.citations.document_title --messages.content.citations.end_char_index --messages.content.citations.start_char_index --messages.content.citations.type --messages.content.citations.end_page_number --messages.content.citations.start_page_number --messages.content.citations.end_block_index --messages.content.citations.start_block_index --messages.content.+citation --messages.content.source.data --messages.content.source.media_type --messages.content.source.type --messages.content.id --messages.content.name --messages.content.tool_use_id --messages.content.content.text --messages.content.content.type --messages.content.content.cache_control.type --messages.content.content.citations.cited_text --messages.content.content.citations.document_index --messages.content.content.citations.document_title --messages.content.content.citations.end_char_index --messages.content.content.citations.start_char_index --messages.content.content.citations.type --messages.content.content.citations.end_page_number --messages.content.content.citations.start_page_number --messages.content.content.citations.end_block_index --messages.content.content.citations.start_block_index --messages.content.content.+citation --messages.content.content.source.data --messages.content.content.source.media_type --messages.content.content.source.type --messages.content.+content --messages.content.is_error --messages.content.source.content --messages.content.source.content.text --messages.content.source.content.type --messages.content.source.content.cache_control.type --messages.content.source.content.citations.cited_text --messages.content.source.content.citations.document_index --messages.content.source.content.citations.document_title --messages.content.source.content.citations.end_char_index --messages.content.source.content.citations.start_char_index --messages.content.source.content.citations.type --messages.content.source.content.citations.end_page_number --messages.content.source.content.citations.start_page_number --messages.content.source.content.citations.end_block_index --messages.content.source.content.citations.start_block_index --messages.content.source.content.+citation --messages.content.source.content.source.data --messages.content.source.content.source.media_type --messages.content.source.content.source.type --messages.content.source.+content --messages.content.citations.enabled --messages.content.context --messages.content.title --messages.content.signature --messages.content.thinking --messages.content.data --messages.+content --messages.role --+message --model --metadata.user_id --stop-sequences --+stop_sequence --system.text --system.type --system.cache_control.type --system.citations.cited_text --system.citations.document_index --system.citations.document_title --system.citations.end_char_index --system.citations.start_char_index --system.citations.type --system.citations.end_page_number --system.citations.start_page_number --system.citations.end_block_index --system.citations.start_block_index --system.+citation --+system --temperature --thinking.budget_tokens --thinking.type --tool-choice.type --tool-choice.disable_parallel_tool_use --tool-choice.name --tools.display_height_px --tools.display_width_px --tools.name --tools.type --tools.cache_control.type --tools.display_number --tools.input_schema.type --tools.description --+tool --top-k --top-p --betas --+beta"
        ;;
      beta.messages.count_tokens)
        flags="--messages.content.text --messages.content.type --messages.content.cache_control.type --messages.content.citations.cited_text --messages.content.citations.document_index --messages.content.citations.document_title --messages.content.citations.end_char_index --messages.content.citations.start_char_index --messages.content.citations.type --messages.content.citations.end_page_number --messages.content.citations.start_page_number --messages.content.citations.end_block_index --messages.content.citations.start_block_index --messages.content.+citation --messages.content.source.data --messages.content.source.media_type --messages.content.source.type --messages.content.id --messages.content.name --messages.content.tool_use_id --messages.content.content.text --messages.content.content.type --messages.content.content.cache_control.type --messages.content.content.citations.cited_text --messages.content.content.citations.document_index --messages.content.content.citations.document_title --messages.content.content.citations.end_char_index --messages.content.content.citations.start_char_index --messages.content.content.citations.type --messages.content.content.citations.end_page_number --messages.content.content.citations.start_page_number --messages.content.content.citations.end_block_index --messages.content.content.citations.start_block_index --messages.content.content.+citation --messages.content.content.source.data --messages.content.content.source.media_type --messages.content.content.source.type --messages.content.+content --messages.content.is_error --messages.content.source.content --messages.content.source.content.text --messages.content.source.content.type --messages.content.source.content.cache_control.type --messages.content.source.content.citations.cited_text --messages.content.source.content.citations.document_index --messages.content.source.content.citations.document_title --messages.content.source.content.citations.end_char_index --messages.content.source.content.citations.start_char_index --messages.content.source.content.citations.type --messages.content.source.content.citations.end_page_number --messages.content.source.content.citations.start_page_number --messages.content.source.content.citations.end_block_index --messages.content.source.content.citations.start_block_index --messages.content.source.content.+citation --messages.content.source.content.source.data --messages.content.source.content.source.media_type --messages.content.source.content.source.type --messages.content.source.+content --messages.content.citations.enabled --messages.content.context --messages.content.title --messages.content.signature --messages.content.thinking --messages.content.data --messages.+content --messages.role --+message --model --system --system.text --system.type --system.cache_control.type --system.citations.cited_text --system.citations.document_index --system.citations.document_title --system.citations.end_char_index --system.citations.start_char_index --system.citations.type --system.citations.end_page_number --system.citations.start_page_number --system.citations.end_block_index --system.citations.start_block_index --system.+citation --+system --thinking.budget_tokens --thinking.type --tool-choice.type --tool-choice.disable_parallel_tool_use --tool-choice.name --tools.display_height_px --tools.display_width_px --tools.name --tools.type --tools.cache_control.type --tools.display_number --tools.input_schema.type --tools.description --+tool --betas --+beta"
        ;;
      beta.messages.batches.create)
        flags="--requests.custom_id --requests.params.max_tokens --requests.params.messages.content.text --requests.params.messages.content.type --requests.params.messages.content.cache_control.type --requests.params.messages.content.citations.cited_text --requests.params.messages.content.citations.document_index --requests.params.messages.content.citations.document_title --requests.params.messages.content.citations.end_char_index --requests.params.messages.content.citations.start_char_index --requests.params.messages.content.citations.type --requests.params.messages.content.citations.end_page_number --requests.params.messages.content.citations.start_page_number --requests.params.messages.content.citations.end_block_index --requests.params.messages.content.citations.start_block_index --requests.params.messages.content.+citation --requests.params.messages.content.source.data --requests.params.messages.content.source.media_type --requests.params.messages.content.source.type --requests.params.messages.content.id --requests.params.messages.content.name --requests.params.messages.content.tool_use_id --requests.params.messages.content.content.text --requests.params.messages.content.content.type --requests.params.messages.content.content.cache_control.type --requests.params.messages.content.content.citations.cited_text --requests.params.messages.content.content.citations.document_index --requests.params.messages.content.content.citations.document_title --requests.params.messages.content.content.citations.end_char_index --requests.params.messages.content.content.citations.start_char_index --requests.params.messages.content.content.citations.type --requests.params.messages.content.content.citations.end_page_number --requests.params.messages.content.content.citations.start_page_number --requests.params.messages.content.content.citations.end_block_index --requests.params.messages.content.content.citations.start_block_index --requests.params.messages.content.content.+citation --requests.params.messages.content.content.source.data --requests.params.messages.content.content.source.media_type --requests.params.messages.content.content.source.type --requests.params.messages.content.+content --requests.params.messages.content.is_error --requests.params.messages.content.source.content --requests.params.messages.content.source.content.text --requests.params.messages.content.source.content.type --requests.params.messages.content.source.content.cache_control.type --requests.params.messages.content.source.content.citations.cited_text --requests.params.messages.content.source.content.citations.document_index --requests.params.messages.content.source.content.citations.document_title --requests.params.messages.content.source.content.citations.end_char_index --requests.params.messages.content.source.content.citations.start_char_index --requests.params.messages.content.source.content.citations.type --requests.params.messages.content.source.content.citations.end_page_number --requests.params.messages.content.source.content.citations.start_page_number --requests.params.messages.content.source.content.citations.end_block_index --requests.params.messages.content.source.content.citations.start_block_index --requests.params.messages.content.source.content.+citation --requests.params.messages.content.source.content.source.data --requests.params.messages.content.source.content.source.media_type --requests.params.messages.content.source.content.source.type --requests.params.messages.content.source.+content --requests.params.messages.content.citations.enabled --requests.params.messages.content.context --requests.params.messages.content.title --requests.params.messages.content.signature --requests.params.messages.content.thinking --requests.params.messages.content.data --requests.params.messages.+content --requests.params.messages.role --requests.params.+message --requests.params.model --requests.params.metadata.user_id --requests.params.stop_sequences --requests.params.+stop_sequence --requests.params.stream --requests.params.system.text --requests.params.system.type --requests.params.system.cache_control.type --requests.params.system.citations.cited_text --requests.params.system.citations.document_index --requests.params.system.citations.document_title --requests.params.system.citations.end_char_index --requests.params.system.citations.start_char_index --requests.params.system.citations.type --requests.params.system.citations.end_page_number --requests.params.system.citations.start_page_number --requests.params.system.citations.end_block_index --requests.params.system.citations.start_block_index --requests.params.system.+citation --requests.params.+system --requests.params.temperature --requests.params.thinking.budget_tokens --requests.params.thinking.type --requests.params.tool_choice.type --requests.params.tool_choice.disable_parallel_tool_use --requests.params.tool_choice.name --requests.params.tools.display_height_px --requests.params.tools.display_width_px --requests.params.tools.name --requests.params.tools.type --requests.params.tools.cache_control.type --requests.params.tools.display_number --requests.params.tools.input_schema.type --requests.params.tools.description --requests.params.+tool --requests.params.top_k --requests.params.top_p --+request --betas --+beta"
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
            COMPREPLY=( $(compgen -W "claude-3-7-sonnet-latest claude-3-7-sonnet-20250219 claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0" -- $cur) )
            ;;
        esac
        ;;
      messages.create)
        case "$prev" in
          --messages.content.type)
            COMPREPLY=( $(compgen -W "text image tool_use tool_result document thinking redacted_thinking" -- $cur) )
            ;;
          --messages.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --messages.content.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --messages.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp application/pdf text/plain" -- $cur) )
            ;;
          --messages.content.source.type)
            COMPREPLY=( $(compgen -W "base64 text content" -- $cur) )
            ;;
          --messages.content.content.type)
            COMPREPLY=( $(compgen -W "text image" -- $cur) )
            ;;
          --messages.content.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --messages.content.content.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --messages.content.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp" -- $cur) )
            ;;
          --messages.content.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --messages.content.source.content.type)
            COMPREPLY=( $(compgen -W "text image" -- $cur) )
            ;;
          --messages.content.source.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --messages.content.source.content.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --messages.content.source.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp" -- $cur) )
            ;;
          --messages.content.source.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --messages.role)
            COMPREPLY=( $(compgen -W "user assistant" -- $cur) )
            ;;
          --model)
            COMPREPLY=( $(compgen -W "claude-3-7-sonnet-latest claude-3-7-sonnet-20250219 claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0" -- $cur) )
            ;;
          --system.type)
            COMPREPLY=( $(compgen -W "text" -- $cur) )
            ;;
          --system.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --system.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --thinking.type)
            COMPREPLY=( $(compgen -W "enabled disabled" -- $cur) )
            ;;
          --tool-choice.type)
            COMPREPLY=( $(compgen -W "auto any tool" -- $cur) )
            ;;
          --tools.name)
            COMPREPLY=( $(compgen -W "bash str_replace_editor" -- $cur) )
            ;;
          --tools.type)
            COMPREPLY=( $(compgen -W "bash_20250124 text_editor_20250124" -- $cur) )
            ;;
          --tools.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
        esac
        ;;
      messages.count_tokens)
        case "$prev" in
          --messages.content.type)
            COMPREPLY=( $(compgen -W "text image tool_use tool_result document thinking redacted_thinking" -- $cur) )
            ;;
          --messages.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --messages.content.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --messages.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp application/pdf text/plain" -- $cur) )
            ;;
          --messages.content.source.type)
            COMPREPLY=( $(compgen -W "base64 text content" -- $cur) )
            ;;
          --messages.content.content.type)
            COMPREPLY=( $(compgen -W "text image" -- $cur) )
            ;;
          --messages.content.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --messages.content.content.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --messages.content.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp" -- $cur) )
            ;;
          --messages.content.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --messages.content.source.content.type)
            COMPREPLY=( $(compgen -W "text image" -- $cur) )
            ;;
          --messages.content.source.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --messages.content.source.content.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --messages.content.source.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp" -- $cur) )
            ;;
          --messages.content.source.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --messages.role)
            COMPREPLY=( $(compgen -W "user assistant" -- $cur) )
            ;;
          --model)
            COMPREPLY=( $(compgen -W "claude-3-7-sonnet-latest claude-3-7-sonnet-20250219 claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0" -- $cur) )
            ;;
          --system.type)
            COMPREPLY=( $(compgen -W "text" -- $cur) )
            ;;
          --system.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --system.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --thinking.type)
            COMPREPLY=( $(compgen -W "enabled disabled" -- $cur) )
            ;;
          --tool-choice.type)
            COMPREPLY=( $(compgen -W "auto any tool" -- $cur) )
            ;;
          --tools.name)
            COMPREPLY=( $(compgen -W "bash str_replace_editor" -- $cur) )
            ;;
          --tools.type)
            COMPREPLY=( $(compgen -W "bash_20250124 text_editor_20250124" -- $cur) )
            ;;
          --tools.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
        esac
        ;;
      messages.batches.create)
        case "$prev" in
          --requests.params.messages.content.type)
            COMPREPLY=( $(compgen -W "text image tool_use tool_result document thinking redacted_thinking" -- $cur) )
            ;;
          --requests.params.messages.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --requests.params.messages.content.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --requests.params.messages.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp application/pdf text/plain" -- $cur) )
            ;;
          --requests.params.messages.content.source.type)
            COMPREPLY=( $(compgen -W "base64 text content" -- $cur) )
            ;;
          --requests.params.messages.content.content.type)
            COMPREPLY=( $(compgen -W "text image" -- $cur) )
            ;;
          --requests.params.messages.content.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --requests.params.messages.content.content.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --requests.params.messages.content.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp" -- $cur) )
            ;;
          --requests.params.messages.content.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --requests.params.messages.content.source.content.type)
            COMPREPLY=( $(compgen -W "text image" -- $cur) )
            ;;
          --requests.params.messages.content.source.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --requests.params.messages.content.source.content.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --requests.params.messages.content.source.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp" -- $cur) )
            ;;
          --requests.params.messages.content.source.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --requests.params.messages.role)
            COMPREPLY=( $(compgen -W "user assistant" -- $cur) )
            ;;
          --requests.params.model)
            COMPREPLY=( $(compgen -W "claude-3-7-sonnet-latest claude-3-7-sonnet-20250219 claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0" -- $cur) )
            ;;
          --requests.params.system.type)
            COMPREPLY=( $(compgen -W "text" -- $cur) )
            ;;
          --requests.params.system.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --requests.params.system.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --requests.params.thinking.type)
            COMPREPLY=( $(compgen -W "enabled disabled" -- $cur) )
            ;;
          --requests.params.tool_choice.type)
            COMPREPLY=( $(compgen -W "auto any tool" -- $cur) )
            ;;
          --requests.params.tools.name)
            COMPREPLY=( $(compgen -W "bash str_replace_editor" -- $cur) )
            ;;
          --requests.params.tools.type)
            COMPREPLY=( $(compgen -W "bash_20250124 text_editor_20250124" -- $cur) )
            ;;
          --requests.params.tools.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
        esac
        ;;
      beta.messages.create)
        case "$prev" in
          --messages.content.type)
            COMPREPLY=( $(compgen -W "text image tool_use tool_result document thinking redacted_thinking" -- $cur) )
            ;;
          --messages.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --messages.content.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --messages.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp application/pdf text/plain" -- $cur) )
            ;;
          --messages.content.source.type)
            COMPREPLY=( $(compgen -W "base64 text content" -- $cur) )
            ;;
          --messages.content.content.type)
            COMPREPLY=( $(compgen -W "text image" -- $cur) )
            ;;
          --messages.content.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --messages.content.content.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --messages.content.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp" -- $cur) )
            ;;
          --messages.content.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --messages.content.source.content.type)
            COMPREPLY=( $(compgen -W "text image" -- $cur) )
            ;;
          --messages.content.source.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --messages.content.source.content.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --messages.content.source.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp" -- $cur) )
            ;;
          --messages.content.source.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --messages.role)
            COMPREPLY=( $(compgen -W "user assistant" -- $cur) )
            ;;
          --model)
            COMPREPLY=( $(compgen -W "claude-3-7-sonnet-latest claude-3-7-sonnet-20250219 claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0" -- $cur) )
            ;;
          --system.type)
            COMPREPLY=( $(compgen -W "text" -- $cur) )
            ;;
          --system.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --system.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --thinking.type)
            COMPREPLY=( $(compgen -W "enabled disabled" -- $cur) )
            ;;
          --tool-choice.type)
            COMPREPLY=( $(compgen -W "auto any tool" -- $cur) )
            ;;
          --tools.name)
            COMPREPLY=( $(compgen -W "computer bash str_replace_editor" -- $cur) )
            ;;
          --tools.type)
            COMPREPLY=( $(compgen -W "computer_20241022 bash_20241022 text_editor_20241022 computer_20250124 bash_20250124 text_editor_20250124 custom" -- $cur) )
            ;;
          --tools.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --tools.input_schema.type)
            COMPREPLY=( $(compgen -W "object" -- $cur) )
            ;;
          --betas)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 computer-use-2025-01-24 pdfs-2024-09-25 token-counting-2024-11-01 token-efficient-tools-2025-02-19 output-128k-2025-02-19" -- $cur) )
            ;;
          --+beta)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 computer-use-2025-01-24 pdfs-2024-09-25 token-counting-2024-11-01 token-efficient-tools-2025-02-19 output-128k-2025-02-19" -- $cur) )
            ;;
        esac
        ;;
      beta.messages.count_tokens)
        case "$prev" in
          --messages.content.type)
            COMPREPLY=( $(compgen -W "text image tool_use tool_result document thinking redacted_thinking" -- $cur) )
            ;;
          --messages.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --messages.content.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --messages.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp application/pdf text/plain" -- $cur) )
            ;;
          --messages.content.source.type)
            COMPREPLY=( $(compgen -W "base64 text content" -- $cur) )
            ;;
          --messages.content.content.type)
            COMPREPLY=( $(compgen -W "text image" -- $cur) )
            ;;
          --messages.content.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --messages.content.content.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --messages.content.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp" -- $cur) )
            ;;
          --messages.content.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --messages.content.source.content.type)
            COMPREPLY=( $(compgen -W "text image" -- $cur) )
            ;;
          --messages.content.source.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --messages.content.source.content.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --messages.content.source.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp" -- $cur) )
            ;;
          --messages.content.source.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --messages.role)
            COMPREPLY=( $(compgen -W "user assistant" -- $cur) )
            ;;
          --model)
            COMPREPLY=( $(compgen -W "claude-3-7-sonnet-latest claude-3-7-sonnet-20250219 claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0" -- $cur) )
            ;;
          --system.type)
            COMPREPLY=( $(compgen -W "text" -- $cur) )
            ;;
          --system.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --system.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --thinking.type)
            COMPREPLY=( $(compgen -W "enabled disabled" -- $cur) )
            ;;
          --tool-choice.type)
            COMPREPLY=( $(compgen -W "auto any tool" -- $cur) )
            ;;
          --tools.name)
            COMPREPLY=( $(compgen -W "computer bash str_replace_editor" -- $cur) )
            ;;
          --tools.type)
            COMPREPLY=( $(compgen -W "computer_20241022 bash_20241022 text_editor_20241022 computer_20250124 bash_20250124 text_editor_20250124 custom" -- $cur) )
            ;;
          --tools.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --tools.input_schema.type)
            COMPREPLY=( $(compgen -W "object" -- $cur) )
            ;;
          --betas)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 computer-use-2025-01-24 pdfs-2024-09-25 token-counting-2024-11-01 token-efficient-tools-2025-02-19 output-128k-2025-02-19" -- $cur) )
            ;;
          --+beta)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 computer-use-2025-01-24 pdfs-2024-09-25 token-counting-2024-11-01 token-efficient-tools-2025-02-19 output-128k-2025-02-19" -- $cur) )
            ;;
        esac
        ;;
      beta.messages.batches.create)
        case "$prev" in
          --requests.params.messages.content.type)
            COMPREPLY=( $(compgen -W "text image tool_use tool_result document thinking redacted_thinking" -- $cur) )
            ;;
          --requests.params.messages.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --requests.params.messages.content.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --requests.params.messages.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp application/pdf text/plain" -- $cur) )
            ;;
          --requests.params.messages.content.source.type)
            COMPREPLY=( $(compgen -W "base64 text content" -- $cur) )
            ;;
          --requests.params.messages.content.content.type)
            COMPREPLY=( $(compgen -W "text image" -- $cur) )
            ;;
          --requests.params.messages.content.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --requests.params.messages.content.content.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --requests.params.messages.content.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp" -- $cur) )
            ;;
          --requests.params.messages.content.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --requests.params.messages.content.source.content.type)
            COMPREPLY=( $(compgen -W "text image" -- $cur) )
            ;;
          --requests.params.messages.content.source.content.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --requests.params.messages.content.source.content.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --requests.params.messages.content.source.content.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp" -- $cur) )
            ;;
          --requests.params.messages.content.source.content.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --requests.params.messages.role)
            COMPREPLY=( $(compgen -W "user assistant" -- $cur) )
            ;;
          --requests.params.model)
            COMPREPLY=( $(compgen -W "claude-3-7-sonnet-latest claude-3-7-sonnet-20250219 claude-3-5-haiku-latest claude-3-5-haiku-20241022 claude-3-5-sonnet-latest claude-3-5-sonnet-20241022 claude-3-5-sonnet-20240620 claude-3-opus-latest claude-3-opus-20240229 claude-3-sonnet-20240229 claude-3-haiku-20240307 claude-2.1 claude-2.0" -- $cur) )
            ;;
          --requests.params.system.type)
            COMPREPLY=( $(compgen -W "text" -- $cur) )
            ;;
          --requests.params.system.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --requests.params.system.citations.type)
            COMPREPLY=( $(compgen -W "char_location page_location content_block_location" -- $cur) )
            ;;
          --requests.params.thinking.type)
            COMPREPLY=( $(compgen -W "enabled disabled" -- $cur) )
            ;;
          --requests.params.tool_choice.type)
            COMPREPLY=( $(compgen -W "auto any tool" -- $cur) )
            ;;
          --requests.params.tools.name)
            COMPREPLY=( $(compgen -W "computer bash str_replace_editor" -- $cur) )
            ;;
          --requests.params.tools.type)
            COMPREPLY=( $(compgen -W "computer_20241022 bash_20241022 text_editor_20241022 computer_20250124 bash_20250124 text_editor_20250124 custom" -- $cur) )
            ;;
          --requests.params.tools.cache_control.type)
            COMPREPLY=( $(compgen -W "ephemeral" -- $cur) )
            ;;
          --requests.params.tools.input_schema.type)
            COMPREPLY=( $(compgen -W "object" -- $cur) )
            ;;
          --betas)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 computer-use-2025-01-24 pdfs-2024-09-25 token-counting-2024-11-01 token-efficient-tools-2025-02-19 output-128k-2025-02-19" -- $cur) )
            ;;
          --+beta)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 computer-use-2025-01-24 pdfs-2024-09-25 token-counting-2024-11-01 token-efficient-tools-2025-02-19 output-128k-2025-02-19" -- $cur) )
            ;;
        esac
        ;;
      beta.messages.batches.retrieve)
        case "$prev" in
          --betas)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 computer-use-2025-01-24 pdfs-2024-09-25 token-counting-2024-11-01 token-efficient-tools-2025-02-19 output-128k-2025-02-19" -- $cur) )
            ;;
          --+beta)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 computer-use-2025-01-24 pdfs-2024-09-25 token-counting-2024-11-01 token-efficient-tools-2025-02-19 output-128k-2025-02-19" -- $cur) )
            ;;
        esac
        ;;
      beta.messages.batches.list)
        case "$prev" in
          --betas)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 computer-use-2025-01-24 pdfs-2024-09-25 token-counting-2024-11-01 token-efficient-tools-2025-02-19 output-128k-2025-02-19" -- $cur) )
            ;;
          --+beta)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 computer-use-2025-01-24 pdfs-2024-09-25 token-counting-2024-11-01 token-efficient-tools-2025-02-19 output-128k-2025-02-19" -- $cur) )
            ;;
        esac
        ;;
      beta.messages.batches.delete)
        case "$prev" in
          --betas)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 computer-use-2025-01-24 pdfs-2024-09-25 token-counting-2024-11-01 token-efficient-tools-2025-02-19 output-128k-2025-02-19" -- $cur) )
            ;;
          --+beta)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 computer-use-2025-01-24 pdfs-2024-09-25 token-counting-2024-11-01 token-efficient-tools-2025-02-19 output-128k-2025-02-19" -- $cur) )
            ;;
        esac
        ;;
      beta.messages.batches.cancel)
        case "$prev" in
          --betas)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 computer-use-2025-01-24 pdfs-2024-09-25 token-counting-2024-11-01 token-efficient-tools-2025-02-19 output-128k-2025-02-19" -- $cur) )
            ;;
          --+beta)
            COMPREPLY=( $(compgen -W "message-batches-2024-09-24 prompt-caching-2024-07-31 computer-use-2024-10-22 computer-use-2025-01-24 pdfs-2024-09-25 token-counting-2024-11-01 token-efficient-tools-2025-02-19 output-128k-2025-02-19" -- $cur) )
            ;;
        esac
        ;;
    esac
}
complete -F _main anthropic-cli