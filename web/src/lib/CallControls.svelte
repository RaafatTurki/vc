<script lang="ts">
  import { LogOut, Mic, MicOff, ScreenShare, ScreenShareOff, SwitchCamera, Video, VideoOff } from "@lucide/svelte"

  interface Props {
    microphoneMuted: boolean
    cameraStopped: boolean
    canShareScreen: boolean
    screenSharing: boolean
    sharingScreen: boolean
    canSwitchCamera: boolean
    switchingCamera: boolean
    onToggleMicrophone: () => void | Promise<void>
    onToggleCamera: () => void | Promise<void>
    onToggleScreenShare: () => void | Promise<void>
    onSwitchCamera: () => void | Promise<void>
    onLeave: () => void
  }

  let {
    microphoneMuted,
    cameraStopped,
    canShareScreen,
    screenSharing,
    sharingScreen,
    canSwitchCamera,
    switchingCamera,
    onToggleMicrophone,
    onToggleCamera,
    onToggleScreenShare,
    onSwitchCamera,
    onLeave,
  }: Props = $props()
</script>

<div class="controls" aria-label="Call controls">
  <button type="button" aria-pressed={microphoneMuted} onclick={onToggleMicrophone}>
    {#if microphoneMuted}<MicOff aria-hidden="true" />{:else}<Mic aria-hidden="true" />{/if}
    <span>{microphoneMuted ? "Unmute" : "Mute"}</span>
  </button>
  <button type="button" aria-pressed={cameraStopped} onclick={onToggleCamera}>
    {#if cameraStopped}<VideoOff aria-hidden="true" />{:else}<Video aria-hidden="true" />{/if}
    <span>{cameraStopped ? "Start video" : "Stop video"}</span>
  </button>
  {#if canShareScreen}
    <button type="button" aria-pressed={screenSharing} disabled={sharingScreen} onclick={onToggleScreenShare}>
      {#if screenSharing}<ScreenShareOff aria-hidden="true" />{:else}<ScreenShare aria-hidden="true" />{/if}
      <span>{sharingScreen ? "Starting…" : screenSharing ? "Stop sharing" : "Share screen"}</span>
    </button>
  {/if}
  {#if canSwitchCamera}
    <button type="button" disabled={switchingCamera} onclick={onSwitchCamera}>
      <SwitchCamera aria-hidden="true" />
      <span>{switchingCamera ? "Switching…" : "Switch camera"}</span>
    </button>
  {/if}
  <button class="danger" type="button" onclick={onLeave}>
    <LogOut aria-hidden="true" />
    <span>Leave</span>
  </button>
</div>

<style>
  .controls {
    position: sticky;
    bottom: 20px;
    z-index: 5;
    display: flex;
    gap: 10px;
    align-items: center;
    justify-content: center;
    width: fit-content;
    margin: 24px auto 0;
    padding: 9px;
    border: 1px solid #2c3342;
    border-radius: 18px;
    background: rgb(13 17 27 / 82%);
    box-shadow: 0 18px 50px rgb(0 0 0 / 38%);
    backdrop-filter: blur(18px);
  }

  button {
    display: flex;
    gap: 7px;
    align-items: center;
    min-height: 46px;
    padding: 0 15px;
    border: 0;
    border-radius: 11px;
    color: #e9ebf2;
    background: #252b39;
    font-weight: 650;
  }

  button[aria-pressed="true"] { color: #ff9aaa; background: #42232d; }
  button.danger { background: #d83c58; }
  button :global(svg) { width: 18px; height: 18px; }

  @media (max-width: 760px) {
    .controls { gap: 4px; width: 100%; padding: 7px; }
    button { flex: 1; flex-direction: column; gap: 2px; min-width: 0; padding: 7px 4px; font-size: 0.72rem; }
  }
</style>
