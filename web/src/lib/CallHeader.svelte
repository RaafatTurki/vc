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
    gap: 24px;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 24px;
  }

  p {
    margin: 0 0 12px;
    color: #8f9cff;
    font-size: 0.76rem;
    font-weight: 800;
    letter-spacing: 0.14em;
    text-transform: uppercase;
  }

  h1 {
    max-width: 650px;
    margin: 0;
    overflow: hidden;
    font-size: clamp(2rem, 5vw, 4.3rem);
    line-height: 1;
    letter-spacing: -0.065em;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  @media (max-width: 760px) {
    header { align-items: stretch; flex-direction: column; }
  }
</style>
