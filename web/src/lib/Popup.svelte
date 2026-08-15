<script lang="ts">
  import { X } from "@lucide/svelte"
  import Button from "./Button.svelte"
  interface Option { value: string; label: string }
  interface Props { open: boolean; title: string; message?: string; options: Option[]; onSelect: (value: string) => void | Promise<void>; onClose: () => void }
  let { open, title, message = "", options, onSelect, onClose }: Props = $props()
  let dialog = $state<HTMLDialogElement>()
  $effect(() => {
    if (!dialog) return
    if (open && !dialog.open) dialog.showModal()
    if (!open && dialog.open) dialog.close()
  })
  function cancel(event: Event) { event.preventDefault(); onClose() }
  function clickBackdrop(event: MouseEvent) { if (event.target === dialog) onClose() }
</script>

<dialog bind:this={dialog} aria-labelledby="popup-title" oncancel={cancel} onclose={onClose} onclick={clickBackdrop}>
  <header><h2 id="popup-title">{title}</h2><Button size="icon" kind="ghost" aria-label="Close" onclick={onClose}><X aria-hidden="true" /></Button></header>
  {#if message}<p class="message">{message}</p>{/if}
  <div class="options">{#each options as option}<Button fullWidth onclick={() => onSelect(option.value)}>{option.label}</Button>{/each}</div>
</dialog>

<style>
  dialog { width: min(23rem, calc(100vw - 2rem)); max-height: calc(100vh - 2rem); padding: var(--space-4); border: 1px solid var(--accent-border); border-radius: .25rem; color: var(--ink); background: var(--surface); }
  dialog::backdrop { background: var(--surface-soft); }
  header { display: flex; align-items: center; justify-content: space-between; gap: var(--space-3); }
  h2 { font-size: 1rem; } .message { margin-top: var(--space-4); color: var(--muted); font-size: .84rem; }
  .options { display: grid; gap: var(--space-2); margin-top: var(--space-4); }
</style>
