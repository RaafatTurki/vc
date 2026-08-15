<script lang="ts">
  import { AudioLines, LogOut, MessageCircle, Mic, MicOff, ScreenShare, ScreenShareOff, SwitchCamera, Video, VideoOff } from "@lucide/svelte"
  import Popup from "./Popup.svelte"

  interface Props {
    microphoneMuted: boolean
    noiseCancellationEnabled: boolean
    cameraStopped: boolean
    canShareScreen: boolean
    screenSharing: boolean
    sharingScreen: boolean
    canSwitchCamera: boolean
    switchingCamera: boolean
    onToggleMicrophone: () => void | Promise<void>
    onToggleNoiseCancellation: () => void | Promise<void>
    onToggleCamera: () => void | Promise<void>
    onToggleScreenShare: () => void | Promise<void>
    onStartScreenShare: (frameRate: 30 | 60 | 120) => void | Promise<void>
    onSwitchCamera: () => void | Promise<void>
    onLeave: () => void
    chatOpen: boolean
    unreadChatMessages: number
    onToggleChat: () => void
  }

  let {
    microphoneMuted,
    noiseCancellationEnabled,
    cameraStopped,
    canShareScreen,
    screenSharing,
    sharingScreen,
    canSwitchCamera,
    switchingCamera,
    onToggleMicrophone,
    onToggleNoiseCancellation,
    onToggleCamera,
    onToggleScreenShare,
    onStartScreenShare,
    onSwitchCamera,
    onLeave,
    chatOpen,
    unreadChatMessages,
    onToggleChat,
  }: Props = $props()

  let frameRatePopupOpen = $state(false)

  async function chooseFrameRate(value: string) {
    frameRatePopupOpen = false
    await onStartScreenShare(Number(value) as 30 | 60 | 120)
  }
</script>

<div class="controls" aria-label="Call controls">
  <button type="button" aria-label={microphoneMuted ? "Unmute microphone" : "Mute microphone"} aria-pressed={microphoneMuted} onclick={onToggleMicrophone}>
    {#if microphoneMuted}<MicOff aria-hidden="true" />{:else}<Mic aria-hidden="true" />{/if}
    <span>{microphoneMuted ? "Unmute" : "Mute"}</span>
  </button>
  <button type="button" aria-label={cameraStopped ? "Start camera" : "Stop camera"} aria-pressed={cameraStopped} onclick={onToggleCamera}>
    {#if cameraStopped}<VideoOff aria-hidden="true" />{:else}<Video aria-hidden="true" />{/if}
    <span>{cameraStopped ? "Start video" : "Stop video"}</span>
  </button>
  <button type="button" aria-label={noiseCancellationEnabled ? "Disable noise cancellation" : "Enable noise cancellation"} aria-pressed={!noiseCancellationEnabled} onclick={onToggleNoiseCancellation}>
    <AudioLines aria-hidden="true" />
    <span>{noiseCancellationEnabled ? "Noise canceling" : "Noise off"}</span>
  </button>
  {#if canShareScreen}
    <div class="screen-tools">
      {#if screenSharing}
        <button type="button" aria-label="Stop screen sharing" aria-pressed="true" disabled={sharingScreen} onclick={onToggleScreenShare}>
          <ScreenShareOff aria-hidden="true" />
          <span>Stop sharing</span>
        </button>
      {:else}
        <button type="button" aria-label="Choose screen share frame rate" disabled={sharingScreen} onclick={() => frameRatePopupOpen = true}>
          <ScreenShare aria-hidden="true" />
          <span>Share screen</span>
        </button>
      {/if}
    </div>
  {/if}
  {#if canSwitchCamera}
    <button type="button" aria-label={switchingCamera ? "Switching camera" : "Switch camera"} disabled={switchingCamera} onclick={onSwitchCamera}>
      <SwitchCamera aria-hidden="true" />
      <span>{switchingCamera ? "Switching…" : "Switch camera"}</span>
    </button>
  {/if}
  <button class="danger" type="button" aria-label="Leave call" onclick={onLeave}>
    <LogOut aria-hidden="true" />
    <span>Leave</span>
  </button>
  <button type="button" aria-label={chatOpen ? "Close chat" : "Open chat"} aria-expanded={chatOpen} onclick={onToggleChat}>
    <MessageCircle aria-hidden="true" />
    <span>{chatOpen ? "Hide chat" : "Chat"}{unreadChatMessages && !chatOpen ? ` · ${unreadChatMessages}` : ""}</span>
  </button>
</div>

<Popup
  open={frameRatePopupOpen}
  title="Screen share quality"
  options={[{ value: "30", label: "Share at 30 FPS" }, { value: "60", label: "Share at 60 FPS" }, { value: "120", label: "Share at 120 FPS" }]}
  onSelect={chooseFrameRate}
  onClose={() => frameRatePopupOpen = false}
/>

<style>
  .controls {
    display: flex;
    gap: var(--space-2);
    align-items: center;
    justify-content: center;
    width: fit-content;
    max-width: 100%;
    overflow-x: auto;
    margin: var(--space-4) auto 0;
    padding: var(--space-2);
    border: 1px solid rgb(var(--white-rgb) / 14%);
    border-radius: 2px;
    background: var(--surface);
  }

  button {
    display: flex;
    gap: 0.4375rem;
    align-items: center;
    height: var(--control-height);
    min-height: var(--control-height);
    padding: 0 0.9375rem;
    font-weight: 650;
    white-space: nowrap;
  }

  button.danger,
  button[aria-pressed="true"] { border-color: var(--danger-border); color: var(--danger); background: var(--danger-bg); }
  button :global(svg) { width: 1.125rem; height: 1.125rem; }

  .screen-tools {
    display: flex;
    gap: var(--space-1);
    align-items: center;
  }

  @media (hover: hover) and (pointer: fine) {
    button:hover:not(:disabled) { color: var(--accent); background: var(--accent-subtle); }
    button.danger:hover:not(:disabled),
    button[aria-pressed="true"]:hover:not(:disabled) {
      color: var(--danger);
      border-color: var(--danger-border);
      background: var(--danger-bg);
    }
  }

  @media (max-width: 47.5em) {
    .controls { gap: var(--space-1); width: fit-content; padding: var(--space-2); justify-content: flex-start; }
    button { flex: 0 0 auto; padding: 0 0.625rem; font-size: 0.72rem; }
  }

  @media (max-width: 38.75em) {
    .controls { width: fit-content; max-width: calc(100% - 0.25rem); justify-content: center; overflow-x: visible; }
    button { width: var(--control-height); padding: 0; justify-content: center; }
    button span { display: none; }
  }
</style>
