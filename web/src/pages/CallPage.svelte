<script lang="ts">
  import CallControls from "../lib/CallControls.svelte"
  import CallHeader from "../lib/CallHeader.svelte"
  import DeviceControls from "../lib/DeviceControls.svelte"
  import VideoGrid from "../lib/VideoGrid.svelte"
  import type { DeviceType, MediaDeviceOption, Peer } from "../lib/types"

  interface Props {
    roomID: string
    copyLabel: string
    localStream: MediaStream | null
    participantName: string
    deviceType: DeviceType
    microphoneMuted: boolean
    cameraStopped: boolean
    cameraFacing: VideoFacingModeEnum
    screenSharing: boolean
    displayStream: MediaStream | null
    peers: Map<string, Peer>
    audioDevices: MediaDeviceOption[]
    videoDevices: MediaDeviceOption[]
    selectedAudioDeviceID?: string
    selectedVideoDeviceID?: string
    switchingAudioDevice: boolean
    switchingVideoDevice: boolean
    canShareScreen: boolean
    sharingScreen: boolean
    canSwitchCamera: boolean
    switchingCamera: boolean
    cameraError: string
    onCopy: () => void | Promise<void>
    onTogglePeerPlayback: (peerID: string) => void
    onAudioDeviceChange: (deviceID: string) => void | Promise<void>
    onVideoDeviceChange: (deviceID: string) => void | Promise<void>
    onToggleMicrophone: () => void | Promise<void>
    onToggleCamera: () => void | Promise<void>
    onToggleScreenShare: () => void | Promise<void>
    onSwitchCamera: () => void | Promise<void>
    onLeave: () => void
  }

  let {
    roomID,
    copyLabel,
    localStream,
    participantName,
    deviceType,
    microphoneMuted,
    cameraStopped,
    cameraFacing,
    screenSharing,
    displayStream,
    peers,
    audioDevices,
    videoDevices,
    selectedAudioDeviceID = $bindable(),
    selectedVideoDeviceID = $bindable(),
    switchingAudioDevice,
    switchingVideoDevice,
    canShareScreen,
    sharingScreen,
    canSwitchCamera,
    switchingCamera,
    cameraError,
    onCopy,
    onTogglePeerPlayback,
    onAudioDeviceChange,
    onVideoDeviceChange,
    onToggleMicrophone,
    onToggleCamera,
    onToggleScreenShare,
    onSwitchCamera,
    onLeave,
  }: Props = $props()
</script>

<section aria-label="Vivid call">
  <CallHeader {roomID} {copyLabel} {onCopy} />
  <VideoGrid
    {localStream}
    {participantName}
    {deviceType}
    {microphoneMuted}
    {cameraStopped}
    {cameraFacing}
    {screenSharing}
    {displayStream}
    {peers}
    {onTogglePeerPlayback}
  />
  <DeviceControls
    {audioDevices}
    {videoDevices}
    bind:selectedAudioDeviceID
    bind:selectedVideoDeviceID
    {switchingAudioDevice}
    {switchingVideoDevice}
    onAudioChange={onAudioDeviceChange}
    onVideoChange={onVideoDeviceChange}
  />
  <CallControls
    {microphoneMuted}
    {cameraStopped}
    {canShareScreen}
    {screenSharing}
    {sharingScreen}
    {canSwitchCamera}
    {switchingCamera}
    {onToggleMicrophone}
    {onToggleCamera}
    {onToggleScreenShare}
    {onSwitchCamera}
    {onLeave}
  />
  {#if cameraError}<p role="alert">{cameraError}</p>{/if}
</section>

<style>
  section { padding-top: clamp(50px, 8vh, 90px); }
  p { margin: 12px auto 0; color: #ff9aaa; font-size: 0.84rem; text-align: center; }
</style>
