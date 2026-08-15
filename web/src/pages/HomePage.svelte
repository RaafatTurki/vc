<script lang="ts">
  import { Download, Mic, MicOff, RefreshCw, Video, VideoOff } from "@lucide/svelte"
  import { createRoomID } from "../lib/call"
  import Button from "../lib/Button.svelte"

  interface Props {
    roomInput?: string
    signalInput?: string
    participantName?: string
    joinWithAudio?: boolean
    joinWithVideo?: boolean
    joining: boolean
    setupError: string
    installHint?: string
    onJoin: (event: SubmitEvent) => void | Promise<void>
    onInstall?: () => void | Promise<void>
  }

  let {
    roomInput = $bindable(),
    signalInput = $bindable(),
    participantName = $bindable(),
    joinWithAudio = $bindable(),
    joinWithVideo = $bindable(),
    joining,
    setupError,
    installHint = "",
    onJoin,
    onInstall = () => {},
  }: Props = $props()
</script>

<section aria-labelledby="setup-title">
  <div>
    <h1 id="setup-title">Simple <span>Secure</span> <strong>Vivid</strong></h1>
    <p class="intro">Create a room, share its link, and talk directly from browser to browser.</p>
    <p class="eyebrow">Private by design</p>
    <div class="install-button">
      <Button onclick={onInstall}>
        <Download size={16} aria-hidden="true" />
        Install web app
      </Button>
    </div>
    {#if installHint}<p class="install-hint" role="status">{installHint}</p>{/if}
  </div>

  <form onsubmit={onJoin} novalidate>
    <div class="form-field">
      <label for="nameInput">Your name</label>
      <input class="ui-field" id="nameInput" bind:value={participantName} maxlength="32" autocomplete="name" placeholder="How others will see you" required>
    </div>

    <div class="form-field">
      <label for="roomInput">Room ID</label>
      <div class="room-entry">
        <input class="ui-field" id="roomInput" bind:value={roomInput} maxlength="6" autocomplete="off" spellcheck="false" required>
        <Button onclick={() => roomInput = createRoomID()}>
          <RefreshCw size={17} aria-hidden="true" />
          New room
        </Button>
      </div>
    </div>

    <fieldset>
      <legend>Join with</legend>
      <button class="media-choice" type="button" aria-pressed={joinWithAudio} onclick={() => joinWithAudio = !joinWithAudio}>
        {#if joinWithAudio}<Mic aria-hidden="true" />{:else}<MicOff aria-hidden="true" />{/if}
        <span>Microphone {joinWithAudio ? "on" : "off"}</span>
      </button>
      <button class="media-choice" type="button" aria-pressed={joinWithVideo} onclick={() => joinWithVideo = !joinWithVideo}>
        {#if joinWithVideo}<Video aria-hidden="true" />{:else}<VideoOff aria-hidden="true" />{/if}
        <span>Camera {joinWithVideo ? "on" : "off"}</span>
      </button>
    </fieldset>

    <details>
      <summary>Connection settings</summary>
      <label for="signalInput">Signaling WebSocket URL</label>
      <input class="ui-field" id="signalInput" bind:value={signalInput} type="url" spellcheck="false" required>
    </details>

    <Button kind="primary" fullWidth type="submit" disabled={joining}>
      {joining ? "Joining…" : "Join call"}
    </Button>
    {#if setupError}<p class="error" role="alert">{setupError}</p>{/if}
  </form>
</section>

<style>
  section {
    display: grid;
    grid-template-columns: minmax(0, 1fr) minmax(18rem, 1fr);
    gap: clamp(3rem, 9vw, 8.125rem);
    align-items: center;
    min-height: calc(100svh - 6.25rem);
  }

  .eyebrow {
    margin: 0 0 var(--space-3);
    color: var(--accent);
    font-size: 0.76rem;
    font-weight: 800;
    letter-spacing: 0.14em;
    text-transform: uppercase;
  }

  h1 {
    max-width: 47.5rem;
    margin: 0;
    font-size: clamp(2.7rem, 7vw, 6.5rem);
    line-height: 0.95;
    letter-spacing: -0.045em;
    font-weight: 500;
  }

  h1 span, h1 strong { display: block; font: inherit; }
  h1 strong { color: var(--accent); }

  .intro {
    max-width: 36.875rem;
    margin: 1.75rem 0 0;
    color: var(--muted);
    font-size: clamp(1rem, 2vw, 1.25rem);
    line-height: 1.65;
  }

  .install-button { margin-top: var(--space-4); }
  .install-hint { margin: var(--space-2) 0 0; color: var(--muted); font-size: 0.78rem; line-height: 1.5; }

  form {
    padding: clamp(1.25rem, 3vw, 1.75rem);
    border: 1px solid var(--line-soft);
    border-radius: 2px;
    background: var(--surface);
  }

  label {
    display: block;
    margin-bottom: 0.5625rem;
    color: var(--muted);
    font-size: 0.84rem;
    font-weight: 650;
  }

  .form-field + .form-field { margin-top: 1.125rem; }

  .room-entry {
    display: flex;
    gap: var(--space-2);
    align-items: center;
  }

  fieldset {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0.5625rem;
    margin: 1.25rem 0 0;
    padding: 0;
    border: 0;
  }

  legend { margin-bottom: 0.5625rem; color: var(--muted); font-size: 0.84rem; font-weight: 650; }

  .media-choice {
    display: flex;
    gap: var(--space-2);
    align-items: center;
    height: var(--control-height);
    min-height: var(--control-height);
    padding: 0 0.6875rem;
    border-color: var(--accent-border);
    background: var(--accent-soft);
    font-size: 0.78rem;
    font-weight: 680;
    white-space: nowrap;
  }

  .media-choice :global(svg) { width: 1rem; height: 1rem; }
  .media-choice[aria-pressed="false"] { color: var(--danger); border-color: var(--danger-border); background: var(--danger-bg); }
  details { margin-top: 1.25rem; color: var(--muted); font-size: 0.86rem; }
  summary { margin-bottom: 0.9375rem; cursor: pointer; }
  .error { margin: 0.9375rem 0 0; color: var(--danger); font-size: 0.86rem; line-height: 1.5; }

  @media (max-width: 47.5em) {
    section {
      grid-template-columns: 1fr;
      gap: clamp(2rem, 8vw, 2.625rem);
      align-content: center;
      padding: clamp(2.5rem, 10vh, 4rem) 0;
    }

    form { padding: clamp(1.125rem, 5vw, 1.5rem); }
    .room-entry { align-items: stretch; flex-direction: column; }
  }
</style>
