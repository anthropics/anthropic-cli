#!/usr/bin/env bash

# This also works for zsh: https://zsh.sourceforge.io/Doc/Release/Completion-System.html#Completion-System
_main()
{
    COMPREPLY=()

    local subcommands="completions.create messages.create"

    if [[ "$COMP_CWORD" -eq 1 ]]
    then
      local cur="${COMP_WORDS[COMP_CWORD]}"
      COMPREPLY=($(compgen -W "$subcommands" -- "$cur"))
      return
    fi

    local subcommand="${COMP_WORDS[1]}"
    local flags
    case "$subcommand" in
      completions.create)
        flags="--max-tokens-to-sample --model --prompt --metadata.user_id --stop-sequences --+stop_sequence --temperature --top-k --top-p"
        ;;
      messages.create)
        flags="--max-tokens --messages.content.text_block_param.text --messages.content.text_block_param.type --messages.content.image_block_param.source.data --messages.content.image_block_param.source.media_type --messages.content.image_block_param.source.type --messages.content.image_block_param.type --messages.content.tool_use_block_param.id --messages.content.tool_use_block_param.name --messages.content.tool_use_block_param.type --messages.content.tool_result_block_param.tool_use_id --messages.content.tool_result_block_param.type --messages.content.tool_result_block_param.content.text_block_param.text --messages.content.tool_result_block_param.content.text_block_param.type --messages.content.tool_result_block_param.content.image_block_param.source.data --messages.content.tool_result_block_param.content.image_block_param.source.media_type --messages.content.tool_result_block_param.content.image_block_param.source.type --messages.content.tool_result_block_param.content.image_block_param.type --messages.content.tool_result_block_param.content.+text_block_param --messages.content.tool_result_block_param.content.+image_block_param --messages.content.tool_result_block_param.is_error --messages.content.+text_block_param --messages.content.+image_block_param --messages.content.+tool_use_block_param --messages.content.+tool_result_block_param --messages.role --+message --model --metadata.user_id --stop-sequences --+stop_sequence --system.text --system.type --+system --temperature --tool-choice.tool_choice_auto.type --tool-choice.tool_choice_auto.disable_parallel_tool_use --tool-choice.tool_choice_any.type --tool-choice.tool_choice_any.disable_parallel_tool_use --tool-choice.tool_choice_tool.name --tool-choice.tool_choice_tool.type --tool-choice.tool_choice_tool.disable_parallel_tool_use --tools.name --tools.description --+tool --top-k --top-p"
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
          --messages.content.text_block_param.type)
            COMPREPLY=( $(compgen -W "text" -- $cur) )
            ;;
          --messages.content.image_block_param.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp" -- $cur) )
            ;;
          --messages.content.image_block_param.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --messages.content.image_block_param.type)
            COMPREPLY=( $(compgen -W "image" -- $cur) )
            ;;
          --messages.content.tool_use_block_param.type)
            COMPREPLY=( $(compgen -W "tool_use" -- $cur) )
            ;;
          --messages.content.tool_result_block_param.type)
            COMPREPLY=( $(compgen -W "tool_result" -- $cur) )
            ;;
          --messages.content.tool_result_block_param.content.text_block_param.type)
            COMPREPLY=( $(compgen -W "text" -- $cur) )
            ;;
          --messages.content.tool_result_block_param.content.image_block_param.source.media_type)
            COMPREPLY=( $(compgen -W "image/jpeg image/png image/gif image/webp" -- $cur) )
            ;;
          --messages.content.tool_result_block_param.content.image_block_param.source.type)
            COMPREPLY=( $(compgen -W "base64" -- $cur) )
            ;;
          --messages.content.tool_result_block_param.content.image_block_param.type)
            COMPREPLY=( $(compgen -W "image" -- $cur) )
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
          --tool-choice.tool_choice_auto.type)
            COMPREPLY=( $(compgen -W "auto" -- $cur) )
            ;;
          --tool-choice.tool_choice_any.type)
            COMPREPLY=( $(compgen -W "any" -- $cur) )
            ;;
          --tool-choice.tool_choice_tool.type)
            COMPREPLY=( $(compgen -W "tool" -- $cur) )
            ;;
        esac
        ;;
    esac
}
complete -F _main anthropic-cli