<script lang="ts">
  import type { Snippet } from "svelte"
  import type { MouseEventHandler } from "svelte/elements"

  interface Props {
    children: Snippet
    kind?: "primary" | "secondary"
    compact?: boolean
    type?: "button" | "reset" | "submit"
    disabled?: boolean
    onclick?: MouseEventHandler<HTMLButtonElement>
  }

  let {
    children,
    kind = "secondary",
    compact = false,
    type = "button",
    disabled = false,
    onclick,
  }: Props = $props()
</script>

<button class:primary={kind === "primary"} class:compact {type} {disabled} {onclick}>
  {@render children()}
</button>

<style>
  button {
    display: inline-flex;
    gap: 8px;
    align-items: center;
    justify-content: center;
    min-height: 46px;
    padding: 0 17px;
    border: 1px solid #303748;
    border-radius: 12px;
    color: #f8f9fd;
    background: #171c28;
    font-weight: 720;
    white-space: nowrap;
  }

  button:hover:not(:disabled) { filter: brightness(1.12); }
  button.compact { min-height: 40px; }

  button.primary {
    width: 100%;
    margin-top: 22px;
    border: 0;
    background: linear-gradient(135deg, #6979ff, #9467ed);
    box-shadow: 0 10px 30px rgb(105 121 255 / 24%);
  }
</style>
