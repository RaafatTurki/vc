<script lang="ts">
  import { ChevronDown } from "@lucide/svelte"

  export interface DropdownOption {
    value: string
    label: string
  }

  interface Props {
    options: DropdownOption[]
    value?: string
    label?: string
    disabled?: boolean
    onchange?: (value: string) => void | Promise<void>
  }

  let {
    options,
    value = $bindable(""),
    label = "Choose an option",
    disabled = false,
    onchange,
  }: Props = $props()

  let open = $state(false)
  let focusedIndex = $state(0)
  let menu: HTMLDivElement

  const selectedLabel = $derived(options.find(option => option.value === value)?.label || label)

  function selectOption(option: DropdownOption): void {
    value = option.value
    open = false
    onchange?.(option.value)
  }

  function toggle(): void {
    if (disabled || options.length < 2) return
    open = !open
    if (open) focusedIndex = Math.max(0, options.findIndex(option => option.value === value))
  }

  function handleKeydown(event: KeyboardEvent): void {
    if (disabled || options.length < 2) return
    if (event.key === "Enter" || event.key === " ") {
      event.preventDefault()
      if (open) selectOption(options[focusedIndex])
      else toggle()
    } else if (event.key === "ArrowDown" || event.key === "ArrowRight") {
      event.preventDefault()
      open = true
      focusedIndex = (focusedIndex + 1) % options.length
    } else if (event.key === "ArrowUp" || event.key === "ArrowLeft") {
      event.preventDefault()
      open = true
      focusedIndex = (focusedIndex - 1 + options.length) % options.length
    } else if (event.key === "Escape") {
      open = false
    }
  }

  $effect(() => {
    if (!open) return
    const close = (event: MouseEvent) => {
      if (menu && !menu.contains(event.target as Node)) open = false
    }
    document.addEventListener("mousedown", close)
    return () => document.removeEventListener("mousedown", close)
  })
</script>

<div class="dropdown" bind:this={menu}>
  <button
    class="trigger"
    class:open
    type="button"
    aria-haspopup="listbox"
    aria-expanded={open}
    {disabled}
    onclick={toggle}
    onkeydown={handleKeydown}
  >
    <span class:selected-placeholder={!options.some(option => option.value === value)}>{selectedLabel}</span>
    <ChevronDown aria-hidden="true" />
  </button>

  {#if open}
    <div class="menu" role="listbox" aria-label={label} tabindex="-1">
      {#each options as option, index (option.value)}
        <button
          class="option"
          class:selected={option.value === value}
          class:focused={index === focusedIndex}
          type="button"
          role="option"
          aria-selected={option.value === value}
          onclick={() => selectOption(option)}
          onmouseenter={() => focusedIndex = index}
        >
          {option.label}
        </button>
      {/each}
    </div>
  {/if}
</div>

<style>
  .dropdown { position: relative; width: 100%; min-width: 0; }

  .trigger {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-2);
    width: 100%;
    height: var(--control-height);
    min-height: var(--control-height);
    min-width: 0;
    padding: 0 0.6875rem;
    border: 1px solid var(--line-strong);
    border-radius: 2px;
    outline: none;
    color: var(--ink);
    background: var(--bg);
    text-align: left;
  }

  .trigger:focus-visible,
  .trigger.open { border-color: var(--accent); box-shadow: 0 0 0 3px var(--accent-12); }
  .trigger:disabled { cursor: default; opacity: 0.65; }
  .trigger span { min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .trigger .selected-placeholder { color: var(--muted); }
  .trigger :global(svg) { flex: 0 0 auto; width: 1rem; height: 1rem; color: var(--accent); transition: transform 160ms ease; }
  .trigger.open :global(svg) { transform: rotate(180deg); }

  .menu {
    position: absolute;
    z-index: 20;
    top: calc(100% + var(--space-1));
    right: 0;
    left: 0;
    display: grid;
    max-height: min(16rem, 40vh);
    overflow-y: auto;
    padding: var(--space-1);
    border: 1px solid var(--accent-border);
    border-radius: 2px;
    background: var(--surface);
    box-shadow: 0 0.75rem 2rem rgb(0 0 0 / 45%);
  }

  .option {
    width: 100%;
    min-height: 2.5rem;
    padding: 0.5rem 0.625rem;
    border: 0;
    border-radius: 1px;
    color: var(--muted);
    background: transparent;
    text-align: left;
    white-space: normal;
  }

  .option.focused { color: var(--ink); background: var(--accent-soft); }
  .option.selected { color: var(--accent); }

  @media (hover: hover) and (pointer: fine) {
    .trigger:hover:not(:disabled) { border-color: var(--accent); box-shadow: 0 0 0 3px var(--accent-12); }
    .option:hover { color: var(--ink); background: var(--accent-soft); }
  }
</style>
