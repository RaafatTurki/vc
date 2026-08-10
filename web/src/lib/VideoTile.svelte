<script lang="ts">
  import { onDestroy, onMount } from "svelte"
  import { Maximize2, MicOff, Minimize2, Monitor, ScreenShare, Smartphone, VideoOff, Volume2, VolumeX } from "@lucide/svelte"
  import type { DeviceType } from "./types"

  interface Props {
    stream: MediaStream | null
    name: string
    device?: DeviceType
    microphoneMuted?: boolean
    cameraStopped?: boolean
    screenSharing?: boolean
    screenOnly?: boolean
    local?: boolean
    mirrored?: boolean
    locallyMuted?: boolean
    onToggleMute?: () => void
  }

  let {
    stream,
    name,
    device = "computer",
    microphoneMuted = false,
    cameraStopped = false,
    screenSharing = false,
    screenOnly = false,
    local = false,
    mirrored = local,
    locallyMuted = false,
    onToggleMute = () => {},
  }: Props = $props()
  let card: HTMLElement
  let video: HTMLVideoElement
  let fullscreen = $state(false)
  let speaking = $state(false)
  let volume = $state(1)

  onMount(() => {
    document.addEventListener("fullscreenchange", updateFullscreenState)
    document.addEventListener("webkitfullscreenchange", updateFullscreenState)
    video?.addEventListener("webkitbeginfullscreen", onVideoFullscreenStart)
    video?.addEventListener("webkitendfullscreen", onVideoFullscreenEnd)
  })

  onDestroy(() => {
    document.removeEventListener("fullscreenchange", updateFullscreenState)
    document.removeEventListener("webkitfullscreenchange", updateFullscreenState)
    video?.removeEventListener("webkitbeginfullscreen", onVideoFullscreenStart)
    video?.removeEventListener("webkitendfullscreen", onVideoFullscreenEnd)
  })

  $effect(() => {
    if (video && video.srcObject !== stream) {
      video.srcObject = stream ?? null
      if (stream) video.play().catch(() => {})
    }
  })

  $effect(() => {
    if (video) video.volume = volume
  })

  $effect(() => {
    const currentStream = stream
    if (!currentStream || screenOnly || currentStream.getAudioTracks().length === 0) {
      speaking = false
      return
    }

    const AudioContextClass = window.AudioContext || window.webkitAudioContext
    if (!AudioContextClass) return
    const context = new AudioContextClass()
    const analyser = context.createAnalyser()
    const source = context.createMediaStreamSource(currentStream)
    analyser.fftSize = 512
    analyser.smoothingTimeConstant = 0.65
    source.connect(analyser)
    context.resume().catch(() => {})

    const samples = new Float32Array(analyser.fftSize)
    let quietTicks = 0
    const interval = window.setInterval(() => {
      analyser.getFloatTimeDomainData(samples)
      let sum = 0
      for (const sample of samples) sum += sample * sample
      const rms = Math.sqrt(sum / samples.length)
      if (rms > 0.025) {
        quietTicks = 0
        speaking = true
      } else if (++quietTicks >= 4) {
        speaking = false
      }
    }, 80)

    return () => {
      window.clearInterval(interval)
      source.disconnect()
      speaking = false
      context.close().catch(() => {})
    }
  })

  function updateFullscreenState() {
    fullscreen = (document.fullscreenElement || document.webkitFullscreenElement) === card
  }

  function onVideoFullscreenStart() {
    fullscreen = true
  }

  function onVideoFullscreenEnd() {
    fullscreen = false
  }

  async function toggleFullscreen() {
    try {
      const activeElement = document.fullscreenElement || document.webkitFullscreenElement
      if (activeElement === card) {
        if (document.exitFullscreen) await document.exitFullscreen()
        else document.webkitExitFullscreen?.()
      } else if (card.requestFullscreen) {
        await card.requestFullscreen()
      } else if (card.webkitRequestFullscreen) {
        card.webkitRequestFullscreen()
      } else {
        video?.webkitEnterFullscreen?.()
      }
    } catch (error) {
      console.error("Could not toggle fullscreen video", error)
    }
  }
</script>

