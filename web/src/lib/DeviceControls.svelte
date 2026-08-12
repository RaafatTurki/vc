<script lang="ts">
  import { Mic, Video } from "@lucide/svelte"
  import type { MediaDeviceOption } from "./types"

  interface Props {
    audioDevices: MediaDeviceOption[]
    videoDevices: MediaDeviceOption[]
    selectedAudioDeviceID?: string
    selectedVideoDeviceID?: string
    switchingAudioDevice: boolean
    switchingVideoDevice: boolean
    onAudioChange: (deviceID: string) => void | Promise<void>
    onVideoChange: (deviceID: string) => void | Promise<void>
  }

  let {
    audioDevices,
    videoDevices,
    selectedAudioDeviceID = $bindable(""),
    selectedVideoDeviceID = $bindable(""),
    switchingAudioDevice,
    switchingVideoDevice,
    onAudioChange,
    onVideoChange,
  }: Props = $props()
</script>

{#if audioDevices.length || videoDevices.length}
  <div class="device-controls" aria-label="Media devices">
    {#if audioDevices.length}
      <label>
        <span><Mic aria-hidden="true" /> Microphone</span>
        <select bind:value={selectedAudioDeviceID} disabled={switchingAudioDevice || audioDevices.length < 2} onchange={() => onAudioChange(selectedAudioDeviceID)}>
          {#each audioDevices as device (device.deviceId)}
            <option value={device.deviceId}>{device.label}</option>
          {/each}
        </select>
      </label>
    {/if}
    {#if videoDevices.length}
      <label>
        <span><Video aria-hidden="true" /> Camera</span>
        <select bind:value={selectedVideoDeviceID} disabled={switchingVideoDevice || videoDevices.length < 2} onchange={() => onVideoChange(selectedVideoDeviceID)}>
          {#each videoDevices as device (device.deviceId)}
            <option value={device.deviceId}>{device.label}</option>
          {/each}
        </select>
      </label>
    {/if}
  </div>
{/if}

<style>
  .device-controls {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: var(--space-3);
    margin: 1.125rem auto 0;
    padding: var(--space-3);
    border: 1px solid var(--line);
    border-radius: 2px;
    background: var(--surface);
  }

  label { min-width: 0; margin: 0; color: var(--muted); font-size: 0.84rem; font-weight: 650; }

  label > span {
    display: flex;
    gap: 0.4375rem;
    align-items: center;
    margin-bottom: 0.4375rem;
  }

  label > span :global(svg) { width: 0.9375rem; height: 0.9375rem; }

  select {
    width: 100%;
    height: var(--control-height);
    min-height: var(--control-height);
    min-width: 0;
    padding: 0.625rem 2.25rem 0.625rem 0.6875rem;
    border: 1px solid var(--line-strong);
    border-radius: 2px;
    outline: none;
    overflow: hidden;
    color: var(--ink);
    background: var(--bg);
    text-overflow: ellipsis;
  }

  select:focus { border-color: var(--accent); box-shadow: 0 0 0 3px var(--accent-12); }
  select:disabled { cursor: default; opacity: 0.65; }

  @media (max-width: 47.5em) {
    .device-controls { grid-template-columns: 1fr; }
  }
</style>
