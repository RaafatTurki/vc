<script lang="ts">
  import { X } from "@lucide/svelte"

  interface Option {
    value: string
    label: string
  }

  interface Props {
    open: boolean
    title: string
    message?: string
    options: Option[]
    onSelect: (value: string) => void | Promise<void>
    onClose: () => void
  }

  let { open, title, message = "", options, onSelect, onClose }: Props = $props()

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === "Escape") onClose()
  }
</script>

<svelte:window onkeydown={event => open && handleKeydown(event)} />

{#if open}
  <div class="backdrop" role="presentation" onclick={onClose}></div>
  <div class="popup" role="dialog" aria-modal="true" aria-labelledby="popup-title">
    <header>
      <h2 id="popup-title">{title}</h2>
      <button class="close" type="button" aria-label="Close" onclick={onClose}>
        <X aria-hidden="true" />
      </button>
    </header>
    {#if message}<p class="message">{message}</p>{/if}
    <div class="options">
      {#each options as option}
        <button type="button" onclick={() => onSelect(option.value)}>{option.label}</button>
      {/each}
    </div>
  </div>
{/if}

<style>
  .backdrop {
    position: fixed;
    inset: 0;
    z-index: 1000;
    background: var(--surface-soft);
    backdrop-filter: blur(0.25rem);
  }

  .popup {
    position: fixed;
    top: 50%;
    left: 50%;
    z-index: 1001;
    display: grid;
    gap: var(--space-4);
    width: min(23rem, calc(100vw - 2rem));
    max-height: calc(100vh - 2rem);
    padding: var(--space-4);
    border: 1px solid var(--accent-border);
    border-radius: 0.25rem;
    background: var(--surface);
    box-shadow: 0 1rem 3rem var(--accent-16);
    transform: translate(-50%, -50%);
    overflow: auto;
  }

  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-3);
  }

  h2 {
    margin: 0;
    color: var(--ink);
    font-size: 1rem;
  }

  .message { margin: 0; color: var(--muted); font-size: 0.84rem; }

  button {
    border: 1px solid var(--line-soft);
    border-radius: 0.2rem;
    color: var(--ink);
    background: var(--surface-faint);
    font: inherit;
    font-weight: 650;
  }

  .close {
    display: grid;
    width: var(--control-height);
    height: var(--control-height);
    min-width: var(--control-height);
    place-items: center;
  }

  .close :global(svg) { width: 1rem; height: 1rem; }

  .options {
    display: grid;
    gap: var(--space-2);
  }

  .options button {
    min-height: var(--control-height);
    padding: 0.65rem 0.8rem;
    text-align: left;
  }

  @media (hover: hover) and (pointer: fine) {
    button:hover { color: var(--accent); border-color: var(--accent-border); background: var(--accent-subtle); }
  }
</style>
