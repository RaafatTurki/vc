<script lang="ts">
  import { Copy } from "@lucide/svelte"
  import Button from "./Button.svelte"

  let { roomID, copyLabel, onCopy }: { roomID: string; copyLabel: string; onCopy: () => void | Promise<void> } = $props()
</script>

<header>
  <div>
    <p>Current room</p>
    <h1>{roomID}</h1>
  </div>
  <Button compact onclick={onCopy}>
    <Copy size={17} aria-hidden="true" />
    {copyLabel}
  </Button>
</header>

<style>
  header {
    display: flex;
    gap: var(--space-6);
    align-items: center;
    justify-content: space-between;
    margin-bottom: var(--space-6);
    padding-bottom: var(--space-4);
    border-bottom: 1px solid var(--line);
  }

  p {
    margin: 0 0 var(--space-3);
    color: var(--accent);
    font-size: 0.76rem;
    font-weight: 800;
    letter-spacing: 0.14em;
    text-transform: uppercase;
  }

  h1 {
    max-width: 40.625rem;
    margin: 0;
    overflow: hidden;
    font-size: clamp(2rem, 5vw, 4.3rem);
    line-height: 1;
    letter-spacing: -0.045em;
    font-weight: 500;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  @media (max-width: 47.5em) {
    header { align-items: stretch; flex-direction: column; }
  }
</style>
