<script lang="ts">
  let { state, text }: { state: "idle" | "connecting" | "connected" | "error"; text: string } = $props()
</script>

<header>
  <a class="brand" href="/" aria-label="Vivid home">
    <span class="brand-mark" aria-hidden="true"></span>
    Vivid
  </a>
  <div class="status" data-state={state} role="status" aria-live="polite">
    <span class="status-dot" aria-hidden="true"></span>
    <span>{text}</span>
  </div>
</header>

<style>
  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .brand {
    display: inline-flex;
    gap: 10px;
    align-items: center;
    color: inherit;
    font-weight: 750;
    letter-spacing: -0.02em;
    text-decoration: none;
  }

  .brand-mark {
    width: 27px;
    height: 20px;
    border-radius: 7px;
    background: linear-gradient(135deg, #7180ff, #a96cff);
    box-shadow: 0 0 24px rgb(113 128 255 / 45%);
  }

  .brand-mark::after {
    display: block;
    width: 0;
    height: 0;
    margin: 5px 0 0 9px;
    border-top: 5px solid transparent;
    border-bottom: 5px solid transparent;
    border-left: 8px solid white;
    content: "";
  }

  .status {
    display: flex;
    gap: 8px;
    align-items: center;
    padding: 8px 12px;
    border: 1px solid #242a38;
    border-radius: 999px;
    color: #a9b0c0;
    background: rgb(13 17 27 / 76%);
    font-size: 0.8rem;
  }

  .status-dot {
    width: 7px;
    height: 7px;
    border-radius: 50%;
    background: #778095;
  }

  .status[data-state="connecting"] .status-dot { background: #f5b942; box-shadow: 0 0 10px #f5b942; }
  .status[data-state="connected"] .status-dot { background: #21d6a1; box-shadow: 0 0 10px #21d6a1; }
  .status[data-state="error"] .status-dot { background: #ff647c; box-shadow: 0 0 10px #ff647c; }
</style>
