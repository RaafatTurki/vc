<script lang="ts">
  import { Mic, Video } from "@lucide/svelte"
  import Dropdown from "./Dropdown.svelte"
  import type { MediaDeviceOption } from "./types"

  interface Props {
    audioDevices: MediaDeviceOption[]
    videoDevices: MediaDeviceOption[]
    selectedAudioDeviceID: string
    selectedVideoDeviceID: string
    switchingAudioDevice: boolean
    switchingVideoDevice: boolean
    onAudioChange: (deviceID: string) => void | Promise<void>
    onVideoChange: (deviceID: string) => void | Promise<void>
  }

  let {
    audioDevices,
    videoDevices,
    selectedAudioDeviceID,
    selectedVideoDeviceID,
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
        <Dropdown
          options={audioDevices.map(device => ({ value: device.deviceId, label: device.label }))}
          value={selectedAudioDeviceID}
          label="Microphone"
          disabled={switchingAudioDevice || audioDevices.length < 2}
          onchange={onAudioChange}
        />
      </label>
    {/if}
    {#if videoDevices.length}
      <label>
        <span><Video aria-hidden="true" /> Camera</span>
        <Dropdown
          options={videoDevices.map(device => ({ value: device.deviceId, label: device.label }))}
          value={selectedVideoDeviceID}
          label="Camera"
          disabled={switchingVideoDevice || videoDevices.length < 2}
          onchange={onVideoChange}
        />
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

  @media (max-width: 47.5em) {
    .device-controls { grid-template-columns: 1fr; }
  }
</style>
