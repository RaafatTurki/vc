<script lang="ts">
  import type { Snippet } from "svelte"
  import type { HTMLButtonAttributes } from "svelte/elements"

  interface Props extends HTMLButtonAttributes {
    children: Snippet
    kind?: "primary" | "secondary" | "danger" | "ghost"
    size?: "regular" | "compact" | "icon"
    pressed?: boolean
    fullWidth?: boolean
  }
  let { children, kind = "secondary", size = "regular", pressed, fullWidth = false, ...attributes }: Props = $props()
</script>

<button class="ui-button" data-kind={kind} data-size={size} data-pressed={pressed} class:full-width={fullWidth} {...attributes}>
  {@render children()}
</button>

<style>
  .ui-button { display: inline-flex; gap: var(--space-2); align-items: center; justify-content: center; min-height: var(--control-height); padding: 0 1.0625rem; border: 1px solid var(--accent-border); border-radius: var(--radius); color: var(--ink); background: var(--accent-soft); font-weight: 720; white-space: nowrap; }
  .ui-button :global(svg) { width: 1.125rem; height: 1.125rem; }
  .ui-button[data-size="compact"] { min-height: var(--control-height-compact); }
  .ui-button[data-size="icon"] { width: var(--control-height); min-width: var(--control-height); padding: 0; }
  .ui-button.full-width { width: 100%; }
  .ui-button[data-kind="primary"] { border-color: var(--accent); color: var(--bg); background: var(--accent); }
  .ui-button[data-kind="danger"], .ui-button[data-pressed="true"] { border-color: var(--danger-border); color: var(--danger); background: var(--danger-bg); }
  .ui-button[data-kind="ghost"] { border-color: transparent; background: transparent; }
  @media (hover: hover) and (pointer: fine) { .ui-button:hover:not(:disabled) { border-color: var(--accent-border-strong); color: var(--accent); background: var(--accent-subtle); } .ui-button[data-kind="primary"]:hover:not(:disabled) { border-color: var(--accent); color: var(--bg); background: var(--accent); } }
</style>
