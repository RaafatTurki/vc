<script lang="ts">
  import { Mic, MicOff, RefreshCw, Video, VideoOff } from "@lucide/svelte"
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
    onJoin: (event: SubmitEvent) => void | Promise<void>
  }

  let {
    roomInput = $bindable(),
    signalInput = $bindable(),
    participantName = $bindable(),
    joinWithAudio = $bindable(),
    joinWithVideo = $bindable(),
    joining,
    setupError,
    onJoin,
  }: Props = $props()
</script>

<section aria-labelledby="setupTitle">
  <div>
    <p class="eyebrow">Private by design</p>
    <h1 id="setupTitle">Start a Vivid call</h1>
    <p class="intro">Create a room, share its link, and talk directly from browser to browser.</p>
  </div>

  <form onsubmit={onJoin} novalidate>
    <div class="form-field">
      <label for="nameInput">Your name</label>
      <input id="nameInput" bind:value={participantName} maxlength="32" autocomplete="name" placeholder="How others will see you" required>
    </div>

    <div class="form-field">
      <label for="roomInput">Room ID</label>
      <div class="room-entry">
        <input id="roomInput" bind:value={roomInput} maxlength="6" autocomplete="off" spellcheck="false" required>
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
      <input id="signalInput" bind:value={signalInput} type="url" spellcheck="false" required>
    </details>

    <Button kind="primary" type="submit" disabled={joining}>
      {joining ? "Joining…" : "Join call"}
    </Button>
    {#if setupError}<p class="error" role="alert">{setupError}</p>{/if}
  </form>
</section>

<style>
  section {
    display: grid;
    grid-template-columns: minmax(0, 1.2fr) minmax(320px, 0.8fr);
    gap: clamp(48px, 9vw, 130px);
    align-items: center;
    min-height: calc(100vh - 100px);
  }

  .eyebrow {
    margin: 0 0 12px;
    color: #8f9cff;
    font-size: 0.76rem;
    font-weight: 800;
    letter-spacing: 0.14em;
    text-transform: uppercase;
  }

  h1 {
    max-width: 760px;
    margin: 0;
    font-size: clamp(2.7rem, 7vw, 6.5rem);
    line-height: 0.95;
    letter-spacing: -0.065em;
  }

  .intro {
    max-width: 590px;
    margin: 28px 0 0;
    color: #a9b0c0;
    font-size: clamp(1rem, 2vw, 1.25rem);
    line-height: 1.65;
  }

  form {
    padding: 28px;
    border: 1px solid #242a38;
    border-radius: 24px;
    background: rgb(17 22 34 / 78%);
    box-shadow: 0 26px 80px rgb(0 0 0 / 28%);
    backdrop-filter: blur(20px);
  }

  label {
    display: block;
    margin-bottom: 9px;
    color: #c8cdda;
    font-size: 0.84rem;
    font-weight: 650;
  }

  input {
    width: 100%;
    min-width: 0;
    padding: 13px 14px;
    border: 1px solid #303748;
    border-radius: 12px;
    outline: none;
    color: #f5f7fb;
    background: #0c1019;
  }

  input:focus { border-color: #7583ff; box-shadow: 0 0 0 3px rgb(117 131 255 / 14%); }
  .form-field + .form-field { margin-top: 18px; }

  .room-entry {
    display: flex;
    gap: 10px;
    align-items: center;
  }

  fieldset {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 9px;
    margin: 20px 0 0;
    padding: 0;
    border: 0;
  }

  legend { margin-bottom: 9px; color: #c8cdda; font-size: 0.84rem; font-weight: 650; }

  .media-choice {
    display: flex;
    gap: 8px;
    align-items: center;
    min-height: 42px;
    padding: 0 11px;
    border: 1px solid #354052;
    border-radius: 10px;
    color: #dfe3ed;
    background: #1b2230;
    font-size: 0.78rem;
    font-weight: 680;
  }

  .media-choice :global(svg) { width: 16px; height: 16px; }
  .media-choice[aria-pressed="false"] { color: #ff9aaa; border-color: #59303a; background: #351f27; }
  details { margin-top: 20px; color: #a9b0c0; font-size: 0.86rem; }
  summary { margin-bottom: 15px; cursor: pointer; }
  .error { margin: 15px 0 0; color: #ff8295; font-size: 0.86rem; line-height: 1.5; }

  @media (max-width: 760px) {
    section {
      grid-template-columns: 1fr;
      gap: 42px;
      align-content: center;
      padding: 55px 0;
    }

    form { padding: 20px; }
    .room-entry { align-items: stretch; flex-direction: column; }
  }
</style>
