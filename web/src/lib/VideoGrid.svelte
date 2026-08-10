<script lang="ts">
  import VideoTile from "./VideoTile.svelte"
  import type { DeviceType, Peer } from "./types"

  interface Props {
    localStream: MediaStream | null
    participantName: string
    deviceType: DeviceType
    microphoneMuted: boolean
    cameraStopped: boolean
    cameraFacing: VideoFacingModeEnum
    screenSharing: boolean
    displayStream: MediaStream | null
    peers: Map<string, Peer>
    onTogglePeerPlayback: (peerID: string) => void
  }

  let {
    localStream,
    participantName,
    deviceType,
    microphoneMuted,
    cameraStopped,
    cameraFacing,
    screenSharing,
    displayStream,
    peers,
    onTogglePeerPlayback,
  }: Props = $props()
</script>

<div class="grid">
  <VideoTile stream={localStream} name={`${participantName} (You)`} device={deviceType} {microphoneMuted} {cameraStopped} local mirrored={cameraFacing !== "environment"} />
  {#if screenSharing && displayStream}
    <VideoTile stream={displayStream} name="Your screen" screenSharing screenOnly local mirrored={false} />
  {/if}
  {#each [...peers.values()] as peer (peer.id)}
    <VideoTile stream={peer.stream} name={peer.name} device={peer.device} microphoneMuted={peer.microphoneMuted} cameraStopped={peer.cameraStopped} locallyMuted={peer.locallyMuted} onToggleMute={() => onTogglePeerPlayback(peer.id)} />
    {#if peer.screenSharing && peer.screenStream}
      <VideoTile stream={peer.screenStream} name={`${peer.name}'s screen`} screenSharing screenOnly locallyMuted={peer.locallyMuted} />
    {/if}
  {/each}
  {#if peers.size === 0}
    <div class="empty-room">
      <span class="waiting-ring" aria-hidden="true"></span>
      <strong>Waiting for someone to join</strong>
      <span>Share the invite link to start the call.</span>
    </div>
  {/if}
</div>

<style>
  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(min(100%, 420px), 1fr));
    gap: 14px;
  }

  .empty-room {
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 8px;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    min-height: 280px;
    border: 1px solid #242a38;
    border-radius: 20px;
    color: #778095;
    background: #111622;
    text-align: center;
    aspect-ratio: 16 / 10;
  }

  .empty-room strong { color: #cbd0dc; }

  .waiting-ring {
    width: 38px;
    height: 38px;
    margin-bottom: 10px;
    border: 3px solid #2d3444;
    border-top-color: #7d89ff;
    border-radius: 50%;
    animation: spin 1.3s linear infinite;
  }

  @keyframes spin { to { transform: rotate(360deg); } }

  @media (max-width: 760px) {
    .empty-room { min-height: 210px; }
  }

  @media (prefers-reduced-motion: reduce) {
    .waiting-ring { animation: none; }
  }
</style>