<article bind:this={card} class:local-card={local} class:remote-card={!local} class:mirrored class:screen-sharing={screenSharing} class:speaking class="video-card">
  <video bind:this={video} autoplay muted={local || locallyMuted} playsinline></video>
  <button
    class="tile-fullscreen"
    type="button"
    aria-label={fullscreen ? `Exit fullscreen for ${name}` : `View ${name} fullscreen`}
    aria-pressed={fullscreen}
    onclick={toggleFullscreen}
  >
    {#if fullscreen}
      <Minimize2 aria-hidden="true" />
    {:else}
      <Maximize2 aria-hidden="true" />
    {/if}
  </button>
  <div class="video-meta" aria-label={`${device === "mobile" ? "Mobile" : "PC"}${microphoneMuted ? ", microphone muted" : ""}`}>
    {#if !screenOnly}
      <span class="video-badge">
        {#if device === "mobile"}
          <Smartphone class="badge-icon" aria-hidden="true" />
          Mobile
        {:else}
          <Monitor class="badge-icon" aria-hidden="true" />
          PC
        {/if}
      </span>
      {#if microphoneMuted}
        <span class="video-badge muted-badge">
          <MicOff class="badge-icon" aria-hidden="true" />
          Mic off
        </span>
      {/if}
      {#if cameraStopped}
        <span class="video-badge muted-badge">
          <VideoOff class="badge-icon" aria-hidden="true" />
          Camera off
        </span>
      {/if}
    {/if}
    {#if screenSharing}
      <span class="video-badge sharing-badge">
        <ScreenShare class="badge-icon" aria-hidden="true" />
        Sharing
      </span>
    {/if}
  </div>
  <div class="video-label">{name}</div>
  {#if !local && !screenOnly}
    <label class="tile-volume" aria-label={`Volume for ${name}`}>
      <Volume2 aria-hidden="true" />
      <input bind:value={volume} type="range" min="0" max="1" step="0.05" disabled={locallyMuted}>
      <span>{Math.round(volume * 100)}%</span>
    </label>
    <button
      class="tile-mute"
      type="button"
      aria-pressed={locallyMuted}
      aria-label={locallyMuted ? `Unmute ${name} for you` : `Mute ${name} for you`}
      onclick={onToggleMute}
    >
      {#if locallyMuted}
        <VolumeX class="tile-mute-icon" aria-hidden="true" />
      {:else}
        <Volume2 class="tile-mute-icon" aria-hidden="true" />
      {/if}
      {locallyMuted ? "Unmute" : "Mute audio"}
    </button>
  {/if}
</article>

<style>
  .video-card {
    position: relative;
    isolation: isolate;
    overflow: hidden;
    min-height: 280px;
    border: 1px solid #242a38;
    border-radius: 20px;
    background: #111622;
    clip-path: inset(0 round 20px);
    aspect-ratio: 16 / 10;
  }

  .video-card.speaking::after {
    position: absolute;
    z-index: 10;
    inset: 0;
    border: 3px solid #21d6a1;
    border-radius: inherit;
    box-shadow: inset 0 0 18px rgb(33 214 161 / 18%);
    pointer-events: none;
    content: "";
  }

  video {
    display: block;
    width: 100%;
    height: 100%;
    border-radius: inherit;
    background: #090c12;
    object-fit: cover;
  }

  .local-card.mirrored video { transform: scaleX(-1); }
  .screen-sharing video { object-fit: contain; }

  .video-card:fullscreen,
  .video-card:-webkit-full-screen {
    width: 100vw;
    height: 100vh;
    min-height: 100vh;
    border: 0;
    border-radius: 0;
    clip-path: none;
    aspect-ratio: auto;
  }

  .tile-fullscreen {
    position: absolute;
    top: 14px;
    right: 14px;
    display: grid;
    width: 34px;
    height: 34px;
    padding: 8px;
    border: 1px solid rgb(255 255 255 / 12%);
    border-radius: 9px;
    color: #e9ebf2;
    background: rgb(25 30 43 / 82%);
    backdrop-filter: blur(12px);
  }

  .tile-fullscreen:hover { background: rgb(42 49 67 / 92%); }
  .tile-fullscreen :global(svg) { width: 100%; height: 100%; }

  .video-label {
    position: absolute;
    bottom: 14px;
    left: 14px;
    max-width: calc(100% - 150px);
    overflow: hidden;
    padding: 6px 10px;
    border: 1px solid rgb(255 255 255 / 12%);
    border-radius: 8px;
    background: rgb(5 7 11 / 66%);
    font-size: 0.76rem;
    font-weight: 720;
    text-overflow: ellipsis;
    white-space: nowrap;
    backdrop-filter: blur(12px);
  }

  .video-meta {
    position: absolute;
    top: 14px;
    left: 14px;
    display: flex;
    gap: 6px;
  }

  .video-badge {
    display: inline-flex;
    gap: 5px;
    align-items: center;
    padding: 5px 8px;
    border: 1px solid rgb(255 255 255 / 12%);
    border-radius: 999px;
    color: #d9deea;
    background: rgb(5 7 11 / 66%);
    font-size: 0.68rem;
    font-weight: 750;
    backdrop-filter: blur(12px);
  }

  .video-badge :global(.badge-icon) { width: 13px; height: 13px; }
  .muted-badge { color: #ffadb9; border-color: rgb(255 100 124 / 30%); background: rgb(74 25 37 / 78%); }
  .sharing-badge { color: #aeb7ff; border-color: rgb(117 131 255 / 32%); background: rgb(34 40 91 / 82%); }

  .tile-mute {
    position: absolute;
    right: 14px;
    bottom: 14px;
    display: inline-flex;
    gap: 6px;
    align-items: center;
    min-height: 31px;
    padding: 0 10px;
    border: 1px solid rgb(255 255 255 / 12%);
    border-radius: 8px;
    color: #e9ebf2;
    background: rgb(25 30 43 / 82%);
    font-size: 0.72rem;
    font-weight: 700;
    backdrop-filter: blur(12px);
  }

  .tile-mute[aria-pressed="true"] { color: #ffadb9; border-color: rgb(255 100 124 / 30%); background: rgb(74 25 37 / 88%); }
  .tile-mute :global(.tile-mute-icon) { width: 14px; height: 14px; }

  .tile-volume {
    position: absolute;
    right: 14px;
    bottom: 54px;
    display: flex;
    gap: 7px;
    align-items: center;
    width: 178px;
    margin: 0;
    padding: 7px 9px;
    border: 1px solid rgb(255 255 255 / 12%);
    border-radius: 9px;
    color: #dfe3ed;
    background: rgb(25 30 43 / 86%);
    font-size: 0.68rem;
    backdrop-filter: blur(12px);
  }

  .tile-volume > :global(svg) { flex: 0 0 auto; width: 14px; height: 14px; }
  .tile-volume input { flex: 1; width: auto; min-width: 0; height: 4px; padding: 0; accent-color: #7d89ff; }
  .tile-volume span { width: 30px; text-align: right; }
  .tile-volume:has(input:disabled) { opacity: 0.55; }

  @media (max-width: 760px) {
    .video-card { min-height: 210px; }
  }
</style>
