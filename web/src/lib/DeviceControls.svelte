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
    gap: 12px;
    margin: 18px auto 0;
    padding: 14px;
    border: 1px solid #242a38;
    border-radius: 16px;
    background: rgb(13 17 27 / 72%);
  }

  label { min-width: 0; margin: 0; color: #c8cdda; font-size: 0.84rem; font-weight: 650; }

  label > span {
    display: flex;
    gap: 7px;
    align-items: center;
    margin-bottom: 7px;
  }

  label > span :global(svg) { width: 15px; height: 15px; }

  select {
    width: 100%;
    min-width: 0;
    padding: 10px 36px 10px 11px;
    border: 1px solid #303748;
    border-radius: 10px;
    outline: none;
    overflow: hidden;
    color: #eef0f6;
    background: #111622;
    text-overflow: ellipsis;
  }

  select:focus { border-color: #7583ff; box-shadow: 0 0 0 3px rgb(117 131 255 / 14%); }
  select:disabled { cursor: default; opacity: 0.65; }

  @media (max-width: 760px) {
    .device-controls { grid-template-columns: 1fr; }
  }
</style>
