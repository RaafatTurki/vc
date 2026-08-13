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
    gap: var(--space-2);
    align-items: center;
    justify-content: center;
    height: var(--control-height);
    min-height: var(--control-height);
    padding: 0 1.0625rem;
    border: 1px solid var(--accent-border);
    border-radius: 2px;
    color: var(--ink);
    background: var(--accent-soft);
    font-weight: 720;
    white-space: nowrap;
  }

  @media (hover: hover) and (pointer: fine) {
    button:hover:not(:disabled) { border-color: var(--accent-border-strong); color: var(--accent); background: var(--accent-subtle); }
    button.primary:hover:not(:disabled) { border-color: var(--accent); color: var(--bg); background: var(--accent); }
  }
  button.compact { min-height: 2.5rem; height: 2.5rem; }

  button.primary {
    width: 100%;
    margin-top: 1.375rem;
    border: 1px solid var(--accent);
    color: var(--bg);
    background: var(--accent);
  }
</style>
